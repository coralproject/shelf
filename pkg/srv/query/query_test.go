package query_test

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/coralproject/shelf/pkg/db"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var context = "testing"

func init() {
	tests.Init()
}

//==============================================================================

// removeSets is used to clear out all the test sets from the collection.
// All test query sets must start with QSTEST in their name.
func removeSets(db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "QTEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(context, "query_sets", f); err != nil {
		return err
	}

	return nil
}

// getFixture retrieves a query record from the filesystem.
func getFixture(filePath string) (*query.Set, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var qs query.Set
	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		return nil, err
	}

	return &qs, nil
}

//==============================================================================

// TestCreateQuery tests if we can create a query record in the db.
func TestCreateQuery(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to save a query set into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			qs2, err := query.GetSetByName(context, db, qs1.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set.", tests.Success)

			if qs1.Name != qs2.Name {
				t.Logf("\t%+v", qs1.Name)
				t.Logf("\t%+v", qs2.Name)
				t.Errorf("\t%s\tShould be able to get back the same Name value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same Name value.", tests.Success)
			}

			if qs1.Enabled != qs2.Enabled {
				t.Logf("\t%+v", qs1.Enabled)
				t.Logf("\t%+v", qs2.Enabled)
				t.Errorf("\t%s\tShould be able to get back the same Enabled value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same Enabled value.", tests.Success)
			}

			if qs1.Description != qs2.Description {
				t.Logf("\t%+v", qs1.Description)
				t.Logf("\t%+v", qs2.Description)
				t.Errorf("\t%s\tShould be able to get back the same Description value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same Description value.", tests.Success)
			}

			if len(qs1.Params) != len(qs2.Params) {
				t.Logf("\t%+v", qs1.Params)
				t.Logf("\t%+v", qs2.Params)
				t.Errorf("\t%s\tShould be able to get back the same number of Param values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same number of Param values.", tests.Success)
			}

			for ind, param1 := range qs1.Params {
				param2 := qs2.Params[ind]
				if !reflect.DeepEqual(param1, param2) {
					t.Logf("\t%+v", param1)
					t.Logf("\t%+v", param2)
					t.Errorf("\t%s\tShould be able to get back the same Param value.", tests.Failed)
				} else {
					t.Logf("\t%s\tShould be able to get back the same Param value.", tests.Success)
				}
			}

			if len(qs1.Queries) != len(qs2.Queries) {
				t.Logf("\t%+v", qs1.Queries)
				t.Logf("\t%+v", qs2.Queries)
				t.Errorf("\t%s\tShould be able to get back the same number of Query values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same number of Query values.", tests.Success)
			}

			for ind, qu := range qs1.Queries {
				qu2 := qs2.Queries[ind]

				if qu.Type != qu2.Type {
					t.Errorf("\t%s\tShould have matching Type for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching Type for query at index %d", tests.Success, ind)
				}

				if qu.Description != qu2.Description {
					t.Errorf("\t%s\tShould have matching description for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching description for query at index %d", tests.Success, ind)
				}

				if qu.Continue != qu2.Continue {
					t.Errorf("\t%s\tShould have matching continue flag for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching continue flag for query at index %d", tests.Success, ind)
				}

				if !reflect.DeepEqual(*qu.SaveOptions, *qu2.SaveOptions) {
					t.Errorf("\t%s\tShould have matching save_options for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching save_options for query at index %d", tests.Success, ind)
				}

				if !reflect.DeepEqual(*qu.ScriptOptions, *qu2.ScriptOptions) {
					t.Errorf("\t%s\tShould have matching script_options for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching script_options for query at index %d", tests.Success, ind)
				}

				if !reflect.DeepEqual(*qu.VarOptions, *qu2.VarOptions) {
					t.Errorf("\t%s\tShould have matching script_options for query at index %d", tests.Failed, ind)
				} else {
					t.Logf("\t%s\tShould have matching script_options for query at index %d", tests.Success, ind)
				}

				for sindex, src := range qu.Scripts {
					csrc := qu2.Scripts[sindex]

					if csrc != src {
						t.Logf("Script Src(Index: %d): %s", sindex, src)
						t.Logf("Script Src(Index: %d): %s", sindex, csrc)
						t.Errorf("\t%s\tShould have matching src for query index: %d at scripts index %d", tests.Failed, ind, sindex)
					} else {
						t.Logf("\t%s\tShould have matching src for query index: %d at scripts index %d", tests.Success, ind, sindex)
					}

				}

			}
		}
	}
}

// TestGetSetNames validates retrieval of query.Set record names.
func TestGetSetNames(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "spending_advice"

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to retrieve a list of query sets.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			qs1.Name += "2"
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a second query set.", tests.Success)

			names, err := query.GetSetNames(context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set names : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the query set names", tests.Success)

			if len(names) != 2 {
				t.Fatalf("\t%s\tShould have two query sets : %s", tests.Failed, names)
			}
			t.Logf("\t%s\tShould have atleast one query record name: %s", tests.Success, names)

			if !strings.Contains(names[0], qsName) || !strings.Contains(names[1], qsName) {
				t.Errorf("\t%s\tShould have \"%s\" in the name.", tests.Failed, qsName)
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, qsName)
			}
		}
	}
}

