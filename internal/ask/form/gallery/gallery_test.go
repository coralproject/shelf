package gallery_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form/gallery"
)

// prefix is what we are looking to delete after the test.
const prefix = "FTEST"

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

	// fms, err := formfix.Get()
	// if err != nil {
	// 	t.Fatalf("%s\tShould be able retrieve form fixture : %s", tests.Failed, err)
	// }
	// t.Logf("%s\tShould be able retrieve form fixture.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	return nil, db
}

func teardown(t *testing.T, db *db.DB) {
	// if err := formfix.Remove(tests.Context, db, prefix); err != nil {
	// 	t.Fatalf("%s\tShould be able to remove the forms : %v", tests.Failed, err)
	// }
	// t.Logf("%s\tShould be able to remove the forms.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}
