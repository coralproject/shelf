package auth_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/srv/auth"
	"github.com/coralproject/shelf/pkg/srv/auth/session"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var context = "testing"

func init() {
	tests.Init()
}

// removeUser is used to clear out all the test user from the collection.
func removeUser(ses *mgo.Session, publicID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		return c.Remove(q)
	}

	if err := mongo.ExecuteDB(context, ses, "users", f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := mongo.ExecuteDB(context, ses, "sessions", f); err != nil {
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
		t.Log("\tWhen using a test user.")
		{
			nu := auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

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

// TestCreateWebToken tests create a web token and a pairing session.
func TestCreateWebToken(t *testing.T) {
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

	t.Log("Given the need to create a web token.")
	{
		t.Log("\tWhen using a new user.")
		{
			nu := auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

			u1, err := auth.CreateUser(context, ses, nu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(context, ses, u1, 250*time.Millisecond)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			u2, err := auth.GetUserByPublicID(context, ses, u1.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			sId, _, err := auth.DecodeWebToken(context, webTok)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to decode the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to decode the web token.", tests.Success)

			s2, err := session.GetBySessionID(context, ses, sId)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Success)

			if u2.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould have the right session for user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the right session for user.", tests.Success)
			}

			u3, err := auth.ValidateWebToken(context, ses, webTok)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to validate the web token : %v", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to validate the web token.", tests.Success)
			}

			if u1.PublicID != u3.PublicID {
				t.Fatalf("\t%s\tShould have the right user for the token.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have the right user for the token.", tests.Success)
			}
		}
	}
}

// TestExpiredWebToken tests create a web token and tests when it expires.
func TestExpiredWebToken(t *testing.T) {
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

	t.Log("Given the need to validate web tokens expire.")
	{
		t.Log("\tWhen using a new user.")
		{
			nu := auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

			u1, err := auth.CreateUser(context, ses, nu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(context, ses, u1, 1*time.Millisecond)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if _, err := auth.ValidateWebToken(context, ses, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the web token : %v", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould Not be able to validate the web token.", tests.Success)
			}
		}
	}
}
