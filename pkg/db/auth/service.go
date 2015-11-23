package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/ArdanStudios/aggserver/auth/crypto"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collection = "user"
const serviceID = ""

// GetUserByEmail retrieves a user record by using the provided email.
func GetUserByEmail(context interface{}, ses *mgo.Session, email string) (*User, error) {
	log.Dev(context, "GetUserByEmail", "Started : Email[%s]", email)

	email = strings.ToLower(email)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"email": email}
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
func GetUserByPublicID(context interface{}, ses *mgo.Session, pid string) (*User, error) {
	log.Dev(context, "GetUserByPublicID", "Started : PID[%s]", pid)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": pid}
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

// Create adds a new user to the database.
func Create(context interface{}, ses *mgo.Session, u *User) error {
	log.Dev(context, "Create", "Started : PublicID[%s]", u.PublicID)

	f := func(col *mgo.Collection) error {
		q := bson.M{"email": u.Email}
		log.Dev(context, "Create", "MGO : db.%s.find(%s).count()", collection, mongo.Query(q))
		count, err := col.Find(q).Count()
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("User Already Exists")
		}

		return nil
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	f = func(col *mgo.Collection) error {
		log.Dev(context, "Create", "MGO : db.%s.insert(%s)", collection, mongo.Query(u))
		return col.Insert(u)
	}

	if err := mongo.ExecuteDB("CONTEXT", ses, collection, f); err != nil {
		log.Error(context, "Create", err, "Completed")
		return err
	}

	log.Dev(context, "Create", "Completed")
	return nil
}

// UpdateName updates an existing user's name in the database.
// Uses the user entity's id as the update parameter.
// Returns a non-nil error, if the operation fails.
func UpdateName(u *User, name string) error {
	log.Dev(u.PublicID, "UpdateName", "Started : Updating User Record")

	ms := time.Now().UTC()

	updateBson := bson.M{"name": name, "email": u.Email, "password": u.Password, "date_modified": ms}

	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateName", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(u.PublicID, "UpdateName", "Started : Mongodb.Update()")
	if err := mongo.ExecuteDB("CONTEXT", ses, collection, f); err != nil {
		log.Error(u.PublicID, "UpdateName", err, "Completed")
		return err
	}

	u.FullName = name
	u.DateModified = ms

	log.Dev(u.PublicID, "UpdateName", "Completed : Updating User Record")
	return nil
}

// UpdateEmail updates the email for a user record.
// Uses the user entity's id as the update parameter.
// Returns a non-nil error, if the operation fails.
func UpdateEmail(u *User, email string) error {
	log.Dev(u.PublicID, "UpdateEmail", "Started : Updating User Record")

	// All emails will be in lowercase.
	email = strings.ToLower(email)

	ms := time.Now().UTC()

	updateBson := bson.M{"name": u.FullName, "email": email, "password": u.Password, "date_modified": ms}

	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateEmail", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(u.PublicID, "UpdateEmail", "Started : Mongodb.Update()")
	if err := mongo.ExecuteDB("CONTEXT", ses, collection, f); err != nil {
		log.Error(u.PublicID, "UpdateEmail", err, "Completed")
		return err
	}

	u.Email = email
	u.DateModified = ms

	log.Dev(u.PublicID, "UpdateEmail", "Completed : Updated User Record")
	return nil
}

// UpdatePassword updates an existing user's password in the database.
// Uses the user entity's id as the update parameter.
// Requires provision of the old password and the new password.
// Returns a non-nil error, if the existingPassword is not a match, or
// the update operation fails.
func UpdatePassword(u *User, existingPassword, newPassword string) error {
	log.Dev(u.PublicID, "UpdatePassword", "Started : Updating User Record")

	log.Dev(u.PublicID, "UpdatePassword", "Started : Validate User Existing Password %s", mongo.Query(existingPassword))
	if !u.IsPasswordValid(existingPassword) {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Validate User Existing Password %s : Error %s", mongo.Query(existingPassword), "Invalid Password")
		return errors.New("Invalid Password")
	}
	log.Dev(u.PublicID, "UpdatePassword", "Compeleted : Validate User Existing Password %s : Success", mongo.Query(existingPassword))

	log.Dev(u.PublicID, "UpdatePassword", "Started : Create New Password %s", mongo.Query(newPassword))
	newPassHash, err := crypto.BcryptHash((u.PrivateID + newPassword))
	if err != nil {
		log.Error(u.PublicID, "UpdatePassword", err, "Completed")
		return err
	}

	log.Dev(u.PublicID, "UpdatePassword", "Completed : Create New Password %s : Success", mongo.Query(newPassword))
	u.Password = newPassHash

	log.Dev(u.PublicID, "UpdatePassword", "Started : User : SetToken")
	if err := u.SetToken(serviceID); err != nil {
		log.Error(u.PublicID, "UpdatePassword", err, "Completed")
		return err
	}
	log.Dev(u.PublicID, "UpdatePassword", "Completed : User : SetToken : Success")

	log.Dev(u.PublicID, "UpdatePassword", "Started : Validate NewUser Password %s", mongo.Query(newPassword))
	if !u.IsPasswordValid(newPassword) {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Validate New User Password %s : Error %s", mongo.Query(newPassword), "Invalid Password")
		return errors.New("Invalid Password")
	}
	log.Dev(u.PublicID, "UpdatePassword", "Completed : Validate New User Password %s : Success", mongo.Query(newPassword))

	ms := time.Now().UTC()
	u.DateModified = ms

	log.Dev(u.PublicID, "UpdatePassword", "Started : Mongodb.UpdateId()")

	updateBson := bson.M{"name": u.FullName, "email": u.Email, "password": newPassHash, "date_modified": ms}
	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdatePassword", "Completed : Mongodb.UpdateId()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	if err = mongo.ExecuteDB("CONTEXT", ses, collection, f); err != nil {
		log.Error(u.PublicID, "UpdatePassword", err, "Completed")
		return nil
	}

	log.Dev(u.PublicID, "UpdatePassword", "Completed : Updated User Record")
	return nil
}

// Delete removes an existing user from the database.
// Returns a non-nil error, if the operation fails.
func Delete(u *User) error {
	log.Dev(u.PublicID, "Delete", "Started : Delete User")

	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "Delete", "Completed : Mongodb.RemoveId()")
		return c.Remove(bson.M{"id": u.ID})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(u.PublicID, "Delete", "Started : Mongodb.RemoveId()")
	if err := mongo.ExecuteDB("CONTEXT", ses, collection, f); err != nil {
		log.Error(u.PublicID, "Delete", err, "Completed")
		return err
	}

	log.Dev(u.PublicID, "Delete", "Completed : Delete User")
	return nil
}
