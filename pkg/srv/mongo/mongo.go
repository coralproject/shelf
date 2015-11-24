package mongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coralproject/shelf/pkg/cfg"

	"gopkg.in/mgo.v2"
)

// Holds global state for mongo access.
var m struct {
	dbName string
	ses    *mgo.Session
}

// InitMGO sets up the MongoDB environment. This expects that the
// cfg package has been initialized first.
func InitMGO() error {
	if m.ses != nil {
		return errors.New("Mongo environment already initialized")
	}

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := mgo.DialInfo{
		Addrs:    []string{cfg.MustString("MONGO_HOST")},
		Timeout:  60 * time.Second,
		Database: cfg.MustString("MONGO_AUTHDB"),
		Username: cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	var err error
	if m.ses, err = mgo.DialWithInfo(&mongoDBDialInfo); err != nil {
		return err
	}

	// Save the database name to use.
	m.dbName = cfg.MustString("MONGO_DB")

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	m.ses.SetMode(mgo.Monotonic, true)

	return nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}

// GetSession returns a copy of the master session for use.
func GetSession() *mgo.Session {
	return m.ses.Copy()
}

// GetDatabase returns a mgo database value based on configuration.
func GetDatabase(session *mgo.Session) *mgo.Database {
	return session.DB(m.dbName)
}

// GetDatabaseName returns the name of the database being used.
func GetDatabaseName() string {
	return m.dbName
}

// GetCollection returns a mgo collection value based on configuration.
func GetCollection(session *mgo.Session, colName string) *mgo.Collection {
	return session.DB(m.dbName).C(colName)
}

// ExecuteDB the MongoDB literal function.
func ExecuteDB(context interface{}, session *mgo.Session, collectionName string, f func(*mgo.Collection) error) error {
	// Capture the specified collection.
	col := session.DB(m.dbName).C(collectionName)
	if col == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}

	// Execute the MongoDB call.
	return f(col)
}

// CollectionExists returns true if the collection name exists in the specified database.
func CollectionExists(context interface{}, session *mgo.Session, useCollection string) bool {
	cols, err := session.DB(m.dbName).CollectionNames()
	if err != nil {
		return false
	}

	for _, col := range cols {
		if col == useCollection {
			return true
		}
	}

	return false
}
