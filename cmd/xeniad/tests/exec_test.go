// Package tests implements users tests for the API layer.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/query/qfix"
)

// TestExec tests the execution of a specific query.
func TestExec(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a specific query.")
	{
		url := "/1.0/exec/" + qPrefix + "_basic?station_id=42021"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`)

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestExecCustom tests the execution of a custom query.
func TestExecCustom(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a custom query.")
	{
		qs, err := qfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		url := "/1.0/exec"
		r := tests.NewRequest("POST", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := tests.IndentJSON(`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`)

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestExecJSONP tests the execution of a specific query using JSONP.
func TestExecJSONP(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a specific query with JSONP output.")
	{
		url := "/1.0/exec/" + qPrefix + "_basic?station_id=42021&callback=handle_data"
		r := tests.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := w.Body.String()

			if !strings.HasPrefix(recv, "handle_data(") {
				t.Fatalf("\t%s\tShould get the expected prefix.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected prefix.", tests.Success)

			recv = strings.TrimPrefix(recv, "handle_data(")

			if !strings.HasSuffix(recv, ")") {
				t.Fatalf("\t%s\tShould get the expected suffix.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected suffix.", tests.Success)

			recv = tests.IndentJSON(strings.TrimSuffix(recv, ")"))

			resp := tests.IndentJSON(`{"results":[{"Name":"Basic","Docs":[{"name":"C14 - Pasco County Buoy, FL"}]}]}`)

			if resp != recv {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// TestExecExplain tests the execution of a custom query with explain.
func TestExecExplain(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to execute a custom query with explain.")
	{
		qs, err := qfix.Get("basic.json")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve the fixture.", tests.Success)

		qs.Explain = true

		qsStrData, err := json.Marshal(&qs)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to marshal the fixture : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to marshal the fixture.", tests.Success)

		url := "/1.0/exec"
		r := tests.NewRequest("POST", url, bytes.NewBuffer(qsStrData))
		w := httptest.NewRecorder()

		a.ServeHTTP(w, r)

		t.Logf("\tWhen calling url : %s", url)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould be able to retrieve the query : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould be able to retrieve the query.", tests.Success)

			recv := tests.IndentJSON(w.Body.String())
			resp := `queryPlanner`

			if !strings.Contains(recv, resp) {
				t.Log(resp)
				t.Log(recv)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}
