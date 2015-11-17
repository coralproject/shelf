// Package db tests the database API.
package db

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/cli/user/db"
	"github.com/coralproject/shelf/cli/user/tests"
)

// TestUserAPI validates the user CRUD API.
func TestUserAPI(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	u, err := db.NewUser("Zhang Luo", "zhang.luo@gmail.com", "Zhu4*20F_M")
	if err != nil {
		t.Fatalf("Creating User Data %s", tests.Failed)
	}
	t.Logf("Creating User Data %s", tests.Success)

	userCreate(*u, t)
	userRecordByName(u.Name, t)
	userRecordByEmail(u.Email, t)
	userRecordByPublicID(u.PublicID, t)
	userUpdate(*u, t)
	userDelete(*u, t)
	tearDown(t)
}

// userCreate tests the addition of a user record into the database.
func userCreate(u db.User, t *testing.T) {
	t.Log("Given the need to create a new User.")
	{
		t.Log("\tWhen giving a new user record")
		{
			if err := db.Create(&u); err != nil {
				t.Errorf("\t\tShould have added new user into the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have added new user into the database %s", tests.Success)
			}
		}
	}
}

// userRecordByName tests the retrieval of a user record from the database, using
// the records "Name".
func userRecordByName(name string, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's name")
		{
			_, err := db.GetUserByName(name)

			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}
		}
	}
}

// userRecordByEmail tests the retrieval of a user record from the database, using
// the records "Email".
func userRecordByEmail(email string, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's email")
		{
			_, err := db.GetUserByEmail(email)

			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}
		}
	}
}

// userRecordByPublicID tests the retrieval of a user record from the database, using
// the records "PublicID".
func userRecordByPublicID(pid string, t *testing.T) {
	t.Log("Given the need to retrieve a User record.")
	{
		t.Log("\tWhen giving a user record's public_id")
		{
			_, err := db.GetUserByPublicID(pid)

			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)
			}
		}
	}
}

// userCreate tests the updating of a user record in the database.
func userUpdate(u db.User, t *testing.T) {
	t.Log("Given the need to update a new User record.")
	{
		t.Log("\tWhen giving a user record with a name update")
		{
			if err := db.UpdateName(&u, "Zhang Shou Luo"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a user record with a email update")
		{
			if err := db.UpdateEmail(&u, "Zhang.Shou.Luo@gmail.com"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a user record with a password update")
		{
			if err := db.UpdatePassword(&u, "Zhu4*20F_M", "Zhu57*sM321"); err != nil {
				t.Errorf("\t\tShould have updated existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated existing User record in the database %s", tests.Success)
			}
		}

		t.Log("\tWhen the need to validate a password change")
		{
			user, err := db.GetUserByEmail("Zhang.Shou.Luo@gmail.com")
			if err != nil {
				t.Errorf("\t\tShould have retrieved existing user record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved existing user record %s", tests.Success)

				if user.IsPasswordValid("Zhu57*sM321") {
					t.Errorf("\t\tShould have new password for existing user %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have new password for existing user %s", tests.Success)
				}
			}
		}
	}
}

// userDelete tests the removal of a user record in the database.
func userDelete(u db.User, t *testing.T) {
	t.Log("Given the need to delete a User record.")
	{
		t.Log("\tWhen giving a user record")
		{
			if err := db.Delete(&u); err != nil {
				t.Errorf("\t\tShould have removed existing User record in the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have removed existing User record in the database %s", tests.Success)
			}
		}
		t.Log("\tWhen giving a deleted user record's email")
		{
			_, err := db.GetUserByEmail("Zhang.Shou.Luo@gmail.com")

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
	err := db.ExecuteDB(db.GetSession(), db.UserCollection, func(c *mgo.Collection) error {
		return c.DropCollection()
	})

	if err != nil {
		t.Errorf("Successfully dropped users collection %s", tests.Failed)
	} else {
		t.Logf("Successfully dropped users collection %s", tests.Success)
	}
}
