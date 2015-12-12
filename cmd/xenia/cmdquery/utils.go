package cmdquery

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/log"
)

// setFromFile serializes the content of a getSet from a file using the
// given file path. Returns the serialized query.getSet, else returns a
// non-nil error if the operation failed.
func setFromFile(context interface{}, path string) (*query.Set, error) {
	log.Dev(context, "setFromFile", "Started : Load getSet : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "setFromFile", err, "Completed : Load getSet : File %s", path)
		return nil, err
	}

	var qs query.Set

	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		log.Error(context, "setFromFile", err, "Completed : Load getSet : File %s", path)
		return nil, err
	}

	log.Dev(context, "setFromFile", "Completed : Load getSet : File %s", path)
	return &qs, nil
}

// LoadDir loadsup a given directory, calling a load function for each valid
// json file found.
func loadDir(dir string, loader func(string) error) error {
	if loader == nil {
		return errors.New("No Loader provided")
	}

	_, err := os.Stat(dir)
	if err != nil && err == os.ErrNotExist {
		return err
	}

	err2 := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		return loader(path)
	})

	if err2 != nil {
		return err2
	}

	return nil
}
