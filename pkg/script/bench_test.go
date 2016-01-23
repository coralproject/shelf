package script

import (
	"testing"
)

//==============================================================================
// PrepareForInsert benchmarks to test the performance and memory alloctions for
// transforming MongoDB $ and . inside of command field names.
/*
	$ go test -run none -bench . -benchtime 3s -benchmem
	PASS
	BenchmarkPrepareForInsert-8	 2000000	      2074 ns/op	      96 B/op	       6 allocs/op
	BenchmarkPrepareForUse-8   	 3000000	      1679 ns/op	      16 B/op	       1 allocs/op
	ok  	github.com/coralproject/xenia/pkg/script	36.236s
*/

// TODO: Review these benchmarks with community. Since these function alter
// the existing map, I am not sure the benchmarks are providing an accurate
// view.

// BenchmarkPrepareForInsert tests the processing of commands for insert.
func BenchmarkPrepareForInsert(b *testing.B) {

	// Generate a set of unique documents to process.
	a := make([]map[string]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = newIns()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prepareForInsert(a[i])
	}
}

// BenchmarkPrepareForUse tests the processing of commands for use.
func BenchmarkPrepareForUse(b *testing.B) {

	// Generate a set of unique documents to process.
	a := make([]map[string]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = newUse()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prepareForUse(a[i])
	}
}

//==============================================================================

func newIns() map[string]interface{} {
	return map[string]interface{}{
		"$group": map[string]interface{}{
			"_id": map[string]interface{}{
				"day": map[string]interface{}{
					"$dayOfMonth": "date_created",
					"month": map[string]interface{}{
						"$month": "date_created",
						"year": map[string]interface{}{
							"$year": "date_created",
							"comm.ents": map[string]interface{}{
								"$sum": 1,
							},
						},
					},
				},
			},
		},
	}
}

func newUse() map[string]interface{} {
	return map[string]interface{}{
		"_$group": map[string]interface{}{
			"_id": map[string]interface{}{
				"day": map[string]interface{}{
					"_$dayOfMonth": "date_created",
					"month": map[string]interface{}{
						"_$month": "date_created",
						"year": map[string]interface{}{
							"_$year": "date_created",
							"comm*ents": map[string]interface{}{
								"_$sum": 1,
							},
						},
					},
				},
			},
		},
	}
}
