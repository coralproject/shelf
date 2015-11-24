// Package db tests the database API.
package auth_test

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/srv/auth"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/tests"
)

var collection = "users"
var u auth.User
var nu = auth.NewUser{
	UserType: 0,
	Status:   auth.StatusActive,
	FullName: "Zhang Luo",
	Email:    "zhang.luo@gmail.com",
	Password: "Zhu4*20F_M",
}

// TestUserAPI validates the user CRUD API.
func TestUserAPI(t *testing.T) {
	// Initialize the test environment.
	tests.Init()

	tests.ResetLog()
	defer tests.DisplayLog()

	userCreate(nu, t)
	userRecordByName(&u, t)
	userRecordByEmail(&u, t)
	userRecordByPublicID(&u, t)
	userUpdate(&u, t)
	userDelete(&u, t)
	tearDown(t)
}

// userCreate tests the addition of a user record into the database.
func userCreate(u auth.NewUser, t *testing.T) {
}

// userRecordByName tests the retrieval of a user record from the database, using
// the records "Name".
func userRecordByName(u *auth.User, t *testing.T) {
}

// userRecordByEmail tests the retrieval of a user record from the database, using
// the records "Email".
func userRecordByEmail(u *auth.User, t *testing.T) {
}

// userRecordByPublicID tests the retrieval of a user record from the database, using
// the records "PublicID".
func userRecordByPublicID(u *auth.User, t *testing.T) {
}

// userCreate tests the updating of a user record in the database.
func userUpdate(u *auth.User, t *testing.T) {
}

// userDelete tests the removal of a user record in the database.
func userDelete(u *auth.User, t *testing.T) {
}

// tearDown tears down the collection being used.
func tearDown(t *testing.T) {
	err := mongo.ExecuteDB("tearDown", mongo.GetSession(), collection, func(c *mgo.Collection) error {
		return c.DropCollection()
	})

	if err != nil {
		t.Errorf("Successfully dropped users collection %s", tests.Failed)
	} else {
		t.Logf("Successfully dropped users collection %s", tests.Success)
	}
}
