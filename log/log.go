package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Level constants that define the supported usable LogLevel.
const (
	DEV int = iota + 1
	USER
)

// l contains a standard logger for all logging.
var l struct {
	*log.Logger
	level func() int
}

// Init must be called to initialize the logging system. This function should
// only be called once.
func Init(w io.Writer, level func() int) {
	l.Logger = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.level = level
}

// Dev logs trace information for developers.
func Dev(context interface{}, funcName string, format string, a ...interface{}) {
	if l.level() == DEV {
		if a != nil {
			format = fmt.Sprintf(format, a...)
		}

		l.Output(2, fmt.Sprintf("%s : %s : DEV : %s", context, funcName, format))
	}
}

// User logs trace information for users.
func User(context interface{}, funcName string, format string, a ...interface{}) {
	if a != nil {
		format = fmt.Sprintf(format, a...)
	}

	l.Output(2, fmt.Sprintf("%s : %s : USER : %s", context, funcName, format))
}

// Fatal logs trace information for users and terminates the app.
func Fatal(context interface{}, funcName string, format string, a ...interface{}) {
	if a != nil {
		format = fmt.Sprintf(format, a...)
	}

	l.Output(2, fmt.Sprintf("%s : %s : USER : %s", context, funcName, format))
	os.Exit(1)
}
