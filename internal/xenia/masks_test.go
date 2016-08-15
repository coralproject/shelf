package xenia

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ardanlabs/kit/tests"
	"github.com/coralproject/shelf/internal/xenia/mask"
	"gopkg.in/mgo.v2/bson"
)

type (
	location struct {
		Type string `json:"type"`
	}

	condition struct {
		TempF float64 `json:"temp_f"`
		Temp  string  `json:"temperature_string"`
	}

	doc struct {
		StationID string    `json:"station_id"`
		Name      string    `json:"name"`
		Location  location  `json:"location"`
		Condition condition `json:"condition"`
	}
)

//==============================================================================

// TestMaskingDelete tests the masking functionality for deletes.
func TestMaskingDelete(t *testing.T) {
	t.Logf("Given the need to mask fields as deletes.")
	{
		masks := map[string]mask.Mask{
			"station_id": {"*", "station_id", mask.MaskRemove},
			"type":       {"*", "type", mask.MaskRemove},
			"wind_dir":   {"*", "wind_dir", mask.MaskRemove},
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
	masks := map[string]mask.Mask{
		"station_id": {"*", "station_id", mask.MaskAll},
		"type":       {"*", "type", mask.MaskAll},
		"temp_f":     {"*", "temp_f", mask.MaskAll},
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

// TestMaskingLeft tests the masking functionality for left.
func TestMaskingLeft(t *testing.T) {
	masks := map[string]mask.Mask{
		"station_id":         {"*", "station_id", mask.MaskLeft},
		"temperature_string": {"*", "temperature_string", mask.MaskLeft},
		"temp_f":             {"*", "temp_f", mask.MaskLeft},
	}

	t.Logf("Given the need to mask fields as left.")
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

			if fin.StationID != "****1" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "****1", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "****1", "station_id")
			}

			if fin.Condition.Temp != "**** F (15.2 C)" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "**** F (15.2 C)", "condition.temperature_string", fin.Condition.Temp)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "**** F (15.2 C)", "condition.temperature_string")
			}

			if fin.Condition.TempF != 0.00 {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "0.00", "condition.temp_f", fin.Condition.TempF)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "0.00", "condition.temp_f")
			}
		}
	}
}

// TestMaskingLeft8 tests the masking functionality for left8.
func TestMaskingLeft8(t *testing.T) {
	masks := map[string]mask.Mask{
		"station_id":         {"*", "station_id", mask.MaskLeft + "8"},
		"temperature_string": {"*", "temperature_string", mask.MaskLeft + "8"},
		"temp_f":             {"*", "temp_f", mask.MaskLeft + "8"},
	}

	t.Logf("Given the need to mask fields as left.")
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

			if fin.StationID != "*****" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "*****", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "*****", "station_id")
			}

			if fin.Condition.Temp != "********15.2 C)" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "********15.2 C)", "condition.temperature_string", fin.Condition.Temp)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "********15.2 C)", "condition.temperature_string")
			}

			if fin.Condition.TempF != 0.00 {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "0.00", "condition.temp_f", fin.Condition.TempF)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "0.00", "condition.temp_f")
			}
		}
	}
}

// TestMaskingRight tests the masking functionality for right.
func TestMaskingRight(t *testing.T) {
	masks := map[string]mask.Mask{
		"station_id":         {"*", "station_id", mask.MaskRight},
		"temperature_string": {"*", "temperature_string", mask.MaskRight},
		"temp_f":             {"*", "temp_f", mask.MaskRight},
	}

	t.Logf("Given the need to mask fields as left.")
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

			if fin.StationID != "4****" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "4****", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "4****", "station_id")
			}

			if fin.Condition.Temp != "59.4 F (15.****" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "59.4 F (15.****", "condition.temperature_string", fin.Condition.Temp)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "59.4 F (15.****", "condition.temperature_string")
			}

			if fin.Condition.TempF != 0.00 {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "0.00", "condition.temp_f", fin.Condition.TempF)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "0.00", "condition.temp_f")
			}
		}
	}
}

// TestMaskingRight8 tests the masking functionality for right8.
func TestMaskingRight8(t *testing.T) {
	masks := map[string]mask.Mask{
		"station_id":         {"*", "station_id", mask.MaskRight + "8"},
		"temperature_string": {"*", "temperature_string", mask.MaskRight + "8"},
		"temp_f":             {"*", "temp_f", mask.MaskRight + "8"},
	}

	t.Logf("Given the need to mask fields as left.")
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

			if fin.StationID != "*****" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "*****", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "*****", "station_id")
			}

			if fin.Condition.Temp != "59.4 F ********" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "59.4 F ********", "condition.temperature_string", fin.Condition.Temp)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "59.4 F ********", "condition.temperature_string")
			}

			if fin.Condition.TempF != 0.00 {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "0.00", "condition.temp_f", fin.Condition.TempF)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "0.00", "condition.temp_f")
			}
		}
	}
}

// TestMaskingEmail tests the masking functionality for email.
func TestMaskingEmail(t *testing.T) {
	masks := map[string]mask.Mask{
		"station_id": {"*", "station_id", mask.MaskEmail},
		"name":       {"*", "name", mask.MaskEmail},
		"temp_f":     {"*", "temp_f", mask.MaskEmail},
	}

	t.Logf("Given the need to mask fields as left.")
	{
		t.Logf("\tWhen using test fixture data.")
		{
			docs, err := fixtures()
			if err != nil {
				t.Fatalf("\t%s\tShould retrieve fixture documents.", tests.Failed)
			}

			docs[0]["station_id"] = "bill.smith@ardanlabs.com"
			docs[0]["name"] = "b@mydomain.com"

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

			if fin.StationID != "******@ardanlabs.com" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "******@ardanlabs.com", "station_id", fin.StationID)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "******@ardanlabs.com", "station_id")
			}

			if fin.Name != "******@mydomain.com" {
				t.Errorf("\t%s\tShould find %q in the document for field %q : %v", tests.Failed, "******@mydomain.com", "name", fin.Condition.Temp)
			} else {
				t.Logf("\t%s\tShould find %q in the document for field %q.", tests.Success, "******@mydomain.com", "name")
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
	path := os.Getenv("GOPATH") + "/src/github.com/coralproject/shelf/tstdata/test_data.json"

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
