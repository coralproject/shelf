package query

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/fixtures/"
}

//==============================================================================

// GetFixture retrieves a query record from the filesystem for testing.
func GetFixture(fileName string) (*Set, error) {
	file, err := os.Open(path + fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var qs Set
	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		return nil, err
	}

	return &qs, nil
}

// AddTestSet inserts a set for testing.
func AddTestSet(db *db.DB, qs *Set) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": qs.Name}
		_, err := c.Upsert(q, qs)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, Collection, f); err != nil {
		return err
	}

	return nil
}

// RemoveTestSets is used to clear out all the test sets from the collection.
// All test query sets must start with QSTEST in their name.
func RemoveTestSets(db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "QTEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "QTEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, CollectionHistory, f); err != nil {
		return err
	}

	return nil
}

// GenerateTestData creates a temp collection with data
// that can be used for testing things.
func GenerateTestData(db *db.DB) error {
	file, err := os.Open(path + "test_data.json")
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var rawDocs []bson.M
	if err := json.Unmarshal(data, &rawDocs); err != nil {
		return err
	}

	var docs []interface{}
	for _, rd := range rawDocs {
		mar, err := json.Marshal(rd)
		if err != nil {
			return err
		}

		doc, err := UmarshalMongoScript(string(mar), &Query{HasDate: true})
		if err != nil {
			return err
		}

		docs = append(docs, doc)
	}

	f := func(c *mgo.Collection) error {
		return c.Insert(docs...)
	}

	if err := db.ExecuteMGO(tests.Context, CollectionExecTest, f); err != nil {
		return err
	}

	return nil
}

// DropTestData drops the temp collection.
func DropTestData() {
	db := db.NewMGO()
	defer db.CloseMGO()

	mongo.GetCollection(db.MGOConn, CollectionExecTest).DropCollection()
}
