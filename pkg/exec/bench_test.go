package exec

import (
	"testing"
)

//==============================================================================
// Parsing benchmarks to test the performance and memory allocations for
// parsing varaibles inside each command.
/*
	$ go test -run none -bench BenchmarkParse -benchtime 3s -benchmem
	PASS
	BenchmarkParseNumber-8	50000000	        79.8 ns/op	      16 B/op	       1 allocs/op
	BenchmarkParseString-8	50000000	        77.8 ns/op	      16 B/op	       1 allocs/op
	BenchmarkParseDate-8  	50000000	        82.0 ns/op	      16 B/op	       1 allocs/op
	ok  	github.com/coralproject/xenia/pkg/exec	12.409s
*/

var parseCmds = []string{
	`{"duration": "#number:duration"}`,
	`{"target": "#string:target"}`,
	`{"$gte": "#date:start"}`,
}

var parseVars = map[string]string{
	"duration": "10",
	"target":   "bill2",
	"start":    "2016-02-15",
}

var parseRes interface{}

func BenchmarkParseNumber(b *testing.B) {
	var res interface{}

	for i := 0; i < b.N; i++ {
		res = parse(parseCmds[0], parseVars)
	}

	parseRes = res
}

func BenchmarkParseString(b *testing.B) {
	var res interface{}

	for i := 0; i < b.N; i++ {
		res = parse(parseCmds[1], parseVars)
	}

	parseRes = res
}

func BenchmarkParseDate(b *testing.B) {
	var res interface{}

	for i := 0; i < b.N; i++ {
		res = parse(parseCmds[2], parseVars)
	}

	parseRes = res
}

//==============================================================================
// PreProcess benchmarks to test the performance and memory alloctions for
// parsing variables inside of commands.
/*
	$ go test -run none -bench BenchmarkPP -benchtime 3s -benchmem
	PASS
	BenchmarkPPNumber-8	50000000	       112 ns/op	       0 B/op	       0 allocs/op
	BenchmarkPPString-8	30000000	       122 ns/op	       0 B/op	       0 allocs/op
	BenchmarkPPDate-8  	30000000	       183 ns/op	       0 B/op	       0 allocs/op
	ok  	github.com/coralproject/xenia/pkg/exec	15.443s
*/

// TODO: Review these benchmarks with community. Since these function alter
// the existing map, I am not sure the benchmarks are providing an accurate
// view.

var ppVars = map[string]string{
	"duration": "10",
	"target":   "bill",
	"start":    "2013-01-16T00:00:00.000Z",
}

func BenchmarkPPNumber(b *testing.B) {
	cmd := map[string]interface{}{
		"$match": map[string]interface{}{"duration": "#number:duration"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PreProcess(cmd, ppVars)
	}
}

func BenchmarkPPString(b *testing.B) {
	cmd := map[string]interface{}{
		"$match": map[string]interface{}{"target": "#string:target"},
	}

	for i := 0; i < b.N; i++ {
		PreProcess(cmd, ppVars)
	}
}

func BenchmarkPPDate(b *testing.B) {
	cmd := map[string]interface{}{
		"$match": map[string]interface{}{"start": map[string]interface{}{"$gte": "#date:start"}},
	}

	for i := 0; i < b.N; i++ {
		PreProcess(cmd, ppVars)
	}
}
