package shelf

import (
	"encoding/json"
	"testing"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/sfix"
)

// TestAddRemoveView tests if we can add/remove a view to/from the db.
func TestAddRemoveView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelManager(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationship manager : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationship manager.", tests.Success)
	}()

	t.Log("Given the need to save a new view into the database.")
	{
		t.Log("\tWhen starting from the relmanager.json test fixture")
		{
			raw, err := sfix.LoadRelManagerData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship manager fixture : %s", tests.Failed, err)
			}
			var rm RelManager
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship manager fixture : %s", tests.Failed, err)
			}
			if err := NewRelManager(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create a relationship manager : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a relationship manager.", tests.Success)
			newView := View{
				Name:      "test_view",
				StartType: "coral_user",
				Path: []PathSegment{PathSegment{
					Level:          1,
					Direction:      "out",
					RelationshipID: "32420143-376e-482a-b4d2-709ea9e31e8e",
				}},
			}
			newViewID, err := AddView(tests.Context, db, newView)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to add a new view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to add a new view.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			var viewIDs []string
			for _, view := range rm.Views {
				viewIDs = append(viewIDs, view.ID)
			}
			if !stringContains(viewIDs, newViewID) {
				t.Errorf("\t%s\tShould be able to get back the same view ID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same view ID.", tests.Success)
			}
			if err = RemoveView(tests.Context, db, newViewID); err != nil {
				t.Fatalf("\t%s\tShould be able to remove a view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove a view.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			viewIDs = []string{}
			for _, view := range rm.Views {
				viewIDs = append(viewIDs, view.ID)
			}
			if stringContains(viewIDs, newViewID) {
				t.Errorf("\t%s\tShould be able to get back the rel. manager without the removed ID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the rel. manager without the removed ID.", tests.Success)
			}

		}
	}
}

// TestUpdateView tests if we can update a view in the db.
func TestUpdateView(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelManager(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationship manager : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationship manager.", tests.Success)
	}()

	t.Log("Given the need to update a view in the database.")
	{
		t.Log("\tWhen starting from the relmanager.json test fixture")
		{
			raw, err := sfix.LoadRelManagerData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship manager fixture : %s", tests.Failed, err)
			}
			var rm RelManager
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship manager fixture : %s", tests.Failed, err)
			}
			if err := NewRelManager(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create a relationship manager : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a relationship manager.", tests.Success)
			updatedView := View{
				ID:        "28d58623-073e-4efe-b8bb-62cfce178752",
				Name:      "updated name",
				StartType: "coral_user",
				Path: []PathSegment{PathSegment{
					Level:          1,
					Direction:      "out",
					RelationshipID: "c9f8df2e-c301-4a90-8abc-02cf201f9cf0",
				}},
			}
			if err := UpdateView(tests.Context, db, updatedView); err != nil {
				t.Fatalf("\t%s\tShould be able to update a view : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a view.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			var names []string
			for _, view := range rm.Views {
				names = append(names, view.Name)
			}
			if !stringContains(names, "updated name") {
				t.Errorf("\t%s\tShould be able to get back the updated view", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the updated view.", tests.Success)
			}
		}
	}
}
