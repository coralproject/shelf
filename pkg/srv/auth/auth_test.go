package auth_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/srv/auth"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collection = "users"
	context    = "testing"
)

func init() {
	tests.Init()
}

// removeUser is used to clear out all the test user from the collection.
func removeUser(ses *mgo.Session, publicID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != mgo.ErrNotFound {
		return err
	}

	return nil
}

// TestCreateUser tests the creation of a user.
func TestCreateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err != nil {
			t.Errorf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to create users in the DB.")
	{
		nu := auth.NewUser{
			UserType: auth.TypeAPI,
			Status:   auth.StatusActive,
			FullName: "Test Kennedy",
			Email:    "bill@ardanlabs.com",
			Password: "_Password124",
		}

		t.Log("\tWhen using a test user.")
		{
			u1, err := auth.CreateUser(context, ses, nu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			u2, err := auth.GetUserByPublicID(context, ses, u1.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			// Remove the objectid to be able to compare the values.
			u2.ID = ""

			// Need to remove the nanoseconds to be able to compare the values.
			u1.DateModified = u1.DateModified.Add(-time.Duration(u1.DateModified.Nanosecond()))
			u1.DateCreated = u1.DateCreated.Add(-time.Duration(u1.DateCreated.Nanosecond()))
			u2.DateModified = u2.DateModified.Add(-time.Duration(u2.DateModified.Nanosecond()))
			u2.DateCreated = u2.DateCreated.Add(-time.Duration(u2.DateCreated.Nanosecond()))

			if !reflect.DeepEqual(*u1, *u2) {
				t.Errorf("\t%s\tShould be able to get back the same user.", tests.Failed)
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u2)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}
