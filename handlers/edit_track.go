package handlers

import (
	"io"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
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
		defer editedAudio.Close()
		c.Header("Content-Type", "audio/mpeg")
		io.Copy(c.Writer, editedAudio)
		return StatusOK, nil
	})
}

func (h *Handler) editAudio(
	cc call.Call,
	editType model.EditType,
	mp3File *multipart.FileHeader,
) (file multipart.File, err error) {
	switch editType {
	case model.EditTypeSlowed, model.EditTypeSpedUp:
		return h.reverbedService.SpeedFilter(cc, editType, mp3File)
	}
	return file, nil
}
