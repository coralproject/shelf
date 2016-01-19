package script_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/coralproject/xenia/pkg/script"
	"github.com/coralproject/xenia/pkg/script/sfix"

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

// TestUpsertCreateScript tests if we can create a script record in the db.
func TestUpsertCreateScript(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

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

			if !reflect.DeepEqual(*scr1, *scr2) {
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
	tests.ResetLog()
	defer tests.DisplayLog()

	scrName := "STEST_basic"

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of scripts.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := *scr1
			scr2.Name += "2"
			if err := script.Upsert(tests.Context, db, &scr2); err != nil {
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
				if name[0:5] == "STEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two scripts : %d", tests.Failed, len(names))
			}
			t.Logf("\t%s\tShould have two scripts.", tests.Success)

			if !strings.Contains(names[0], scrName) || !strings.Contains(names[1], scrName) {
				t.Errorf("\t%s\tShould have \"%s\" in the name : %s", tests.Failed, scrName, names[0])
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, scrName)
			}
		}
	}
}

// TestGetScripts validates retrieval of all Script records.
func TestGetScripts(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of scripts.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr1.Name += "2"
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second script.", tests.Success)

			scrs, err := script.GetScripts(tests.Context, db, nil)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the scripts : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the scripts", tests.Success)

			var count int
			for _, scr := range scrs {
				if scr.Name[0:5] == "STEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two scripts : %d", tests.Failed, count)
			}
			t.Logf("\t%s\tShould have two scripts.", tests.Success)
		}
	}
}

// TestGetScriptByNames validates retrieval of Script records by a set of names.
func TestGetScriptByNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of script values.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := *scr1
			scr2.Name += "2"
			if err := script.Upsert(tests.Context, db, &scr2); err != nil {
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
				if scr.Name[0:5] == "STEST" {
					count++
				}
			}

			if count != 2 {
				t.Fatalf("\t%s\tShould have two scripts : %d", tests.Failed, len(scripts))
			}
			t.Logf("\t%s\tShould have two scripts.", tests.Success)

			if scripts[0].Name != scr1.Name || scripts[1].Name != scr2.Name {
				t.Errorf("\t%s\tShould have retrieve the correct scripts.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have retrieve the correct scripts.", tests.Success)
			}
		}
	}
}

// TestGetLastScriptHistoryByName validates retrieval of Script from the history
// collection.
func TestGetLastScriptHistoryByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	scrName := "STEST_basic"

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

	t.Log("Given the need to retrieve a script from history.")
	{
		t.Log("\tWhen using script", scr1)
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr1.Commands = append(scr1.Commands, "Command 4")

			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2, err := script.GetLastHistoryByName(tests.Context, db, scrName)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the last script from history : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the last script from history.", tests.Success)

			if !reflect.DeepEqual(*scr1, *scr2) {
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
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

	t.Log("Given the need to update a script into the database.")
	{
		t.Log("\tWhen using two scripts")
		{
			if err := script.Upsert(tests.Context, db, scr1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a script : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a script.", tests.Success)

			scr2 := *scr1
			scr2.Commands = append(scr2.Commands, "Command 4")

			if err := script.Upsert(tests.Context, db, &scr2); err != nil {
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
	tests.ResetLog()
	defer tests.DisplayLog()

	scrName := "STEST_basic"
	scrBadName := "STEST_basic_advice"

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := sfix.Remove(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the scripts : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the scripts.", tests.Success)
	}()

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
	tests.ResetLog()
	defer tests.DisplayLog()

	scrName := "STEST_unknown"

	const fixture = "basic.json"
	scr1, err := sfix.Get(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load script record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load script record from file.", tests.Success)

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

			_, err = script.GetByName(tests.Context, nil, scrName)
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
