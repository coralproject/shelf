package gallery_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
	"github.com/coralproject/shelf/internal/ask/form/gallery/galleryfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "FGTEST"

func TestMain(m *testing.M) {
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

	os.Exit(m.Run())
}

func setup(t *testing.T, fixture string) ([]gallery.Gallery, *db.DB) {
	tests.ResetLog()

	gs, err := galleryfix.Get()
	if err != nil {
		t.Fatalf("%s\tShould be able retrieve gallery fixture : %s", tests.Failed, err)
	}
	t.Logf("%s\tShould be able retrieve gallery fixture.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	return gs, db
}

func teardown(t *testing.T, db *db.DB) {
	if err := galleryfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the galleries : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the galleries.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func Test_CreateDelete(t *testing.T) {
	gs, db := setup(t, "gallery")
	defer teardown(t, db)

	t.Log("Given the need to upsert and delete galleries.")
	{
		t.Log("\tWhen starting from an empty galleries collection")
		{
			//----------------------------------------------------------------------
			// Upsert the gallery.

			g, err := gallery.Create(tests.Context, db, gs[0].FormID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to upsert a gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to upsert a gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Get the gallery.

			rg, err := gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the gallery by id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the gallery by id.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the gallery we expected.

			if rg.ID.Hex() != g.ID.Hex() {
				t.Fatalf("\t%s\tShould be able to get back the same gallery.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the gallery.

			if err := gallery.Delete(tests.Context, db, g.ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the gallery : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the gallery.", tests.Success)

			//----------------------------------------------------------------------
			// Get the gallery.

			_, err = gallery.Retrieve(tests.Context, db, g.ID.Hex())
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a gallery with the deleted id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a gallery with the deleted id.", tests.Success)
		}
	}
}
