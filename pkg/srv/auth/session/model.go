package session

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Session denotes a user's session within the system.
type Session struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	SessionID   string        `bson:"session_id" json:"session_id"`
	PublicID    string        `bson:"public_id" json:"public_id"`
	DateExpires time.Time     `bson:"date_expires" json:"date_expires"`
	DateCreated time.Time     `bson:"date_created" json:"date_created"`
}
