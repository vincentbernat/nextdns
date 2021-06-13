package proxy

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	bufferPools = []struct {
		size int
		p    *sync.Pool
	}{
		{maxUDPSize, &sync.Pool{}},
		{maxDNS0Size, &sync.Pool{}},
		{maxDNS0Size * 2, &sync.Pool{}},
		{maxDNS0Size * 8, &sync.Pool{}},
		{maxTCPSize, &sync.Pool{}},
	}
	countBufferPoolAllocs = false
	bufferPoolAllocs      = uint64(0)
)

func bufferPoolGet(size int) *[]byte {
	for _, bp := range bufferPools {
		if size <= bp.size {
			buf := bp.p.Get()
			if buf == nil {
				if countBufferPoolAllocs {
					atomic.AddUint64(&bufferPoolAllocs, 1)
				}
				b := make([]byte, bp.size)
				return &b
			}
			return buf.(*[]byte)
		}
	}
	// Never ask for a too big buffer
	panic(fmt.Sprintf("too big buffer len=%d", size))
}

func bufferPoolPut(p *[]byte) {
	for _, bp := range bufferPools {
		if len(*p) == bp.size {
			bp.p.Put(p)
			return
		}
	}
	// Never release a buffer of the wrong size
	panic(fmt.Sprintf("unexpected buffer len=%v", len(*p)))
}
