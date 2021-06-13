package proxy

import (
	"runtime/debug"
	"testing"

	"github.com/nextdns/nextdns/internal/race"
)

func TestBufferPool(t *testing.T) {
	countBufferPoolAllocs = true
	origGCPercent := debug.SetGCPercent(-1)
	defer func() {
		countBufferPoolAllocs = false
		debug.SetGCPercent(origGCPercent)
	}()
	sizes := []int{10, 100, 1000, 2000, 3000, 8000, 16000, 32000, maxTCPSize}
	routines := 5
	wait := make(chan int)
	for r := 0; r < routines; r++ {
		go func() {
			for i := 0; i < 2000; i++ {
				for _, s := range sizes {
					buf := bufferPoolGet(s)
					if len(*buf) < s {
						t.Fatalf("bufferPoolGet(%d) length too short: %d",
							s, len(*buf))
					}
					bufferPoolPut(buf)
				}
			}
			wait <- 1
		}()
	}
	for r := 0; r < routines; r++ {
		<-wait
	}
	if !race.Enabled && bufferPoolAllocs > uint64(routines*len(sizes)) {
		t.Errorf("bufferPoolGet() did allocate too many buffers: %d",
			bufferPoolAllocs)
	}
}
