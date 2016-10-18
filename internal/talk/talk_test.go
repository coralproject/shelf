package talk_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/talk"
)

func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Initialize the configuration and logging systems. Plus anything
	// else the web app layer needs.
	tests.Init("TALK")

	return m.Run()
}

// setup initializes for each indivdual test.
func setup(t *testing.T) *httptest.Server {
	tests.ResetLog()

	// Initialization of stub server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var itm item.Item

		// check that the row is what we want it to be
		switch r.RequestURI {
		case "/v1/item/33":
			itm = item.Item{ID: "33", Type: "comment", Data: map[string]interface{}{"body": "Something."}}
			itm.Data["flagged_by"] = []string{"11"}
			b, _ := json.Marshal([]item.Item{itm})
			w.Write(b)
		case "/v1/item/notfoundtarget":
			b, _ := json.Marshal([]item.Item{})
			w.Write(b)
		default:
			err = errors.New("Bad Request")
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err == nil {
			w.WriteHeader(http.StatusOK)
		}
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintln(w, err)

	}))

	return server
}

func teardown(t *testing.T, server *httptest.Server) {
	// teardown deinitializes for each indivdual test.
	tests.DisplayLog()
}

//==============================================================================

// TestAddAction tests that the action was added correctly to the target item.
func TestAddAction(t *testing.T) {
	server := setup(t)
	defer teardown(t, server)

	usr := item.Item{ID: "11", Type: "user", Data: map[string]interface{}{"name": "Maria"}}
	itm := item.Item{ID: "33", Type: "comment", Data: map[string]interface{}{"body": "Something."}}
	itm.Data["flagged_by"] = []interface{}{"11"}

	// Build our table of the different test sets.
	actionSets := []struct {
		url            string
		user           item.Item
		action         string
		targetID       string
		expectedTarget item.Item
		expectedError  error
	}{
		{url: server.URL, user: item.Item{}, action: "liked_by", targetID: "notfoundtarget", expectedTarget: item.Item{}, expectedError: talk.ErrItemNotFound},
		{url: server.URL, user: item.Item{}, action: "wrong action", targetID: "33", expectedTarget: item.Item{}, expectedError: talk.ErrActionNotAllowed},
		{url: server.URL, user: usr, action: "flagged_by", targetID: "33", expectedTarget: itm, expectedError: nil},
	}

	// Iterate over all the different test sets.
	for _, actionSet := range actionSets {

		t.Logf("Given the need to add action %s to target %s.", actionSet.action, actionSet.targetID)
		{

			a, err := talk.AddAction(actionSet.url, actionSet.user, actionSet.action, actionSet.targetID)

			if err != actionSet.expectedError {
				t.Errorf("\t%s\tShould be able to return error %v but got : %v.", tests.Failed, actionSet.expectedError, err)
				return
			}
			t.Logf("\t%s\tShould be able to return error: %v", tests.Success, actionSet.expectedError)

			for f, p := range a.Data {
				if !reflect.DeepEqual(p, actionSet.expectedTarget.Data[f]) {
					t.Errorf("\t%s\tShould be able to return target '%v' but got : '%v'", tests.Failed, actionSet.expectedTarget.Data[f], a.Data[f])
					return
				}
				t.Logf("\t%s\tShould be able to return target: %v", tests.Success, actionSet.expectedTarget)
			}
		}
	}
}
