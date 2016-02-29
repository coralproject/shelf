package exec

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/coralproject/xenia/pkg/mask"

	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2/bson"
)

//==============================================================================

// TestMaskingDelete tests the masking functionality for deletes.
func TestMaskingDelete(t *testing.T) {
	t.Logf("Given the need to mask fields as deletes.")
	{
		masks := map[string]mask.Mask{
			"station_id": {"test_xenia_data", "station_id", mask.MaskRemove},
			"type":       {"test_xenia_data", "type", mask.MaskRemove},
			"wind_dir":   {"test_xenia_data", "wind_dir", mask.MaskRemove},
		}

		docs, err := fixtures()
		if err != nil {
			t.Fatalf("\t%s\tShould retrieve fixture documents.", tests.Failed)
		}

		t.Logf("\tWhen using test data fixtures")
		{
			if _, exists := docs[0]["station_id"]; !exists {
				t.Fatalf("\t%s\tShould find %q in the document.", tests.Failed, "station_id")
			}
			t.Logf("\t%s\tShould find %q in the document.", tests.Success, "station_id")

			if _, exists := docs[0]["location"].(map[string]interface{})["type"]; !exists {
				t.Fatalf("\t%s\tShould find %q in the document.", tests.Failed, "type")
			}
			t.Logf("\t%s\tShould find %q in the document.", tests.Success, "type")

			if _, exists := docs[0]["condition"].(map[string]interface{})["wind_dir"]; !exists {
				t.Fatalf("\t%s\tShould find %q in the document.", tests.Failed, "wind_dir")
			}
			t.Logf("\t%s\tShould find %q in the document.", tests.Success, "wind_dir")

			if err := matchMaskField(tests.Context, masks, docs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to mask fields : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to mask fields.", tests.Success)

			if _, exists := docs[0]["station_id"]; exists {
				t.Errorf("\t%s\tShould not find %q in the document.", tests.Failed, "station_id")
			} else {
				t.Logf("\t%s\tShould not find %q in the document.", tests.Success, "station_id")
			}

			if _, exists := docs[0]["location"].(map[string]interface{})["type"]; exists {
				t.Errorf("\t%s\tShould not find %q in the document.", tests.Failed, "type")
			} else {
				t.Logf("\t%s\tShould not find %q in the document.", tests.Success, "type")
			}

			if _, exists := docs[0]["condition"].(map[string]interface{})["wind_dir"]; exists {
				t.Errorf("\t%s\tShould not find %q in the document.", tests.Failed, "wind_dir")
			} else {
				t.Logf("\t%s\tShould not find %q in the document.", tests.Success, "wind_dir")
			}
		}
	}
}

// TestMaskingAll tests the masking functionality for all.
func TestMaskingAll(t *testing.T) {
	type location struct {
		Type string `json:"type"`
	}

	type condition struct {
		TempF float64 `json:"temp_f"`
	}

	type doc struct {
		StationID string    `json:"station_id"`
		Location  location  `json:"location"`
		Condition condition `json:"condition"`
	}

	masks := map[string]mask.Mask{
		"station_id": {"test_xenia_data", "station_id", mask.MaskAll},
		"type":       {"test_xenia_data", "type", mask.MaskAll},
		"temp_f":     {"test_xenia_data", "temp_f", mask.MaskAll},
	}

	t.Logf("Given the need to mask fields as all.")
	{
		t.Logf("\tWhen using test fixture data.")
		{
			docs, err := fixtures()
			if err != nil {
				t.Fatalf("\t%s\tShould retrieve fixture documents.", tests.Failed)
			}

			if err := matchMaskField(tests.Context, masks, docs[0]); err != nil {
				t.Fatalf("\t%s\tShould be able to mask fields : %s", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to mask fields.", tests.Success)

			data, err := json.Marshal(docs[0])
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal document.", tests.Failed)
			}

			var fin doc
			if err := json.Unmarshal(data, &fin); err != nil {
				t.Fatalf("\t%s\tShould unmarshal document.", tests.Failed)
			}

			if fin.StationID != "******" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "******", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "******", "station_id")
			}

			if fin.Location.Type != "******" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "******", "location.type", fin.Location.Type)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "******", "location.type")
			}

			if fin.Condition.TempF != 0.00 {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "0.00", "condition.temp_f", fin.Condition.TempF)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "0.00", "condition.temp_f")
			}
		}
	}
}

//==============================================================================

// fixtures reads the test data fixture for documents to use for this testing.
func fixtures() ([]bson.M, error) {
	path := os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/tstdata/test_data.json"

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var docs []bson.M
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}
