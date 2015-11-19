// Package servicestest implements users tests for the services layer.
// USE THIS AS A MODEL FOR NOW.
package services

import (
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/srv/xenia/app"
	"github.com/coralproject/shelf/srv/xenia/models"
	"github.com/coralproject/shelf/srv/xenia/services"
	"github.com/coralproject/shelf/srv/xenia/tests"
)

// TestUsers validates a user can be created, retrieved and
// then removed from the system.
func TestUsers(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	c := &app.Context{
		Session:   mongo.GetSession(),
		SessionID: "TESTING",
	}
	defer c.Session.Close()

	now := time.Now()

	u := models.User{
		UserType:     1,
		FirstName:    "Bill",
		LastName:     "Kennedy",
		Email:        "bill@ardanstugios.com",
		Company:      "Ardan Labs",
		DateModified: &now,
		DateCreated:  &now,
		Addresses: []models.UserAddress{
			{
				Type:         1,
				LineOne:      "12973 SW 112th ST",
				LineTwo:      "Suite 153",
				City:         "Miami",
				State:        "FL",
				Zipcode:      "33172",
				Phone:        "305-527-3353",
				DateModified: &now,
				DateCreated:  &now,
			},
		},
	}

	t.Log("Given the need to add a new user, retrieve and remove that user from the system.")
	{
		if _, err := u.Validate(); err != nil {
			t.Fatal("\tShould be able to validate the user data.", tests.Failed)
		}
		t.Log("\tShould be able to validate the user data.", tests.Succeed)

		if _, err := services.Users.Create(c, &u); err != nil {
			t.Fatal("\tShould be able to create a user in the system.", tests.Failed)
		}
		t.Log("\tShould be able to create a user in the system.", tests.Succeed)

		if u.UserID == "" {
			t.Fatal("\tShould have an UserID for the user.", tests.Failed)
		}
		t.Log("\tShould have an UserID for the user.", tests.Succeed)

		ur, err := services.Users.Retrieve(c, u.UserID)
		if err != nil {
			t.Fatal("\tShould be able to retrieve the user back from the system.", tests.Failed)
		}
		t.Log("\tShould be able to retrieve the user back from the system.", tests.Succeed)

		if _, err := u.Compare(ur); err != nil {
			t.Fatal("\tShould find both the original and retrieved value are identical.", tests.Failed)
		}
		t.Log("\tShould find both the original and retrieved value are identical.", tests.Succeed)

		if ur == nil || u.UserID != ur.UserID {
			t.Fatal("\tShould have a match between the created user and the one retrieved.", tests.Failed)
		}
		t.Log("\tShould have a match between the created user and the one retrieved.", tests.Succeed)

		if err := services.Users.Delete(c, u.UserID); err != nil {
			t.Fatal("\tShould be able to remove the user from the system.", tests.Failed)
		}
		t.Log("\tShould be able to remove the user from the system", tests.Succeed)

		if _, err := services.Users.Retrieve(c, u.UserID); err == nil {
			t.Fatal("\tShould NOT be able to retrieve the user back from the system.", tests.Failed)
		}
		t.Log("\tShould NOT be able to retrieve the user back from the system.", tests.Succeed)
	}
}
