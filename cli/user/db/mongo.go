package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/coralproject/shelf/log"
	"gopkg.in/mgo.v2"
)

// MongoDB connection information.
// TODO: Need to read configuration.
const (
	mongoDBHosts = "ds039441.mongolab.com:39441"
	authDatabase = "gotraining"
	authUserName = "got"
	authPassword = "got2015"
	database     = "gotraining"
)

// session maintains the master session
var session *mgo.Session

// InitMGO sets up the MongoDB environment.
func InitMGO() {
	log.Dev("mongodb", "InitMGO", "Started : Host[%s] Database[%s]", mongoDBHosts, database)

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := mgo.DialInfo{
		Addrs:    []string{mongoDBHosts},
		Timeout:  60 * time.Second,
		Database: authDatabase,
		Username: authUserName,
		Password: authPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	var err error
	if session, err = mgo.DialWithInfo(&mongoDBDialInfo); err != nil {
		log.Fatal("mongodb", "InitMGO", "MongoDB Dial : %v", err.Error())
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	session.SetMode(mgo.Monotonic, true)

	log.Dev("mongodb", "InitMGO", "Completed")
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
	return session.Copy()
}

// ExecuteDB the MongoDB literal function.
func ExecuteDB(session *mgo.Session, collectionName string, f func(*mgo.Collection) error) error {
	log.Dev("mongodb", "ExecuteDB", "Started : Collection[%s]", collectionName)

	// Capture the specified collection.
	collection := session.DB(database).C(collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		log.Dev("mongodb", "ExecuteDB", "Completed : ERROR %s ", err.Error())
		return err
	}

	// Execute the MongoDB call.
	if err := f(collection); err != nil {
		log.Dev("mongodb", "ExecuteDB", "Completed : ERROR %s", err.Error())
		return err
	}

	log.Dev("mongodb", "ExecuteDB", "Completed")
	return nil
}
