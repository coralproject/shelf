package query

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/pkg/db/query"
	"github.com/coralproject/shelf/pkg/log"
)

// RuleSetFromReader serializes the content of a RuleSet from a io.Reader.
// Returns the serialized RuleSet pointer, else returns a non-nil error if
// the operation failed.
func RuleSetFromReader(context interface{}, r io.Reader) (*query.RuleSet, error) {
	log.Dev(context, "RuleSetFromReader", "Started : Load RuleSet")
	var rs query.RuleSet

	err := json.NewDecoder(r).Decode(&rs)
	if err != nil {
		log.Error(context, "RuleSetFromReader", err, "Completed : Load RuleSet")
		return nil, err
	}

	log.Dev(context, "RuleSetFromReader", "Completed : Load RuleSet")
	return &rs, nil
}

// RuleSetFromFile serializes the content of a RuleSet from a file using the
// given file path.
// Returns the serialized query.RuleSet, else returns a non-nil error if
// the operation failed.
func RuleSetFromFile(context interface{}, path string) (*query.RuleSet, error) {
	log.Dev(context, "RuleSetFromFile", "Started : Load RuleSet : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	var rs query.RuleSet

	err = json.NewDecoder(file).Decode(&rs)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	log.Dev(context, "RuleSetFromFile", "Completed : Load RuleSet : File %s", path)
	return &rs, nil
}

// RulesFromPaths loads sets of rules from the giving array of file paths.
// Returns a list of query.Rule, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func RulesFromPaths(context interface{}, ruleFilePaths []string) ([]query.Rule, error) {
	log.Dev(context, "RuleFromPaths", "Started : Paths %s", ruleFilePaths)

	var rules []query.Rule

	for _, file := range ruleFilePaths {
		ruleFile, err := os.Open(file)
		if err != nil {
			log.Error(context, "RuleFromPaths", err, "Completed : Paths %s", ruleFilePaths)
			return nil, err
		}

		var r query.Rule
		err = json.NewDecoder(ruleFile).Decode(&r)
		if err != nil {
			log.Error(context, "RuleFromPaths", err, "Completed : Paths %s", ruleFilePaths)
			return nil, err
		}

		rules = append(rules, r)
	}

	log.Dev(context, "RuleFromPaths", "Completed : Paths %s", ruleFilePaths)
	return rules, nil
}

// RulesFromDir loads sets of rules from the giving files in the directory path,
// only reading the current directory level and not sub-directories.
// Returns a list of Rule pointers, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func RulesFromDir(context interface{}, dirPath string) ([]query.Rule, error) {
	log.Dev(context, "RulesFromDir", "Started : Load Rules : Dir %s", dirPath)

	stat, err := os.Stat(dirPath)
	if err != nil {
		log.Error(context, "RulesFromDir", err, "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	if !stat.IsDir() {
		log.Error(context, "RulesFromDir", fmt.Errorf("Path[%s] is not a Directory", dirPath), "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	//open up the filepath since its a directory, read and sort
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Error(context, "RulesFromDir", err, "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	filesInfo, err := dir.Readdir(0)
	if err != nil {
		log.Error(context, "RulesFromDir", err, "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	dir.Close()

	var files []string

	for _, info := range filesInfo {
		if info.IsDir() {
			continue
		}

		files = append(files, filepath.Join(dirPath, info.Name()))
	}

	rules, err := RulesFromPaths(context, files)
	if err != nil {
		log.Error(context, "RulesFromDir", err, "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	log.Dev(context, "RulesFromDir", "Completed : Load Rules : Dir %s", dirPath)
	return rules, nil
}
