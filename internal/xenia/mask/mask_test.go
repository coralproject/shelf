package mask_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/mask"
	"github.com/coralproject/shelf/internal/xenia/mask/mfix"
)

// collection is what we are looking to delete after the test.
const collection = "test_xenia_data"

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
func setup(t *testing.T, fixture string) ([]mask.Mask, *db.DB) {
	tests.ResetLog()

	masks, err := mfix.Get(fixture)
	if err != nil {
		t.Fatalf("%s\tShould load query mask record from file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load query mask record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	return masks, db
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	if err := mfix.Remove(db, collection); err != nil {
		t.Fatalf("%s\tShould be able to remove the query mask : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the query mask.", tests.Success)

	db.CloseMGO(tests.Context)

	tests.DisplayLog()
}

//==============================================================================

// TestUpsertCreateMask tests if we can create a query mask record in the db.
func TestUpsertCreateMask(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to save a query mask into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := mask.Upsert(tests.Context, db, masks[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query mask.", tests.Success)

			if _, err := mask.GetLastHistoryByName(tests.Context, db, masks[0].Collection, masks[0].Field); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query mask from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query mask from history.", tests.Success)

			msk, err := mask.GetByName(tests.Context, db, masks[0].Collection, masks[0].Field)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query mask.", tests.Success)

			if !reflect.DeepEqual(masks[0], msk) {
				t.Logf("\t%+v", masks[0])
				t.Logf("\t%+v", msk)
				t.Errorf("\t%s\tShould be able to get back the same query mask values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query mask values.", tests.Success)
			}
		}
	}
}

// TestGetMasks validates retrieval of all query mask records.
func TestGetMasks(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to retrieve a list of query masks.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			for _, msk := range masks {
				if err := mask.Upsert(tests.Context, db, msk); err != nil {
					t.Fatalf("\t%s\tShould be able to create a query mask : %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to create a query mask.", tests.Success)
			}

			msks, err := mask.GetAll(tests.Context, db, nil)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query masks : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query masks", tests.Success)

			var count int
			for _, msk := range msks {
				if msk.Collection == collection || (msk.Collection == "*" && msk.Field == "test") {
					count++
				}
			}

			if count != len(masks) {
				t.Fatalf("\t%s\tShould have two query masks : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two query masks.", tests.Success)
		}
	}
}

// TestGetMaskByCollection validates retrieval of all query mask records by collection.
func TestGetMaskByCollection(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to retrieve a list of query masks by collection.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			for _, msk := range masks {
				if err := mask.Upsert(tests.Context, db, msk); err != nil {
					t.Fatalf("\t%s\tShould be able to create a query mask : %s", tests.Failed, err)
				}
				t.Logf("\t%s\tShould be able to create a query mask.", tests.Success)
			}

			msks, err := mask.GetByCollection(tests.Context, db, masks[0].Collection)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query masks : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query masks", tests.Success)

			var count int
			for _, msk := range msks {
				if msk.Collection == collection || (msk.Collection == "*" && msk.Field == "test") {
					count++
				}
			}

			if count != len(masks) {
				t.Fatalf("\t%s\tShould have two query masks : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two query masks.", tests.Success)
		}
	}
}

// TestGetLastMaskHistoryByName validates retrieval of Mask from the history
// collection.
func TestGetLastMaskHistoryByName(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to retrieve a mask from history.")
	{
		t.Log("\tWhen using mask", masks[0])
		{
			if err := mask.Upsert(tests.Context, db, masks[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a mask.", tests.Success)

			if err := mask.Upsert(tests.Context, db, masks[1]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a mask.", tests.Success)

			msk, err := mask.GetLastHistoryByName(tests.Context, db, masks[1].Collection, masks[1].Field)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the last mask from history : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the last mask from history.", tests.Success)

			if !reflect.DeepEqual(masks[1], msk) {
				t.Logf("\t%+v", masks[1])
				t.Logf("\t%+v", msk)
				t.Errorf("\t%s\tShould be able to get back the same mask values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same mask values.", tests.Success)
			}
		}
	}
}

// TestUpsertUpdateMask validates update operation of a given mask.
func TestUpsertUpdateMask(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to update a mask into the database.")
	{
		t.Log("\tWhen using two masks")
		{
			if err := mask.Upsert(tests.Context, db, masks[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a mask.", tests.Success)

			masks[0].Type = mask.MaskAll

			if err := mask.Upsert(tests.Context, db, masks[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to update a mask record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a mask record.", tests.Success)

			if _, err := mask.GetLastHistoryByName(tests.Context, db, masks[0].Collection, masks[0].Field); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the mask from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the mask from history.", tests.Success)

			updMsk, err := mask.GetByName(tests.Context, db, masks[0].Collection, masks[0].Field)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a mask record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a mask record.", tests.Success)

			if updMsk.Type != masks[0].Type {
				t.Logf("\t%+v", updMsk.Type)
				t.Logf("\t%+v", masks[0].Type)
				t.Errorf("\t%s\tShould have an updated mask record.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have an updated mask record.", tests.Success)
			}
		}
	}
}

// TestDeleteMask validates the removal of a mask from the database.
func TestDeleteMask(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to delete a mask in the database.")
	{
		t.Log("\tWhen using mask", masks[0])
		{
			if err := mask.Upsert(tests.Context, db, masks[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a mask.", tests.Success)

			if err := mask.Delete(tests.Context, db, masks[0].Collection, masks[0].Field); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a mask : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete a mask.", tests.Success)

			if err := mask.Delete(tests.Context, db, "collection", "field"); err == nil {
				t.Fatalf("\t%s\tShould not be able to delete a mask that does not exist.", tests.Failed)
			}
			t.Logf("\t%s\tShould not be able to delete a mask that does not exist.", tests.Success)

			if _, err := mask.GetByName(tests.Context, db, masks[0].Collection, masks[0].Field); err == nil {
				t.Fatalf("\t%s\tShould be able to validate mask does not exists: %s", tests.Failed, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate mask does not exists.", tests.Success)
		}
	}
}

// TestAPIFailureMasks validates the failure of the api using a nil session.
func TestAPIFailureMasks(t *testing.T) {
	const fixture = "basic.json"
	masks, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to validate failure of API with bad session.")
	{
		t.Log("When giving a nil session")
		{
			err := mask.Upsert(tests.Context, nil, masks[0])
			if err == nil {
				t.Fatalf("\t%s\tShould be refused create by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused create by api with bad session: %s", tests.Success, err)

			_, err = mask.GetAll(tests.Context, nil, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = mask.GetByCollection(tests.Context, nil, masks[0].Collection)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = mask.GetByName(tests.Context, nil, masks[0].Collection, masks[0].Field)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = mask.GetLastHistoryByName(tests.Context, nil, masks[0].Collection, masks[0].Field)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			err = mask.Delete(tests.Context, nil, masks[0].Collection, masks[0].Field)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused delete by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused delete by api with bad session: %s", tests.Success, err)
		}
	}
}
