package query

import (
	"testing"

	"github.com/coralproject/shelf/pkg/tests"
)

// TestRenderScript validates the process of subsituting variables within
// source scripts.
func TestRenderScript(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Logf("Given the need to subsitute variables into a script source.")
	{
		script := "{ \"$match\" : { \"user_id\" : \"#user_id#\", \"category\" : \"#category#\" }}"
		t.Logf("\tWhen giving a script source and a variable map")
		{
			model := map[string]string{"user_id": "10", "category": "petrol"}
			wanted := "{ \"$match\" : { \"user_id\" : \"10\", \"category\" : \"petrol\" }}"
			rendered := renderScript(script, model)
			if rendered != wanted {
				t.Logf("\tModel: %+v", model)
				t.Logf("\tRender: %+v", rendered)
				t.Logf("\tExpected: %+v", wanted)
				t.Errorf("\t%s\tShould have matched expected output from script src rendering", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched expected output from script src rendering", tests.Success)
			}

			model2 := map[string]string{"category": "diesel"}
			wanted2 := "{ \"$match\" : { \"user_id\" : \"#user_id#\", \"category\" : \"diesel\" }}"
			rendered2 := renderScript(script, model2)
			if rendered2 != wanted2 {
				t.Logf("\tModel: %+v", model2)
				t.Logf("\tRender: %+v", rendered2)
				t.Logf("\tExpected: %+v", wanted2)
				t.Errorf("\t%s\tShould have matched expected output with partial data", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched expected output with partial data", tests.Success)
			}

			model3 := map[string]string{"name": "fieldsteign", "age": "30"}
			rendered3 := renderScript(script, model3)
			if rendered3 != script {
				t.Logf("\tModel: %+v", model3)
				t.Logf("\tRender: %+v", rendered3)
				t.Logf("\tExpected: %+v", script)
				t.Errorf("\t%s\tShould have matched input script source when model lacks proper data", tests.Failed)
			} else {
				t.Logf("\t%s\tShould have matched input script source when model lacks proper data", tests.Success)
			}
		}
	}
}
