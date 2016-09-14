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
func LoadQuadParams(context interface{}, path string) ([]wire.QuadParam, error) {
	log.Dev(context, "LoadQuadParams", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadQuadParams", err, "Completed")
		return nil, err
	}
	defer file.Close()

	var quadParams []wire.QuadParam
	if err = json.NewDecoder(file).Decode(&quadParams); err != nil {
		log.Error(context, "LoadQuadParams", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadQuadParams", "Completed")
	return quadParams, nil
}

// LoadItem serializes the content of an item from a file using the
// given file path. Returns the serialized item values.
func LoadItem(context interface{}, path string) (map[string]interface{}, error) {
	log.Dev(context, "LoadItem", "Started : File %s", path)

	var item map[string]interface{}

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadItem", err, "Completed")
		return item, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&item); err != nil {
		log.Error(context, "LoadItem", err, "Completed")
		return item, err
	}

	log.Dev(context, "LoadItem", "Completed")
	return item, nil
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
