package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/log"

	"github.com/astaxie/beego/validation"
	"gopkg.in/mgo.v2/bson"
)

// UserEntity model denotes a user entity for a tenant.
type UserEntity struct {
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

// LoginUser model used for when a user logs in.
type LoginUser struct {
	Email    string `json:"email" form:"email" valid:"Required;Email;MaxSize(150)"`
	Password string `json:"password" form:"email" valid:"Required;"`
}

// AuthenticateToken authenticates a User entities token.
func (user *UserEntity) AuthenticateToken(token string) error {
	if err := IsTokenValidForEntity(user, token); err != nil {
		return err
	}

	return nil
}

// SetToken sets the Token for the User.
func (user *UserEntity) SetToken(sID string) error {
	t, err := TokenforEntity(user)
	if err != nil {
		return err
	}

	user.Token = base64.StdEncoding.EncodeToString([]byte(sID + ":" + base64.StdEncoding.EncodeToString(t)))
	return nil
}

// IsPasswordValid compares Password string to Stored Password.
func (user *UserEntity) IsPasswordValid(pwd string) bool {
	if user.Password == "" {
		return false
	}

	// Hashed Password comes first, then the plain text version.
	if err := CompareBcryptHashPassword([]byte(user.Password), []byte(user.PrivateID+pwd)); err != nil {
		return false
	}

	return true
}

// Pwd implements the secure entity interface.
func (user *UserEntity) Pwd() ([]byte, error) {
	if user.Password == "" {
		return nil, errors.New("User PWD is Blank")
	}

	return []byte(user.Password), nil
}

// Salt implements the secure entity interface.
func (user *UserEntity) Salt() ([]byte, error) {
	if (user.PublicID == "" ||
		user.PrivateID == "" ||
		user.DateCreated == time.Time{}) {

		return nil, errors.New("Unable to Generate user Token, Missing Data")
	}

	s := user.PublicID + ":" + user.PrivateID + ":" + fmt.Sprintf("%d", user.DateCreated.UTC().Unix())

	return []byte(s), nil
}

// NewUser model is used for creating a new user.
type NewUser struct {
	FullName        string `msgpack:"first_name" json:"first_name"`
	Email           string `msgpack:"email" json:"email"`                       // error:"Email is required."`
	Password        string `msgpack:"password" json:"password"`                 // error:"Password is required and must be at least 6 characters."`
	PasswordConfirm string `msgpack:"password_confirm" json:"password_confirm"` // error:"Password confirmation is required and must match password."`
	PostalCode      string `msgpack:"postal_code" json:"postal_code"`
}

// Validate is the interface method for a model to specify what needs validation.
func (nu *NewUser) Validate(v *validation.Validation) {
	// Fullname
	v.Required(nu.FullName, "full_name")
	v.AlphaNumeric(nu.FullName, "full_name")
	v.MinSize(nu.FullName, 2, "full_name")

	// Email
	v.Required(nu.Email, "email")
	v.Email(nu.Email, "email")
	v.MaxSize(nu.Email, 100, "email")

	// Password
	v.Required(nu.Password, "password")
	v.MinSize(nu.Password, 8, "password")

	// Password Confirm
	v.Required(nu.PasswordConfirm, "password_confirm")
	v.MinSize(nu.PasswordConfirm, 8, "password_confirm")

	if nu.Password != nu.PasswordConfirm {
		v.SetError("password", "Does not match password confirm")
	}
}

// Create creates a New UserEntity value from a NewUser value.
func (nu *NewUser) Create(context interface{}) (*UserEntity, error) {
	log.User(context, "Create", "Started : Email[%s]", nu.Email)

	user, err := nu.toUserEntity()
	if err != nil {
		log.Error(context, "Create", err, "Creating Entity User")
		return nil, err
	}

	user.UserType = entityUserType.Customer
	if nu.IsAcctHolder {
		user.UserType = entityUserType.CustomerAccountHolder
	}
	user.Status = entityStatus.Active

	log.User(context, "Create", "Completed : Email[%s], CompanyId[%s]", user.Email, user.CompanyID)
	return user, nil
}

// CreateSysadmin creates a New Customer Entity type from New User type.
func (nu *NewUser) CreateSysadmin() (*SysadminEntity, error) {
	log.Startf("authentication", "NewUser.CreateSysadmin - Email[%s], CompanyId[%s]", nu.Email, nu.CompanyID)

	user, err := nu.toUserEntity()
	if err != nil {
		log.Errf(err, "NewUser", "NewUser.CreateSysadmin - Creating Sysadmin User")
		return nil, err
	}

	user.PostalCode = ""
	user.CompanyID = defaultMasterDB
	user.UserType = entityUserType.Sysadmin
	user.Status = entityStatus.Active

	sysadmin := SysadminEntity{UserEntity: *user}

	log.Completef("authentication", "NewUser.CreateSysadmin - Email[%s], CompanyId[%s]", user.Email, user.CompanyID)
	return &sysadmin, nil
}

// toUserEntity creates a userentity from a new user type. Used when creating a new user.
func (nu *NewUser) toUserEntity() (*UserEntity, error) {
	log.Start("authentication")

	publicIDBytes, err := uuid.NewV4()
	if err != nil {
		log.Errf(err, "authentication", "NewUser.toUserEntity - Generating Public Id")
		return nil, err
	}

	privateIDBytes, err2 := uuid.NewV4()
	if err2 != nil {
		log.Errf(err2, "authentication", "NewUser.toUserEntity - Generating Private Id")
		return nil, err2
	}

	publicID := publicIDBytes.String()
	privateID := privateIDBytes.String()

	if nu.Language == language.UnknownLanguage {
		nu.Language = language.USEnglish
	}

	if nu.Currency == currency.UnknownCurrency {
		nu.Currency = currency.USDollar
	}

	entityUser := UserEntity{
		PublicID:     publicID,
		PrivateID:    privateID,
		Status:       entityStatus.Active,
		Language:     nu.Language,
		Currency:     nu.Currency,
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		FullName:     fmt.Sprintf("%s %s", strings.ToLower(nu.FirstName), strings.ToLower(nu.LastName)),
		Email:        strings.ToLower(nu.Email),
		PostalCode:   nu.PostalCode,
		CompanyID:    nu.CompanyID,
		Roles:        nu.Roles,
		DateModified: time.Now(),
		DateCreated:  time.Now(),
		IsDeleted:    false,
	}

	entityUser.Password, err = crypto.BcryptPassword(entityUser.PrivateID + nu.Password)
	if err != nil {
		log.Errf(err, "authentication", "NewUser.toUserEntity - Generating Password Hash")
		return nil, err
	}

	log.Complete("authentication")
	return &entityUser, nil
}

// UserUpdate is the model to update a user.
type UserUpdate struct {
	FirstName  string `msgpack:"first_name,omitempty" bson:"first_name" json:"first_name,omitempty"`
	LastName   string `msgpack:"last_name,omitempty" bson:"last_name" json:"last_name,omitempt"`
	Email      string `msgpack:"email,omitempty" bson:"email" json:"email,omitempty"`
	PostalCode string `msgpack:"postal_code,omitempty" bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	Roles      []int  `msgpack:"roles,omitempty" json:"roles,omitempty" bson:"roles"`
	Status     int    `msgpack:"status,omitempty" json:"status,omitempty" bson:"status"`
	MaxRole    int
}

// Name returns the formatted name of the user.
func (up *UserUpdate) Name() string {
	return fmt.Sprintf("%s %s", strings.ToTitle(up.FirstName), strings.ToTitle(up.LastName))
}

// Validate implements the Validation interface.
func (up *UserUpdate) Validate(v *validation.Validation) {
	// FirstName
	if up.FirstName != "" {
		v.AlphaNumeric(up.FirstName, "first_name")
		v.MinSize(up.FirstName, 2, "first_name")
	}

	// LastName
	if up.LastName != "" {
		v.AlphaNumeric(up.LastName, "last_name")
		v.MinSize(up.LastName, 2, "last_name")
	}
	//Email
	if up.Email != "" {
		v.Email(up.Email, "email")
		v.MaxSize(up.Email, 100, "email")
	}

	if len(up.Roles) > 0 {
		for _, role := range up.Roles {
			if !(role%2 == 0 && role > 0 && role <= up.MaxRole) {
				v.SetError("Roles", "has an invalid value")
				break
			}
		}
	}

	if up.Status != 0 {
		if up.Status != entityStatus.Active && up.Status != entityStatus.Disabled {
			v.SetError("Status", "has an invalid value")
		}
	}
}
