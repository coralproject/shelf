package submission_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/ask/form/submission"
	"github.com/coralproject/shelf/internal/ask/form/submission/submissionfix"
)

// prefix is what we are looking to delete after the test.
const prefix = "FSTEST"

func TestMain(m *testing.M) {
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

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		log.Fatalf("Should be able to get a Mongo session : %v", err)
	}
	defer db.CloseMGO(tests.Context)

	// We need the database indexes setup before we can call anything, so do this
	// first. This is fairly important, so we want to fail the entire test suite
	// if we can't setup the indexes.
	if err := submission.EnsureIndexes(tests.Context, db); err != nil {
		log.Fatal("Can't ensure the database indexes")
	}

	os.Exit(m.Run())
}

func setup(t *testing.T, fixture string) ([]submission.Submission, *db.DB) {
	tests.ResetLog()

	subs, err := submissionfix.GetMany("submission.json")
	if err != nil {
		t.Fatalf("%s\tShould be able retrieve submission fixture : %s", tests.Failed, err)
	}
	t.Logf("%s\tShould be able retrieve submission fixture.", tests.Success)

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("Should be able to get a Mongo session : %v", err)
	}

	return subs, db
}

func teardown(t *testing.T, db *db.DB) {
	if err := submissionfix.Remove(tests.Context, db, prefix); err != nil {
		t.Fatalf("%s\tShould be able to remove the submissions : %v", tests.Failed, err)
	}
	t.Logf("%s\tShould be able to remove the submissions.", tests.Success)

	db.CloseMGO(tests.Context)
	tests.DisplayLog()
}

