package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/auth/crypto"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collection = "user"

//==============================================================================

// Create adds a new user to the database.
func Create(context interface{}, ses *mgo.Session, nu NewUser) (*User, error) {
	log.Dev(context, "Create", "Started : Email[%s]", nu.Email)

	if err := nu.validate(context); err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	u, err := nu.create(context)
	if err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	f := func(col *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", collection, mongo.Query(&u))
		if err := col.Insert(&u); err != nil {
			return err
		}

		return nil
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return nil, err
	}

	log.Dev(context, "Create", "Completed")
	return u, nil
}

//==============================================================================

// UpdateName updates an existing user's full Name in the database.
func UpdateName(context interface{}, ses *mgo.Session, publicID string, fullName string) error {
	log.Dev(context, "UpdateName", "Started : PublicID[%v] FullName[%s]", publicID, fullName)

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		upd := bson.M{"$set": bson.M{"full_name": fullName, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdateName", "MGO : db.%s.update(%s, %s)", collection, mongo.Query(q), mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateName", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateName", "Completed")
	return nil
}

// UpdateEmail updates an existing user's email in the database.
func UpdateEmail(context interface{}, ses *mgo.Session, publicID string, email string) error {
	log.Dev(context, "UpdateEmail", "Started : PublicID[%v] Email[%s]", publicID, email)

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		upd := bson.M{"$set": bson.M{"email": email, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdateEmail", "MGO : db.%s.update(%s, %s)", collection, mongo.Query(q), mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateEmail", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateEmail", "Completed")
	return nil
}

// UpdatePassword updates an existing user's password in the database.
func UpdatePassword(context interface{}, ses *mgo.Session, publicID, privateID string, password string) error {
	log.Dev(context, "UpdatePassword", "Started : PublicID[%s]", publicID)

	if len(password) < 8 {
		err := errors.New("Invalid password length")
		log.Error(context, "UpdatePassword", err, "Completed")
		return err
	}

	newPassHash, err := crypto.BcryptPassword(privateID + password)
	if err != nil {
		log.Error(context, "UpdatePassword", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		upd := bson.M{"$set": bson.M{"password": newPassHash, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdatePassword", "MGO : db.%s.update(%s, %s)", collection, mongo.Query(q), mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err = mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdatePassword", err, "Completed")
		return nil
	}

	log.Dev(context, "UpdatePassword", "Completed")
	return nil
}

//==============================================================================

// Delete removes an existing user from the database.
func Delete(context interface{}, ses *mgo.Session, publicID string) error {
	log.Dev(context, "Delete", "Started : PublicID[%s]", publicID)

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		log.Dev(context, "Delete", "MGO : db.%s.remove(%s)", collection, mongo.Query(q))
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Delete", err, "Completed")
		return err
	}

	log.Dev(context, "Delete", "Completed")
	return nil
}

//==============================================================================

// GetUserByEmail retrieves a user record by using the provided email.
func GetUserByEmail(context interface{}, ses *mgo.Session, email string) (*User, error) {
	log.Dev(context, "GetUserByEmail", "Started : Email[%s]", email)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"email": strings.ToLower(email)}
		log.Dev(context, "GetUserByEmail", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetUserByEmail", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByEmail", "Completed")
	return &user, nil
}

// GetUserByName retrieves a user record by using the provided name.
func GetUserByName(context interface{}, ses *mgo.Session, fullName string) (*User, error) {
	log.Dev(context, "GetUserByName", "Started : FullName[%s]", fullName)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"full_name": fullName}
		log.Dev(context, "GetUserByName", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetUserByName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByName", "Completed")
	return &user, nil
}

// GetUserByPublicID retrieves a user record by using the provided PublicID.
func GetUserByPublicID(context interface{}, ses *mgo.Session, publicID string) (*User, error) {
	log.Dev(context, "GetUserByPublicID", "Started : PID[%s]", publicID)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		log.Dev(context, "GetUserByName", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetUserByPublicID", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByPublicID", "Completed")
	return &user, nil
}
