package comment

import (
	"time"
	"strings"

	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/vendor/github.com/pborman/uuid"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collectionUser = "user"
const collectionComment = "comment"

// GetUser retrieves an individual user resource
func GetCommentById(context interface{}, session *mgo.Session, id string) (*User, error) {
	log.Dev(context, "GetCommentById", "Started : Id[%s]", id)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetCommentById", "MGO : db.%s.findOne(%s)", collectionComment, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, session, collectionComment, f); err != nil {
		log.Error(context, "GetUserById", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetCommentById", "Completed")
	return &user, nil
}

// CreateComment creates a new comment
func CreateComment(context interface{}, session *mgo.Session, comment Comment) (*Comment, error) {

	if comment.Id == "" {
		comment.Id = uuid.New()
	}
	comment.CreatedDate = time.Now()
	comment.Status = "New"

	// Write the user to mongo
	err1 := mongo.GetCollection(session, collectionComment).Insert(comment)
	if err1 != nil {
		return nil, err1
	}

	return &comment, nil
}

// GetUser retrieves an individual user resource
func GetUserById(context interface{}, session *mgo.Session, id string) (*User, error) {
	log.Dev(context, "GetUserById", "Started : Id[%s]", id)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetUserById", "MGO : db.%s.findOne(%s)", collectionUser, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, session, collectionUser, f); err != nil {
		log.Error(context, "GetUserById", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserById", "Completed")
	return &user, nil
}

// GetUser retrieves an individual user resource
func GetUserByEmail(context interface{}, session *mgo.Session, email string) (*User, error) {
	log.Dev(context, "GetUserByEmail", "Started : Email[%s]", email)

	email = strings.ToLower(email)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"email": email}
		log.Dev(context, "GetUserByEmail", "MGO : db.%s.findOne(%s)", collectionUser, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := mongo.ExecuteDB(context, session, collectionUser, f); err != nil {
		log.Error(context, "GetUserByEmail", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByEmail", "Completed")
	return &user, nil
}

// CreateUser creates a new user resource
func CreateUser(context interface{}, session *mgo.Session, user User) (*User, error) {
	log.Dev(context, "CreateUser", "Started : User: ", user)

	dbUser, err := GetUserByEmail(context, session, user.Email)
	if dbUser != nil {
		log.Error(context, "CreateUser", err, "User exists")
		return dbUser, nil
	}

	if user.Id == "" {
		user.Id = uuid.New()
	}
	user.MemberSince = time.Now()

	// Write the user to mongo
	err1 := mongo.GetCollection(session, collectionUser).Insert(user)
	if err1 != nil {
		return nil, err1
	}

	log.Dev(context, "CreateUser", "Completed")
	return &user, nil
}
