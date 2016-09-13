package disk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/sponge/item"
)

// LoadItems serializes the content of Items from a file using the
// given file path. Returns the serialized Item value.
func LoadItems(context interface{}, path string) ([]item.Item, error) {
	log.Dev(context, "LoadItems", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadItems", err, "Completed")
		return nil, err
	}
	defer file.Close()

	var items []item.Item
	if err = json.NewDecoder(file).Decode(&items); err != nil {
		log.Error(context, "LoadItems", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadItems", "Completed")
	return items, nil
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
