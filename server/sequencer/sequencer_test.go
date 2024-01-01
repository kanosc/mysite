package sequencer

import (
	"testing"
)

func TestGenerateId(t *testing.T) {
	ids := make(chan string, 100)
	for i := 0; i < 10; i++ {
		go func(i int, out chan string) {
			for j := 0; j < 10; j++ {
				id := DefaultIdGenerator().GenerateId()
				out <- id
			}
		}(i, ids)
	}
	for i := 0; i < 100; i++ {
		t.Log(<-ids)
	}
}
