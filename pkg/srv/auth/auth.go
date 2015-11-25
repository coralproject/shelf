// Package auth provides CRUD methods for the auth user API.
package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/auth/crypto"
	"github.com/coralproject/shelf/pkg/srv/auth/session"
	"github.com/coralproject/shelf/pkg/srv/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collection = "users"

//==============================================================================

// CreateUser adds a new user to the database.
func CreateUser(context interface{}, ses *mgo.Session, nu NewUser) (*User, error) {
	log.Dev(context, "CreateUser", "Started : Email[%s]", nu.Email)

	if err := nu.validate(context); err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return nil, err
	}

	u, err := nu.new(context)
	if err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return nil, err
	}

	f := func(col *mgo.Collection) error {
		log.Dev(context, "CreateUser", "MGO : db.%s.insert(CAN'T SHOW)", collection)
		return col.Insert(u)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return nil, err
	}

	log.Dev(context, "CreateUser", "Completed : PublicID[%s]", u.PublicID)
	return u, nil
}

// CreateWebToken return a token and session that can be used to authenticate a user.
func CreateWebToken(context interface{}, ses *mgo.Session, u *User, expires time.Duration) (string, error) {
	log.Dev(context, "CreateWebToken", "Started : PublicID[%s]", u.PublicID)

	// Do we have a valid session right now?
	s, err := session.GetByLatest(context, ses, u.PublicID)
	if err != nil && err != mgo.ErrNotFound {
		log.Error(context, "CreateUser", err, "Completed")
		return "", err
	}

	// If we don't have one or it has been expired create
	// a new one.
	if err == mgo.ErrNotFound || s.IsExpired(context) {
		if s, err = session.Create(context, ses, u.PublicID, expires); err != nil {
			log.Error(context, "CreateUser", err, "Completed")
			return "", err
		}
	}

	// Set the return arguments though we will explicitly
	// return them. Don't want any confusion.
	token, err := u.WebToken(s.SessionID)
	if err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return "", err
	}

	log.Dev(context, "CreateWebToken", "Completed : WebToken[%s]", token)
	return token, nil
}

//==============================================================================

// DecodeWebToken breaks a web token into its parts.
func DecodeWebToken(context interface{}, webToken string) (sessionID string, token string, err error) {
	log.Dev(context, "DecodeWebToken", "Started : WebToken[%s]", webToken)

	// Decode the web token to break it into its parts.
	data, err := base64.StdEncoding.DecodeString(webToken)
	if err != nil {
		log.Error(context, "DecodeWebToken", err, "Completed")
		return "", "", err
	}

	// Split the web token.
	str := strings.Split(string(data), ":")
	if len(str) != 2 {
		err := errors.New("Invalid token")
		log.Error(context, "DecodeWebToken", err, "Completed")
		return "", "", err
	}

	// Pull out the session and token.
	sessionID = str[0]
	token = str[1]

	log.Dev(context, "DecodeWebToken", "Completed : SessionID[%s] Token[%s]", sessionID, token)
	return sessionID, token, nil
}

// ValidateWebToken accepts a web token and validates its credibility. Returns
// a User value is the token is valid.
func ValidateWebToken(context interface{}, ses *mgo.Session, webToken string) (*User, error) {
	log.Dev(context, "ValidateWebToken", "Started : WebToken[%s]", webToken)

	// Extract the sessionID and token from the web token.
	sessionID, token, err := DecodeWebToken(context, webToken)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Find the session in the database.
	s, err := session.GetBySessionID(context, ses, sessionID)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Validate the session has not expired.
	if s.IsExpired(context) {
		err := errors.New("Expired token")
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Pull the user for this session.
	u, err := GetUserByPublicID(context, ses, s.PublicID)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Validate the token against this user.
	if err := crypto.IsTokenValid(u, token); err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	log.Dev(context, "ValidateWebToken", "Completed : PublicID[%s]", u.PublicID)
	return u, nil
}

//==============================================================================

// GetUserByPublicID retrieves a user record by using the provided PublicID.
func GetUserByPublicID(context interface{}, ses *mgo.Session, publicID string) (*User, error) {
	log.Dev(context, "GetUserByPublicID", "Started : PID[%s]", publicID)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		log.Dev(context, "GetUserByPublicID", "MGO : db.%s.findOne(%s)", collection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "GetUserByPublicID", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByPublicID", "Completed")
	return &user, nil
}

//==============================================================================

// Update updates an existing user to the database.
func Update(context interface{}, ses *mgo.Session, uu UpdUser) error {
	log.Dev(context, "Update", "Started : PublicID[%s]", uu.PublicID)

	if err := uu.validate(context); err != nil {
		log.Error(context, "Update", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": uu.PublicID}
		upd := bson.M{"$set": bson.M{"full_name": uu.FullName, "email": uu.Email, "type": uu.UserType, "status": uu.Status, "modified_at": time.Now().UTC()}}
		log.Dev(context, "Update", "MGO : db.%s.update(%s)", collection, mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "Update", err, "Completed")
		return err
	}

	log.Dev(context, "Update", "Completed")
	return nil
}

// UpdateCredentials updates an existing user's password and token in the database.
func UpdateCredentials(context interface{}, ses *mgo.Session, u *User, password string) error {
	log.Dev(context, "UpdateCredentials", "Started : PublicID[%s]", u.PublicID)

	if len(password) < 8 {
		err := errors.New("Invalid password length")
		log.Error(context, "UpdateCredentials", err, "Completed")
		return err
	}

	newPassHash, err := crypto.BcryptPassword(u.PrivateID + password)
	if err != nil {
		log.Error(context, "UpdateCredentials", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": u.PublicID}
		upd := bson.M{"$set": bson.M{"password": newPassHash, "modified_at": time.Now().UTC()}}
		log.Dev(context, "Create", "MGO : db.%s.Update(%s, CAN'T SHOW)", collection, mongo.Query(q))
		return c.Update(q, upd)
	}

	if err = mongo.ExecuteDB(context, ses, collection, f); err != nil {
		log.Error(context, "UpdateCredentials", err, "Completed")
		return nil
	}

	log.Dev(context, "UpdateCredentials", "Completed")
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
