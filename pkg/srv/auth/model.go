// Package auth provides CRUD methods for the auth user API.
package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/srv/auth/crypto"

	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"
)

// Set of user status codes.
const (
	StatusUnknown = iota
	StatusActive
	StatusDisabled
	StatusDeleted
	StatusInvalid
)

// Set of user type codes.
const (
	TypeAPI = iota + 1
)

//==============================================================================

// LoginUser model used for when a user logs in.
type LoginUser struct {
	Email    string `json:"email" form:"email" valid:"Required;Email;MaxSize(150)"`
	Password string `json:"password" form:"email" valid:"Required;"`
}

//==============================================================================

// User model denotes a user entity for a tenant.
type User struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"-"`
	PublicID     string        `bson:"public_id" json:"public_id"`
	PrivateID    string        `bson:"private_id" json:"-"`
	UserType     int           `bson:"type" json:"type"`
	Status       int           `bson:"status" json:"status"`
	FullName     string        `bson:"full_name" json:"full_name"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"password" json:"-"`
	IsDeleted    bool          `bson:"is_deleted" json:"-"`
	DateModified time.Time     `bson:"date_modified" json:"-"`
	DateCreated  time.Time     `bson:"date_created" json:"-"`
}

// Pwd implements the secure entity interface.
func (u *User) Pwd() ([]byte, error) {
	if u.Password == "" {
		return nil, errors.New("User password is blank")
	}

	return []byte(u.Password), nil
}

// Salt implements the secure entity interface.
func (u *User) Salt() ([]byte, error) {
	if (u.PublicID == "" || u.PrivateID == "" || u.DateCreated == time.Time{}) {
		return nil, errors.New("Unable to generate user token, missing data")
	}

	s := u.PublicID + ":" + u.PrivateID + ":" + fmt.Sprintf("%d", u.DateCreated.UTC().Unix())

	return []byte(s), nil
}

// WebToken returns a token ready for web use.
func (u *User) WebToken(sessionID string) (string, error) {
	t, err := crypto.GenerateToken(u)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString([]byte(sessionID + ":" + base64.StdEncoding.EncodeToString(t)))
	return token, nil
}

// IsPasswordValid compares the user provided password with what is in the db.
func (u *User) IsPasswordValid(password string) bool {
	if u.Password == "" {
		return false
	}

	// Hashed Password comes first, then the plain text version.
	if err := crypto.CompareBcryptHashPassword([]byte(u.Password), []byte(u.PrivateID+password)); err != nil {
		return false
	}

	return true
}

//==============================================================================

// NewUser is provided to create new users in the system.
type NewUser struct {
	UserType int    `bson:"type" json:"type"`
	Status   int    `bson:"status" json:"status"`
	FullName string `bson:"full_name" json:"full_name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
}

// validate performs validation on a NewUser value before it is processed.
func (nu *NewUser) validate(context interface{}) error {
	var v validation.Validation

	v.Required(nu.FullName, "FullName")
	v.MinSize(nu.FullName, 2, "FullName")

	v.Required(nu.Email, "Email")
	v.Email(nu.Email, "Email")
	v.MaxSize(nu.Email, 100, "Email")

	v.Required(nu.Password, "Password")
	v.MinSize(nu.Password, 8, "Password")

	if v.HasErrors() {
		return fmt.Errorf("%v", v.ErrorsMap)
	}

	return nil
}

// new takes a new user and creates a valid User value.
func (nu *NewUser) new(context interface{}) (*User, error) {
	u := User{
		PublicID:     uuid.New(),
		PrivateID:    uuid.New(),
		Status:       StatusActive,
		FullName:     nu.FullName,
		Email:        strings.ToLower(nu.Email),
		DateModified: time.Now(),
		DateCreated:  time.Now(),
		IsDeleted:    false,
	}

	var err error
	if u.Password, err = crypto.BcryptPassword(u.PrivateID + nu.Password); err != nil {
		return nil, err
	}

	return &u, nil
}

//==============================================================================

// UpdUser is provided to update an existing user in the system.
type UpdUser struct {
	PublicID string `bson:"public_id" json:"public_id"`
	UserType int    `bson:"type" json:"type"`
	Status   int    `bson:"status" json:"status"`
	FullName string `bson:"full_name" json:"full_name"`
	Email    string `bson:"email" json:"email"`
}

// validate performs validation on a NewUser value before it is processed.
func (uu *UpdUser) validate(context interface{}) error {
	var v validation.Validation

	v.Required(uu.PublicID, "public_id")
	v.MinSize(uu.PublicID, 2, "public_id")

	v.Required(uu.FullName, "full_name")
	v.AlphaNumeric(uu.FullName, "full_name")
	v.MinSize(uu.FullName, 2, "full_name")

	v.Required(uu.Email, "email")
	v.Email(uu.Email, "email")
	v.MaxSize(uu.Email, 100, "email")

	if v.HasErrors() {
		return fmt.Errorf("%v", v.Errors)
	}

	return nil
}
