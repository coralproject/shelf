package session_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/srv/auth/session"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/tests"
)

var (
	collection = "sessions"
	publicID   = "6dcda2da-92c3-11e5-8994-feff819cdc9f"
)

func init() {
	tests.Init()
}

// removeSessions is used to clear out all the test sessions that are
// created from tests.
func removeSessions() error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		return c.Remove(q)
	}

	return mongo.ExecuteDB(context, ses, collection, f)
}

// retrieveSession is used to validate sessions are being saved
// correctly.
func retrieveSession(sessionID string) (*session.Session, error) {
	var s session.Session
	f := func(c *mgo.Collection) error {
		q := bson.M{"session_id": sessionID}
		return c.Find(q).One(&s)
	}

	if err := mongo.ExecuteDB(context, ses, collection, f); err != nil {
		return nil, err
	}

	return &s, nil
}

// TestCreate tests the creation of sessions.
func TestCreate(t *testing.T) {
	t.Log("Given the need to create sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Succeed)

			ses := mongo.GetSession()
			defer ses.Close()

			s1, err := session.Create("context", ses, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a session.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Succeed)

			s2, err := retrieveSession(s1.SessionID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Succeed)

			if s1.SessionID != s2.SessionID {
				t.Fatalf("\t%s\tShould be able to get back the same session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same session.", tests.Succeed)
			}

			if s1.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Succeed)
			}
		}
	}
}
