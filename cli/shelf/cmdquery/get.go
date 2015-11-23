package cmdquery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coralproject/shelf/pkg/db/query"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/spf13/cobra"
)

var getLong = `Retrieves a query record from the system using the supplied name.
Example:

		user get -n user_advice

`

// get contains the state for this command.
var get struct {
	name string
}

// addGet handles the retrival users records, displayed in json formatted response.
func addGet() {
	cmd := &cobra.Command{
		Use:   "get [-n name]",
		Short: "Retrieves a query record",
		Long:  getLong,
		Run:   runGet,
	}

	cmd.Flags().StringVarP(&get.name, "name", "n", "", "name of the user record")

	queryCmd.AddCommand(cmd)
}

// runGet is the code that implements the get command.
func runGet(cmd *cobra.Command, args []string) {
}

// getQuerySetFromFile serializes the content of a RuleSet from a file using the
// given file path.
// Returns the serialized query.RuleSet, else returns a non-nil error if
// the operation failed.
func getQuerySetFromFile(context interface{}, path string) (*query.QuerySet, error) {
	log.Dev(context, "RuleSetFromFile", "Started : Load RuleSet : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	var qs query.QuerySet

	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	log.Dev(context, "RuleSetFromFile", "Completed : Load RuleSet : File %s", path)
	return &qs, nil
}

// getQueriesFromPaths loads sets of rules from the giving array of file paths.
// Returns a list of query.Rule, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func getQueriesFromPaths(context interface{}, ruleFilePaths []string) ([]query.Query, error) {
	log.Dev(context, "RuleFromPaths", "Started : Paths %s", ruleFilePaths)

	var queries []query.Query

	for _, file := range ruleFilePaths {
		ruleFile, err := os.Open(file)
		if err != nil {
			log.Error(context, "RuleFromPaths", err, "Completed : Paths %s", ruleFilePaths)
			return nil, err
		}

		var q query.Query
		err = json.NewDecoder(ruleFile).Decode(&q)
		if err != nil {
			log.Error(context, "RuleFromPaths", err, "Completed : Paths %s", ruleFilePaths)
			return nil, err
		}

		queries = append(queries, q)
	}

	log.Dev(context, "RuleFromPaths", "Completed : Paths %s", ruleFilePaths)
	return queries, nil
}

// queryFromDir loads sets of rules from the giving files in the directory path,
// only reading the current directory level and not sub-directories.
// Returns a list of Rule pointers, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func queryFromDir(context interface{}, dirPath string) ([]query.Query, error) {
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

	rules, err := getQueriesFromPaths(context, files)
	if err != nil {
		log.Error(context, "RulesFromDir", err, "Completed : Load Rules : Dir %s", dirPath)
		return nil, err
	}

	log.Dev(context, "RulesFromDir", "Completed : Load Rules : Dir %s", dirPath)
	return rules, nil
}
