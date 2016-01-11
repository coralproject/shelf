package query_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
)

func init() {
	tests.Init("XENIA")

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

// TestUpsertCreateQuery tests if we can create a query record in the db.
func TestUpsertCreateQuery(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to save a query set into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			if _, err = query.Sets.GetLastHistoryByName(tests.Context, db, set1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set from history.", tests.Success)

			set2, err := query.Sets.GetByName(tests.Context, db, set1.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set.", tests.Success)

			if !reflect.DeepEqual(*set1, *set2) {
				t.Logf("\t%+v", set1)
				t.Logf("\t%+v", set2)
				t.Errorf("\t%s\tShould be able to get back the same query values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query values.", tests.Success)
			}
		}
	}
}

// TestGetSetNames validates retrieval of query Set record names.
func TestGetSetNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_basic"

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of query sets.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			set1.Name += "2"
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second query set.", tests.Success)

			names, err := query.Sets.GetNames(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set names", tests.Success)

			var count int
			for _, name := range names {
				if name[0:5] == "QTEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two query sets : %d", tests.Failed, len(names))
			}
			t.Logf("\t%s\tShould have two query sets.", tests.Success)

			if !strings.Contains(names[0], qsName) || !strings.Contains(names[1], qsName) {
				t.Errorf("\t%s\tShould have \"%s\" in the name.", tests.Failed, qsName)
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, qsName)
			}
		}
	}
}

// TestGetLastSetHistoryByName validates retrieval of query Set from the history
// collection.
func TestGetLastSetHistoryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_basic"

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to retrieve a query set from history.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			set1.Description = "Next Version"

			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			set2, err := query.Sets.GetLastHistoryByName(tests.Context, db, qsName)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the last query set from history : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the last query set from history.", tests.Success)

			if !reflect.DeepEqual(*set1, *set2) {
				t.Logf("\t%+v", set1)
				t.Logf("\t%+v", set2)
				t.Errorf("\t%s\tShould be able to get back the same query values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query values.", tests.Success)
			}
		}
	}
}

// TestUpsertUpdateQuery validates update operation of a given query Set.
func TestUpsertUpdateQuery(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to update a query set into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			set2 := *set1
			set2.Params = append(set2.Params, query.Param{
				Name:    "group",
				Default: "1",
				Desc:    "provides the group number for the query script",
			})

			if err := query.Sets.Upsert(tests.Context, db, &set2); err != nil {
				t.Fatalf("\t%s\tShould be able to update a query set record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a query set record.", tests.Success)

			if _, err = query.Sets.GetLastHistoryByName(tests.Context, db, set1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set from history.", tests.Success)

			updSet, err := query.Sets.GetByName(tests.Context, db, set2.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a query set record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a query set record.", tests.Success)

			if updSet.Name != set1.Name {
				t.Errorf("\t%s\tShould be able to get back the same query set name.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query set name.", tests.Success)
			}

			if lendiff := len(updSet.Params) - len(set1.Params); lendiff != 1 {
				t.Errorf("\t%s\tShould have one more parameter in set record: %d", tests.Failed, lendiff)
			} else {
				t.Logf("\t%s\tShould have one more parameter in set record.", tests.Success)
			}

			if !reflect.DeepEqual(set2.Params[0], updSet.Params[0]) {
				t.Logf("\t%+v", set2.Params[0])
				t.Logf("\t%+v", updSet.Params[0])
				t.Errorf("\t%s\tShould be abe to validate the query param values in db.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be abe to validate the query param values in db.", tests.Success)
			}
		}
	}
}

// TestDeleteSet validates the removal of a query from the database.
func TestDeleteSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_basic"
	qsBadName := "QTEST_brod_advice"

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to delete a query set in the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			if err := query.Sets.Delete(tests.Context, db, qsName); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a query set using its name[%s]: %s", tests.Failed, qsName, err)
			}
			t.Logf("\t%s\tShould be able to delete a query set using its name[%s]:", tests.Success, qsName)

			if err := query.Sets.Delete(tests.Context, db, qsBadName); err == nil {
				t.Fatalf("\t%s\tShould not be able to delete a query set using wrong name name[%s]", tests.Failed, qsBadName)
			}
			t.Logf("\t%s\tShould not be able to delete a query set using wrong name name[%s]", tests.Success, qsBadName)

			if _, err := query.Sets.GetByName(tests.Context, db, qsName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate query set with Name[%s] does not exists: %s", tests.Failed, qsName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate query set with Name[%s] does not exists:", tests.Success, qsName)
		}
	}
}

// TestUnknownName validates the behaviour of the query API when using a invalid/
// unknown query name.
func TestUnknownName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_unknown"

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := query.RemoveTestData(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to validate bad query name response.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.Sets.Upsert(tests.Context, db, set1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			if _, err := query.Sets.GetByName(tests.Context, db, qsName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate query set with Name[%s] does not exists: %s", tests.Failed, qsName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate query set with Name[%s] does not exists.", tests.Success, qsName)

			if err := query.Sets.Delete(tests.Context, db, qsName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate query set with Name[%s] can not be deleted: %s", tests.Failed, qsName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate query set with Name[%s] can not be deleted.", tests.Success, qsName)
		}
	}
}

// TestAPIFailureSet validates the failure of the api using a nil session.
func TestAPIFailureSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_unknown"

	const fixture = "basic.json"
	set1, err := query.GetFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	t.Log("Given the need to validate failure of API with bad session.")
	{
		t.Log("When giving a nil session")
		{
			err := query.Sets.Upsert(tests.Context, nil, set1)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused create by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused create by api with bad session: %s", tests.Success, err)

			_, err = query.Sets.GetNames(tests.Context, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = query.Sets.GetByName(tests.Context, nil, qsName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = query.Sets.GetLastHistoryByName(tests.Context, nil, qsName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			err = query.Sets.Delete(tests.Context, nil, qsName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused delete by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused delete by api with bad session: %s", tests.Success, err)
		}
	}
}
