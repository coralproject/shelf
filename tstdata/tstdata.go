package tstdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/xenia"
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

// Docs reads the fixture and returns the documents.
func Docs() ([]map[string]interface{}, error) {
	file, err := os.Open(path + "test_data.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var docs []map[string]interface{}
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}

	for i := range docs {
		xenia.ProcessVariables("", docs[i], map[string]string{}, nil)
	}

	return docs, nil
}

// Generate creates a temp collection with data
// that can be used for testing things.
func Generate(db *db.DB) error {
	docs, err := Docs()
	if err != nil {
		return err
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
func Drop(db *db.DB) {
	col, err := db.CollectionMGO(tests.Context, CollectionExecTest)
	if err != nil {
		fmt.Printf("***********> Should be able to get a Mongo session : %v\n", err)
		return
	}

	col.DropCollection()
}
