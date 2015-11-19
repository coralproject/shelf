package db

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection sets the collection to store expressions in.
var QueryCollection = "queries"

// Expression represents a condition query to be executed against a given collection.
type Expression struct {
	Collection string   `bson:"collection" json:"collection"`
	Queries    []string `bson:"queries" json:"queries"`
}

// Compare compares the fields of one expression against itself. Returns a non-nil
// error if any did not matche.
func (e *Expression) Compare(eu *Expression) error {
	if e.Collection != eu.Collection {
		return errors.New("Collection is not a match")
	}

	// match the Queries of both using their index.
	for i, qs := range e.Queries {
		if eu.Queries[i] == qs {
			continue
		}
		return errors.New("Queries are not a match")
	}

	return nil
}

// Query represents a set of expressions to be runned against a Collection,
// with a binary mindset on the result of its Test.
type Query struct {
	ID         bson.ObjectId `bson:"id" json:"id,omitempty"`
	Name       string        `bson:"name" json:"name"`
	Test       *Expression   `bson:"test" json:"test"`
	Failed     *Expression   `bon:"failed" json:"failed"`
	Passed     *Expression   `bson:"passed" json:"passed"`
	CreatedAt  *time.Time    `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt *time.Time    `bson:"modified_at" json:"modified_at,omitempty"`
}

// QueryFromFile loads the query structure from a giving file path. It expects
// the file to be a json document. Returns the query object if successfull else
// returns a non-nil error.
func QueryFromFile(file string) (*Query, error) {
	var q Query

	if err := q.LoadFile(file); err != nil {
		return nil, err
	}

	return &q, nil
}

// NewQuery returns a new query object.
func NewQuery(name string) *Query {
	ms := time.Now()
	q := Query{
		ID:         bson.NewObjectId(),
		Name:       name,
		CreatedAt:  &ms,
		ModifiedAt: &ms,
	}
	return &q
}

// Compare compares the fields of one query against itself. Returns a non-nil
// error if any did not matche.
func (q *Query) Compare(qu *Query) error {
	if q.Name != qu.Name {
		return errors.New("Name is not a match")
	}

	if err := q.Test.Compare(qu.Test); err != nil {
		return err
	}

	if err := q.Passed.Compare(qu.Passed); err != nil {
		return err
	}

	if err := q.Failed.Compare(qu.Failed); err != nil {
		return err
	}

	return nil
}

// LoadFile loads the query structure from a giving file path. It expects
// the file to be a json document.
// Returns a non-nil error if the operation fails.
func (q *Query) LoadFile(file string) error {
	log.Dev("Query", "LoadFile", "Started : Query : LoadFile : File[%s]", file)
	inputFile, err := os.Open(file)
	if err != nil {
		log.Error("Query", "LoadFile", err, "Completed : Query : LoadFile : File[%s]", file)
		return err
	}

	defer inputFile.Close()

	if err := json.NewDecoder(inputFile).Decode(q); err != nil {
		log.Error("Query", "LoadFile", err, "Completed : Query : LoadFile : File[%s]", file)
		return err
	}

	if q.Name == "" {
		_, fileName := filepath.Split(file)
		ext := filepath.Ext(fileName)
		q.Name = strings.Replace(fileName, ext, "", -1)
	}

	log.Dev("Query", "LoadFile", "Completed : Query : LoadFile : File[%s]", file)
	return nil
}

// GetByID retrieves the query object by its ID, and returns the query document.
// Returns a non-nil error if the operation fails.
func GetByID(id string) (*Query, error) {
	log.Dev(id, "GetByID", "Started : Query : Get Query")

	if !bson.IsObjectIdHex(id) {
		err := errors.New("Invalid bson.ObjectID value")
		log.Error(id, "GetByID", err, "Completed : Query : Get Query")
		return nil, err
	}

	session := mongo.GetSession()
	defer session.Close()

	var q Query

	var getQuery = func(c *mgo.Collection) error {
		log.Dev(id, "GetByID", "Completed : Query : MongoDb.Find().One()")
		return c.Find(bson.M{"id": bson.ObjectIdHex(id)}).One(&q)
	}

	log.Dev(id, "GetByID", "Started : Query : MongoDb.Find().One()")
	err := mongo.ExecuteDB("CONTEXT", session, QueryCollection, getQuery)
	if err != nil {
		log.Error(id, "GetByID", err, "Completed : Query : Get Query")
		return nil, err
	}

	log.Dev(id, "GetByID", "Completed : Query : Get Query")
	return &q, nil
}

// GetByName retrieves the query object by its Name, and returns the query document.
// Returns a non-nil error if the operation fails.
func GetByName(name string) (*Query, error) {
	log.Dev(name, "GetByName", "Started : Query : Get Query")
	session := mongo.GetSession()
	defer session.Close()

	var q Query

	var getQuery = func(c *mgo.Collection) error {
		log.Dev(name, "GetByName", "Completed : Query : MongoDb.Find().One()")
		return c.Find(bson.M{"name": name}).One(&q)
	}

	log.Dev(name, "GetByName", "Started : Query : MongoDb.Find().One()")
	err := mongo.ExecuteDB("CONTEXT", session, QueryCollection, getQuery)
	if err != nil {
		log.Error(name, "GetByName", err, "Completed : Query : Get Query")
		return nil, err
	}

	log.Dev(name, "GetByName", "Completed : Query : Get Query")
	return &q, nil
}

// Create adds the given query into the database collection.
// Returns a non-nil error if the operation fails.
func Create(q *Query) error {
	log.Dev(q.Name, "Create", "Started : Query : Create Query")
	session := mongo.GetSession()
	defer session.Close()

	var createQuery = func(c *mgo.Collection) error {
		log.Dev(q.Name, "Create", "Complete : Query : MongoDb.Insert()")
		return c.Insert(q)
	}

	log.Dev(q.Name, "Create", "Started : Query : MongoDb.Insert()")
	err := mongo.ExecuteDB("CONTEXT", session, QueryCollection, createQuery)
	if err != nil {
		log.Error(q.Name, "Create", err, "Completed : Query : Create Query")
		return err
	}

	log.Dev(q.Name, "Create", "Completed : Query : Create Query")
	return nil
}

// Update modifies the internal structure of a existing query in the database.
// Returns a non-nil error if the operation fails.
func Update(q *Query) error {
	log.Dev(mongo.Query(q.ID), "Update", "Started : Query : Update Query")
	session := mongo.GetSession()
	defer session.Close()

	ms := time.Now()
	updates := bson.M{
		"name": q.Name,
		"test": bson.M{
			"collection": q.Test.Collection,
			"queries":    q.Test.Queries,
		},
		"passed": bson.M{
			"collection": q.Passed.Collection,
			"queries":    q.Passed.Queries,
		},
		"failed": bson.M{
			"collection": q.Failed.Collection,
			"queries":    q.Failed.Queries,
		},
		"modified_at": &ms,
	}

	var updateQuery = func(c *mgo.Collection) error {
		log.Dev(mongo.Query(q.ID), "Update", "Completed : Query : MongoDb.Update()")
		return c.Update(bson.M{"id": q.ID}, bson.M{"$set": updates})
	}

	log.Dev(mongo.Query(q.ID), "Update", "Started : Query : MongoDb.Update()")
	err := mongo.ExecuteDB("CONTEXT", session, QueryCollection, updateQuery)
	if err != nil {
		log.Error(mongo.Query(q.ID), "Update", err, "Completed : Query : Update Query")
		return err
	}

	log.Dev(mongo.Query(q.ID), "Update", "Completed : Query : Update Query")
	return nil
}

// Delete removes the giving query from the database collection
// Returns a non-nil error if the operation fails.
func Delete(q *Query) error {
	log.Dev(q.Name, "Delete", "Started : Query : Delete Query")
	session := mongo.GetSession()
	defer session.Close()

	var deleteQuery = func(c *mgo.Collection) error {
		log.Dev(q.Name, "Delete", "Completed : Query : MongoDb.Remove()")
		return c.Remove(bson.M{"name": q.Name})
	}

	log.Dev(q.Name, "Delete", "Started : Query : MongoDb.Remove()")
	err := mongo.ExecuteDB("CONTEXT", session, QueryCollection, deleteQuery)
	if err != nil {
		log.Error(q.Name, "Delete", err, "Completed : Query : Delete Query")
		return err
	}

	log.Dev(q.Name, "Delete", "Completed : Query : Delete Query")
	return nil
}
