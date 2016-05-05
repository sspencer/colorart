package colorart

import (
	"runtime"
	"sync/atomic"
)

type countedFn func(ch chan countedSet, pmin, pmax int)

// parallelize data processing if 'enabled' is true
func parallelize(datamin, datamax int, fn countedFn) countedSet {
	datasize := datamax - datamin
	partsize := datasize

	numGoroutines := 1
	numProcs := runtime.GOMAXPROCS(0)

	if numProcs > 1 {
		numGoroutines = numProcs
		partsize = partsize / numGoroutines
		if partsize < 1 {
			partsize = 1
		}

		// if partsize had a fraction, bump it by 1 so entire image is covered
		if partsize*numGoroutines < datasize {
			partsize++
		}
	}

	idx := int64(datamin)
	ch := make(chan countedSet, numGoroutines)

	for p := 0; p < numGoroutines; p++ {
		go func() {
			for {
				pmin := int(atomic.AddInt64(&idx, int64(partsize))) - partsize
				if pmin >= datamax {
					break
				}
				pmax := pmin + partsize
				if pmax > datamax {
					pmax = datamax
				}
				fn(ch, pmin, pmax)
			}
		}()
	}

	colors := newCountedSet(10000, 0)

	for p := 0; p < numGoroutines; p++ {
		colors.merge(<-ch)
	}

	close(ch)

	return colors
}
