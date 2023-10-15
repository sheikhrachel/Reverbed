package call

import (
	"fmt"
	"log"
)

const loggingPrefix = "reverbedService"

// getLoggingPrefix returns the logging prefix for the current service and environment
func (cc *Call) getLoggingPrefix() string {
	return fmt.Sprintf("%+v:[%+v]", loggingPrefix, cc.Env)
}

// InfoF logs an info message
func (cc *Call) InfoF(stringRaw string, a ...any) {
	log.Println(cc.getLoggingPrefix(), fmt.Sprintf(stringRaw, a...))
}

// TraceF logs a trace message
func (cc *Call) TraceF(stringRaw string, a ...any) {
	log.Println(cc.getLoggingPrefix(), fmt.Sprintf(stringRaw, a...))
}

// ErrorF logs an error message
func (cc *Call) ErrorF(stringRaw string, a ...any) error {
	log.Println(cc.getLoggingPrefix(), fmt.Sprintf(stringRaw, a...))
	return fmt.Errorf(stringRaw, a...)
}
