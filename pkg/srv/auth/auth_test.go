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

//==============================================================================

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

//==============================================================================

// TestCreateUser tests the creation of a user.
func TestCreateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
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
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u2)
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}

// TestCreateUserValidation tests the creation of a user that is not valid.
func TestCreateUserValidation(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err == nil {
			t.Fatalf("\t%s\tShould Not be able to remove the test user", tests.Failed)
		}
		t.Logf("\t%s\tShould Not be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to make sure only valid users are created in the DB.")
	{
		t.Log("\tWhen using a test user.")
		{
			nu := auth.NewUser{
				UserType: 0,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

			if _, err := auth.CreateUser(context, ses, nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid UserType", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid UserType.", tests.Success)
			}

			nu = auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   0,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

			if _, err := auth.CreateUser(context, ses, nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Status", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Status.", tests.Success)
			}

			nu = auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "1234567",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			}

			if _, err := auth.CreateUser(context, ses, nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid FullName", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid FullName.", tests.Success)
			}

			nu = auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill",
				Password: "_Password124",
			}

			if _, err := auth.CreateUser(context, ses, nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Email", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Email.", tests.Success)
			}

			nu = auth.NewUser{
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "1234567",
			}

			if _, err := auth.CreateUser(context, ses, nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Password", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Password.", tests.Success)
			}
		}
	}
}

// TestUpdateUser tests we can update user information.
func TestUpdateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
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

			uu := auth.UpdUser{
				PublicID: publicID,
				UserType: auth.TypeUSER,
				Status:   auth.StatusInvalid,
				FullName: "Update Kennedy",
				Email:    "upt@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

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

			// Update the fields that changed
			u1.UserType = u2.UserType
			u1.Status = u2.Status
			u1.FullName = u2.FullName
			u1.Email = u2.Email

			if !reflect.DeepEqual(*u1, *u2) {
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u2)
				t.Errorf("\t%s\tShould be able to get back the same user with changes.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user with changes.", tests.Success)
			}
		}
	}
}

// TestUpdateUserValidation tests the update of a user that is not valid.
func TestUpdateUserValidation(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err == nil {
			t.Fatalf("\t%s\tShould Not be able to remove the test user", tests.Failed)
		}
		t.Logf("\t%s\tShould Not be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to make sure only valid users are created in the DB.")
	{
		t.Log("\tWhen using a test user.")
		{
			uu := auth.UpdUser{
				PublicID: "asdasdasd",
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid PublicID", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid PublicID.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				UserType: 0,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid UserType", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid UserType.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				UserType: auth.TypeAPI,
				Status:   0,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid Status", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid Status.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "1234567",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid FullName", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid FullName.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid Email", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid Email.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				UserType: auth.TypeAPI,
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}
		}
	}
}

// TestUpdateUserPassword tests we can update user password.
func TestUpdateUserPassword(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
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

			webTok, err := auth.CreateWebToken(context, ses, u1, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if err := auth.UpdateUserPassword(context, ses, u1, "_Password567"); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

			if _, err := auth.ValidateWebToken(context, ses, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the org web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould Not be able to validate the new org token.", tests.Success)

			u2, err := auth.GetUserByPublicID(context, ses, u1.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			webTok2, err := auth.CreateWebToken(context, ses, u2, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a new web token.", tests.Success)

			if webTok == webTok2 {
				t.Fatalf("\t%s\tShould have different web tokens after the update.", tests.Failed)
			}
			t.Logf("\t%s\tShould have different web tokens after the update.", tests.Success)

			u3, err := auth.ValidateWebToken(context, ses, webTok2)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to validate the new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the new web token.", tests.Success)

			if u1.PublicID != u3.PublicID {
				t.Log(u2.PublicID)
				t.Log(u3.PublicID)
				t.Fatalf("\t%s\tShould have the right user for the new token.", tests.Failed)
			}
			t.Logf("\t%s\tShould have the right user for the new token.", tests.Success)
		}
	}
}

// TestDeleteUser test the deleting of a user.
func TestDeleteUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
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

			u2, err := auth.GetUserByPublicID(context, ses, u1.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			if err := auth.DeleteUser(context, ses, u2.PublicID); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the user.", tests.Success)

			if _, err := auth.GetUserByPublicID(context, ses, u1.PublicID); err == nil {
				t.Fatalf("\t%s\tShould Not be able to retrieve the user by PublicID.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the user by PublicID.", tests.Success)
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
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
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

			webTok, err := auth.CreateWebToken(context, ses, u1, time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

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

			u2, err := auth.GetUserByPublicID(context, ses, u1.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			if u2.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould have the right session for user.", tests.Failed)
				t.Log(u2.PublicID)
				t.Log(s2.PublicID)
			}
			t.Logf("\t%s\tShould have the right session for user.", tests.Success)

			webTok2, err := u2.WebToken(sId)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web new token.", tests.Success)

			if webTok != webTok2 {
				t.Log(webTok)
				t.Log(webTok2)
				t.Fatalf("\t%s\tShould be able to create the same web token.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to create the same web token.", tests.Success)

			u3, err := auth.ValidateWebToken(context, ses, webTok2)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to validate the new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the new web token.", tests.Success)

			if u1.PublicID != u3.PublicID {
				t.Log(u1.PublicID)
				t.Log(u3.PublicID)
				t.Fatalf("\t%s\tShould have the right user for the token.", tests.Failed)
			}
			t.Logf("\t%s\tShould have the right user for the token.", tests.Success)
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
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
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
			}
			t.Logf("\t%s\tShould Not be able to validate the web token.", tests.Success)
		}
	}
}

// TestInvalidWebTokens tests create an invalid web token and tests it fails.
func TestInvalidWebTokens(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	tokens := []string{
		"",
		"6dcda2da-92c3-11e5-8994-feff819cdc9f",
		"OGY4OGI3YWQtZjc5Ny00ODI1LWI0MmUtMjIwZTY5ZDQxYjMzOmFKT2U1b0pFZlZ4cWUrR0JONEl0WlhmQTY0K3JsN2VGcmM2MVNQMkV1WVE9",
	}

	t.Log("Given the need to validate bad web tokens don't validate.")
	{
		for _, token := range tokens {
			t.Logf("\tWhen using token [%s]", token)
			{
				if _, err := auth.ValidateWebToken(context, ses, token); err == nil {
					t.Errorf("\t%s\tShould Not be able to validate the web token : %v", tests.Failed, err)
				} else {
					t.Logf("\t%s\tShould Not be able to validate the web token.", tests.Success)
				}
			}
		}
	}
}

// TestInvalidWebTokenUpdateEmail tests a token becomes invalid after an update.
func TestInvalidWebTokenUpdateEmail(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	ses := mongo.GetSession()
	defer ses.Close()

	var publicID string
	defer func() {
		if err := removeUser(ses, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to validate web tokens don't work after user update.")
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

			webTok, err := auth.CreateWebToken(context, ses, u1, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if _, err := auth.ValidateWebToken(context, ses, webTok); err != nil {
				t.Fatalf("\t%s\tShould be able to validate the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the web token.", tests.Success)

			uu := auth.UpdUser{
				PublicID: publicID,
				UserType: auth.TypeUSER,
				Status:   auth.StatusInvalid,
				FullName: "Update Kennedy",
				Email:    "change@ardanlabs.com",
			}

			if err := auth.UpdateUser(context, ses, uu); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

			if _, err := auth.ValidateWebToken(context, ses, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the org web token.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to validate the org web token.", tests.Success)
		}
	}
}
