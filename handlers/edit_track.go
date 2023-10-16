package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	timeout "github.com/s-wijaya/gin-timeout"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/model"
)

const (
	editTypeParam     = "edit_type"
	fileParam         = "file"
	audioDir          = "audio"
	inputFile         = "input-%s.mp3"
	outputFile        = "output-%s.mp3"
	contentTypeHeader = "Content-Type"
	audioContentType  = "audio/mpeg"
)

// EditTrack is used to edit a track
// Path - /edit_track/:edit_type
func (h *Handler) EditTrack(c *gin.Context) {
	timeout.APIWrapper(c, func(c *gin.Context) (int, interface{}) {
		cc := call.New(h.appEnv, h.appRegion)
		editType := c.Param(editTypeParam)
		mp3File, err := c.FormFile(fileParam)
		if errutil.HandleError(cc, err) {
			cc.InfoF("No mp3 file provided")
			return StatusBadRequest, err.Error()
		}
		editTypeInt, err := strconv.Atoi(editType)
		if errutil.HandleError(cc, err) {
			cc.InfoF("Invalid edit type: %s", editType)
			return StatusBadRequest, err.Error()
		}
		var (
			editedAudio       multipart.File
			uniqueID          = uuid.New().String()[0:8]
			tempMP3FilePath   = filepath.Join(audioDir, fmt.Sprintf(inputFile, uniqueID))
			outputMp3FilePath = filepath.Join(audioDir, fmt.Sprintf(outputFile, uniqueID))
		)
		editedAudio, err = h.reverbedService.Filter(cc, model.GetEditType(editTypeInt), mp3File, tempMP3FilePath, outputMp3FilePath)
		if errutil.HandleError(cc, err) {
			cc.InfoF("Error editing audio for id: %+v: %+v", uniqueID, err)
			return StatusInternalServerError, err.Error()
		}
		defer editedAudio.Close()
		c.Header(contentTypeHeader, audioContentType)
		io.Copy(c.Writer, editedAudio)
		go cleanupTempFiles(tempMP3FilePath, outputMp3FilePath)
		return StatusOK, nil
	})
}

func cleanupTempFiles(tempMP3FilePath, outputMp3FilePath string) {
	os.Remove(tempMP3FilePath)
	os.Remove(outputMp3FilePath)
}
