package shelf

import (
	"encoding/json"
	"reflect"
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

// TestNewRelsAndViews tests if we can create a new relationships and views in the db.
func TestNewRelsAndViews(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelsAndViews(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships and views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships and views.", tests.Success)
	}()

	t.Log("Given the need to save a new set of relationships and views into the database.")
	{
		t.Log("\tWhen using the relsandviews.json test fixture")
		{
			raw, err := sfix.LoadRelAndViewData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship and view fixture : %s", tests.Failed, err)
			}
			var rm1 RelsAndViews
			if err := json.Unmarshal(raw, &rm1); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship and view fixture : %s", tests.Failed, err)
			}
			if err := NewRelsAndViews(tests.Context, db, rm1); err != nil {
				t.Fatalf("\t%s\tShould be able to create new relationships and views : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create new relationships and views.", tests.Success)
			rm2, err := GetRelsAndViews(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve relationships and views : %s", tests.Failed, err)
			}
			if !reflect.DeepEqual(rm1, rm2) {
				t.Logf("\t%+v", rm1)
				t.Logf("\t%+v", rm2)
				t.Errorf("\t%s\tShould be able to get back the same relationships and views.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same relationships and views.", tests.Success)
			}
		}
	}
}
