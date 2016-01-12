package tstdata

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CollectionExecTest contains the name of the collection that is
// going to store the xenia test data.
const CollectionExecTest = "test_xenia_data"

var path string

func init() {
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/tstdata/"
}

//==============================================================================

// Generate creates a temp collection with data
// that can be used for testing things.
func Generate(db *db.DB) error {
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

		doc, err := query.UmarshalMongoScript(string(mar), &query.Query{HasDate: true})
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

// Drop drops the temp collection.
func Drop() {
	db := db.NewMGO()
	defer db.CloseMGO()

	mongo.GetCollection(db.MGOConn, CollectionExecTest).DropCollection()
}
