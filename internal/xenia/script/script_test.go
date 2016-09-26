package script_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/script"
	"github.com/coralproject/shelf/internal/xenia/script/sfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "STEST_O"

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
func setup(t *testing.T, fixture string) (script.Script, *db.DB) {
	tests.ResetLog()

	scr, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("%s\tShould load query mask record from file : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould load query mask record from file.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}

	return scr, db
}

// teardown deinitializes for each indivdual test.
func teardown(t *testing.T, db *db.DB) {
	if err := sfix.Remove(db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the query mask : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the query mask.", tests.Success)

	db.CloseMGO(tests.Context)

	tests.DisplayLog()
}

//==============================================================================

// TestUpsertCreateScript tests if we can create a script record in the db.
func TestUpsertCreateScript(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to save a script into the database.")
	{
		t.Log("\tWhen using script", scr1)
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			if _, err := script.GetLastHistoryByName(tests.Context, db, scr1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the script from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the script from history.", tests.Success)

			scr2, err := script.GetByName(tests.Context, db, scr1.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the script.", tests.Success)

			if !reflect.DeepEqual(scr1, scr2) {
				t.Logf("\t%+v", scr1)
				t.Logf("\t%+v", scr2)
				t.Errorf("\t%s\tShould be able to get back the same script values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same script values.", tests.Success)
			}
		}
	}
}

// TestGetScriptNames validates retrieval of Script record names.
func TestGetScriptNames(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	scrName := prefix + "_basic"

	t.Log("Given the need to retrieve a list of scripts.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := scr1
			scr2.Name += "2"
			if err := script.Upsert(tests.Context, db, scr2); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second script.", tests.Success)

			names, err := script.GetNames(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the script names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the script names", tests.Success)

			var count int
			for _, name := range names {
				if len(name) > len(prefix) && name[0:len(prefix)] == prefix {
					count++
				}
			}

			// When tests are running in parallel with the query and exec package, we could
			// have more scripts.

			if count < 2 {
				t.Fatalf("\t%s\tShould have at least two scripts : %d : %v", tests.Failed, len(names), names)
			}
			t.Logf("\t%s\tShould have at least two scripts.", tests.Success)

			var found bool
			for _, n := range names {
				if strings.Contains(n, scrName) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("\t%s\tShould have \"%s\" in the name : %s", tests.Failed, scrName, names[0])
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, scrName)
			}
		}
	}
}

// TestGetScripts validates retrieval of all Script records.
func TestGetScripts(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to retrieve a list of scripts.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := scr1
			scr2.Name += "2"
			if err := script.Upsert(tests.Context, db, scr2); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second script.", tests.Success)

			scripts, err := script.GetAll(tests.Context, db, nil)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the scripts : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the scripts", tests.Success)

			var count int
			for _, scr := range scripts {
				if len(scr.Name) > len(prefix) && scr.Name[0:len(prefix)] == prefix {
					count++
				}
			}

			// When tests are running in parallel with the query and exec package, we could
			// have more scripts.

			if count < 2 {
				t.Fatalf("\t%s\tShould have at least two scripts : %d : %v", tests.Failed, len(scripts), scripts)
			}
			t.Logf("\t%s\tShould have at least two scripts.", tests.Success)

			var found int
			for _, s := range scripts {
				if s.Name == scr1.Name || s.Name == scr2.Name {
					found++
				}
			}

			if found != 2 {
				t.Errorf("\t%s\tShould have retrieve the correct scripts : found[%d]", tests.Failed, found)
			} else {
				t.Logf("\t%s\tShould have retrieve the correct scripts.", tests.Success)
			}
		}
	}
}

// TestGetScriptByNames validates retrieval of Script records by a set of names.
func TestGetScriptByNames(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to retrieve a list of script values.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := scr1
			scr2.Name += "2"
			if err := script.Upsert(tests.Context, db, scr2); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second script.", tests.Success)

			scripts, err := script.GetByNames(tests.Context, db, []string{scr1.Name, scr2.Name})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the scripts by names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the scripts by names", tests.Success)

			var count int
			for _, scr := range scripts {
				if len(scr.Name) > len(prefix) && scr.Name[0:len(prefix)] == prefix {
					count++
				}
			}

			// When tests are running in parallel with the query and exec package, we could
			// have more scripts.

			if count < 2 {
				t.Fatalf("\t%s\tShould have at least two scripts : %d : %v", tests.Failed, len(scripts), scripts)
			}
			t.Logf("\t%s\tShould have at least two scripts.", tests.Success)

			var found int
			for _, s := range scripts {
				if s.Name == scr1.Name || s.Name == scr2.Name {
					found++
				}
			}

			if found != 2 {
				t.Errorf("\t%s\tShould have retrieve the correct scripts : found[%d]", tests.Failed, found)
			} else {
				t.Logf("\t%s\tShould have retrieve the correct scripts.", tests.Success)
			}
		}
	}
}

