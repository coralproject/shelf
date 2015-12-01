package cmdquery

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/query"
)

// setFromFile serializes the content of a getSet from a file using the
// given file path.
// Returns the serialized query.getSet, else returns a non-nil error if
// the operation failed.
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

// queriesFromList loads sets of getQuerys from the giving array of file paths.
// Returns a list of query.getQuery, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func queriesFromList(context interface{}, getQueryFilePaths []string) ([]query.Query, error) {
	log.Dev(context, "queriesFromList", "Started : Paths %s", getQueryFilePaths)

	var queries []query.Query

	for _, file := range getQueryFilePaths {
		getQueryFile, err := os.Open(file)
		if err != nil {
			log.Error(context, "queriesFromList", err, "Completed : Paths %s", getQueryFilePaths)
			return nil, err
		}

		var q query.Query
		err = json.NewDecoder(getQueryFile).Decode(&q)
		if err != nil {
			log.Error(context, "queriesFromList", err, "Completed : Paths %s", getQueryFilePaths)
			return nil, err
		}

		queries = append(queries, q)
	}

	log.Dev(context, "queriesFromList", "Completed : Paths %s", getQueryFilePaths)
	return queries, nil
}

// LoadDir loadsup a given directory, calling a load function for each valid
// json file found.
func loadDir(cdir string, ses *mgo.Session, loader func(string, *mgo.Session) error) error {
	if loader == nil {
		return errors.New("No Loader provided")
	}

	var dir string

	// If we have a empty directory argument.
	if cdir == "" {
		if envdir, err := cfg.String(envKey); err == nil {
			dir = envdir
		} else {
			dir = "rules"
		}
	} else {
		dir = cdir
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

		ext := filepath.Ext(path)

		if ext != ".json" {
			return nil
		}

		return loader(path, ses)
	})

	if err2 != nil {
		return err2
	}

	return nil
}
