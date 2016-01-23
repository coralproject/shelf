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
	BenchmarkPPNumber-8	50000000	       113 ns/op	       0 B/op	       0 allocs/op
	BenchmarkPPString-8	30000000	       126 ns/op	       0 B/op	       0 allocs/op
	BenchmarkPPDate-8  	30000000	       175 ns/op	       0 B/op	       0 allocs/op
	ok  	github.com/coralproject/xenia/pkg/exec	15.359s
*/

var ppCmds = []map[string]interface{}{
	{"$match": map[string]interface{}{"duration": "#number:duration"}},
	{"$match": map[string]interface{}{"target": "#string:target"}},
	{"$match": map[string]interface{}{"start": map[string]interface{}{"$gte": "#date:start"}}},
}

var ppVars = map[string]string{
	"duration": "10",
	"target":   "bill",
	"start":    "2013-01-16T00:00:00.000Z",
}

var ppRes map[string]interface{}

func BenchmarkPPNumber(b *testing.B) {
	var res map[string]interface{}

	for i := 0; i < b.N; i++ {
		res = PreProcess(ppCmds[0], ppVars)
	}

	ppRes = res
}

func BenchmarkPPString(b *testing.B) {
	var res map[string]interface{}

	for i := 0; i < b.N; i++ {
		res = PreProcess(ppCmds[1], ppVars)
	}

	ppRes = res
}

func BenchmarkPPDate(b *testing.B) {
	var res map[string]interface{}

	for i := 0; i < b.N; i++ {
		res = PreProcess(ppCmds[2], ppVars)
	}

	ppRes = res
}
