package query_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"
)

var context = "testing"

func init() {
	tests.Init()
}

//==============================================================================

// removeSets is used to clear out all the test sets from the collection.
// All test query sets must start with QSTEST in their name.
func removeSets(ses *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": bson.RegEx{Pattern: "QTEST"}}
		_, err := c.RemoveAll(q)
		return err
	}

	err := mongo.ExecuteDB(context, ses, "query_sets", f)
	if err != mgo.ErrNotFound {
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
	} else {
		t.Logf("\t%s\tShould load query record from file.", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSets(ses); err != nil {
			t.Errorf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		} else {
			t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
		}
	}()

	t.Log("Given the need to save a query set into the database.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, ses, *qs1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to create a query set.", tests.Success)
			}

			qs2, err := query.GetSetByName(context, ses, qs1.Name)
			if err != nil {
				t.Errorf("\t%s\tShould be able to retrieve the query set : %s", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to retrieve the query set.", tests.Success)
			}

			// TODO: WE NEED TO CHECK THE ENTIRE VALUE
			// Try using reflect.DeepEqual but you have pointers so you might
			// need to check the larger parts. This is a place holder.
			if qs1.Name != qs2.Name {
				t.Errorf("\t%s\tShould be able to get back the same query set.", tests.Failed)
				t.Logf("\t%+v", *qs1)
				t.Logf("\t%+v", *qs2)
			} else {
				t.Logf("\t%s\tShould be able to get back the same query set.", tests.Success)
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
	qs, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	} else {
		t.Logf("\t%s\tShould load query record from file.", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSets(ses); err != nil {
			t.Errorf("\t%s\tShould be able to remove the query set : %v", tests.Failed, err)
		} else {
			t.Logf("\t%s\tShould be able to remove the query set.", tests.Success)
		}
	}()

	t.Log("Given the need to retrieve a list of query sets.")
	{
		t.Log("\tWhen using fixture", fixture)
		{
			if err := query.CreateSet(context, ses, *qs); err != nil {
				t.Fatalf("\t%s\tShould be able to create a query set : %s", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to create a query set.", tests.Success)
			}

			qs.Name = qs.Name + "2"
			if err := query.CreateSet(context, ses, *qs); err != nil {
				t.Fatalf("\t%s\tShould be able to create a second query set : %s", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to create a second query set.", tests.Success)
			}

			names, err := query.GetSetNames(context, ses)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the query set names : %v", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to retrieve the query set names", tests.Success)
			}

			if len(names) != 2 {
				t.Errorf("\t%s\tShould have two query sets : %s", tests.Failed, names)
			} else {
				t.Logf("\t%s\tShould have atleast one query record name: %s", tests.Success, names)
			}

			if !strings.Contains(names[0], qsName) || !strings.Contains(names[1], qsName) {
				t.Errorf("\t%s\tShould have \"%s\" in the name.", tests.Failed, qsName)
			} else {
				t.Logf("\t%s\tShould have \"%s\" in the name.", tests.Success, qsName)
			}
		}
	}
}

/*
DON'T THINK WE NEED THIS TEST. TRY TO USE THE API EXCEPT FOR REMOVING EVERYTHING IN THE END.
THAT IS MY FAULT BECAUSE I WAS WRONG BEFORE. I THOUGHT I HAD TOLD YOU THIS IN THE MORNING.

// TestSetNamesList validates the accuracy of the names list returned.
func TestSetNamesList(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	q, err := getQuery()
	if err != nil {
		t.Fatalf("%s\t\tShould load query record from file : %v", tests.Failed, err)
	} else {
		t.Logf("%s\t\tShould load query record from file", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSession(ses); err != nil {
			t.Errorf("%s\t\tShould have removed query record from dbs : %v", tests.Failed, err)
		} else {
			t.Logf("%s\t\tShould have removed query record from db ", tests.Success)
		}
	}()

	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving the need to retrieve all names")
		{

			err := query.CreateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			var names []bson.M

			f := func(c *mgo.Collection) error {
				q := bson.M{"name": 1}
				return c.Find(nil).Select(q).Sort("name").All(&names)
			}

			if err := mongo.ExecuteDB("Test", ses, collection, f); err != nil {
				t.Errorf("%s\t\tShould have retrieved query record names successfully %s", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have retrieved query record names successfully", tests.Success)

				var prefixed int
				var unprefixed int

				for _, item := range names {
					if !strings.HasPrefix(item["name"].(string), "test") {
						unprefixed++
						continue
					}
					prefixed++
				}

				if unprefixed == len(names) {
					t.Errorf("%s\t\tShould find names are without 'test' prefix in list", tests.Failed)
				} else {
					t.Logf("%s\t\tShould find names are without 'test' prefix in list", tests.Success)
				}

				if prefixed == 0 {
					t.Errorf("%s\t\tShould find name with 'test' prefix in names list", tests.Failed)
				} else {
					t.Logf("%s\t\tShould find name with 'test' prefix in names list", tests.Success)
				}
			}

		}
	}
}
*/

/*
DON'T THINK WE NEED THIS TEST. TRY TO USE THE API EXCEPT FOR REMOVING EVERYTHING IN THE END.
THAT IS MY FAULT BECAUSE I WAS WRONG BEFORE. I THOUGHT I HAD TOLD YOU THIS IN THE MORNING.

// TestGetSetByName validates the retrieval of a query.Set record using its name.
func TestGetSetByName(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	q, err := getQuery()
	if err != nil {
		t.Fatalf("%s\t\tShould load query record from file : %v", tests.Failed, err)
	} else {
		t.Logf("%s\t\tShould load query record from file", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSession(ses); err != nil {
			t.Errorf("%s\t\tShould have removed query record from dbs : %v", tests.Failed, err)
		} else {
			t.Logf("%s\t\tShould have removed query record from db ", tests.Success)
		}
	}()

	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving a query record's name")
		{

			err := query.CreateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %s", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			ses := mongo.GetSession()
			defer ses.Close()

			qs, err := query.GetSetByName("Tests", ses, record)
			if err != nil {
				t.Errorf("%s\t\tShould have retrieved query record name[%s] successfully : %v", record, tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have retrieved query record name[%s] successfully", record, tests.Success)

				if qs.Description != q.Description {
					t.Errorf("%s\t\tShould have matching description with retrieved record ", tests.Failed)
				} else {
					t.Logf("%s\t\tShould have matching description with retrieved record ", tests.Success)
				}

				if len(qs.Params) != len(q.Params) {
					t.Errorf("%s\t\tShould have matching param size with retrieved record ", tests.Failed)
				} else {
					t.Logf("%s\t\tShould have matching param size with retrieved record", tests.Success)
				}

				if len(qs.Queries) != len(q.Queries) {
					t.Errorf("%s\t\tShould have matching rule size with retrieved record", tests.Failed)
				} else {
					t.Logf("%s\t\tShould have matching rule size with retrieved record", tests.Success)
				}

				if qs.Name != q.Name {
					t.Errorf("%s\t\tShould have matching name with retrieved record", tests.Failed)
				} else {
					t.Logf("%s\t\tShould have matching name with retrieved record", tests.Success)
				}

				if qs.Enabled != q.Enabled {
					t.Errorf("%s\t\tShould match run 'enabled' flag with retrieved record", tests.Failed)
				} else {
					t.Logf("%s\t\tShould have match run 'enabled' flag with retrieved record", tests.Success)
				}
			}

			if _, err := query.GetSetByName("Tests", ses, "advice_spending"); err == nil {
				t.Errorf("%s\t\tShould have failed to retrieve query record name[%s] successfully : %v", "advice_spending", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have failed to retrieve query record name[%s] successfully", "advice_spending", tests.Success)
			}
		}
	}
}
*/

/*
LET'S FIX THIS TESTS TO FOLLOW THE PATTERNS IN THE FIRST TWO TESTS. STUDY THE CODE.
ALSO STUDY THE CODE IN AUTH_TEST.

// TestUpdateSet set validates update operation of a given record.
func TestUpdateSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	q, err := getQuery()
	if err != nil {
		t.Fatalf("%s\t\tShould load query record from file : %v", tests.Failed, err)
	} else {
		t.Logf("%s\t\tShould load query record from file", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSession(ses); err != nil {
			t.Errorf("%s\t\tShould have removed query record from dbs : %v", tests.Failed, err)
		} else {
			t.Logf("%s\t\tShould have removed query record from db ", tests.Success)
		}
	}()

	t.Log("Given the need to update a query record")
	{
		t.Log("\tWhen giving the a query record")
		{

			err := query.UpdateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			if err := recordThoseExists(ses); err != nil {
				t.Errorf("%s\t\tShould have found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have found a query record ", tests.Success)
			}

			if err := recordThoseNotExists(ses); err == nil {
				t.Errorf("%s\t\tShould have not found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have not found a query record ", tests.Success)
			}

		}
	}
}
*/

/*
LET'S FIX THIS TESTS TO FOLLOW THE PATTERNS IN THE FIRST TWO TESTS. STUDY THE CODE.
ALSO STUDY THE CODE IN AUTH_TEST.

// TestDeleteSet validates the removal of a query from the database.
func TestDeleteSet(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	q, err := getQuery()
	if err != nil {
		t.Fatalf("%s\t\tShould load query record from file : %v", tests.Failed, err)
	} else {
		t.Logf("%s\t\tShould load query record from file", tests.Success)
	}

	ses := mongo.GetSession()
	defer ses.Close()

	defer func() {
		if err := removeSession(ses); err != nil {
			t.Errorf("%s\t\tShould have removed query record from dbs : %v", tests.Failed, err)
		} else {
			t.Logf("%s\t\tShould have removed query record from db ", tests.Success)
		}
	}()

	t.Log("Given the need to remove a query from the database")
	{
		t.Log("\tWhen giving the record's name")
		{
			err := query.CreateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			if err := query.DeleteSet("Tests", ses, record); err != nil {
				t.Errorf("%s\t\tShould have added new query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			if err := recordThoseExists(ses); err == nil {
				t.Errorf("%s\t\tShould have found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have found a query record ", tests.Success)
			}

			if err := recordThoseNotExists(ses); err != nil {
				t.Errorf("%s\t\tShould have not found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have not found a query record ", tests.Success)
			}

		}
	}
}
*/

/*
LET'S FIX THIS TESTS TO FOLLOW THE PATTERNS IN THE FIRST TWO TESTS. STUDY THE CODE.
ALSO STUDY THE CODE IN AUTH_TEST.

// TestNoSession tests the when a nil session is used.
func TestNoSession(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	var q query.Set

	t.Log("Given the need to validate the use of a nil session")
	{
		t.Log("\tWhen giving the no mongo session")
		{
			err := query.CreateSet("Tests", nil, q)
			if err == nil {
				t.Errorf("%s\t\tShould not be able to create a record", tests.Failed)
			} else {
				t.Logf("%s\t\tShould not be able to create a record", tests.Success)
			}

			if _, err := query.GetSetNames("Test", nil); err == nil {
				t.Errorf("%s\t\tShould not be able to retrieve record names", tests.Failed)
			} else {
				t.Logf("%s\t\tShould not be able to retrieve record names", tests.Success)
			}

			if _, err := query.GetSetByName("Test", nil, record); err == nil {
				t.Errorf("%s\t\tShould not be able to retrieve a record", tests.Failed)
			} else {
				t.Logf("%s\t\tShould not be able to retrieve a record", tests.Success)
			}

			if err := query.UpdateSet("Test", nil, q); err == nil {
				t.Errorf("%s\t\tShould not be able to update record", tests.Failed)
			} else {
				t.Logf("%s\t\tShould not be able to update record", tests.Success)
			}

			if err := query.DeleteSet("Tests", nil, record); err == nil {
				t.Errorf("%s\t\tShould not be able to delete a record", tests.Failed)
			} else {
				t.Logf("%s\t\tShould not be able to delete a record", tests.Success)
			}

		}
	}
}
*/
