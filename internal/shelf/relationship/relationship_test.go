package relationship_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/relationship"
	"github.com/coralproject/xenia/internal/shelf/relationship/relationshipfix"
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

	// Initialize the logging system.
	logLevel := func() int {
		return log.DEV
	}
	log.Init(os.Stdout, logLevel, log.Ldefault)
}

// prefix is what we are looking to delete after the test.
const prefix = "RTEST_"

// TestUpsertDelete tests if we can add/remove a relationship to/from the db.
func TestUpsertDelete(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := relationshipfix.Remove(tests.Context, db, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships.", tests.Success)
	}()

	t.Log("Given the need to upsert and delete relationships.")
	{
		t.Log("\tWhen starting from an empty relationships collection")
		{
			rels, err := relationshipfix.Get()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship fixture : %s", tests.Failed, err)
			}
			if err := relationship.Upsert(tests.Context, db, &rels[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a relationship.", tests.Success)
			rel, err := relationship.GetByPredicate(tests.Context, db, rels[0].Predicate)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the relationship by predicate : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the relationship by predicate.", tests.Success)
			if !reflect.DeepEqual(rels[0], *rel) {
				t.Logf("\t%+v", rels[0])
				t.Logf("\t%+v", rel)
				t.Fatalf("\t%s\tShould be able to get back the same relationship.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationship.", tests.Success)
			if err := relationship.Delete(tests.Context, db, rels[0].Predicate); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the relationship.", tests.Success)
			rel, err = relationship.GetByPredicate(tests.Context, db, rels[0].Predicate)
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a relationship with the deleted predicate : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a relationship with the deleted predicate.", tests.Success)

		}
	}
}

// TestGetAll tests if we can get all relationships from the db.
func TestGetAll(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := relationshipfix.Remove(tests.Context, db, prefix); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships.", tests.Success)
	}()

	t.Log("Given the need to get all the relationships in the database.")
	{
		t.Log("\tWhen starting from an empty relationships collection")
		{
			rels1, err := relationshipfix.Get()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship fixture : %s", tests.Failed, err)
			}
			for _, rel := range rels1 {
				if err := relationship.Upsert(tests.Context, db, &rel); err != nil {
					t.Fatalf("\t%s\tShould be able to upsert a relationship : %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to upsert a relationship.", tests.Success)
			}
			rels2, err := relationship.GetAll(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get all relationships : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get all relationships.", tests.Success)
			if !reflect.DeepEqual(rels1, rels2) {
				t.Logf("\t%+v", rels1)
				t.Logf("\t%+v", rels2)
				t.Fatalf("\t%s\tShould be able to get back the same relationships.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationships.", tests.Success)
		}
	}
}