// TestGetLastScriptHistoryByName validates retrieval of Script from the history
// collection.
func TestGetLastScriptHistoryByName(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	scrName := prefix + "_basic"

	t.Log("Given the need to retrieve a script from history.")
	{
		t.Log("\tWhen using script", scr1)
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr1.Commands = append(scr1.Commands, map[string]interface{}{"command": 4})

			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2, err := script.GetLastHistoryByName(tests.Context, db, scrName)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the last script from history : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the last script from history.", tests.Success)

			if !reflect.DeepEqual(scr1, scr2) {
				t.Logf("\t%+v", scr1)
				t.Logf("\t%+v", scr2)
				t.Errorf("\t%s\tShould be able to get back the same script values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same script values.", tests.Success)
			}
		}
	}
}

// TestUpsertUpdateScript validates update operation of a given Script.
func TestUpsertUpdateScript(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	t.Log("Given the need to update a script into the database.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := scr1
			scr2.Commands = append(scr2.Commands, map[string]interface{}{"command": 4})

			if err := script.Upsert(tests.Context, db, scr2); err != nil {
				t.Fatalf("\t%s\tShould be able to update a script record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a script record.", tests.Success)

			if _, err := script.GetLastHistoryByName(tests.Context, db, scr1.Name); err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the script from history: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the script from history.", tests.Success)

			updScr, err := script.GetByName(tests.Context, db, scr2.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a script record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a script record.", tests.Success)

			if updScr.Name != scr1.Name {
				t.Errorf("\t%s\tShould be able to get back the same script name.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same script name.", tests.Success)
			}

			if lendiff := len(updScr.Commands) - len(scr1.Commands); lendiff != 1 {
				t.Errorf("\t%s\tShould have one more parameter in script record: %d", tests.Failed, lendiff)
			} else {
				t.Logf("\t%s\tShould have one more parameter in script record.", tests.Success)
			}

			if !reflect.DeepEqual(scr2.Commands[0], updScr.Commands[0]) {
				t.Logf("\t%+v", scr2.Commands[0])
				t.Logf("\t%+v", updScr.Commands[0])
				t.Errorf("\t%s\tShould be abe to validate the script param values in db.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be abe to validate the script param values in db.", tests.Success)
			}
		}
	}
}

// TestDeleteScript validates the removal of a script from the database.
func TestDeleteScript(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	scrName := prefix + "_basic"
	scrBadName := prefix + "_basic_advice"

	t.Log("Given the need to delete a script in the database.")
	{
		t.Log("\tWhen using script", scr1)
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			if err := script.Delete(tests.Context, db, scrName); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a script using its name[%s]: %s", tests.Failed, scrName, err)
			}
			t.Logf("\t%s\tShould be able to delete a script using its name[%s]:", tests.Success, scrName)

			if err := script.Delete(tests.Context, db, scrBadName); err == nil {
				t.Fatalf("\t%s\tShould not be able to delete a script using wrong name name[%s]", tests.Failed, scrBadName)
			}
			t.Logf("\t%s\tShould not be able to delete a script using wrong name name[%s]", tests.Success, scrBadName)

			if _, err := script.GetByName(tests.Context, db, scrName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate script with Name[%s] does not exists: %s", tests.Failed, scrName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate script  with Name[%s] does not exists:", tests.Success, scrName)
		}
	}
}

// TestAPIFailureScripts validates the failure of the api using a nil session.
func TestAPIFailureScripts(t *testing.T) {
	const fixture = "basic.json"
	scr1, db := setup(t, fixture)
	defer teardown(t, db)

	scrName := prefix + "_unknown"

	t.Log("Given the need to validate failure of API with bad session.")
	{
		t.Log("When giving a nil session")
		{
			err := script.Upsert(tests.Context, nil, scr1)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused create by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused create by api with bad session: %s", tests.Success, err)

			_, err = script.GetNames(tests.Context, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = script.GetAll(tests.Context, nil, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = script.GetByName(tests.Context, nil, scrName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = script.GetByNames(tests.Context, nil, nil)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			_, err = script.GetLastHistoryByName(tests.Context, nil, scrName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			err = script.Delete(tests.Context, nil, scrName)
			if err == nil {
				t.Fatalf("\t%s\tShould be refused delete by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused delete by api with bad session: %s", tests.Success, err)
		}
	}
}
