package worker

import (
	"context"
	"fmt"
	"testing"
)

var terms = []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}

func BenchmarkConcurrent(b *testing.B) {
	dd := New(3).Start(context.Background(), func(j Job) {}) // start up worker pool
	for n := 0; n < b.N; n++ {
		for i := range terms {
			dd.Submit(Job{
				ID:   uint64(i),
				Data: []byte(fmt.Sprintf("JobID::%d", i)),
			})
		}
	}
}
