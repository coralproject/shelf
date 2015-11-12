package cfg_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ArdanStudios/aggserver/cfg"
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

// TestExists validates the ability to load configuration values
// using the OS-level environment variables and read them back.
func TestExists(t *testing.T) {
	t.Log("Given the need to read environment variables.")
	{
		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")

		cfg.Init("myapp")

		t.Log("\tWhen given a namspace key to search for that exists.")
		{
			if cfg.Int("proc_id") != 322 {
				t.Errorf("\t\t%s Should have key %q with value %d", failed, "proc_id", 322)
			} else {
				t.Logf("\t\t%s Should have key %q with value %d", succeed, "proc_id", 322)
			}

			if cfg.String("socket") != "./tmp/sockets.po" {
				t.Errorf("\t\t%s Should have key %q with value %q", failed, "socket", "./tmp/sockets.po")
			} else {
				t.Logf("\t\t%s Should have key %q with value %q", succeed, "socket", "./tmp/sockets.po")
			}

			if cfg.Int("port") != 4034 {
				t.Errorf("\t\t%s Should have key %q with value %d", failed, "port", 4034)
			} else {
				t.Logf("\t\t%s Should have key %q with value %d", succeed, "port", 4034)
			}
		}
	}
}

// TestNotExists validates the ability to load configuration values
// using the OS-level environment variables and panic when something
// is missing.
func TestNotExists(t *testing.T) {
	t.Log("Given the need to panic when environment variables are missing.")
	{
		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")

		cfg.Init("myapp")

		t.Log("\tWhen given a namspace key to search for that does NOT exist.")
		{
			shouldPanic(t, "stamp", func() {
				cfg.Time("stamp")
			})

			shouldPanic(t, "pid", func() {
				cfg.Int("pid")
			})

			shouldPanic(t, "dest", func() {
				cfg.String("dest")
			})

		}
	}
}

// shouldPanic receives a context string and a function to run, if the function
// panics, it is considered a success else a failure.
func shouldPanic(t *testing.T, context string, fx func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("\t\t%s Should paniced when giving unknown key %q.", failed, context)
		} else {
			t.Logf("\t\t%s Should paniced when giving unknown key %q.", succeed, context)
		}
	}()

	fx()
}
