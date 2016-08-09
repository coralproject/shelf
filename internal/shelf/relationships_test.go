package shelf

import (
	"encoding/json"
	"testing"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/xenia/internal/shelf/sfix"
)

// TestAddRemoveRelationship tests if we can add/remove a relationship to/from the db.
func TestAddRemoveRelationship(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelsAndViews(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships and views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships and views.", tests.Success)
	}()

	t.Log("Given the need to save a new relationship into the database.")
	{
		t.Log("\tWhen starting from the relsandviews.json test fixture")
		{
			raw, err := sfix.LoadRelAndViewData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship and view fixture : %s", tests.Failed, err)
			}
			var rm RelsAndViews
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship and view fixture : %s", tests.Failed, err)
			}
			if err := NewRelsAndViews(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create relationships and views : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create relationships and views.", tests.Success)
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
			rm, err = GetRelsAndViews(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve a relationships : %s", tests.Failed, err)
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
			rm, err = GetRelsAndViews(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve relationships : %s", tests.Failed, err)
			}
			relIDs = []string{}
			for _, rel := range rm.Relationships {
				relIDs = append(relIDs, rel.ID)
			}
			if stringContains(relIDs, newRelID) {
				t.Errorf("\t%s\tShould be able to get back the relationships without the removed ID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the relationships without the removed ID.", tests.Success)
			}

		}
	}
}

// TestAddRelationshipFail tests if we can properly throw an error for an illegal relationship.
func TestAddRelationshipFail(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := ClearRelsAndViews(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships and views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships and views.", tests.Success)
	}()

	t.Log("Given the need to save a new relationship into the database.")
	{
		t.Log("\tWhen starting from the relsandviews.json test fixture")
		{
			raw, err := sfix.LoadRelAndViewData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship and view fixture : %s", tests.Failed, err)
			}
			var rm RelsAndViews
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship and view fixture : %s", tests.Failed, err)
			}
			if err := NewRelsAndViews(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create relationships and views : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create relationships and views.", tests.Success)
			newRel := Relationship{
				SubjectTypes: []string{"coral_user"},
				Predicate:    "authored",
				ObjectTypes:  []string{"coral_comment"},
			}
			_, err = AddRelationship(tests.Context, db, newRel)
			if err == nil {
				t.Fatalf("\t%s\tShould be able to throw error on preexisting predicate : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to throw error on preexisting predicate.", tests.Success)
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
		if err := ClearRelsAndViews(tests.Context, db); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the relationships and views : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the relationships and views.", tests.Success)
	}()

	t.Log("Given the need to update a relationship in the database.")
	{
		t.Log("\tWhen starting from the relsandviews.json test fixture")
		{
			raw, err := sfix.LoadRelAndViewData()
			if err != nil {
				t.Fatalf("\t%s\tShould be able retrieve relationship and view fixture : %s", tests.Failed, err)
			}
			var rm RelsAndViews
			if err := json.Unmarshal(raw, &rm); err != nil {
				t.Fatalf("\t%s\tShould be able unmarshal relationship and view fixture : %s", tests.Failed, err)
			}
			if err := NewRelsAndViews(tests.Context, db, rm); err != nil {
				t.Fatalf("\t%s\tShould be able to create relationships and views : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create relationships and views.", tests.Success)
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
			rm, err = GetRelsAndViews(tests.Context, db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve relationships : %s", tests.Failed, err)
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
