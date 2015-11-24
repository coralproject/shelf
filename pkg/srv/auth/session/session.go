package session

import (
	"time"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collection = "session"

//==============================================================================

// Create adds a new session for the specified user to the database.
func Create(context interface{}, ses *mgo.Session, publicID string, expires time.Duration) (*Session, error) {
	log.Dev(context, "Create", "Started : PublicID[%s]", publicID)

	s := Session{
		SessionID:   uuid.New(),
		PublicID:    publicID,
		DateExpires: time.Now().Add(expires),
		DateCreated: time.Now(),
	}

	f := func(col *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", collection, mongo.Query(&s))
		return col.Insert(&s)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Create", "Completed")
	return &s, nil
}

//==============================================================================

// Get retrieves a session from the session store.
func Get(context interface{}, ses *mgo.Session, sessionID string) (*Session, error) {
	log.Dev(context, "Get", "Started : SessionID[%s]", sessionID)

	var s Session
	f := func(c *mgo.Collection) error {
		q := bson.M{"session_id": sessionID}
		log.Dev(context, "Get", "MGO : db.%s.find(%s).sort({\"date_created\": 1}).limit(1)", collection, mongo.Query(q))
		return c.Find(q).Sort("date_created").Limit(1).One(&s)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Get", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Get", "Completed")
	return &s, nil
}
