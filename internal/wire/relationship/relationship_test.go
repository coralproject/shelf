package relationship_test

import (
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/wire/relationship"
	"github.com/coralproject/shelf/internal/wire/relationship/relationshipfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "RTEST_"

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

//==============================================================================

// setup initializes for each indivdual test.
func setup(t *testing.T) ([]relationship.Relationship, *db.DB) {
	tests.ResetLog()

	rels, err := relationshipfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould load relationship records from file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load relationship records from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	return rels, db
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	if err := relationshipfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the relationship records : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the relationship records.", tests.Success)

	db.CloseMGO(tests.Context)

	tests.DisplayLog()
}

//==============================================================================

// TestUpsertDelete tests if we can add/remove a relationship to/from the db.
func TestUpsertDelete(t *testing.T) {
	rels, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to upsert and delete relationships.")
	{
		t.Log("\tWhen starting from an empty relationships collection")
		{

			//----------------------------------------------------------------------
			// Upsert the relationship.

			if err := relationship.Upsert(tests.Context, db, &rels[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a relationship.", tests.Success)

			//----------------------------------------------------------------------
			// Get the relationship.

			rel, err := relationship.GetByPredicate(tests.Context, db, rels[0].Predicate)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the relationship by predicate : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the relationship by predicate.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the relationship we expected.

			if !reflect.DeepEqual(rels[0], *rel) {
				t.Logf("\t%+v", rels[0])
				t.Logf("\t%+v", rel)
				t.Fatalf("\t%s\tShould be able to get back the same relationship.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationship.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the relationship.

			if err := relationship.Delete(tests.Context, db, rels[0].Predicate); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the relationship.", tests.Success)

			//----------------------------------------------------------------------
			// Get the relationship.

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
	rels, db := setup(t)
	defer teardown(t, db)

	t.Log("Given the need to get all the relationships in the database.")
	{
		t.Log("\tWhen starting from an empty relationships collection")
		{

			for _, rel := range rels {
				if err := relationship.Upsert(tests.Context, db, &rel); err != nil {
					t.Fatalf("\t%s\tShould be able to upsert a relationships : %s", tests.Failed, err)
				}
			}
			t.Logf("\t%s\tShould be able to upsert relationships.", tests.Success)

			rels2, err := relationship.GetAll(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get all relationships : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get all relationships.", tests.Success)

			if !reflect.DeepEqual(rels, rels2) {
				t.Logf("\t%+v", rels)
				t.Logf("\t%+v", rels2)
				t.Fatalf("\t%s\tShould be able to get back the same relationships.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same relationships.", tests.Success)
		}
	}
}
