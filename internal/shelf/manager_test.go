package shelf

import (
	"encoding/json"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/sfix"
)

func init() {
	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("XENIA")

	// Initialize MongoDB using the `tests.TestSession` as the name of the
	// master session.
	cfg := mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}
	tests.InitMongo(cfg)
}

// TestNewRelManager tests if we can create a new relationship manager in the db.
func TestNewRelManager(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelManager(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationship manager : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationship manager.", tests.Success)
	}()

	t.Log("Given the need to save a relationship manager into the database.")
	{
		t.Log("\tWhen using the default relationship manager")
		{
			raw, err := sfix.LoadRelManagerData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship manager fixture : %s", tests.Failed, err)
			}
			var rm RelManager
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship manager fixture : %s", tests.Failed, err)
			}
			if err := NewRelManager(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create a relationship manager : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a relationship manager.", tests.Success)
		}
	}
}
