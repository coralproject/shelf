package query_test

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/shelf/pkg/log"
	"github.com/coralproject/shelf/pkg/srv/mongo"
	"github.com/coralproject/shelf/pkg/srv/query"
	"github.com/coralproject/shelf/pkg/tests"
)

// collection used for testing the query CRUD API
var collection = "queries"
var record = "spending_advice"
var testRecord = "test_spending_advice"

func init() {

	// Initialize the test environment.
	tests.Init()
}

func removeSession(session *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": record}
		log.Dev("Test", "removeSession", "db.queries.remove(%s)", mongo.Query(q))
		return c.Remove(q)
	}

	err := mongo.ExecuteDB("Tests", session, collection, f)
	if err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func removeTestSession(session *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": testRecord}
		log.Dev("Test", "removeSession", "db.queries.remove(%s)", mongo.Query(q))
		return c.Remove(q)
	}

	err := mongo.ExecuteDB("Tests", session, collection, f)
	if err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func recordThoseNotExists(session *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": record}
		log.Dev("Test", "recordThoseNotExists", "db.queries.remove(%s)", mongo.Query(q))
		n, err := c.Find(q).Count()
		if err != nil {
			return err
		}

		if n != 0 {
			return errors.New("Record Found")
		}

		return nil
	}

	err := mongo.ExecuteDB("Tests", session, collection, f)
	if err != mgo.ErrNotFound {
		return err
	}

	return nil
}

func recordThoseExists(session *mgo.Session) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"name": record}
		log.Dev("Test", "recordThoseExists", "db.queries.find(%s).count()", mongo.Query(q))

		n, err := c.Find(q).Count()
		if err != nil {
			return err
		}

		if n == 0 {
			return errors.New("No Record Found")
		}

		return nil
	}

	err := mongo.ExecuteDB("Tests", session, collection, f)
	if err != mgo.ErrNotFound {
		return err
	}

	return nil
}

// getQuery retrieves a query record from the filesystem.
func getQuery() (query.Set, error) {
	qFile := "./fixtures/spending_advice.json"

	var qs query.Set

	file, err := os.Open(qFile)
	if err != nil {
		return qs, err
	}

	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		return qs, err
	}

	return qs, nil
}

// TestCreateQuery tests if we can create a query record in the db.
func TestCreateQuery(t *testing.T) {
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

	t.Log("Given the need to save a query into the database")
	{
		t.Log("\tWhen giving a query object to save")
		{

			err := query.CreateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %s", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			if err := recordThoseExists(ses); err != nil {
				t.Errorf("%s\t\tShould have found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have found a query record ", tests.Success)
			}

			if err := recordThoseNotExists(ses); err == nil {
				t.Errorf("%s\t\tShould have found a query record : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have found a query record ", tests.Success)
			}

		}
	}
}

// TestGetSetNames validates retrieval of query.Set record names.
func TestGetSetNames(t *testing.T) {
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

		if err := removeTestSession(ses); err != nil {
			t.Errorf("%s\t\tShould have removed query test record from dbs : %v", tests.Failed, err)
		} else {
			t.Logf("%s\t\tShould have removed query test record from db ", tests.Success)
		}
	}()

	t.Log("Given the need to retrieve a query from the database")
	{
		t.Log("\tWhen giving the need to retrieve all names")
		{

			err := query.CreateSet("Tests", ses, q)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record : %s", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record ", tests.Success)
			}

			testQuery := q
			testQuery.Name = testRecord
			err = query.CreateSet("Tests", ses, testQuery)
			if err != nil {
				t.Errorf("%s\t\tShould have added new query record with name %s : %s", testQuery.Name, tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have added new query record with name %s : %s", testQuery.Name, tests.Success, err)
			}

			ses := mongo.GetSession()
			defer ses.Close()

			names, err := query.GetSetNames("Test", ses)
			if err != nil {
				t.Errorf("%s\t\tShould have retrieved query record names successfully : %v", tests.Failed, err)
			} else {
				t.Logf("%s\t\tShould have retrieved query record names successfully", tests.Success)

				if len(names) == 0 {
					t.Errorf("%s\t\tShould have atleast one query record name: %s", tests.Failed, names)
				} else {
					t.Logf("%s\t\tShould have atleast one query record name: %s", tests.Success, names)
				}

				expectedName := "spending_advice"
				if names[0] != expectedName {
					t.Errorf("%s\t\tShould have first name equal %q", tests.Failed, expectedName)
				} else {
					t.Logf("%s\t\tShould have first name equal %q", tests.Success, expectedName)
				}
			}

		}
	}
}

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
