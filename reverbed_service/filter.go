package reverbed_service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/model"
)

const (
	ffmpegHermitPath = "./bin/ffmpeg"
	inputFlag        = "-i"
	filterFlag       = "-filter:a"
	yesFlag          = "-y"
)

func (r *ReverbedService) Filter(
	cc call.Call,
	editType model.EditType,
	MP3File *multipart.FileHeader,
	tempMP3FilePath, outputMP3FilePath string,
) (modifiedMP3File multipart.File, err error) {
	if MP3File == nil {
		return nil, errors.New("no file provided")
	}
	var file multipart.File
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
	modifiedMP3File, err = os.Create(outputMP3FilePath)
	if errutil.HandleError(cc, err) {
		cc.InfoF("Error creating output mp3 file: %+v", err)
		modifiedMP3File.Close()
		return nil, err
	}
	if err = exec.Command(ffmpegHermitPath, inputFlag, tempMP3FilePath, filterFlag, editType.GetFilter(), yesFlag, outputMP3FilePath).Run(); errutil.HandleError(cc, err) {
		cc.InfoF("Error running ffmpeg command: %+v", err)
		modifiedMP3File.Close()
		return nil, err
	}
	return modifiedMP3File, nil
}
