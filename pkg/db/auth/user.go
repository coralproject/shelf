package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/log"

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

// LoginUser model used for when a user logs in.
type LoginUser struct {
	Email    string `json:"email" form:"email" valid:"Required;Email;MaxSize(150)"`
	Password string `json:"password" form:"email" valid:"Required;"`
}

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
	Token        string        `bson:"-" json:"token"`
	IsDeleted    bool          `bson:"is_deleted" json:"-"`
	DateModified time.Time     `bson:"date_modified" json:"-"`
	DateCreated  time.Time     `bson:"date_created" json:"-"`
}

// AuthenticateToken authenticates a User entities token.
func (u *User) AuthenticateToken(token string) error {
	if err := IsTokenValidForEntity(u, token); err != nil {
		return err
	}

	return nil
}

// SetToken sets the Token for the User.
func (u *User) SetToken(sID string) error {
	t, err := TokenforEntity(u)
	if err != nil {
		return err
	}

	u.Token = base64.StdEncoding.EncodeToString([]byte(sID + ":" + base64.StdEncoding.EncodeToString(t)))
	return nil
}

// IsPasswordValid compares Password string to Stored Password.
func (u *User) IsPasswordValid(pwd string) bool {
	if u.Password == "" {
		return false
	}

	// Hashed Password comes first, then the plain text version.
	if err := CompareBcryptHashPassword([]byte(u.Password), []byte(u.PrivateID+pwd)); err != nil {
		return false
	}

	return true
}

// Pwd implements the secure entity interface.
func (u *User) Pwd() ([]byte, error) {
	if u.Password == "" {
		return nil, errors.New("User PWD is Blank")
	}

	return []byte(u.Password), nil
}

// Salt implements the secure entity interface.
func (u *User) Salt() ([]byte, error) {
	if (u.PublicID == "" || u.PrivateID == "" || u.DateCreated == time.Time{}) {
		return nil, errors.New("Unable to Generate user Token, Missing Data")
	}

	s := u.PublicID + ":" + u.PrivateID + ":" + fmt.Sprintf("%d", u.DateCreated.UTC().Unix())

	return []byte(s), nil
}

// NewUser is used for creating a properly initialized user value.
type NewUser struct {
	FullName        string `msgpack:"first_name" json:"first_name"`
	Email           string `msgpack:"email" json:"email"`
	Password        string `msgpack:"password" json:"password"`
	PasswordConfirm string `msgpack:"password_confirm" json:"password_confirm"`
	PostalCode      string `msgpack:"postal_code" json:"postal_code"`
}

// Validate performs validation on a NewUser value before it is processed.
func (nu *NewUser) Validate(context interface{}, v *validation.Validation) {
	log.Dev(context, "NewUser.Validate", "Started")

	v.Required(nu.FullName, "full_name")
	v.AlphaNumeric(nu.FullName, "full_name")
	v.MinSize(nu.FullName, 2, "full_name")

	v.Required(nu.Email, "email")
	v.Email(nu.Email, "email")
	v.MaxSize(nu.Email, 100, "email")

	v.Required(nu.Password, "password")
	v.MinSize(nu.Password, 8, "password")

	v.Required(nu.PasswordConfirm, "password_confirm")
	v.MinSize(nu.PasswordConfirm, 8, "password_confirm")

	if nu.Password != nu.PasswordConfirm {
		v.SetError("password", "Does not match password confirm")
	}

	log.Dev(context, "NewUser.Validate", "Completed : HasErrors[%v]", v.HasErrors())
}

// Create creates a User value from a NewUser value.
func (nu *NewUser) Create(context interface{}) (*User, error) {
	log.Dev(context, "NewUser.Create", "Started : Email[%s]", nu.Email)

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
	if u.Password, err = BcryptPassword(u.PrivateID + nu.Password); err != nil {
		log.Error(context, "NewUser.Create", err, "Generating Password Hash")
		return nil, err
	}

	log.Dev(context, "NewUser.Create", "Completed")
	return &u, nil
}

// UserUpdate is the model to update a user.
type UserUpdate struct {
	FullName string `bson:"full_name" json:"full_name,omitempty"`
	Email    string `bson:"email" json:"email,omitempty"`
	Status   int    `bson:"status" json:"status,omitempty"`
}

// Validate implements the Validation interface.
func (uu *UserUpdate) Validate(context interface{}, v *validation.Validation) {
	log.Dev(context, "UserUpdate.Validate", "Started")

	if uu.FullName != "" {
		v.AlphaNumeric(uu.FullName, "full_name")
		v.MinSize(uu.FullName, 2, "full_name")
	}

	if uu.Email != "" {
		v.Email(uu.Email, "email")
		v.MaxSize(uu.Email, 100, "email")
	}

	if uu.Status != 0 {
		if uu.Status != StatusActive && uu.Status != StatusDisabled {
			v.SetError("Status", "has an invalid value")
		}
	}

	log.Dev(context, "UserUpdate.Validate", "Completed : HasErrors[%v]", v.HasErrors())
}
