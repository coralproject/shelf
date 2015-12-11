package data

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

// collection contains the name of the comments collection.
const (
	commentCollection = "comments"
	userCollection    = "users"
)

//==============================================================================
//==================================== Comment =================================
//==============================================================================

// CreateComment adds a new comment in the database.
func CreateComment(context interface{}, db *db.DB, com *Comment) error {
	if com.CommentID == "" {
		com.CommentID = uuid.New()
	}
	com.DateCreated = time.Now()
	com.Status = "New"

	f := func(col *mgo.Collection) error {
		log.Dev(context, "CreateComment", "MGO: db.%s.insert()", commentCollection)
		return col.Insert(com)
	}

	if err := db.ExecuteMGO(context, commentCollection, f); err != nil {
		log.Error(context, "CreateComment", err, "Completed")
		return err
	}

	log.Dev(context, "CreateComment", "Completed")

	return nil
}

// GetCommentByID retrieves an individual comment by ID
func GetCommentByID(context interface{}, db *db.DB, id string) (*Comment, error) {
	log.Dev(context, "GetCommentById", "Started : Id[%s]", id)

	var comment Comment
	f := func(c *mgo.Collection) error {
		q := bson.M{"comment_id": id}
		log.Dev(context, "GetCommentById", "MGO : db.%s.find(CommentId: '%s')", commentCollection, mongo.Query(q))
		return c.Find(q).One(&comment)
	}

	if err := db.ExecuteMGO(context, commentCollection, f); err != nil {
		log.Error(context, "GetUserById", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetCommentById", "Completed")
	return &comment, nil
}

//==============================================================================
//==================================== User ====================================
//==============================================================================

// GetUserByID retrieves an individual user by ID
func GetUserByID(context interface{}, db *db.DB, id string) (*User, error) {
	log.Dev(context, "GetUserById", "Started : Id[%s]", id)

	var user User
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": id}
		log.Dev(context, "GetUserById", "MGO : db.%s.findOne(%s)", userCollection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, userCollection, f); err != nil {
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
		log.Dev(context, "GetUserByUserName", "MGO : db.%s.findOne(%s)", userCollection, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, userCollection, f); err != nil {
		log.Error(context, "GetUserByUserName", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByUserName", "Completed")
	return &user, nil
}

// CreateUser creates a new user resource
func CreateUser(context interface{}, db *db.DB, user *User) error {
	log.Dev(context, "CreateUser", "Started : User: ", user)

	/* This error condition should be performed by DB indexes
	dbUser, err := GetUserByUserName(context, db, user.UserName)
	if dbUser != nil {
		log.Error(context, "CreateUser", err, "User exists")
		return "CreateUser: user with same UserName already exists"
	}
	*/

	// set defaults, may want to move this to a factory method on the User struct
	if user.UserID == "" {
		user.UserID = uuid.New()
	}
	user.MemberSince = time.Now()

	f := func(col *mgo.Collection) error {
		log.Dev(context, "CreateUser", "MGO: db.%s.insert()", userCollection)
		return col.Insert(user)
	}

	// Write the user to mongo
	if err := db.ExecuteMGO(context, userCollection, f); err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return err
	}

	log.Dev(context, "CreateUser", "Completed")
	return nil
}

/*

Todo: bulk methods should call individual inserts in a loop

// AddUsers adds an array of users to user collectionollection
func AddUsers(context interface{}, db *db.DB, users []User) error {
	return mongo.GetCollection(db.MGOConn, collection).Insert(users)
}

// AddComments adds an array of comments to comment collection
func AddComments(context interface{}, db *db.DB, comments []Comment) error {
	return mongo.GetCollection(db.MGOConn, collection).Insert(comments)
}
*/
