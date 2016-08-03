package sfix

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

// LoadDefaultRelManager loads the default relationship manager based on default.json.
func LoadDefaultRelManager() ([]byte, error) {
	raw, err := ioutil.ReadFile("sfix/default.json")
	if err != nil {
		return nil, errors.Wrap(err, "Could not load default.json")
	}
	return raw, nil
}
