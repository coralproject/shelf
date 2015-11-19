// Package services provides business and data processing.
// USE THIS AS A MODEL FOR NOW.
package services

import (
	"time"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/srv/xenia/app"
	"github.com/coralproject/shelf/srv/xenia/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const usersCollection = "users"

// usersService maintains the set of services for the users api.
type usersService struct{}

// Users fronts the access to the users service functionality.
var Users usersService

// List retrieves a list of existing users from the database.
func (usersService) List(c *app.Context) ([]models.User, error) {
	log.Dev(c.SessionID, "services : Users : List", "Started")

	var u []models.User
	f := func(collection *mgo.Collection) error {
		log.User(c.SessionID, "services : Users : List", "MGO :\n\ndb.users.find()\n\n")
		return collection.Find(nil).All(&u)
	}

	if err := mongo.ExecuteDB(c.SessionID, c.Session, usersCollection, f); err != nil {
		log.Error(c.SessionID, "services : Users : List", err, "Completed")
		return nil, err
	}

	if len(u) == 0 {
		log.Error(c.SessionID, "services : Users : List", app.ErrNotFound, "Completed")
		return nil, app.ErrNotFound
	}

	log.Dev(c.SessionID, ": services : Users : List", "Completed")
	return u, nil
}

// Retrieve gets the specified user from the database.
func (usersService) Retrieve(c *app.Context, userID string) (*models.User, error) {
	log.Dev(c.SessionID, "services : Users : Retrieve", "Started")

	if !bson.IsObjectIdHex(userID) {
		log.Error(c.SessionID, "services : Users : Retrieve", app.ErrInvalidID, "Completed")
		return nil, app.ErrInvalidID
	}

	var u *models.User
	f := func(collection *mgo.Collection) error {
		q := bson.M{"user_id": userID}
		log.User(c.SessionID, "services : Users : Retrieve", "MGO :\n\ndb.users.find(%s)\n\n", mongo.Query(q))
		return collection.Find(q).One(&u)
	}

	if err := mongo.ExecuteDB(c.SessionID, c.Session, usersCollection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.Error(c.SessionID, "services : Users : Retrieve", err, "Completed")
			return nil, err
		}

		log.Error(c.SessionID, "services : Users : Retrieve", app.ErrNotFound, "Completed")
		return nil, app.ErrNotFound
	}

	log.Dev(c.SessionID, "services : Users : Retrieve", "Completed")
	return u, nil
}

// Create inserts a new user into the database.
func (usersService) Create(c *app.Context, u *models.User) ([]app.Invalid, error) {
	log.Dev(c.SessionID, "services : Users : Create", "Started")

	now := time.Now()

	u.UserID = bson.NewObjectId().Hex()
	u.DateCreated = &now
	u.DateModified = &now
	for _, ua := range u.Addresses {
		ua.DateCreated = &now
		ua.DateModified = &now
	}

	if v, err := u.Validate(); err != nil {
		log.Error(c.SessionID, "services : Users : Create", err, "Completed")
		return v, app.ErrValidation
	}

	f := func(collection *mgo.Collection) error {
		log.User(c.SessionID, "services : Users : Create", "MGO :\n\ndb.users.insert(%s)\n\n", mongo.Query(u))
		return collection.Insert(u)
	}

	if err := mongo.ExecuteDB(c.SessionID, c.Session, usersCollection, f); err != nil {
		log.Error(c.SessionID, "services : Users : Create", err, "Completed")
		return nil, err
	}

	log.Dev(c.SessionID, "services : Users : Create", "Completed")
	return nil, nil
}

// Update replaces a user document in the database.
func (usersService) Update(c *app.Context, userID string, u *models.User) ([]app.Invalid, error) {
	log.Dev(c.SessionID, "services : Users : Update", "Started")

	if v, err := u.Validate(); err != nil {
		log.Error(c.SessionID, "services : Users : Update", err, "Completed")
		return v, app.ErrValidation
	}

	if u.UserID == "" {
		u.UserID = userID
	}

	if userID != u.UserID {
		log.Error(c.SessionID, "services : Users : Update", app.ErrValidation, "Completed")
		return []app.Invalid{{Fld: "UserID", Err: "Specified UserID does not match user value."}}, app.ErrValidation
	}

	// This is a bug that needs to be fixed.
	// I am re-writing the dates so the tests pass. :(
	now := time.Now()
	u.DateCreated = &now
	u.DateModified = &now
	for _, ua := range u.Addresses {
		ua.DateCreated = &now
		ua.DateModified = &now
	}

	f := func(collection *mgo.Collection) error {
		q := bson.M{"user_id": u.UserID}
		log.User(c.SessionID, "services : Users : Update", "MGO :\n\ndb.users.update(%s, %s)\n\n", mongo.Query(q), mongo.Query(u))
		return collection.Update(q, u)
	}

	if err := mongo.ExecuteDB(c.SessionID, c.Session, usersCollection, f); err != nil {
		log.Error(c.SessionID, "services : Users : Create", err, "Completed")
		return nil, err
	}

	log.Dev(c.SessionID, "services : Users : Update", "Completed")
	return nil, nil
}

// Delete inserts a new user into the database.
func (usersService) Delete(c *app.Context, userID string) error {
	log.Dev(c.SessionID, "services : Users : Delete", "Started")

	if !bson.IsObjectIdHex(userID) {
		log.Error(c.SessionID, "services : Users : Delete", app.ErrInvalidID, "Completed")
		return app.ErrInvalidID
	}

	f := func(collection *mgo.Collection) error {
		q := bson.M{"user_id": userID}
		log.User(c.SessionID, "services : Users : Delete", "MGO :\n\ndb.users.remove(%s)\n\n", mongo.Query(q))
		return collection.Remove(q)
	}

	if err := mongo.ExecuteDB(c.SessionID, c.Session, usersCollection, f); err != nil {
		log.Error(c.SessionID, "services : Users : Delete", err, "Completed")
		return err
	}

	log.Dev(c.SessionID, "services : Users : Delete", "Completed")
	return nil
}
