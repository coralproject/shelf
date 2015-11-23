package query

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the rules collection.
const collection = "rules"

// GetQuerySetNames retrieves a list of rule names.
func GetQuerySetNames(context interface{}, ses *mgo.Session) ([]string, error) {
	log.Dev(context, "GetQuerySetNames", "Started")

	var names []bson.M
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": 1}
		log.Dev(context, "GetQuerySetNames", "MGO : db.%s.find({}, %s).sort([\"name\"])", collection, mongo.Query(q))
		return c.Find(nil).Select(q).Sort("name").All(&names)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetQuerySetNames", err, "Completed")
		return nil, err
	}

	var rsn []string
	for _, doc := range names {
		name := doc["name"].(string)
		if strings.HasPrefix(name, "test") {
			continue
		}

		rsn = append(rsn, name)
	}

	log.Dev(context, "GetQuerySetNames", "Completed : RSN[%+v]", rsn)
	return rsn, nil
}

// GetQuerySet retrieves the configuration for the specified QuerySet.
func GetQuerySet(context interface{}, ses *mgo.Session, name string) (*QuerySet, error) {
	log.Dev(context, "GetQuerySet", "Started : Name[%s]", name)

	var rs QuerySet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}
		log.Dev(context, "GetQuerySet", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&rs)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetQuerySet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetQuerySet", "Completed : RS[%+v]", rs)
	return &rs, nil
}

// UpdateQuerySet is used to create or update existing QuerySet documents.
func UpdateQuerySet(context interface{}, ses *mgo.Session, rs *QuerySet) error {
	log.Dev(context, "UpdateQuerySet", "Started : Name[%s]", rs.Name)

	f := func(c *mgo.Collection) error {
		q := bson.M{"name": rs.Name}

		log.Dev(context, "UpdateQuerySet", "MGO :\n\ndb.%s.upsert(%s, %s)\n", collection, mongo.Query(q), mongo.Query(rs))
		_, err := c.Upsert(q, rs)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateQuerySet", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateQuerySet", "Completed")
	return nil
}

// RemoveQuerySet is used to remove an existing QuerySet documents.
func RemoveQuerySet(context interface{}, ses *mgo.Session, name string) (*QuerySet, error) {
	log.Dev(context, "RemoveQuerySet", "Started : Name[%s]", name)

	var rs QuerySet
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": name}

		log.Dev(context, "RemoveQuerySet", "MGO :\n\ndb.%s.remove(%s)\n", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "RemoveQuerySet", err, "Completed")
		return nil, err
	}

	log.Dev(context, "RemoveQuerySet", "Completed")
	return &rs, nil
}

// QuerySetFromReader serializes the content of a RuleSet from a io.Reader.
// Returns the serialized RuleSet pointer, else returns a non-nil error if
// the operation failed.
func QuerySetFromReader(context interface{}, r io.Reader) (*QuerySet, error) {
	log.Dev(context, "RuleSetFromReader", "Started : Load RuleSet")
	var rs QuerySet

	err := json.NewDecoder(r).Decode(&rs)
	if err != nil {
		log.Error(context, "RuleSetFromReader", err, "Completed : Load RuleSet")
		return nil, err
	}

	log.Dev(context, "RuleSetFromReader", "Completed : Load RuleSet")
	return &rs, nil
}

// QuerySetFromFile serializes the content of a RuleSet from a file using the
// given file path.
// Returns the serialized query.RuleSet, else returns a non-nil error if
// the operation failed.
func QuerySetFromFile(context interface{}, path string) (*QuerySet, error) {
	log.Dev(context, "RuleSetFromFile", "Started : Load RuleSet : File %s", path)

	file, err := os.Open(path)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	var rs QuerySet

	err = json.NewDecoder(file).Decode(&rs)
	if err != nil {
		log.Error(context, "RuleSetFromFile", err, "Completed : Load RuleSet : File %s", path)
		return nil, err
	}

	log.Dev(context, "RuleSetFromFile", "Completed : Load RuleSet : File %s", path)
	return &rs, nil
}

// queryFromPaths loads sets of rules from the giving array of file paths.
// Returns a list of query.Rule, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func queryFromPaths(context interface{}, ruleFilePaths []string) ([]Query, error) {
	log.Dev(context, "RuleFromPaths", "Started : Paths %s", ruleFilePaths)

	var rules []Query

	for _, file := range ruleFilePaths {
		ruleFile, err := os.Open(file)
		if err != nil {
			log.Error(context, "RuleFromPaths", err, "Completed : Paths %s", ruleFilePaths)
			return nil, err
		}

		var r Query
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

// queryFromDir loads sets of rules from the giving files in the directory path,
// only reading the current directory level and not sub-directories.
// Returns a list of Rule pointers, each serialized with the contents of it's file.
// If any of the paths are invalid or there was a failure to load their content,
// a non-nil error is returned.
func queryFromDir(context interface{}, dirPath string) ([]Query, error) {
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
