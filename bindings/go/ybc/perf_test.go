package ybc

import (
	"sync"
	"testing"
)

type benchmark_InitFunc func(cache TinyCacher)
type benchmark_IterationFunc func(cache TinyCacher)

func runBenchmark(isTinyCache bool, initFunc benchmark_InitFunc, iterationFunc benchmark_IterationFunc, workersCount int, b *testing.B) {
	b.StopTimer()

	config := newConfig()
	var cache TinyCacher
	var err error
	if isTinyCache {
		cache, err = config.OpenTinyCache(8, true)
	} else {
		cache, err = config.OpenCache(true)
	}
	if err != nil {
		b.Fatal(err)
	}
	defer cache.Close()

	initFunc(cache)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	batchSize := 1000
	loopsCount := b.N
	ch := make(chan int, b.N/batchSize+1)
	for {
		if batchSize > loopsCount {
			ch <- loopsCount
			break
		}
		loopsCount -= batchSize
		ch <- batchSize
	}
	defer close(ch)

	benchLoop := func() {
		defer wg.Done()
		for loopsCount := range ch {
			for i := 0; i < loopsCount; i++ {
				iterationFunc(cache)
			}
		}
	}

	b.StartTimer()
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go benchLoop()
	}
}

func BenchmarkCache_GetHit_1(b *testing.B) {
	benchmarkCache_GetHit(b, 1, false)
}

func BenchmarkCache_GetHit_2(b *testing.B) {
	benchmarkCache_GetHit(b, 2, false)
}

func BenchmarkCache_GetHit_4(b *testing.B) {
	benchmarkCache_GetHit(b, 4, false)
}

func BenchmarkCache_GetHit_8(b *testing.B) {
	benchmarkCache_GetHit(b, 8, false)
}

func BenchmarkCache_GetHit_16(b *testing.B) {
	benchmarkCache_GetHit(b, 16, false)
}

func BenchmarkCache_GetHitTiny_1(b *testing.B) {
	benchmarkCache_GetHit(b, 1, true)
}

func BenchmarkCache_GetHitTiny_2(b *testing.B) {
	benchmarkCache_GetHit(b, 2, true)
}

func BenchmarkCache_GetHitTiny_4(b *testing.B) {
	benchmarkCache_GetHit(b, 4, true)
}

func BenchmarkCache_GetHitTiny_8(b *testing.B) {
	benchmarkCache_GetHit(b, 8, true)
}

func BenchmarkCache_GetHitTiny_16(b *testing.B) {
	benchmarkCache_GetHit(b, 16, true)
}

func benchmarkCache_GetHit(b *testing.B, workersCount int, isTinyCache bool) {
	key := []byte("12345")
	value := []byte("value")

	initFunc := func(cache TinyCacher) {
		if err := cache.Set(key, value, MaxTtl); err != nil {
			b.Fatalf("Cannot set item with key=[%s], value=[%s]: [%s]", key, value, err)
		}
	}

	iterationFunc := func(cache TinyCacher) {
		if _, err := cache.Get(key); err != nil {
			b.Fatalf("Cannot find item with key=[%s]: [%s]", key, err)
		}
	}

	runBenchmark(isTinyCache, initFunc, iterationFunc, workersCount, b)
}

func BenchmarkCache_GetMiss_1(b *testing.B) {
	benchmarkCache_GetMiss(b, 1, false)
}

func BenchmarkCache_GetMiss_2(b *testing.B) {
	benchmarkCache_GetMiss(b, 2, false)
}

func BenchmarkCache_GetMiss_4(b *testing.B) {
	benchmarkCache_GetMiss(b, 4, false)
}

func BenchmarkCache_GetMiss_8(b *testing.B) {
	benchmarkCache_GetMiss(b, 8, false)
}

func BenchmarkCache_GetMiss_16(b *testing.B) {
	benchmarkCache_GetMiss(b, 16, false)
}

func BenchmarkCache_GetMissTiny_1(b *testing.B) {
	benchmarkCache_GetMiss(b, 1, true)
}

func BenchmarkCache_GetMissTiny_2(b *testing.B) {
	benchmarkCache_GetMiss(b, 2, true)
}

func BenchmarkCache_GetMissTiny_4(b *testing.B) {
	benchmarkCache_GetMiss(b, 4, true)
}

func BenchmarkCache_GetMissTiny_8(b *testing.B) {
	benchmarkCache_GetMiss(b, 8, true)
}

func BenchmarkCache_GetMissTiny_16(b *testing.B) {
	benchmarkCache_GetMiss(b, 16, true)
}

func benchmarkCache_GetMiss(b *testing.B, workersCount int, isTinyCache bool) {
	key := []byte("12345")

	initFunc := func(cache TinyCacher) {}

	iterationFunc := func(cache TinyCacher) {
		if _, err := cache.Get(key); err != ErrCacheMiss {
			b.Fatalf("Unexpected item found for key=[%s]", key)
		}
	}

	runBenchmark(isTinyCache, initFunc, iterationFunc, workersCount, b)
}

func BenchmarkCache_Set_1(b *testing.B) {
	benchmarkCache_Set(b, 1, false)
}

func BenchmarkCache_Set_2(b *testing.B) {
	benchmarkCache_Set(b, 2, false)
}

func BenchmarkCache_Set_4(b *testing.B) {
	benchmarkCache_Set(b, 4, false)
}

func BenchmarkCache_Set_8(b *testing.B) {
	benchmarkCache_Set(b, 8, false)
}

func BenchmarkCache_Set_16(b *testing.B) {
	benchmarkCache_Set(b, 16, false)
}

func BenchmarkCache_SetTiny_1(b *testing.B) {
	benchmarkCache_Set(b, 1, true)
}

func BenchmarkCache_SetTiny_2(b *testing.B) {
	benchmarkCache_Set(b, 2, true)
}

func BenchmarkCache_SetTiny_4(b *testing.B) {
	benchmarkCache_Set(b, 4, true)
}

func BenchmarkCache_SetTiny_8(b *testing.B) {
	benchmarkCache_Set(b, 8, true)
}

func BenchmarkCache_SetTiny_16(b *testing.B) {
	benchmarkCache_Set(b, 16, true)
}

func benchmarkCache_Set(b *testing.B, workersCount int, isTinyCache bool) {
	key := []byte("12345")
	value := []byte("value")

	initFunc := func(cache TinyCacher) {}

	iterationFunc := func(cache TinyCacher) {
		if err := cache.Set(key, value, MaxTtl); err != nil {
			b.Fatalf("Error in Cache.Set(key=[%s], value=[%s])=[%s]", key, value, err)
		}
	}

	runBenchmark(isTinyCache, initFunc, iterationFunc, workersCount, b)
}
