package query

import (
	"testing"
	"time"

	"github.com/coralproject/shelf/pkg/tests"
)

/*
	const name = "QTEST_spending_advice"
	const fixture = "./fixtures/spending_advice.json"
	qs1, err := getFixture(fixture)
	if err != nil {
		t.Fatalf("\t%s\tShould load query record from file : %v", tests.Failed, err)
	}
	t.Logf("\t%s\tShould load query record from file.", tests.Success)
*/

// compareBsonMaps compares two bson maps for equivalence.
func compareBsonMaps(m1 bson.M, m2 bson.M) bool {

}

//==============================================================================

// TestUmarshalMongoScript tests the ability to convert string based Mongo
// commands into a bson map for processing.
func TestUmarshalMongoScript(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	scripts := []struct {
		text string
		so   ScriptOption
		cmp  bson.M
	}{
		{"{\"name\":\"bill\"}", nil, bson.M{"name": "bill"}},
	}

	t.Logf("Given the need to convert mongo commands.")
	{
		for _, script := range scripts {
			t.Logf("\tWhen using %s with %+v", script.text, script.so)
			{
				doc, err := umarshalMongoScript(script.text, script.so)
				if err != nil {
					t.Errorf("\t%s\tShould be able to convert without an error : %v", tests.Failed, err)
					continue
				}
				t.Logf("\t%s\tShould be able to convert without an error.", tests.Success)

			}
		}
	}
}

// TestIsoDate tests the support function to convert an ISODate string
// to a Go Time value.
func TestIsoDate(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	have := "ISODate('2013-01-16T00:00:00.000Z')"
	expected, _ := time.Parse("2006-01-02T15:04:05.999Z", "2013-01-16T00:00:00.000Z")

	t.Logf("Given the need to convert ISODate strings.")
	{
		t.Logf("\tWhen using %s", have)
		{
			tm := isoDate(have)
			if tm != expected {
				t.Fatalf("\t%s\tShould have the proper time value : %v", tests.Failed, tm)
			}
			t.Logf("\t%s\tShould have the proper time value.", tests.Success)
		}
	}
}
