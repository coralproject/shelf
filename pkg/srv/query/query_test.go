package query_test

import (
	"encoding/json"
	"os"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"
)

// collection used for testing the query CRUD API
var collection = "query"

// TestQueryAPI validates the operations of the query database and file loading API.
func TestQueryAPI(t *testing.T) {
	// Initialize the test environment.
	tests.Init()

	tests.ResetLog()
	defer tests.DisplayLog()

	qFile := "./fixtures/spending_advice.json"

	file, err := os.Open(qFile)
	if err != nil {
		t.Fatalf("\t\tShould open File[%q] without error %s", qFile, tests.Failed)
	} else {
		t.Logf("\t\tShould open File[%q] without error %s", qFile, tests.Success)
	}

	var qs query.Set

	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		t.Fatalf("\t\tShould load File[%q] without error %s", qFile, tests.Failed)
	} else {
		t.Logf("\t\tShould load File[%q] without error %s", qFile, tests.Success)
	}

	queryCreate(&qs, t)
	queryGetNames(t)
	queryGetByName(qs.Name, &qs, t)
	queryUpdate(&qs, t)
	queryDelete(&qs, t)
	tearDown(t)
}

// queryCreate validates the creation of a query in the databae.
func queryCreate(q *query.Set, t *testing.T) {
	t.Log("Given the need to save a query into the database")
	{
		t.Log("\tWhen giving a query object to save")
		{

			ses := mongo.GetSession()
			defer ses.Close()

			err := query.Create("Tests", ses, q)
			if err != nil {
				t.Errorf("\t\tShould have added new query record %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have added new query record %s", tests.Success)
			}
		}
	}
}

// queryGetNames validates the retrieval of a query using its name.
func queryGetNames(t *testing.T) {
	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving a query record's name")
		{

			ses := mongo.GetSession()
			defer ses.Close()

			names, err := query.GetNames("Test", ses)
			if err != nil {
				t.Errorf("\t\tShould have retrieved query record names successfully %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved query record names successfully %s", tests.Success)
			}

			if len(names) == 0 {
				t.Errorf("\t\tShould have atleast one query record name %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have atleast one query record name %s", tests.Success)
			}

			expectedName := "spending_advice"
			if names[0] != expectedName {
				t.Errorf("\t\tShould have first name equal %q %s", expectedName, tests.Failed)
			} else {
				t.Logf("\t\tShould have first name equal %q %s", expectedName, tests.Success)
			}
		}
	}
}

// queryGetByName validates the retrieval of a query using its name.
func queryGetByName(name string, q *query.Set, t *testing.T) {
	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving a query record's name")
		{
			ses := mongo.GetSession()
			defer ses.Close()

			qs, err := query.Get("Tests", ses, name)
			if err != nil {
				t.Errorf("\t\tShould have retrieved query record name[%s] successfully %s", name, tests.Failed)
			} else {
				t.Logf("\t\tShould have retrieved query record name[%s] successfully %s", name, tests.Success)

				if qs.Description != q.Description {
					t.Errorf("\t\tShould have matching description with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have matching description with retrieved record %s", tests.Success)
				}

				if len(qs.Params) != len(q.Params) {
					t.Errorf("\t\tShould have matching param size with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have matching param size with retrieved record %s", tests.Success)
				}

				if len(qs.Rules) != len(q.Rules) {
					t.Errorf("\t\tShould have matching rule size with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have matching rule size with retrieved record %s", tests.Success)
				}

				if qs.Name != q.Name {
					t.Errorf("\t\tShould have matching name with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have matching name with retrieved record %s", tests.Success)
				}

				if qs.Enabled != q.Enabled {
					t.Errorf("\t\tShould match run 'enabled' flag with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have match run 'enabled' flag with retrieved record %s", tests.Success)
				}
			}

		}
	}
}

// queryUpdate validates the updating of a query's content in the database.
func queryUpdate(q *query.Set, t *testing.T) {
	t.Log("Given the need to update a query record in the database")
	{
		t.Log("\tWhen giving an updated query")
		{

			//disable the run state of the query set
			q.Enabled = false

			ses := mongo.GetSession()
			defer ses.Close()

			err := query.Update("Tests", ses, q)

			if err != nil {
				t.Errorf("\t\tShould have updated query record name[%s] successfully %s", q.Name, tests.Failed)
			} else {
				t.Logf("\t\tShould have updated query record name[%s] successfully %s", q.Name, tests.Success)

				getSes := mongo.GetSession()
				defer getSes.Close()

				qs, err := query.Get("Tests", getSes, q.Name)
				if err != nil {
					t.Errorf("\t\tShould have retrieved query record name[%s] successfully %s", q.Name, tests.Failed)
				} else {
					t.Logf("\t\tShould have retrieved query record name[%s] successfully %s", q.Name, tests.Success)
				}

				if qs.Enabled != q.Enabled && qs.Enabled == false {
					t.Errorf("\t\tShould match run 'enabled' flag with retrieved record %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have match run 'enabled' flag with retrieved record %s", tests.Success)
				}

			}

		}
	}
}

// queryDelete validates the removal of a query from the database.
func queryDelete(q *query.Set, t *testing.T) {
	t.Log("Given the need to remove a query record in the database")
	{
		t.Log("\tWhen giving an query")
		{
			ses := mongo.GetSession()
			defer ses.Close()

			_, err := query.Delete("Tests", ses, q.Name)
			if err != nil {
				t.Errorf("\t\tShould have removed query record name[%s] successfully %s", q.Name, tests.Failed)
			} else {
				t.Logf("\t\tShould have removed query record name[%s] successfully %s", q.Name, tests.Success)
			}

		}
	}
}

// tearDown tears down the collection being used.
func tearDown(t *testing.T) {
	err := mongo.ExecuteDB("tearDown", mongo.GetSession(), collection, func(c *mgo.Collection) error {
		return c.DropCollection()
	})

	if err != nil {
		t.Errorf("Successfully dropped query collection [Error: %s] %s", err, tests.Failed)
	} else {
		t.Logf("Successfully dropped query collection %s", tests.Success)
	}
}
