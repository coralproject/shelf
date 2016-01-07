package cfg_test

import (
	"testing"

	"github.com/ardanlabs/kit/cfg"
)

// Success and failure markers.
var (
	success = "\u2713"
	failed  = "\u2717"
)

//==============================================================================

// TestExists validates the ability to load configuration values
// using the OS-level environment variables and read them back.
func TestExists(t *testing.T) {
	t.Log("Given the need to read environment variables.")
	{
		uStr := "postgres://root:root@127.0.0.1:8080/postgres?sslmode=disable"

		cfg.Init(cfg.MapProvider{
			Map: map[string]string{
				"PROC_ID": "322",
				"SOCKET":  "./tmp/sockets.po",
				"PORT":    "4034",
				"FLAG":    "on",
				"DSN":     uStr,
			},
		})

		t.Log("\tWhen given a namspace key to search for that exists.")
		{
			proc, err := cfg.Int("PROC_ID")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "PROC_ID")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", success, "PROC_ID")

				if proc != 322 {
					t.Errorf("\t\t%s Should have key %q with value %d", failed, "PROC_ID", 322)
				} else {
					t.Logf("\t\t%s Should have key %q with value %d", success, "PROC_ID", 322)
				}
			}

			socket, err := cfg.String("SOCKET")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "SOCKET")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", success, "SOCKET")

				if socket != "./tmp/sockets.po" {
					t.Errorf("\t\t%s Should have key %q with value %q", failed, "SOCKET", "./tmp/sockets.po")
				} else {
					t.Logf("\t\t%s Should have key %q with value %q", success, "SOCKET", "./tmp/sockets.po")
				}
			}

			port, err := cfg.Int("PORT")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "PORT")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", success, "PORT")

				if port != 4034 {
					t.Errorf("\t\t%s Should have key %q with value %d", failed, "PORT", 4034)
				} else {
					t.Logf("\t\t%s Should have key %q with value %d", success, "PORT", 4034)
				}
			}

			flag, err := cfg.Bool("FLAG")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "FLAG")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", success, "FLAG")

				if flag == false {
					t.Errorf("\t\t%s Should have key %q with value %v", failed, "FLAG", true)
				} else {
					t.Logf("\t\t%s Should have key %q with value %v", success, "FLAG", true)
				}
			}

			u, err := cfg.URL("DSN")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "DSN")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", success, "DSN")

				if u.String() != uStr {
					t.Errorf("\t\t%s Should have key %q with value %v", failed, "DSN", true)
				} else {
					t.Logf("\t\t%s Should have key %q with value %v", success, "DSN", true)
				}
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

		cfg.Init(cfg.MapProvider{
			Map: map[string]string{
				"PROC_ID": "322",
				"SOCKET":  "./tmp/sockets.po",
				"PORT":    "4034",
				"FLAG":    "on",
			},
		})

		t.Log("\tWhen given a namspace key to search for that does NOT exist.")
		{
			shouldPanic(t, "STAMP", func() {
				cfg.MustTime("STAMP")
			})

			shouldPanic(t, "PID", func() {
				cfg.MustInt("PID")
			})

			shouldPanic(t, "DEST", func() {
				cfg.MustString("DEST")
			})

			shouldPanic(t, "ACTIVE", func() {
				cfg.MustBool("ACTIVE")
			})

			shouldPanic(t, "SOCKET_DSN", func() {
				cfg.MustURL("SOCKET_DSN")
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
			t.Logf("\t\t%s Should paniced when giving unknown key %q.", success, context)
		}
	}()

	fx()
}
