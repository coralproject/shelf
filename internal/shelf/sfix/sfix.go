package sfix

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/internal/shelf/sfix/"
}

// LoadRelManagerData loads the default relationship manager based on default.json.
func LoadRelManagerData() ([]byte, error) {
	raw, err := ioutil.ReadFile(path + "relmanager.json")
	if err != nil {
		return nil, errors.Wrap(err, "Could not load relmanager.json")
	}
	return raw, nil
}
