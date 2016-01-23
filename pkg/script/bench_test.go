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
	BenchmarkPrepareForInsert-8	10000000	       722 ns/op	       0 B/op	       0 allocs/op
	BenchmarkPrepareForUse-8   	10000000	       763 ns/op	       0 B/op	       0 allocs/op
	ok  	github.com/coralproject/xenia/pkg/script	16.433s
*/

// TODO: Review these benchmarks with community. Since these function alter
// the existing map, I am not sure the benchmarks are providing an accurate
// view.

func BenchmarkPrepareForInsert(b *testing.B) {
	var cmd = map[string]interface{}{
		"$group": map[string]interface{}{
			"_id": map[string]interface{}{
				"day": map[string]interface{}{
					"$dayOfMonth": "$date_created",
					"month": map[string]interface{}{
						"$month": "$date_created",
						"year": map[string]interface{}{
							"$year": "$date_created",
							"comm.ents": map[string]interface{}{
								"$sum": 1,
							},
						},
					},
				},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prepareForInsert(cmd)
	}
}

func BenchmarkPrepareForUse(b *testing.B) {
	var cmd = map[string]interface{}{
		"_$group": map[string]interface{}{
			"_id": map[string]interface{}{
				"day": map[string]interface{}{
					"_$dayOfMonth": "_$date_created",
					"month": map[string]interface{}{
						"_$month": "_$date_created",
						"year": map[string]interface{}{
							"_$year": "_$date_created",
							"comm~ents": map[string]interface{}{
								"_$sum": 1,
							},
						},
					},
				},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prepareForUse(cmd)
	}
}