// TestUpdateSet set validates update operation of a given record.
func TestUpdateSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to update a query set into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			qs2 := *qs1
			qs2.Params = append(qs2.Params, query.SetParam{
				Name:    "group",
				Default: "1",
				Desc:    "provides the group number for the query script",
			})

			if err := query.UpdateSet(context, db, &qs2); err != nil {
				t.Fatalf("\t%s\tShould be able to update a query set record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a query set record.", tests.Success)

			updSet, err := query.GetSetByName(context, db, qs2.Name)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a query set record: %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve a query set record.", tests.Success)

			if updSet.Name != qs1.Name {
				t.Errorf("\t%s\tShould be able to get back the same query set name.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query set name.", tests.Success)
			}

			if lendiff := len(updSet.Params) - len(qs1.Params); lendiff != 1 {
				t.Errorf("\t%s\tShould have one more parameter in set record: %d", tests.Failed, lendiff)
			} else {
				t.Logf("\t%s\tShould have one more parameter in set record.", tests.Success)
			}

			oparam := qs1.Params[0]
			uparam := updSet.Params[0]

			if !reflect.DeepEqual(oparam, uparam) {
				t.Logf("\t%+v", oparam)
				t.Logf("\t%+v", uparam)
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

	qsName := "QTEST_spending_advice"
	qsBadName := "QTEST_brod_advice"

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to delete a query set in the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			if err := query.DeleteSet(context, db, qsName); err != nil {
				t.Fatalf("\t%s\tShould be able to delete a query set using its name[%s]: %s", tests.Failed, qsName, err)
			}
			t.Logf("\t%s\tShould be able to delete a query set using its name[%s]:", tests.Success, qsName)

			if err := query.DeleteSet(context, db, qsBadName); err == nil {
				t.Fatalf("\t%s\tShould not be able to delete a query set using wrong name name[%s]", tests.Failed, qsBadName)
			}
			t.Logf("\t%s\tShould not be able to delete a query set using wrong name name[%s]", tests.Success, qsBadName)

			if _, err := query.GetSetByName(context, db, qsName); err == nil {
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

	qsName := "QTEST_spending_desire"

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	db := db.NewMGO()
	defer db.CloseMGO()

	defer func() {
		if err := removeSets(db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
	}()

	t.Log("Given the need to validate bad query name response.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, db, qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a query set.", tests.Success)

			if _, err := query.GetSetByName(context, db, qsName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate query set with Name[%s] does not exists: %s", tests.Failed, qsName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate query set with Name[%s] does not exists.", tests.Success, qsName)

			if err := query.DeleteSet(context, db, qsName); err == nil {
				t.Fatalf("\t%s\tShould be able to validate query set with Name[%s] can not be deleted: %s", tests.Failed, qsName, errors.New("Record Exists"))
			}
			t.Logf("\t%s\tShould be able to validate query set with Name[%s] can not be deleted.", tests.Success, qsName)
		}
	}
}

// TestAPIFailure validates the failure of the api using a nil session.
func TestAPIFailure(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	qsName := "QTEST_spending_desire"

	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)

	t.Log("Given the need to to validate failure of API with bad session.")
	{
		t.Log("When giving a nil session")
		{
			if err := query.CreateSet(context, nil, qs1); err == nil {
				t.Fatalf("\t%s\tShould be refused create by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused create by api with bad session: %s", tests.Success, err)

			if err := query.UpdateSet(context, nil, qs1); err == nil {
				t.Fatalf("\t%s\tShould be refused update by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused update by api with bad session: %s", tests.Success, err)

			if _, err := query.GetSetByName(context, nil, qsName); err == nil {
				t.Fatalf("\t%s\tShould be refused get request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused get request by api with bad session: %s", tests.Success, err)

			if _, err := query.GetSetNames(context, nil); err == nil {
				t.Fatalf("\t%s\tShould be refused names request by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused names request by api with bad session: %s", tests.Success, err)

			if err := query.DeleteSet(context, nil, qsName); err == nil {
				t.Fatalf("\t%s\tShould be refused delete by api with bad session", tests.Failed)
			}
			t.Logf("\t%s\tShould be refused delete by api with bad session: %s", tests.Success, err)

		}
	}
}

// TestVariableSubsitution validates the process of subsituting variables within
// source scripts.
func TestVariableSubsitution(t *testing.T) {
	t.Logf("Given the need to subsitute variables into a script source.")
	{

		script := "{ \"$match\" : { \"user_id\" : \"#user_id#\", \"category\" : \"#category#\" }}"

		t.Logf("\tWhen giving a script source and a variable map")
		{

			model := map[string]string{"user_id": "10", "category": "petrol"}
			wanted := "{ \"$match\" : { \"user_id\" : \"10\", \"category\" : \"petrol\" }}"
			rendered := query.RenderScript(script, model)

			if rendered != wanted {
				t.Logf("\tModel: %+v", model)
				t.Logf("\tRender: %+v", rendered)
				t.Logf("\tExpected: %+v", wanted)
				t.Errorf("\t%s\tShould have matched expected output from script src rendering", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched expected output from script src rendering", tests.Success)
			}

			model2 := map[string]string{"category": "diesel"}
			wanted2 := "{ \"$match\" : { \"user_id\" : \"#user_id#\", \"category\" : \"diesel\" }}"
			rendered2 := query.RenderScript(script, model2)

			if rendered2 != wanted2 {
				t.Logf("\tModel: %+v", model2)
				t.Logf("\tRender: %+v", rendered2)
				t.Logf("\tExpected: %+v", wanted2)
				t.Errorf("\t%s\tShould have matched expected output with partial data", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched expected output with partial data", tests.Success)
			}

			model3 := map[string]string{"name": "fieldsteign", "age": "30"}
			rendered3 := query.RenderScript(script, model3)

			if rendered3 != script {
				t.Logf("\tModel: %+v", model3)
				t.Logf("\tRender: %+v", rendered3)
				t.Logf("\tExpected: %+v", script)
				t.Errorf("\t%s\tShould have matched input script source when model lacks proper data", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched input script source when model lacks proper data", tests.Success)
			}

		}
	}
}
