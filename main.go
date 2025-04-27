package main

import (
	"context"
	"fmt"
	"time"

	"github.com/souvik150/fluxcache/pkg/fluxcache"

	concurrentcorrectness "fluxcache-test/concurrent"
)


func main() {
	ctx := context.Background()

	opts := fluxcache.Options{
		RedisAddr:    "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		MemoryCap:     1_000_000,
		SyncAtStart:   false, 
		SyncInterval:  5,              
		KeyPattern:    "user:*",
	}

	cache, err := fluxcache.NewFluxCache(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("failed to initialize FluxCache: %w", err))
	}

	const totalUsers = 100_000
	// const concurrency = 100

	const totalOps = 500_000
	const concurrency = 200

	// set.SetBenchmark(ctx, cache, totalUsers, concurrency) 

	// get.GetBenchmark(ctx, cache, totalUsers, concurrency)   

	// get.GetBenchmark(ctx, cache, totalUsers, concurrency)   

	// mixed.MixedSetGetBenchmark(ctx, cache, totalOps, concurrency)

	concurrentcorrectness.MultiReaderWriterCorrectness(ctx, cache, 10000, 50, 200, 100*time.Second)
}
