package disk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/pkg/query"
	"github.com/coralproject/xenia/pkg/regex"
	"github.com/coralproject/xenia/pkg/script"

	"github.com/ardanlabs/kit/log"
)

// LoadSet serializes the content of a Set from a file using the
// given file path. Returns the serialized Set value.
func LoadSet(context interface{}, path string) (*query.Set, error) {
	log.Dev(context, "LoadSet", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadSet", err, "Completed")
		return nil, err
	}

	var set query.Set
	if err = json.NewDecoder(file).Decode(&set); err != nil {
		log.Error(context, "LoadSet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadSet", "Completed")
	return &set, nil
}

// LoadScript serializes the content of a Script from a file using the
// given file path. Returns the serialized Script value.
func LoadScript(context interface{}, path string) (*script.Script, error) {
	log.Dev(context, "LoadScript", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadScript", err, "Completed")
		return nil, err
	}

	var scr script.Script
	if err = json.NewDecoder(file).Decode(&scr); err != nil {
		log.Error(context, "LoadScript", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadScript", "Completed")
	return &scr, nil
}

// LoadRegex serializes the content of a regex from a file using the
// given file path. Returns the serialized regex value.
func LoadRegex(context interface{}, path string) (*regex.Regex, error) {
	log.Dev(context, "LoadRegex", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadRegex", err, "Completed")
		return nil, err
	}

	var rgx regex.Regex
	if err = json.NewDecoder(file).Decode(&rgx); err != nil {
		log.Error(context, "LoadRegex", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoadRegex", "Completed")
	return &rgx, nil
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
