package mongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
)

// Holds global state for mongo access.
var m struct {
	dbName string
	ses    *mgo.Session
}

// InitMGO sets up the MongoDB environment.
func InitMGO(hostKey, authDBKey, userKey, passKey, dbKey string) error {
	if m.ses != nil {
		return errors.New("Mongo environment already initialized")
	}

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := mgo.DialInfo{
		Addrs:    []string{cfg.MustString(hostKey)},
		Timeout:  60 * time.Second,
		Database: cfg.MustString(authDBKey),
		Username: cfg.MustString(userKey),
		Password: cfg.MustString(passKey),
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	var err error
	if m.ses, err = mgo.DialWithInfo(&mongoDBDialInfo); err != nil {
		return err
	}

	// Save the database name to use.
	m.dbName = cfg.MustString(dbKey)

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
	log.Dev(context, "ExecuteDB", "Started : Collection[%s]", collectionName)

	// Capture the specified collection.
	col := session.DB(m.dbName).C(collectionName)
	if col == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		log.Error(context, "ExecuteDB", err, "Completed")
		return err
	}

	// Execute the MongoDB call.
	if err := f(col); err != nil {
		log.Error(context, "ExecuteDB", err, "Completed")
		return err
	}

	log.Dev(context, "ExecuteDB", "Completed")
	return nil
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
