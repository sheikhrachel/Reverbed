package reverbed_service

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hajimehoshi/go-mp3"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/model"
)

const (
	slowFilter   = "atempo=0.85"
	spedUpFilter = "atempo=2.0"
)

func (r *ReverbedService) SpeedFilter(cc call.Call, editType model.EditType, fileBytes *[]byte, uniqueID uuid.UUID) (mp3bytes []byte, err error) {
	if fileBytes == nil || len(*fileBytes) == 0 {
		return nil, errors.New("no file provided")
	}
	cc.InfoF("Editing audio with type: %s", editType)
	cc.InfoF("File bytes: %d", len(*fileBytes))
	decoder, err := mp3.NewDecoder(bytes.NewReader(*fileBytes))
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	var (
		tempDir, outputDir = generateUniqueDirs(uniqueID)
		tempWavFilePath    = filepath.Join(tempDir, "temp_audio.wav")
		outputMp3FilePath  = filepath.Join(outputDir, "output_audio.mp3")
		filter             string
	)
	defer os.RemoveAll(tempDir)
	defer os.RemoveAll(outputDir)
	tempWavFile, err := os.Create(tempWavFilePath)
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	defer tempWavFile.Close()
	writeTempWavFile(cc, tempWavFile, decoder)
	outputMp3File, err := os.Create(outputMp3FilePath)
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	defer outputMp3File.Close()
	switch editType {
	case model.EditTypeSlowed:
		filter = slowFilter
	case model.EditTypeSpedUp:
		filter = spedUpFilter
	}
	ffmpegHelpCmd := exec.Command("ffmpeg", "-h")
	if err = ffmpegHelpCmd.Run(); errutil.HandleError(cc, err) {
		return nil, err
	}
	ffmpegCmd := exec.Command("ffmpeg", "-i", tempWavFilePath, "-filter:a", filter, "-y", outputMp3FilePath)
	if err = ffmpegCmd.Run(); errutil.HandleError(cc, err) {
		return nil, err
	}
	return os.ReadFile(outputMp3FilePath)
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
func generateUniqueDirs(uniqueID uuid.UUID) (tempDir, outputDir string) {
	tempDir = "temp_audio_" + uniqueID.String()
	outputDir = "output_audio_" + uniqueID.String()
	os.MkdirAll(tempDir, os.ModePerm)
	os.MkdirAll(outputDir, os.ModePerm)
	return tempDir, outputDir
}
