package exec_test

import (
	"strconv"
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
		time  bool
		doc   map[string]interface{}
		vars  map[string]string
		after map[string]interface{}
	}{
		{
			false,
			map[string]interface{}{"{field_name}": "bill"},
			map[string]string{"field_name": "name"},
			map[string]interface{}{"name": "bill"},
		},
		{
			false,
			map[string]interface{}{"statstics.comments.{dimension}.{commentStatus}.{value}": "bill"},
			map[string]string{"dimension": "dim", "commentStatus": "cstat", "value": "v"},
			map[string]interface{}{"statstics.comments.dim.cstat.v": "bill"},
		},
		{
			false,
			map[string]interface{}{"{dimension}": map[string]interface{}{"{commentStatus}": map[string]interface{}{"{value}": "bill"}}},
			map[string]string{"dimension": "dim", "commentStatus": "cstat", "value": "v"},
			map[string]interface{}{"dim": map[string]interface{}{"cstat": map[string]interface{}{"v": "bill"}}},
		},
		{
			false,
			map[string]interface{}{"field_name": "#string:name"},
			map[string]string{"name": "bill"},
			map[string]interface{}{"field_name": "bill"},
		},
		{
			false,
			map[string]interface{}{"field_name": "#number:value"},
			map[string]string{"value": "10"},
			map[string]interface{}{"field_name": 10},
		},
		{
			false,
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16T00:00:00.000Z"},
			map[string]interface{}{"field_name": time1},
		},
		{
			false,
			map[string]interface{}{"field_name": "#date:value"},
			map[string]string{"value": "2013-01-16"},
			map[string]interface{}{"field_name": time2},
		},
		{
			false,
			map[string]interface{}{"field_name": "#date:2013-01-16T00:00:00.000Z"},
			map[string]string{},
			map[string]interface{}{"field_name": time1},
		},
		{
			false,
			map[string]interface{}{"field_name": "#objid:value"},
			map[string]string{"value": "5660bc6e16908cae692e0593"},
			map[string]interface{}{"field_name": bson.ObjectIdHex("5660bc6e16908cae692e0593")},
		},
		{
			true,
			map[string]interface{}{"t": "#time:0"},
			map[string]string{"dur": "0"},
			map[string]interface{}{"t": time.Now().UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:3600"},
			map[string]string{"dur": strconv.Itoa(3600 * int(time.Second))},
			map[string]interface{}{"t": time.Now().Add(3600 * time.Second).UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:-3600"},
			map[string]string{"dur": strconv.Itoa(-3600 * int(time.Second))},
			map[string]interface{}{"t": time.Now().Add(-3600 * time.Second).UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:3s"},
			map[string]string{"dur": strconv.Itoa(3 * int(time.Second))},
			map[string]interface{}{"t": time.Now().Add(3 * time.Second).UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:3ns"},
			map[string]string{"dur": strconv.Itoa(3 * int(time.Nanosecond))},
			map[string]interface{}{"t": time.Now().Add(3 * time.Nanosecond).UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:-3s"},
			map[string]string{"dur": strconv.Itoa(-3 * int(time.Second))},
			map[string]interface{}{"t": time.Now().Add(-3 * time.Second).UTC()},
		},
		{
			true,
			map[string]interface{}{"t": "#time:-5m"},
			map[string]string{"dur": strconv.Itoa(-5 * int(time.Minute))},
			map[string]interface{}{"t": time.Now().Add(-5 * time.Minute).UTC()},
		},
	}

	t.Logf("Given the need to preprocess commands.")
	{
		for _, cmd := range commands {
			t.Logf("\tWhen using %+v with %+v", cmd.doc, cmd.vars)
			{
				err := exec.ProcessVariables("", cmd.doc, cmd.vars, nil)

				if !cmd.time {

					if eq := compareBson(cmd.doc, cmd.after); !eq {
						t.Log(cmd.doc)
						t.Log(cmd.after)

						t.Errorf("\t%s\tShould get back the expected document.", tests.Failed)
						continue
					}
					t.Logf("\t%s\tShould get back the expected document.", tests.Success)

					continue
				}

				v, _ := strconv.Atoi(cmd.vars["dur"])
				dur := time.Duration(v)

				t.Log(time.Now().UTC())
				t.Log(cmd.after["t"])

				dt, ok := cmd.doc["t"].(time.Time)
				if !ok {
					t.Errorf("\t%s\tShould get back a time value within %v of difference : %v", tests.Failed, dur, err)
					continue
				}

				if eq := compareTime(dt, cmd.after["t"].(time.Time)); !eq {
					t.Errorf("\t%s\tShould get back a time value within %v of difference : %v", tests.Failed, dur, err)
					continue
				}
				t.Logf("\t%s\tShould get back a time value within %v of difference.", tests.Success, dur)
			}
		}
	}
}

// compareTime compares two bson maps for equivalence. This is based
// on a percent of difference since we are dealing with time.
func compareTime(t1 time.Time, t2 time.Time) bool {
	diff := t1.Sub(t2)
	if diff > time.Second {
		return false
	}

	return true
}

// compareBson compares two bson maps for equivalence.
func compareBson(m1 bson.M, m2 bson.M) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if bv, ok := v.(bson.M); ok {
			compareBson(bv, m2)
			continue
		}

		if bv, ok := v.(map[string]interface{}); ok {
			compareBson(bv, m2)
			continue
		}

		if m2[k] != v {
			return false
		}
	}

	for k, v := range m2 {
		if bv, ok := v.(bson.M); ok {
			compareBson(m1, bv)
			continue
		}

		if bv, ok := v.(map[string]interface{}); ok {
			compareBson(m1, bv)
			continue
		}

		if m1[k] != v {
			return false
		}
	}

	return true
}
