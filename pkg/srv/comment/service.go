package comment

import (
	"strings"
	"time"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/log"
	"github.com/pborman/uuid"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collections contains the name of the user collection.
const collectionUser = "user"
const collectionComment = "comment"

// GetCommentByID retrieves an individual comment by ID
func GetCommentByID(context interface{}, db *db.DB, id string) (*User, error) {
	log.Dev(context, "GetCommentById", "Started : Id[%s]", id)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetCommentById", "MGO : db.%s.findOne(%s)", collectionComment, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, collectionComment, f); err != nil {
		log.Error(context, "GetUserById", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetCommentById", "Completed")
	return &user, nil
}

// CreateComment creates a new comment
func CreateComment(context interface{}, db *db.DB, comment Comment) (*Comment, error) {
	if comment.ID == "" {
		comment.ID = uuid.New()
	}
	comment.DateCreated = time.Now()
	comment.Status = "New"

	// Write the user to mongo
	err1 := mongo.GetCollection(db.MGOConn, collectionComment).Insert(comment)
	if err1 != nil {
		return nil, err1
	}

	return &comment, nil
}

// GetUserByID retrieves an individual user by ID
func GetUserByID(context interface{}, db *db.DB, id string) (*User, error) {
	log.Dev(context, "GetUserById", "Started : Id[%s]", id)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetUserById", "MGO : db.%s.findOne(%s)", collectionUser, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, collectionUser, f); err != nil {
		log.Error(context, "GetUserById", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserById", "Completed")
	return &user, nil
}

// GetUserByUserName retrieves an individual user by email
func GetUserByUserName(context interface{}, db *db.DB, userName string) (*User, error) {
	log.Dev(context, "GetUserByUserName", "Started : User[%s]", userName)

	userName = strings.ToLower(userName)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"user_name": userName}
		log.Dev(context, "GetUserByUserName", "MGO : db.%s.findOne(%s)", collectionUser, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, collectionUser, f); err != nil {
		log.Error(context, "GetUserByUserName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByUserName", "Completed")
	return &user, nil
}

// CreateUser creates a new user resource
func CreateUser(context interface{}, db *db.DB, user User) (*User, error) {
	log.Dev(context, "CreateUser", "Started : User: ", user)

	dbUser, err := GetUserByUserName(context, db, user.UserName)
	if dbUser != nil {
		log.Error(context, "CreateUser", err, "User exists")
		return dbUser, nil
	}

	if user.ID == "" {
		user.ID = uuid.New()
	}
	user.MemberSince = time.Now()

	// Write the user to mongo
	err1 := mongo.GetCollection(db.MGOConn, collectionUser).Insert(user)
	if err1 != nil {
		return nil, err1
	}

	log.Dev(context, "CreateUser", "Completed")
	return &user, nil
}

// AddUsers adds an array of users to user collection
func AddUsers(context interface{}, db *db.DB, users []User) error {
	return mongo.GetCollection(db.MGOConn, collectionUser).Insert(users)
}

// AddComments adds an array of comments to comment collection
func AddComments(context interface{}, db *db.DB, comments []Comment) error {
	return mongo.GetCollection(db.MGOConn, collectionComment).Insert(comments)
}
