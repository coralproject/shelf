package tstdata

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/coralproject/xenia/pkg/exec"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2"
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

	var docs []map[string]interface{}
	if err := json.Unmarshal(data, &docs); err != nil {
		return err
	}

	for i := range docs {
		docs[i] = exec.PreProcess(docs[i], map[string]string{})
	}

	// The Insert calls requires this converstion.
	var insDocs []interface{}
	for _, doc := range docs {
		insDocs = append(insDocs, doc)
	}

	f := func(c *mgo.Collection) error {
		return c.Insert(insDocs...)
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
