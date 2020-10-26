package uuid

import (
	"sync"
	"testing"
	"time"
)

func TestSnowflake_Generate(t *testing.T) {
	w := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			t.Log(GetTime(GetID()).Format(time.RFC3339))
			w.Done()
		}()
	}
	w.Wait()
}

func BenchmarkSnowflake_Generate(b *testing.B) {
	var id int64 = 0
	for i := 0; i < b.N; i++ {
		id = GetID()
	}
	b.Log("last:", id)
}
