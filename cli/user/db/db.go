package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User represents an entity for a user document.
type User struct {
	ID         bson.ObjectId `json:"id" bson:"id"`
	Name       string        `json:"name,omitempty" bson:"name"`
	Password   string        `json:"password" bson:"password" `
	Token      string        `json:"token,omitempty" bson:"-"`
	PublicID   string        `json:"public_id,omitempty" bson:"public_id"`
	PrivateID  string        `json:"private_id,omitempty" bson:"private_id"`
	ModifiedAt *time.Time    `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time    `json:"created_at,omitempty" bson:"created_at"`
}

// Create adds a new user to the database.
func Create(u *User) error {
	// TODO: Model from code in xenia that talks to the DB.
	return nil
}

// Update modifies an existing user in the database.
func Update(u *User) error {
	// TODO: Model from code in xenia that talks to the DB.
	return nil
}

// Delete removes an existing user from the database.
func Delete(u *User) error {
	// TODO: Model from code in xenia that talks to the DB.
	return nil
}
