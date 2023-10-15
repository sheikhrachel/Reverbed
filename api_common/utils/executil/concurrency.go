package executil

import (
	"runtime"
	"sync"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

// RecoverGoroutine recovers from a panic in a goroutine and logs the stack trace
func RecoverGoroutine(cc call.Call, wg *sync.WaitGroup) {
	defer wg.Done()
	if r := recover(); r != nil {
		for i := 0; i < 10; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			cc.InfoF("RecoverGoroutine, skip: %+v. %+v recovered: %+v, file: %+v, line: %+v", i, runtime.FuncForPC(pc).Name(), r, file, line)
		}
	}
}
