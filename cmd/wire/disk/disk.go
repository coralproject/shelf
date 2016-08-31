package disk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/wire"
)

// LoadQuadParams serializes the content of a set of QuadParams from a file using the
// given file path. Returns the serialized QuadParams value.
func LoadQuadParams(context interface{}, path string) ([]wire.QuadParams, error) {
	log.Dev(context, "LoadQuadParams", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadQuadParams", err, "Completed")
		return nil, err
	}

	var quadParams []wire.QuadParams
	if err = json.NewDecoder(file).Decode(&quadParams); err != nil {
		log.Error(context, "LoadQuadParams", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadQuadParams", "Completed")
	return quadParams, nil
}

// LoadDir loadsup a given directory, calling a load function for each valid
// json file found.
func LoadDir(dir string, loader func(string) error) error {
	if loader == nil {
		return errors.New("No Loader provided")
	}

	if _, err := os.Stat(dir); err != nil {
		return err
	}

	f := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		return loader(path)
	}

	if err := filepath.Walk(dir, f); err != nil {
		return err
	}

	return nil
}
