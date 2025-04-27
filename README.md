# FluxCache Benchmarking Suite

Benchmarking scripts, concurrency tests, and result files for [FluxCache](https://github.com/souvik150/fluxcache).

## 📊 Tests Included
- Concurrent Multi-Reader Multi-Writer Stress Test
- Random read/write correctness validation
- Per-second ops/sec tracking
- Full CSV output logging

## 📊 Benchmark Results

| Benchmark | Result |
|:---|:---|
| **SET Benchmark** | 1,000,000 users in 3m41.76s (4509 ops/sec) |
| **GET Benchmark** | 1,000,000 users in 1m23.23s (12014 ops/sec) |
| **MIXED Benchmark** | 500,000 ops in 1m13.63s (6790 ops/sec) — 0 errors |
| **Concurrent Correctness** | 994,277 reads, 59,089 writes, 0 mismatches ✅ |

### Notes:
- ✅ No data races detected under full `-race` mode.
- ✅ Memory + Redis hybrid caching worked flawlessly.
- ✅ Benchmark conducted under concurrency of 200.
