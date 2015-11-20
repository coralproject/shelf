package db

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/shelf/pkg/tests"
)

// TestQueryPI validates the operations of the query database and file loading API.
func TestQueryPI(t *testing.T) {
	// Initialize the test environment.
	tests.Init()

	tests.ResetLog()
	defer tests.DisplayLog()

	q := NewQuery("")

	queryLoadFile(q, t)
	queryCreate(q, t)
	queryGetByName(q.Name, q, t)
	queryGetByID(q.ID.Hex(), q, t)
	queryUpdate(q, t)
	queryDelete(q, t)
	tearDown(t)
}

// queryLoadFile validates the loading of a query from a file.
func queryLoadFile(q *Query, t *testing.T) {
	t.Log("Given the need to load a query from a file")
	{
		t.Log("\tWhen giving a query file path")
		{

			qFile := "./fixture/user_advice.json"
			err := q.LoadFile(qFile)

			if err != nil {
				t.Fatalf("\t\tShould load File[%q] without error %s", qFile, tests.Failed)
			} else {
				t.Logf("\t\tShould load File[%q] without error %s", qFile, tests.Success)

				qName := "user_advice"
				if q.Name != qName {
					t.Errorf("\t\tShould have name[%q] for loaded query %s", qName, tests.Failed)
				} else {
					t.Logf("\t\tShould have name[%q] for loaded query %s", qName, tests.Success)
				}

				qTestCollection := "user_transactions"
				if q.Test.Collection != qTestCollection {
					t.Errorf("\t\tShould have name[%q] for loaded query.Test expression collection %s", qTestCollection, tests.Failed)
				} else {
					t.Logf("\t\tShould have name[%q] for loaded query.Test expression collection %s", qTestCollection, tests.Success)
				}

				if len(q.Test.Queries) != 3 {
					t.Errorf("\t\tShould have a length of 3 for the query.Test conditions %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have a length of 3 for the query.Test conditions %s", tests.Success)
				}
			}
		}
	}
}

// queryCreate validates the creation of a query in the databae.
func queryCreate(q *Query, t *testing.T) {
	t.Log("Given the need to save a query into the database")
	{
		t.Log("\tWhen giving a query object to save")
		{
			err := Create(q)
			if err != nil {
				t.Errorf("\t\tShould save query into the database %s", tests.Failed)
			} else {
				t.Logf("\t\tShould save query into the database %s", tests.Success)
			}
		}

		t.Log("\tWhen query record has been saved and we need to validate")
		{
			query, err := GetByName(q.Name)
			if err != nil {
				t.Errorf("\t\tShould get query from the database with the Name[%q] %s", q.Name, tests.Failed)
			} else {
				t.Logf("\t\tShould save query from the database with the Name[%q] %s", q.Name, tests.Success)

				if err := q.Compare(query); err != nil {
					t.Errorf("\t\tShould have query matching the query from the database %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have query matching the query from the database %s", tests.Success)
				}
			}

		}
	}
}

// queryGetByName validates the retrieval of a query using its name.
func queryGetByName(name string, q *Query, t *testing.T) {
	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving a query record's name")
		{
			query, err := GetByName(name)
			if err != nil {
				t.Errorf("\t\tShould get query from the database with the Name[%q] %s", name, tests.Failed)
			} else {
				t.Logf("\t\tShould save query from the database with the Name[%q] %s", name, tests.Success)

				if err := q.Compare(query); err != nil {
					t.Errorf("\t\tShould have query matching the query from the database %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have query matching the query from the database %s", tests.Success)
				}
			}

		}
	}
}

// queryGetByID validates the retrieval of a query using its id.
func queryGetByID(id string, q *Query, t *testing.T) {
	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving a query record's ID")
		{
			query, err := GetByID(id)
			if err != nil {
				t.Errorf("\t\tShould get query from the database with the ID[%q] %s", id, tests.Failed)
			} else {
				t.Logf("\t\tShould save query from the database with the ID[%q] %s", id, tests.Success)

				if err := q.Compare(query); err != nil {
					t.Errorf("\t\tShould have query matching the query from the database %s", tests.Failed)
				} else {
					t.Logf("\t\tShould have query matching the query from the database %s", tests.Success)
				}
			}

		}
	}
}

// queryUpdate validates the updating of a query's content in the database.
func queryUpdate(q *Query, t *testing.T) {
	t.Log("Given the need to update a query record in the database")
	{
		t.Log("\tWhen giving an updated query")
		{

			queryUpdate := []string{
				"{ \"$match\" : { \"user_id\" : \"#userId#\", \"category\" : \"gas\" }}",
				"{ \"$group\" : { \"_id\" : { \"category\" : \"$category\" }, \"amount\" : { \"$sum\" : \"$amount\" }}}",
				"{ \"$match\" : { \"amount\" : { \"$gt\" : 10.00}}}",
			}

			q.Test.Queries = queryUpdate

			if err := Update(q); err != nil {
				t.Errorf("\t\tShould update query record in db successfully %s", tests.Failed)
			} else {
				t.Logf("\t\tShould update query record in db successfully %s", tests.Success)

				updatedQuery, err := GetByName(q.Name)
				if err != nil {
					t.Errorf("\t\tShould retrieve updated record from the db successfully %s", tests.Failed)
				} else {
					t.Logf("\t\tShould retrieve updated record from the db successfully %s", tests.Success)

					if err := q.Compare(updatedQuery); err != nil {
						t.Errorf("\t\tShould have query matching the query from the database %s", tests.Failed)
					} else {
						t.Logf("\t\tShould have query matching the query from the database %s", tests.Success)
					}
				}
			}

		}
	}
}

// queryDelete validates the removal of a query from the database.
func queryDelete(q *Query, t *testing.T) {
	t.Log("Given the need to remove a query record in the database")
	{
		t.Log("\tWhen giving an query")
		{
			if err := Delete(q); err != nil {
				t.Errorf("\t\tShould update query record in db successfully %s", tests.Failed)
			} else {
				t.Logf("\t\tShould update query record in db successfully %s", tests.Success)

				_, err := GetByName(q.Name)
				if err == nil {
					t.Errorf("\t\tShould not successfully retrieve deleted record from the db %s", tests.Failed)
				} else {
					t.Logf("\t\tShould not successfully retrieve deleted record from the db %s", tests.Success)
				}
			}

		}
	}
}

// tearDown tears down the collection being used.
func tearDown(t *testing.T) {
	err := mongo.ExecuteDB("tearDown", mongo.GetSession(), QueryCollection, func(c *mgo.Collection) error {
		return c.DropCollection()
	})

	if err != nil {
		t.Errorf("Successfully dropped query collection %s", tests.Failed)
	} else {
		t.Logf("Successfully dropped query collection %s", tests.Success)
	}
}
