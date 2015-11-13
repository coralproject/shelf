package log_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/coralproject/shelf/log"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// logdest implements io.Writer and is the log package destination.
var logdest bytes.Buffer

// resetLog can be called at the beginning of a test or example.
func resetLog() { logdest.Reset() }

// displayLog can be called at the end of a test or example.
// It only prints the log contents if the -test.v flag is set.
func displayLog() {
	if !testing.Verbose() {
		return
	}
	logdest.WriteTo(os.Stdout)
}

func init() {
	log.Init(&logdest, func() int { return log.USER })
}

// TestLogLevelUSER tests the basic functioning of the logger in USER mode.
func TestLogLevelUSER(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to USER.")
		{
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/1/2 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:52: context : FuncName : USER : Message 2 with format: A, B\n", dt)

			log.Dev("context", "FuncName", "Message 1 no format")
			log.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")

			if logdest.String() == log1 {
				t.Logf("\t\t%v : Should log the expected trace line.", succeed)
			} else {
				t.Errorf("\t\t%v : Should log the expected trace line.", failed)
			}
		}
	}
}

// TestLogLevelDEV tests the basic functioning of the logger in DEV mode.
func TestLogLevelDEV(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to DEV.")
		{
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/1/2 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:77: context : FuncName : DEV : Message 1 no format\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:78: context : FuncName : USER : Message 2 with format: A, B\n", dt)

			log.Dev("context", "FuncName", "Message 1 no format")
			log.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")

			if logdest.String() == log1+log2 {
				t.Logf("\t\t%v : Should log the expected trace line.", succeed)
			} else {
				t.Errorf("\t\t%v : Should log the expected trace line.", failed)
			}
		}
	}
}
