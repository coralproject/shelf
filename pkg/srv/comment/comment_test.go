package comment_test

import (
	"testing"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/comment"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	tests.InitMGO()
}

//==============================================================================

func TestCreateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db := db.NewMGO()
	defer db.CloseMGO()

	var ID string
	defer func() {
		if err := removeUser(db, ID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	u1 := comment.User{
		UserName: "David",
		Avatar:   "https://picture.of/david.jpg",
	}

	if err := comment.CreateUser(tests.Context, db, &u1); err != nil {
		t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to create a user", tests.Success)

	// set ID for the deferred removeUser method
	ID = u1.UserID

	if err := comment.CreateUser(tests.Context, db, &u1); err == nil {
		t.Fatalf("\t%s\tShould not be able to create a duplicate user", tests.Failed)
	}
	t.Logf("\t%s\tShould not be able to create a duplicate user", tests.Success)

}

func TestCreateComment(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db := db.NewMGO()
	defer db.CloseMGO()

	var CommentID string
	defer func() {
		if err := removeComment(db, CommentID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test comment : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test comment.", tests.Success)
	}()

	c1 := comment.Comment{
		UserId: "4",
		Body:   "Wonderful story!  The world is going in the right direction!",
	}

	if err := comment.CreateComment(tests.Context, db, &c1); err != nil {
		t.Fatalf("\t%s\tShould be able to create a comment : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to create a comment", tests.Success)

	CommentID = c1.CommentID

	t.Logf("\t%s\tYeah, ok, this works.", tests.Success)
}

//==============================================================================

// removeUser is used to clear out all the test user from the collection.
func removeUser(db *db.DB, UserID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"user_id": UserID}
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(tests.Context, "users", f); err != nil {
		return err
	}

	return nil
}

// removeComment is used to clear out all the test user from the collection.
func removeComment(db *db.DB, CommentID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"comment_id": CommentID}
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(tests.Context, "comments", f); err != nil {
		return err
	}

	return nil
}
