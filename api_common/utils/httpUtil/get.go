package httpUtil

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

// GetRespBody performs a get http request, sets the headers, and returns the response body
func GetRespBody(cc call.Call, apiKey, url string, bodyJSON []byte) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(bodyJSON))
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Add("X-API-KEY", apiKey)
	}
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if errutil.HandleError(cc, err) {
		return nil, err
	}
	if body != nil && strings.Contains(string(body), "error") {
		cc.InfoF("GetRespBody: error %+v", string(body))
		return nil, nil
	}
	return body, nil
}
