package exec_test

import (
	"testing"
	"time"

	"github.com/coralproject/xenia/pkg/exec"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"
	"gopkg.in/mgo.v2/bson"
)

func init() {
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
}

//==============================================================================

// TestPreProcessing tests the ability to preprocess json documents.
func TestPreProcessing(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	time1, _ := time.Parse("2006-01-02T15:04:05.999Z", "2013-01-16T00:00:00.000Z")
	time2, _ := time.Parse("2006-01-02", "2013-01-16")

	commands := []struct {
		doc   map[string]interface{}
		vars  map[string]string
		after map[string]interface{}
	}{
		{
			map[string]interface{}{"field_name": "#string:name"},
			map[string]string{"name": "bill"},
			map[string]interface{}{"field_name": "bill"},
		},
		{
			map[string]interface{}{"field_name": "#number:value"},
			map[string]string{"value": "10"},
			map[string]interface{}{"field_name": 10},
		},
		{
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16T00:00:00.000Z"},
			map[string]interface{}{"field_name": time1},
		},
		{
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16"},
			map[string]interface{}{"field_name": time2},
		},
		{
			map[string]interface{}{"field_name": "#date:2013-01-16T00:00:00.000Z"},
			map[string]string{},
			map[string]interface{}{"field_name": time1},
		},
		{
			map[string]interface{}{"field_name": "#objid:value"},
			map[string]string{"value": "5660bc6e16908cae692e0593"},
			map[string]interface{}{"field_name": bson.ObjectIdHex("5660bc6e16908cae692e0593")},
		},
	}

	t.Logf("Given the need to preprocess commands.")
	{
		for _, cmd := range commands {
			t.Logf("\tWhen using %+v with %+v", cmd.doc, cmd.vars)
			{
				exec.ProcessVariables("", cmd.doc, cmd.vars, nil)

				if eq := compareBson(cmd.doc, cmd.after); !eq {
					t.Log(cmd.doc)
					t.Log(cmd.after)
					t.Errorf("\t%s\tShould get back the expected document.", tests.Failed)
					continue
				}
				t.Logf("\t%s\tShould get back the expected document.", tests.Success)
			}
		}
	}
}

// compareBson compares two bson maps for equivalence.
func compareBson(m1 bson.M, m2 bson.M) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	for k, v := range m2 {
		if m1[k] != v {
			return false
		}
	}

	return true
}
