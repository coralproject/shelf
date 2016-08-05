package shelf

import (
	"encoding/json"
	"testing"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/sfix"
)

// TestAddRemoveUpdateRelationships tests if we can add/remove a relationship to/from the db.
func TestAddRemoveRelationship(t *testing.T) {
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

	t.Log("Given the need to save a new relationship into the database.")
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
			newRel := Relationship{
				SubjectTypes: []string{"coral_user"},
				Predicate:    "tested",
				ObjectTypes:  []string{"coral_comment"},
			}
			newRelID, err := AddRelationship(tests.Context, db, newRel)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to add a new relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to add a new relationship.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			var relIDs []string
			for _, rel := range rm.Relationships {
				relIDs = append(relIDs, rel.ID)
			}
			if !stringContains(relIDs, newRelID) {
				t.Errorf("\t%s\tShould be able to get back the same relationship ID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same relationship ID.", tests.Success)
			}
			if err = RemoveRelationship(tests.Context, db, newRelID); err != nil {
				t.Fatalf("\t%s\tShould be able to remove a relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove a relationship.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			relIDs = []string{}
			for _, rel := range rm.Relationships {
				relIDs = append(relIDs, rel.ID)
			}
			if stringContains(relIDs, newRelID) {
				t.Errorf("\t%s\tShould be able to get back the rel. manager without the removed ID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the rel. manager without the removed ID.", tests.Success)
			}

		}
	}
}

// TestUpdateRelationship tests if we can update a relationship in the db.
func TestUpdateRelationship(t *testing.T) {
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

	t.Log("Given the need to save a new relationship into the database.")
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
			updatedRel := Relationship{
				ID:           "4c25fa45-d8ff-47ec-8379-aed63e731604",
				SubjectTypes: []string{"coral_comment"},
				Predicate:    "on_test",
				ObjectTypes:  []string{"coral_asset"},
			}
			if err := UpdateRelationship(tests.Context, db, updatedRel); err != nil {
				t.Fatalf("\t%s\tShould be able to update a relationship : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a relationship.", tests.Success)
			rm, err = GetRelManager(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationship manager : %s", tests.Failed, err)
			}
			var predicates []string
			for _, rel := range rm.Relationships {
				predicates = append(predicates, rel.Predicate)
			}
			if !stringContains(predicates, "on_test") {
				t.Errorf("\t%s\tShould be able to get back the updated relationship.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the updated relationship.", tests.Success)
			}
		}
	}
}
