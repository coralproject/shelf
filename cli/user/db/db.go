package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coralproject/shelf/cli/user/crypto"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"

	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Global variables pointing to the user's collection in the database.
// TODO: read this from configuration.
var UserCollection = "users"

// User represents an entity for a user document.
type User struct {
	ID         bson.ObjectId `json:"id" bson:"id"`
	Name       string        `json:"name,omitempty" bson:"name"`
	Email      string        `json:"email" bson:"email"`
	Password   string        `json:"password" bson:"password" `
	Token      string        `json:"token,omitempty" bson:"-"`
	PublicID   string        `json:"public_id,omitempty" bson:"public_id"`
	PrivateID  string        `json:"private_id,omitempty" bson:"private_id"`
	ModifiedAt *time.Time    `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time    `json:"created_at,omitempty" bson:"created_at"`
}

// NewUser creates a new user entity.
func NewUser(name, email, password string) (*User, error) {
	// Generate the users PublicID and PrivateID
	pubID := uuid.New()
	privID := uuid.New()

	// Generate the user encrypted password using the supplied password
	pass, err := crypto.BcryptHash(privID + password)
	if err != nil {
		return nil, err
	}

	var user User

	user.ID = bson.NewObjectId()
	user.Name = name
	user.Email = strings.ToLower(email)
	user.Password = pass
	user.PublicID = pubID
	user.PrivateID = privID

	created := time.Now().UTC()
	mod := time.Now().UTC()

	user.CreatedAt = &created
	user.ModifiedAt = &mod

	if err := user.SetToken(); err != nil {
		return nil, err
	}

	return &user, nil
}

// Salt returns the user's password salt.
func (u *User) Salt() ([]byte, error) {
	if !u.hasCrendentials() {
		return nil, errors.New("Invalid User Entity")
	}

	return []byte(u.PublicID + u.PrivateID + fmt.Sprintf("%v", u.CreatedAt.UTC())), nil
}

// Pwd returns the users password.
func (u *User) Pwd() ([]byte, error) {
	if !u.hasCrendentials() {
		return nil, errors.New("Invalid Crendentials")
	}

	return []byte(u.Password), nil
}

// IsPasswordValid validates if the given password matches the user password.
func (u *User) IsPasswordValid(pass string) bool {
	if !u.hasCrendentials() {
		return false
	}

	if err := crypto.CompareBase64BcryptHash(u.Password, (u.PrivateID + pass)); err != nil {
		return false
	}

	return true
}

// IsTokenValid validates if the given token matches the user entities token.
func (u *User) IsTokenValid(token string) bool {
	if err := crypto.IsTokenValidForEntity(u, token); err != nil {
		return false
	}

	return true
}

// SetToken sets the users authentication token
func (u *User) SetToken() error {
	token, err := crypto.TokenForEntity(u)
	if err != nil {
		return err
	}

	u.Token = crypto.Base64Token(token)
	return nil
}

// hasCrendentials returns true/false if authentication required fields are
// blank or unloaded.
func (u *User) hasCrendentials() bool {
	if u.Email == "" || u.PrivateID == "" || u.PublicID == "" {
		return false
	}

	return true
}

// GetUserByEmail retrieves a user record by using the provided email.
// It returns a non-nil error if no record was found.
func GetUserByEmail(email string) (*User, error) {
	log.Dev(email, "GetUserByEmail", "Started : Find User")
	var user User

	// All emails will be in lowercase.
	email = strings.ToLower(email)

	f := func(c *mgo.Collection) error {
		log.Dev(email, "GetUserByEmail", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"email": email}).One(&user)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(email, "GetUserByEmail", "Started : User : Mongodb.Find().One()")
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(email, "GetUserByEmail", err, "Completed")
		return nil, err
	}

	log.Dev(email, "GetUserByEmail", "Completed : Find User")
	return &user, nil
}

// GetUserByName retrieves a user record by using the provided name.
// It returns a non-nil error if no record was found.
func GetUserByName(name string) (*User, error) {
	log.Dev(name, "GetUserByName", "Started : Find User")
	var user User

	f := func(c *mgo.Collection) error {
		log.Dev(name, "GetUserByName", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"name": name}).One(&user)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(name, "GetUserByName", "Started : User : Mongodb.Find().One()")
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(name, "GetUserByName", err, "Completed")
		return nil, err
	}

	log.Dev(name, "GetUserByName", "Completed : Find User")
	return &user, nil
}

// GetUserByPublicID retrieves a user record by using the provided PublicID.
// It returns a non-nil error if no record was found.
func GetUserByPublicID(pid string) (*User, error) {
	log.Dev(pid, "GetUserByPublicID", "Started : Find User")
	var user User

	f := func(c *mgo.Collection) error {
		log.Dev(pid, "GetUserByPublicID", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"public_id": pid}).One(&user)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(pid, "GetUserByPublicID", "Started : User : Mongodb.Find().One()")
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(pid, "GetUserByPublicID", err, "Completed")
		return nil, err
	}

	log.Dev(pid, "GetUserByPublicID", "Completed : Find User")
	return &user, nil
}

// Create adds a new user to the database.
// Returns a non-nil error, if the operation fails.
func Create(u *User) error {
	log.Dev(u.PublicID, "Create", "Started : Create User")

	meta := fmt.Sprintf("{ Email: %q Name: %q}", u.Email, u.Name)
	log.Dev(u.PublicID, "Create", "Started : Check User Exists : User %s", UserCollection, meta)

	// All emails will be in lowercase.
	u.Email = strings.ToLower(u.Email)

	f := func(col *mgo.Collection) error {
		count, err := col.Find(bson.M{"email": u.Email}).Count()
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("User Already Exists")
		}

		return nil
	}

	ses := mongo.GetSession()
	defer ses.Close()

	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(u.PublicID, "Create", err, "Completed")
		return err
	}

	log.Dev(u.PublicID, "Create", "Completed : Check User Exists : User %s ", UserCollection, meta)

	f = func(col *mgo.Collection) error {
		return col.Insert(u)
	}

	log.Dev(u.PublicID, "Create", "Add User : Collection %s", UserCollection)
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(u.PublicID, "Create", err, "Completed")
		return err
	}

	log.Dev(u.PublicID, "Create", "Completed : Create User ")
	return nil
}

// UpdateName updates an existing user's name in the database.
// Uses the user entity's id as the update parameter.
// Returns a non-nil error, if the operation fails.
func UpdateName(u *User, name string) error {
	log.Dev(u.PublicID, "UpdateName", "Started : Updating User Record")

	ms := time.Now().UTC()

	updateBson := bson.M{"name": name, "email": u.Email, "password": u.Password, "modified_at": &ms}

	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateName", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(u.PublicID, "UpdateName", "Started : Mongodb.Update()")
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(u.PublicID, "UpdateName", err, "Completed")
		return err
	}

	u.Name = name
	u.ModifiedAt = &ms

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

	updateBson := bson.M{"name": u.Name, "email": email, "password": u.Password, "modified_at": &ms}

	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateEmail", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	log.Dev(u.PublicID, "UpdateEmail", "Started : Mongodb.Update()")
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(u.PublicID, "UpdateEmail", err, "Completed")
		return err
	}

	u.Email = email
	u.ModifiedAt = &ms

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
	if err := u.SetToken(); err != nil {
		log.Error(u.PublicID, "UpdatePassword", err, "Completed")
		return err
	}
	log.Dev(u.PublicID, "UpdatePassword", "Completed : User : SetToken : Success")

	log.Dev(u.PublicID, "UpdatePassword", "Started : Validate NewUser Password %s", mongo.Query(newPassword))
	if !u.IsPasswordValid(newPassword) {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Validate New User Password %s : Error %s", mongo.Query(newPassword), "Invalid Password")
		return errors.New("Invalid Password")
	}
	log.Dev(u.PublicID, "UpdatePassword", "Compeleted : Validate New User Password %s : Success", mongo.Query(newPassword))

	ms := time.Now().UTC()
	u.ModifiedAt = &ms

	log.Dev(u.PublicID, "UpdatePassword", "Started : Mongodb.UpdateId()")

	updateBson := bson.M{"name": u.Name, "email": u.Email, "password": newPassHash, "modified_at": &ms}
	f := func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdatePassword", "Completed : Mongodb.UpdateId()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	}

	ses := mongo.GetSession()
	defer ses.Close()

	if err = mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
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
	if err := mongo.ExecuteDB("CONTEXT", ses, UserCollection, f); err != nil {
		log.Error(u.PublicID, "Delete", err, "Completed")
		return err
	}

	log.Dev(u.PublicID, "Delete", "Completed : Delete User")
	return nil
}
