package disk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/coralproject/xenia/internal/mask"
	"github.com/coralproject/xenia/internal/query"
	"github.com/coralproject/xenia/internal/regex"
	"github.com/coralproject/xenia/internal/script"

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
func LoadScript(context interface{}, path string) (script.Script, error) {
	log.Dev(context, "LoadScript", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadScript", err, "Completed")
		return script.Script{}, err
	}

	var scr script.Script
	if err = json.NewDecoder(file).Decode(&scr); err != nil {
		log.Error(context, "LoadScript", err, "Completed")
		return script.Script{}, err
	}

	log.Dev(context, "LoadScript", "Completed")
	return scr, nil
}

// LoadRegex serializes the content of a regex from a file using the
// given file path. Returns the serialized regex value.
func LoadRegex(context interface{}, path string) (regex.Regex, error) {
	log.Dev(context, "LoadRegex", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadRegex", err, "Completed")
		return regex.Regex{}, err
	}

	var rgx regex.Regex
	if err = json.NewDecoder(file).Decode(&rgx); err != nil {
		log.Error(context, "LoadRegex", err, "Completed")
		return regex.Regex{}, err
	}

	log.Dev(context, "LoadRegex", "Completed")
	return rgx, nil
}

// LoadMask serializes the content of a Mask from a file using the
// given file path. Returns the serialized Mask value.
func LoadMask(context interface{}, path string) (mask.Mask, error) {
	log.Dev(context, "LoadMask", "Started : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "LoadMask", err, "Completed")
		return mask.Mask{}, err
	}

	var msk mask.Mask
	if err = json.NewDecoder(file).Decode(&msk); err != nil {
		log.Error(context, "LoadMask", err, "Completed")
		return mask.Mask{}, err
	}

	log.Dev(context, "LoadMask", "Completed")
	return msk, nil
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