func Test_CreateDelete(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to create and delete submissions.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{
			//----------------------------------------------------------------------
			// Create the submission.

			if err := submission.Create(tests.Context, db, subs[0].FormID.Hex(), &subs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a submission.", tests.Success)

			//----------------------------------------------------------------------
			// Get the submission.

			sub, err := submission.Retrieve(tests.Context, db, subs[0].ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the submission by id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the submission by id.", tests.Success)

			//----------------------------------------------------------------------
			// Check that we got the submission we expected.

			if subs[0].ID.Hex() != sub.ID.Hex() {
				t.Fatalf("\t%s\tShould be able to get back the same submission.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to get back the same submission.", tests.Success)

			//----------------------------------------------------------------------
			// Delete the submission.

			if err := submission.Delete(tests.Context, db, subs[0].ID.Hex()); err != nil {
				t.Fatalf("\t%s\tShould be able to delete the submission : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete the submission.", tests.Success)

			//----------------------------------------------------------------------
			// Get the submission.

			_, err = submission.Retrieve(tests.Context, db, subs[0].ID.Hex())
			if err == nil {
				t.Fatalf("\t%s\tShould generate an error when getting a submission with the deleted id : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould generate an error when getting a submission with the deleted id.", tests.Success)
		}
	}
}

func matchSubmissions(t *testing.T, got, expected []submission.Submission) {
	if len(got) != len(expected) {
		t.Fatalf("\t%s\tShould be able to list all the submissions : Only found %d results, expected %d", tests.Failed, len(got), len(expected))
	}
	t.Logf("\t%s\tShould be able to list all the submissions", tests.Success)

	for _, fsub := range expected {
		var found bool

		for _, dsub := range got {
			if dsub.ID.Hex() == fsub.ID.Hex() {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("\t%s\tShould be able to find all submissions in resuts : Could not find submission %s in result set", tests.Failed, fsub.ID.Hex())
		}
	}
	t.Logf("\t%s\tShould be able to find all submissions in results", tests.Success)
}

func Test_Search(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to create and delete submissions.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{

			//----------------------------------------------------------------------
			// Create the submissions.

			for _, sub := range subs {
				if err := submission.Create(tests.Context, db, sub.FormID.Hex(), &sub); err != nil {
					t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
				}
			}
			t.Logf("\t%s\tShould be able to create submissions.", tests.Success)

			//----------------------------------------------------------------------
			// Search the submissions.

			results, err := submission.Search(tests.Context, db, subs[0].FormID.Hex(), len(subs), 0, submission.SearchOpts{})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list submissions", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the counts make sense.

			if results.Counts.TotalSearch != len(subs) {
				t.Fatalf("\t%s\tShould have the same total search count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSearch, len(subs))
			}
			t.Logf("\t%s\tShould have the same total search count", tests.Success)

			if results.Counts.TotalSubmissions != len(subs) {
				t.Fatalf("\t%s\tShould have the same total submission count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSubmissions, len(subs))
			}
			t.Logf("\t%s\tShould have the same total submission count", tests.Success)

			//----------------------------------------------------------------------
			// Verify that the docs exist inside the results.

			matchSubmissions(t, subs, results.Submissions)

			//----------------------------------------------------------------------
			// Search the submissions with a query.

			results, err = submission.Search(tests.Context, db, subs[0].FormID.Hex(), len(subs), 0, submission.SearchOpts{
				Query: "Option 1",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list submissions", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the counts make sense.

			if results.Counts.TotalSearch != len(subs) {
				t.Fatalf("\t%s\tShould have the same total search count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSearch, len(subs))
			}
			t.Logf("\t%s\tShould have the same total search count", tests.Success)

			if results.Counts.TotalSubmissions != len(subs) {
				t.Fatalf("\t%s\tShould have the same total submission count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSubmissions, len(subs))
			}
			t.Logf("\t%s\tShould have the same total submission count", tests.Success)

			//----------------------------------------------------------------------
			// Verify that the docs exist inside the results.

			matchSubmissions(t, subs, results.Submissions)

			//----------------------------------------------------------------------
			// Search the submissions with a filter.

			results, err = submission.Search(tests.Context, db, subs[0].FormID.Hex(), len(subs), 0, submission.SearchOpts{
				FilterBy: "flagged",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list submissions", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the counts make sense.

			if results.Counts.TotalSearch != len(subs) {
				t.Fatalf("\t%s\tShould have the same total search count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSearch, len(subs))
			}
			t.Logf("\t%s\tShould have the same total search count", tests.Success)

			if results.Counts.TotalSubmissions != len(subs) {
				t.Fatalf("\t%s\tShould have the same total submission count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSubmissions, len(subs))
			}
			t.Logf("\t%s\tShould have the same total submission count", tests.Success)

			//----------------------------------------------------------------------
			// Verify that the docs exist inside the results.

			matchSubmissions(t, subs, results.Submissions)

			//----------------------------------------------------------------------
			// Search the submissions with a negating filter.

			results, err = submission.Search(tests.Context, db, subs[0].FormID.Hex(), len(subs), 0, submission.SearchOpts{
				FilterBy: "-flagged",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list submissions", tests.Success)

			//----------------------------------------------------------------------
			// Verify that the results do not contain any results (as all results
			// contain the flagged flag).

			if results.Counts.TotalSearch != 0 {
				t.Fatalf("\t%s\tShould have the same total search count : Only found %d results, expected %d", tests.Failed, results.Counts.TotalSearch, 0)
			}
			t.Logf("\t%s\tShould have the same total search count", tests.Success)

			if len(results.Submissions) != 0 {
				t.Fatalf("\t%s\tShould not list any results unmatched : Found %d, expected 0", tests.Failed, len(results.Submissions))
			}
			t.Logf("\t%s\tShould not list any results unmatched", tests.Failed)
		}
	}
}

func Test_RetrieveMany(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to retrieve many submissions.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{

			//----------------------------------------------------------------------
			// Create the submissions.

			ids := make([]string, 0, len(subs))

			for _, sub := range subs {
				ids = append(ids, sub.ID.Hex())

				if err := submission.Create(tests.Context, db, sub.FormID.Hex(), &sub); err != nil {
					t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
				}
			}
			t.Logf("\t%s\tShould be able to create submissions.", tests.Success)

			rsubs, err := submission.RetrieveMany(tests.Context, db, ids)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to list submissions : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to list submissions", tests.Success)

			//----------------------------------------------------------------------
			// Verify that the docs exist inside the results.

			matchSubmissions(t, subs, rsubs)
		}
	}
}

func Test_Flags(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to add and remove flags.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{

			//----------------------------------------------------------------------
			// Create the submission.

			if err := submission.Create(tests.Context, db, subs[0].FormID.Hex(), &subs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a submission.", tests.Success)

			//----------------------------------------------------------------------
			// Ensure the submission has been added to the database.

			sub, err := submission.Retrieve(tests.Context, db, subs[0].ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			newFlag := time.Now().String()

			//----------------------------------------------------------------------
			// Ensure that the new flag is not already on the submission.

			// Ensure that the new flag is not in the current flags.
			for _, flag := range sub.Flags {
				if flag == newFlag {
					t.Fatalf("\t%s\tShould not have test flag already : Flag already exists", tests.Failed)
				}
			}
			t.Logf("\t%s\tShould not have test flag already.", tests.Success)

			//----------------------------------------------------------------------
			// Add the flag to the submission.

			nsub, err := submission.AddFlag(tests.Context, db, sub.ID.Hex(), newFlag)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to add the flag : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to add the flag", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the new flag was added to the submission returned.

			var found bool
			for _, flag := range nsub.Flags {
				if flag == newFlag {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("\t%s\tShould have test flag on returned object : Flag does not exist", tests.Failed)
			}
			t.Logf("\t%s\tShould have test flag on returned object.", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the new flag was added to the store's submission.

			rsub, err := submission.Retrieve(tests.Context, db, sub.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			found = false
			for _, flag := range rsub.Flags {
				if flag == newFlag {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("\t%s\tShould have test flag on database object : Flag does not exist", tests.Failed)
			}
			t.Logf("\t%s\tShould have test flag on database object.", tests.Success)

			//----------------------------------------------------------------------
			// Remove the new flag from the submission.

			nsub, err = submission.RemoveFlag(tests.Context, db, sub.ID.Hex(), newFlag)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to remove the flag : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to remove the flag", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the new flag was removed from the submission returned.

			found = false
			for _, flag := range nsub.Flags {
				if flag == newFlag {
					found = true
					break
				}
			}
			if found {
				t.Fatalf("\t%s\tShould not have test flag on returned object : Flag found", tests.Failed)
			}
			t.Logf("\t%s\tShould not have test flag on returned object.", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the new flag was removed from the submission in the store.

			rsub, err = submission.Retrieve(tests.Context, db, sub.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			found = false
			for _, flag := range rsub.Flags {
				if flag == newFlag {
					found = true
					break
				}
			}
			if found {
				t.Fatalf("\t%s\tShould not have test flag on database object : Flag found", tests.Failed)
			}
			t.Logf("\t%s\tShould not have test flag on database object.", tests.Success)
		}
	}
}

func Test_UpdateStatus(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to update the status of a submission.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{

			//----------------------------------------------------------------------
			// Create the submission.

			if err := submission.Create(tests.Context, db, subs[0].FormID.Hex(), &subs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a submission.", tests.Success)

			sub, err := submission.Retrieve(tests.Context, db, subs[0].ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			newStatus := time.Now().String()

			//----------------------------------------------------------------------
			// Update the submission's status.

			nsub, err := submission.UpdateStatus(tests.Context, db, sub.ID.Hex(), newStatus)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to update the submission status : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to update the submission status", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the status was updated on the returned submission.

			if nsub.Status != newStatus {
				t.Fatalf("\t%s\tShould be able to update the submission status : Expected %s, got %s", tests.Failed, newStatus, nsub.Status)
			}
			t.Logf("\t%s\tShould be able to update the submission status", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the status was updated on the stored submission.

			rsub, err := submission.Retrieve(tests.Context, db, sub.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			if rsub.Status != newStatus {
				t.Fatalf("\t%s\tShould be able to update the submission status in the store : Expected %s, got %s", tests.Failed, newStatus, nsub.Status)
			}
			t.Logf("\t%s\tShould be able to update the submission status in the store", tests.Success)
		}
	}
}

func Test_UpdateAnswer(t *testing.T) {
	subs, db := setup(t, "submission")
	defer teardown(t, db)

	t.Log("Given the need to update an answer on a submission.")
	{
		t.Log("\tWhen starting from an empty submissions collection")
		{

			//----------------------------------------------------------------------
			// Create the submisison.

			if err := submission.Create(tests.Context, db, subs[0].FormID.Hex(), &subs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to create a submission : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a submission.", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that we have created it.

			sub, err := submission.Retrieve(tests.Context, db, subs[0].ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			newAnswer := time.Now().String()

			//----------------------------------------------------------------------
			// Update the question's answer.

			nsub, err := submission.UpdateAnswer(tests.Context, db, sub.ID.Hex(), submission.AnswerInput{
				WidgetID: sub.Answers[0].WidgetID,
				Answer:   newAnswer,
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to update the submission answer : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to update the submission answer", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the answer was updated on the returned submission.

			if nsub.Answers[0].EditedAnswer != newAnswer {
				t.Fatalf("\t%s\tShould be able to update the submission answer : Expected %s, got %s", tests.Failed, newAnswer, nsub.Answers[0].EditedAnswer)
			}
			t.Logf("\t%s\tShould be able to update the submission answer", tests.Success)

			//----------------------------------------------------------------------
			// Ensure that the answer was updated on the stored submission.

			rsub, err := submission.Retrieve(tests.Context, db, sub.ID.Hex())
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the submission : %s", tests.Failed, err.Error())
			}
			t.Logf("\t%s\tShould be able to retrieve the submission", tests.Success)

			if rsub.Answers[0].EditedAnswer != newAnswer {
				t.Fatalf("\t%s\tShould be able to update the submission answer in the store : Expected %s, got %s", tests.Failed, newAnswer, nsub.Answers[0].EditedAnswer)
			}
			t.Logf("\t%s\tShould be able to update the submission answer in the store", tests.Success)
		}
	}
}
