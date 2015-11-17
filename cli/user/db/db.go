package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coralproject/shelf/cli/user/crypto"
	"github.com/coralproject/shelf/log"
	"github.com/satori/go.uuid"

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
	pubID := uuid.NewV4()
	privID := uuid.NewV4()

	// Generate the user encrypted password using the supplied password
	pass, err := crypto.BcryptHash(privID.String() + password)
	if err != nil {
		return nil, err
	}

	var user User

	user.ID = bson.NewObjectId()
	user.Name = name
	user.Email = strings.ToLower(email)
	user.Password = pass
	user.PublicID = pubID.String()
	user.PrivateID = privID.String()

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

// SetToken sets the users authentication token
func (u *User) SetToken() error {
	token, err := crypto.TokenForEntity(u)
	if err != nil {
		return err
	}

	u.Token = crypto.Base64Token(token)
	return nil
}

// UserHeader represents a user entity's publically allowed fields, for safe
// transmission over the wire.
type UserHeader struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	PublicID string `json:"public_id"`
}

// SanitizedUser returns a User struct which only contain the safe and publically
// allowed fields. Useful for creating a JWT Header.
func (u *User) SanitizedUser() (*UserHeader, error) {
	if !u.hasCrendentials() {
		return nil, errors.New("Invalid User Entity")
	}

	return &UserHeader{
		ID:       u.ID.String(),
		Name:     u.Name,
		Email:    u.Email,
		Token:    u.Token,
		PublicID: u.PublicID,
	}, nil
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

	log.Dev(email, "GetUserByEmail", "Started : User : Mongodb.Find().One()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(email, "GetUserByEmail", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"email": email}).One(&user)
	})

	if err != nil {
		log.Dev(email, "GetUserByEmail", "Completed Error : Error %s", err.Error())
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

	log.Dev(name, "GetUserByName", "Started : User : Mongodb.Find().One()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(name, "GetUserByName", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"name": name}).One(&user)
	})

	if err != nil {
		log.Dev(name, "GetUserByName", "Completed Error : Error %s", err.Error())
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

	log.Dev(pid, "GetUserByPublicID", "Started : User : Mongodb.Find().One()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(pid, "GetUserByPublicID", "Completed : User : Mongodb.Find().One()")
		return c.Find(bson.M{"public_id": pid}).One(&user)
	})

	if err != nil {
		log.Dev(pid, "GetUserByPublicID", "Completed Error : Error %s", err.Error())
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

	err := ExecuteDB(GetSession(), UserCollection, func(col *mgo.Collection) error {
		count, err := col.Find(bson.M{"email": u.Email}).Count()
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("User Already Exists")
		}

		return nil
	})

	if err != nil {
		log.Dev(u.PublicID, "Create", "Completed Error : Check User Exists :  User %s : Error %q", meta, err.Error())
		return err
	}

	log.Dev(u.PublicID, "Create", "Completed : Check User Exists : User %s ", UserCollection, meta)

	log.Dev(u.PublicID, "Create", "Add User : Collection %s", UserCollection)
	err2 := ExecuteDB(GetSession(), UserCollection, func(col *mgo.Collection) error {
		return col.Insert(u)
	})

	if err2 != nil {
		log.Dev(u.PublicID, "Create", "Completed Error : Add User : Error %s", err2.Error())
		return err2
	}

	log.Dev(u.PublicID, "Create", "Completed : Add User : Collection %s", UserCollection)

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

	log.Dev(u.PublicID, "UpdateName", "Started : Mongodb.Update()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateName", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	})

	if err != nil {
		log.Dev(u.PublicID, "UpdateName", "Completed Error : Updating User Record : Error %s", err.Error())
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

	log.Dev(u.PublicID, "UpdateEmail", "Started : Mongodb.Update()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdateEmail", "Completed : Mongodb.Update()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	})

	if err != nil {
		log.Dev(u.PublicID, "UpdateEmail", "Completed Error : Updating User Record : Error %s", err.Error())
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

	log.Dev(u.PublicID, "UpdatePassword", "Started : Validate User Existing Password %s", Query(existingPassword))
	if !u.IsPasswordValid(existingPassword) {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Validate User Existing Password %s : Error %s", Query(existingPassword), "Invalid Password")
		return errors.New("Invalid Password")
	}

	log.Dev(u.PublicID, "UpdatePassword", "Compeleted : Validate User Existing Password %s : Success", Query(existingPassword))

	log.Dev(u.PublicID, "UpdatePassword", "Started : Create New Password %s", Query(newPassword))
	newPassHash, err := crypto.BcryptHash((u.PrivateID + newPassword))
	if err != nil {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Create New Password %s : Error %s", Query(newPassword), err.Error())
		return err
	}

	log.Dev(u.PublicID, "UpdatePassword", "Completed : Create New Password %s : Success", Query(newPassword))
	u.Password = newPassHash

	log.Dev(u.PublicID, "UpdatePassword", "Started : User : SetToken")
	if err := u.SetToken(); err != nil {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : User : SetToken : Error %s", err.Error())
		return err
	}
	log.Dev(u.PublicID, "UpdatePassword", "Completed : User : SetToken : Success")

	ms := time.Now().UTC()
	u.ModifiedAt = &ms

	log.Dev(u.PublicID, "UpdatePassword", "Started : Mongodb.UpdateId()")

	updateBson := bson.M{"name": u.Name, "email": u.Email, "password": newPassHash, "modified_at": &ms}
	err = ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "UpdatePassword", "Completed : Mongodb.UpdateId()")
		return c.Update(bson.M{"id": u.ID}, bson.M{"$set": updateBson})
	})

	if err != nil {
		log.Dev(u.PublicID, "UpdatePassword", "Completed Error : Updating User Record : Error %s", err.Error())
		return nil
	}

	log.Dev(u.PublicID, "UpdatePassword", "Completed : Updating User Record")
	return nil
}

// Delete removes an existing user from the database.
// Returns a non-nil error, if the operation fails.
func Delete(u *User) error {
	log.Dev(u.PublicID, "Delete", "Started : Delete User")

	log.Dev(u.PublicID, "Delete", "Started : Mongodb.RemoveId()")
	err := ExecuteDB(GetSession(), UserCollection, func(c *mgo.Collection) error {
		log.Dev(u.PublicID, "Delete", "Completed : Mongodb.RemoveId()")
		return c.Remove(bson.M{"id": u.ID})
	})

	if err != nil {
		log.Dev(u.PublicID, "Delete", "Completed : Delete User : Error %s", err.Error())
		return err
	}

	log.Dev(u.PublicID, "Delete", "Completed : Delete User")
	return nil
}
