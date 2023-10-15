package handlers

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	timeout "github.com/s-wijaya/gin-timeout"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/model"
)

// EditTrack is used to edit a track
// Path - /edit_track/:edit_type
func (h *Handler) EditTrack(c *gin.Context) {
	timeout.APIWrapper(c, func(c *gin.Context) (int, interface{}) {
		cc := call.New(h.appEnv, h.appRegion)
		editType := c.Param("edit_type")
		mp3File, err := c.FormFile("file")
		if errutil.HandleError(cc, err) {
			cc.InfoF("No mp3 file provided")
			return StatusBadRequest, err.Error()
		}
		editTypeInt, err := strconv.Atoi(editType)
		if errutil.HandleError(cc, err) {
			cc.InfoF("Invalid edit type: %s", editType)
			return StatusBadRequest, err.Error()
		}
		editedAudio, err := h.editAudio(cc, model.GetEditType(editTypeInt), mp3File)
		if errutil.HandleError(cc, err) {
			cc.InfoF("Error editing audio: %s", err.Error())
			return StatusInternalServerError, err.Error()
		}
		return StatusOK, EditTrackResponse{EditedAudio: editedAudio}
	})
}

func (h *Handler) editAudio(
	cc call.Call,
	editType model.EditType,
	mp3File *multipart.FileHeader,
) (mp3bytes []byte, err error) {
	//mp3File.Header.Set("Content-Type", "audio/mpeg")
	file, err := mp3File.Open()
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	defer file.Close()
	mp3Bytes, err := io.ReadAll(file)
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	testFile, err := os.Create("testfile.mp3")
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	defer testFile.Close()
	testFile.Write(mp3Bytes)
	switch editType {
	case model.EditTypeSlowed, model.EditTypeSpedUp:
		return h.reverbedService.SpeedFilter(cc, editType, &mp3Bytes, uuid.New())
	}
	return mp3Bytes, nil
}

type EditTrackResponse struct {
	EditedAudio []byte `json:"edited_audio"`
}
