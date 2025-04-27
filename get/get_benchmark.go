package get

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/souvik150/fluxcache/pkg/fluxcache"

	"fluxcache-test/userpb"
)

func GetBenchmark(ctx context.Context, cache *fluxcache.FluxCache, totalUsers, concurrency int) {
	fmt.Println("Starting GET benchmark...")

	start := time.Now()

	var wg sync.WaitGroup
	getCh := make(chan int, totalUsers)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range getCh {
				var user userpb.User
				err := cache.GetProto(fmt.Sprintf("user:%d", id), &user)
				if err != nil {
					fmt.Printf("Error getting user %d: %v\n", id, err)
				}
			}
		}()
	}

	for i := 0; i < totalUsers; i++ {
		getCh <- i
	}
	close(getCh)

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("GET completed: %d users in %v (%.2f ops/sec)\n", totalUsers, duration, float64(totalUsers)/duration.Seconds())
}
