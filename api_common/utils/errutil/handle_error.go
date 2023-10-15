package errutil

import (
	"runtime"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

// HandleError is a util wrapper for exposing error context throughout the project
// to make triaging and maintaining code more easily via the logs
func HandleError(cc call.Call, err error) (b bool) {
	if err != nil {
		/*
			for runtime.Caller():
			skip = 0: this function
			skip = 1: the caller function this helper function is used in
		*/
		_, filename, line, _ := runtime.Caller(1)
		cc.InfoF("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
}
