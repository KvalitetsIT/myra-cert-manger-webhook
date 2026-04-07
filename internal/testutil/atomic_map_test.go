package testutil_test

import (
	"sync"
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil"
)

func TestAtomicMap_GetAll_Race(t *testing.T) {
	m := testutil.NewAtomicMap[int, int]("test")

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Set(i, i)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.GetAll()
		}()
	}

	wg.Wait()
}
