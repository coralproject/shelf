// Package db tests the database API.
package db

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/tests"
)

// TestUserAPI validates the user CRUD API.
func TestUserAPI(t *testing.T) {
	// Initialize the test environment.
	tests.Init()

	tests.ResetLog()
	defer tests.DisplayLog()

	u, err := NewUser("Zhang Luo", "zhang.luo@gmail.com", "Zhu4*20F_M")
	if err != nil {
		t.Fatalf("Creating User Data %s", tests.Failed)
	}
	t.Logf("Creating User Data %s", tests.Success)

	userCreate(u, t)
	userRecordByName(u.Name, u, t)
	userRecordByEmail(u.Email, u, t)
	userRecordByPublicID(u.PublicID, u, t)
	userUpdate(u, t)
	userDelete(u, t)
	tearDown(t)
}

// userCreate tests the addition of a user record into the database.
func userCreate(u *User, t *testing.T) {
	t.Log("Given the need to create a new User.")
	{
		t.Log("\tWhen giving a new user record")
		{
			if err := Create(u); err != nil {
				t.Errorf("\t\tShould have added new user into the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have added new user into the database %s", tests.Success)
			}
		}
		t.Log("\tWhen user record has been saved and we need to validate")
		{
			user, err := GetUserByEmail(u.Email)
			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}

			if err := CompareAll(u, user); err != nil {
				t.Errorf("\t\tShould have similar field values %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have similar field values %s", tests.Success)
			}

		}
	}
}

// userRecordByName tests the retrieval of a user record from the database, using
// the records "Name".
func userRecordByName(name string, u *User, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's name")
		{
			user, err := GetUserByName(name)
			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}

			if err := CompareAll(u, user); err != nil {
				t.Errorf("\t\tShould have similar field values %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have similar field values %s", tests.Success)
			}

		}
	}
}

// userRecordByEmail tests the retrieval of a user record from the database, using
// the records "Email".
func userRecordByEmail(email string, u *User, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's email")
		{
			user, err := GetUserByEmail(email)
			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}

			if err := CompareAll(u, user); err != nil {
				t.Errorf("\t\tShould have similar field values %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have similar field values %s", tests.Success)
			}
		}
	}
}

// userRecordByPublicID tests the retrieval of a user record from the database, using
// the records "PublicID".
func userRecordByPublicID(pid string, u *User, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's public_id")
		{
			user, err := GetUserByPublicID(pid)

			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}

			if err := CompareAll(u, user); err != nil {
				log.Printf("Compared: %s", err)
				t.Errorf("\t\tShould have similar field values %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have similar field values %s", tests.Success)
			}
		}
	}
}

// userCreate tests the updating of a user record in the database.
func userUpdate(u *User, t *testing.T) {
	t.Log("Given the need to update a new User record.")
	{
		t.Log("\tWhen giving a user record with a name update")
		{
			if err := UpdateName(u, "Zhang Shou Luo"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a user record with a email update")
		{
			if err := UpdateEmail(u, "Zhang.Shou.Luo@gmail.com"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a user record with a password update")
		{
			if err := UpdatePassword(u, "Zhu4*20F_M", "Zhu57*sM321"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}

		t.Log("\tWhen the need to validate a password change")
		{
			if !u.IsPasswordValid("Zhu57*sM321") {
				t.Errorf("\t\tShould have new password for existing user %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have new password for existing user %s", tests.Success)
			}

			user, err := GetUserByEmail("Zhang.Shou.Luo@gmail.com")
			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)

				if !user.IsPasswordValid("Zhu57*sM321") {
					t.Errorf("\t\tShould have new password for retrieved user %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have new password for retrieved user %s", tests.Success)
				}
			}

		}
	}
}

// userDelete tests the removal of a user record in the database.
func userDelete(u *User, t *testing.T) {
	t.Log("Given the need to delete a User record.")
	{
		t.Log("\tWhen giving a user record")
		{
			if err := Delete(u); err != nil {
				t.Errorf("\t\tShould have removed existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have removed existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a deleted user record's email")
		{
			_, err := GetUserByEmail("Zhang.Shou.Luo@gmail.com")

			if err == nil {
				t.Errorf("\t\tShould not receive a user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould not receive a user record %s", tests.Success)
			}
		}
	}
}

// tearDown tears down the collection being used.
func tearDown(t *testing.T) {
	err := mongo.ExecuteDB("tearDown", mongo.GetSession(), UserCollection, func(c *mgo.Collection) error {
		return c.DropCollection()
	})

	if err != nil {
		t.Errorf("Successfully dropped users collection %s", tests.Failed)
	} else {
		t.Logf("Successfully dropped users collection %s", tests.Success)
	}
}
