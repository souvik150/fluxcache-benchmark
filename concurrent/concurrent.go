package concurrentcorrectness

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"fluxcache-test/userpb"
	"github.com/souvik150/fluxcache/pkg/fluxcache"
)

type record struct {
	Timestamp string
	Writes    int64
	Reads     int64
}

func createUser(id int) *userpb.User {
	return &userpb.User{
		Id:    fmt.Sprintf("user-%d", id),
		Name:  fmt.Sprintf("Name-%d", id),
		Email: fmt.Sprintf("user%d@example.com", id),
	}
}

func usersEqual(u1, u2 *userpb.User) bool {
	if u1 == nil || u2 == nil {
		return false
	}
	return u1.Id == u2.Id && u1.Name == u2.Name && u1.Email == u2.Email
}

func MultiReaderWriterCorrectness(ctx context.Context, cache *fluxcache.FluxCache, totalKeys, writerCount, readerCount int, testDuration time.Duration) {
	fmt.Println("Starting MULTI-READER-WRITER CONSISTENCY benchmark...")

	ctx, cancel := context.WithTimeout(ctx, testDuration)
	defer cancel()

	var (
		expectedData sync.Map
		totalErrors  int64
		totalReads   int64
		totalWrites  int64
		wg           sync.WaitGroup
	)

	// --- Writers ---
	for i := 0; i < writerCount; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(writerID)))
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				id := r.Intn(totalKeys)
				key := fmt.Sprintf("user:%d", id)
				user := createUser(id)

				err := cache.SetProto(key, user)
				if err == nil {
					expectedData.Store(key, user)
					atomic.AddInt64(&totalWrites, 1)
				}
			}
		}(i)
	}

	// --- Readers ---
	for i := 0; i < readerCount; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(readerID*1000)))
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				id := r.Intn(totalKeys)
				key := fmt.Sprintf("user:%d", id)

				var got userpb.User
				err := cache.GetProto(key, &got)
				if err != nil {
					continue
				}

				value, ok := expectedData.Load(key)
				if !ok {
					continue
				}
				expected := value.(*userpb.User)

				if !usersEqual(expected, &got) {
					atomic.AddInt64(&totalErrors, 1)
					fmt.Printf("[Reader %d] MISMATCH key=%s\n", readerID, key)
				}
				atomic.AddInt64(&totalReads, 1)
			}
		}(i)
	}

	// --- Live Ops/sec Tracker ---
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var csvRecords []record

	go func() {
		startTime := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				writes := atomic.LoadInt64(&totalWrites)
				reads := atomic.LoadInt64(&totalReads)
				now := time.Since(startTime).Truncate(time.Second)

				fmt.Printf("[Live] Time=%v | Total Writes=%d, Total Reads=%d\n", now, writes, reads)

				csvRecords = append(csvRecords, record{
					Timestamp: now.String(),
					Writes:    writes,
					Reads:     reads,
				})
			}
		}
	}()

	// --- Wait for all goroutines to exit after test duration ---
	wg.Wait()

	// --- Final Report ---
	fmt.Println("MULTI-READER-WRITER CONSISTENCY benchmark finished")
	finalWrites := atomic.LoadInt64(&totalWrites)
	finalReads := atomic.LoadInt64(&totalReads)
	finalErrors := atomic.LoadInt64(&totalErrors)

	fmt.Printf("Total writes: %d\n", finalWrites)
	fmt.Printf("Total reads: %d\n", finalReads)
	fmt.Printf("Total mismatches (errors): %d\n", finalErrors)

	// Save to CSV
	saveResultsToCSV(csvRecords)

	if finalErrors > 0 {
		fmt.Println("❌ Test Failed: Found mismatches!")
		panic(errors.New("test failed"))
	} else {
		fmt.Println("✅ Test Passed: No mismatches!")
	}
}

func saveResultsToCSV(records []record) {
	file, err := os.Create(fmt.Sprintf("benchmark_result_%d.csv", time.Now().Unix()))
	if err != nil {
		fmt.Printf("Failed to create CSV: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Timestamp", "Writes", "Reads"})

	// Write all records
	for _, rec := range records {
		writer.Write([]string{rec.Timestamp, fmt.Sprintf("%d", rec.Writes), fmt.Sprintf("%d", rec.Reads)})
	}

	fmt.Println("📈 Benchmark results saved to CSV successfully.")
}
