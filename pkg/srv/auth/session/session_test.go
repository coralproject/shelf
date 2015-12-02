package session_test

import (
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/auth/session"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	publicID = "6dcda2da-92c3-11e5-8994-feff819cdc9f"
	context  = "testing"
)

func init() {
	tests.Init()
}

// removeSessions is used to clear out all the test sessions that are
// created from tests.
func removeSessions(db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, "sessions", f); err != nil {
		return err
	}

	return nil
}

// TestIsExpired tests that we can properly identify when a
// session is expired or not.
func TestIsExpired(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to validate a session has expired.")
	{
		s := session.Session{
			DateExpires: time.Now().Add(-time.Hour),
		}

		t.Log("\tWhen using an expired session.")
		{
			if !s.IsExpired(context) {
				t.Fatalf("\t%s\tShould be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould be expired", tests.Success)
		}

		s = session.Session{
			DateExpires: time.Now().Add(time.Hour),
		}

		t.Log("\tWhen using an valid session")
		{
			if s.IsExpired(context) {
				t.Fatalf("\t%s\tShould Not be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be expired", tests.Success)
		}
	}
}

// TestCreate tests the creation of sessions.
func TestCreate(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSessions(db); err != nil {
			t.Errorf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)
	}()

	t.Log("Given the need to create sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(db); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)

			s1, err := session.Create(context, db, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			s2, err := session.GetBySessionID(context, db, s1.SessionID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Success)

			if s1.SessionID != s2.SessionID {
				t.Fatalf("\t%s\tShould be able to get back the same session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same session.", tests.Success)
			}

			if s1.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}

// TestGetLatest tests the retrieval of the latest session.
func TestGetLatest(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSessions(db); err != nil {
			t.Errorf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)
	}()

	t.Log("Given the need to get the latest sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(db); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)

			if _, err := session.Create(context, db, publicID, 10*time.Second); err != nil {
				t.Fatalf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			time.Sleep(time.Second)

			s2, err := session.Create(context, db, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create another session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create another session.", tests.Success)

			s3, err := session.GetByLatest(context, db, publicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the latest session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the latest session.", tests.Success)

			if s2.SessionID != s3.SessionID {
				t.Errorf("\t%s\tShould be able to get back the latest session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the latest session.", tests.Success)
			}
		}
	}
}

// TestGetNotFound tests when a session is not found.
func TestGetNotFound(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db := db.NewMGO()
	defer db.CloseMGO()

	t.Log("Given the need to test finding a session and it is not found.")
	{
		t.Logf("\tWhen using SessionID %s", "NOT EXISTS")
		{
			if _, err := session.GetBySessionID(context, db, "NOT EXISTS"); err == nil {
				t.Fatalf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
		}
	}
}

// TestNoSession tests when a nil session is used.
func TestNoSession(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test calls with a bad session.")
	{
		t.Log("\tWhen using a nil session")
		{
			if _, err := session.Create(context, nil, publicID, 10*time.Second); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a session.", tests.Success)
			}

			if _, err := session.GetBySessionID(context, nil, "NOT EXISTS"); err == nil {
				t.Errorf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
			}

			if _, err := session.GetByLatest(context, nil, publicID); err == nil {
				t.Errorf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
			}
		}
	}
}
