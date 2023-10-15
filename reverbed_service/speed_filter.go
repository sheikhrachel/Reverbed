package reverbed_service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hajimehoshi/go-mp3"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/model"
)

func (r *ReverbedService) SpeedFilter(
	cc call.Call,
	editType model.EditType,
	MP3File *multipart.FileHeader,
) (modifiedMP3File multipart.File, err error) {
	if MP3File == nil {
		return nil, errors.New("no file provided")
	}
	var (
		file               multipart.File
		tempDir, outputDir = generateUniqueDirs()
		tempMP3FilePath    = filepath.Join(tempDir, "temp_audio.mp3")
		outputMp3FilePath  = filepath.Join(outputDir, "output_audio.mp3")
		filter             = editType.GetFilter()
	)
	defer os.RemoveAll(tempDir)
	defer os.RemoveAll(outputDir)
	file, err = MP3File.Open()
	if errutil.HandleError(cc, err) {
		cc.InfoF("Error opening mp3 file: %+v", err)
		return nil, err
	}
	defer file.Close()
	tempMP3File, err := os.Create(tempMP3FilePath)
	if errutil.HandleError(cc, err) {
		cc.InfoF("Error creating temp mp3 file: %+v", err)
		return nil, err
	}
	defer tempMP3File.Close()
	if _, err = io.Copy(tempMP3File, file); errutil.HandleError(cc, err) {
		cc.InfoF("Error copying mp3 file to temp mp3 file: %+v", err)
		return nil, err
	}
	outputMp3File, err := os.Create(outputMp3FilePath)
	if errutil.HandleError(cc, err) {
		cc.InfoF("Error creating output mp3 file: %+v", err)
		return nil, err
	}
	ffmpegCmd := exec.Command("ffmpeg", "-i", tempMP3FilePath, "-filter:a", filter, "-y", outputMp3FilePath)
	bytesOut, err := ffmpegCmd.CombinedOutput()
	if errutil.HandleError(cc, err) {
		cc.InfoF("Error running ffmpeg command: bytes: %+v, err: %+v", string(bytesOut), err)
		outputMp3File.Close()
		return nil, err
	}
	cc.InfoF("Successfully ran ffmpeg command: %s", string(bytesOut))
	//if err = ffmpegCmd.Run(); errutil.HandleError(cc, err) {
	//	return nil, err
	//}
	return outputMp3File, nil
}

func writeTempWavFile(cc call.Call, tempWavFile *os.File, decoder *mp3.Decoder) {
	var (
		n   int
		err error
		buf = make([]byte, 4096)
	)
	for {
		n, err = decoder.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else if errutil.HandleError(cc, err) {
				return
			}
		}
		if n == 0 {
			break
		}
		if _, err = tempWavFile.Write(buf[:n]); errutil.HandleError(cc, err) {
			return
		}
	}
}

// Generate unique directory names based on a unique identifier.
func generateUniqueDirs() (tempDir, outputDir string) {
	uniqueID := uuid.New().String()
	tempDir = "temp_audio_" + uniqueID
	outputDir = "output_audio_" + uniqueID
	os.MkdirAll(tempDir, os.ModePerm)
	os.MkdirAll(outputDir, os.ModePerm)
	return tempDir, outputDir
}
