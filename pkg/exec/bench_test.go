package exec

import (
	"testing"
)

//==============================================================================
// PreProcess benchmarks to test the performance and memory alloctions for
// parsing variables inside of commands.
/*
	$ go test -run none -bench BenchmarkPP -benchtime 3s -benchmem
	PASS
	BenchmarkPPNumber-8	20000000	       815 ns/op	      16 B/op	       1 allocs/op
	BenchmarkPPString-8	20000000	       945 ns/op	      16 B/op	       1 allocs/op
	BenchmarkPPDate-8  	 5000000	      1129 ns/op	      32 B/op	       1 allocs/op
	ok  	github.com/coralproject/xenia/pkg/exec	123.666s
*/

// TODO: Review these benchmarks with community. Since these function alter
// the existing map, I am not sure the benchmarks are providing an accurate
// view.

var ppVars = map[string]string{
	"duration": "10",
	"target":   "bill",
	"start":    "2013-01-16T00:00:00.000Z",
}

// BenchmarkPPNumber tests the processing of numbers.
func BenchmarkPPNumber(b *testing.B) {

	// Generate a set of unique documents to process.
	a := make([]map[string]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = map[string]interface{}{
			"$match": map[string]interface{}{"duration": "#number:duration"},
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ProcessVariables("", a[i], ppVars, nil)
	}
}

// BenchmarkPPNumber tests the processing of strings.
func BenchmarkPPString(b *testing.B) {

	// Generate a set of unique documents to process.
	a := make([]map[string]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = map[string]interface{}{
			"$match": map[string]interface{}{"target": "#string:target"},
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ProcessVariables("", a[i], ppVars, nil)
	}
}

// BenchmarkPPNumber tests the processing of dates.
func BenchmarkPPDate(b *testing.B) {

	// Generate a set of unique documents to process.
	a := make([]map[string]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = map[string]interface{}{
			"$match": map[string]interface{}{"start": map[string]interface{}{"$gte": "#date:start"}},
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ProcessVariables("", a[i], ppVars, nil)
	}
}
