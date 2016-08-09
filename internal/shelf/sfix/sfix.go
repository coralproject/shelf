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

// LoadRelAndViewData loads the default relationship manager based on default.json.
func LoadRelAndViewData() ([]byte, error) {
	raw, err := ioutil.ReadFile(path + "relsandviews.json")
	if err != nil {
		return nil, errors.Wrap(err, "Could not load relmanager.json")
	}
	return raw, nil
}
