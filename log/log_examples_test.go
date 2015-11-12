package log_test

import (
	"errors"
	"os"
	"testing"

	"github.com/coralproject/shelf/log"
)

// ExampleDev shows how to use the log package.
func ExampleDev(t *testing.T) {
	// Init the log package for stdout. Hardcode the logging level
	// function to use USER level logging.
	log.Init(os.Stdout, func() int { return log.USER })

	// Write a simple log line with no formatting.
	log.User("context", "ExampleDev", "This is a simple line with no formatting")

	// Write a simple log line with formatting.
	log.User("context", "ExampleDev", "This is a simple line with no formatting %d", 10)

	// Write a message error for the user.
	log.User("context", "ExampleDev", "ERROR: %v", errors.New("A user error"))

	// Write a message error for the developer only.
	log.Dev("context", "ExampleDev", "ERROR: %v", errors.New("An developer error"))
}
