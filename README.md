# workerpool

[![en](https://img.shields.io/badge/lang-en-blue.svg)](https://github.com/oseau/workerpool/blob/main/README.md)
[![zh-Hans](https://img.shields.io/badge/lang-zh--Hans-blue.svg)](https://github.com/oseau/workerpool/blob/main/README.zh-Hans.md)

[![Go](https://github.com/oseau/workerpool/actions/workflows/codecov.yml/badge.svg)](https://github.com/oseau/workerpool/actions/workflows/codecov.yml)
[![codecov](https://codecov.io/gh/oseau/workerpool/branch/main/graph/badge.svg)](https://codecov.io/gh/oseau/workerpool)

workerpool is a simple and efficient worker pool implementation in Go.

## Installation

```bash
go get github.com/oseau/workerpool
```

## Usage

Check out the examples in the `examples` directory for more information.

## Features & Design choices

- Zero dependencies
- 100% test coverage
- Simple and composable API
  - `New()` for creating a new pool with defaults, optional configuration may be provided by composing options
    - `WithPoolSize()` for setting the pool size, fixed worker count leads to predictable resource usage
    - `WithQueueSize()` for setting the queue size
  - `Add()` for adding tasks to the pool, non-blocking, fails fast if the pool is closed or the queue is full
  - `Wait()` for waiting for all tasks to complete
  - `Stop()` for stopping the pool
- Context-based cancellation - follows Go idioms for graceful shutdown
- Decoupled by `Task` interface, any type that implements `Task` can be added to the pool
- Efficient worker pool implementation (check benchmarks below)
  - Zero allocations by design - all allocations happen at pool creation

## Benchmarks

```bash
> go test -bench=. -benchmem .
goos: linux
goarch: arm64
pkg: github.com/oseau/workerpool
BenchmarkComparison/small_load/workerpool_only_exec                89955             13018 ns/op               0 B/op          0 allocs/op
BenchmarkComparison/small_load/workerpool                          43704             32827 ns/op            4896 B/op         32 allocs/op
BenchmarkComparison/small_load/raw_goroutines                      35260             35770 ns/op            2420 B/op        102 allocs/op
BenchmarkComparison/small_load/unbuffered_pool                     44536             29641 ns/op            1908 B/op        109 allocs/op
BenchmarkComparison/small_load/buffered_pool                      126960             13456 ns/op            2804 B/op        110 allocs/op
BenchmarkComparison/medium_load/workerpool_only_exec                8961            121890 ns/op               1 B/op          0 allocs/op
BenchmarkComparison/medium_load/workerpool                          8308            156800 ns/op           19067 B/op         44 allocs/op
BenchmarkComparison/medium_load/raw_goroutines                      3288            393416 ns/op           24020 B/op       1002 allocs/op
BenchmarkComparison/medium_load/unbuffered_pool                     4467            277537 ns/op           16404 B/op       1013 allocs/op
BenchmarkComparison/medium_load/buffered_pool                      12555             91724 ns/op           24596 B/op       1014 allocs/op
BenchmarkComparison/high_load/workerpool_only_exec                   984           1213473 ns/op             168 B/op          0 allocs/op
BenchmarkComparison/high_load/workerpool                             954           1442740 ns/op          168692 B/op         76 allocs/op
BenchmarkComparison/high_load/raw_goroutines                         324           3961079 ns/op          240020 B/op      10002 allocs/op
BenchmarkComparison/high_load/unbuffered_pool                        483           2257318 ns/op          160596 B/op      10021 allocs/op
BenchmarkComparison/high_load/buffered_pool                         1617            870150 ns/op          242516 B/op      10022 allocs/op
PASS
ok      github.com/oseau/workerpool     24.072s
```

Running on a MacBook Air M1, 8GB RAM, with [OrbStack](https://www.orbstack.dev/) as the container runtime throttled to 100% CPU and 1GB memory limit.

Since our `workerpool` implementation is rather simple and it's not fair to compare it against other libraries with more complex implementations, we're only comparing it against raw goroutines to demonstrate the performance difference.

This benchmark is not meant to be comprehensive, but rather to give you a general idea of the performance. Our `workerpool` scales well across small(100 tasks), medium(1000 tasks) and high(10000 tasks) loads. It's comparable to other implementations when the load is small, but it outperforms raw goroutines and unbuffered channels by a large margin when the load is high.

## Error handling

`workerpool` provides clear error types for common scenarios:

- `ErrPoolClosed`: Returned when trying to add tasks to a stopped pool
- `ErrQueueFull`: Returned when the task queue is full (non-blocking design)

These errors are returned by the `Add()` method, and are meant to be handled by the caller. We could remove errors in the future (check `Possible improvements` section below).

## Possible improvements

- The initial requirement was to have a pool that listens to incoming tasks and processes them in a FIFO manner, and remains waiting for new tasks after all tasks are completed until the pool is stopped. This provides a simple way to control the total number of tasks that can be processed concurrently. However, if this control is not needed, it's possible to have a simpler implementation, which accepts the tasks when initializing the pool and get rid of the `Add()` method as well as the `queueSize` option. And it kinda make more sense as the user of this library might only interested in how many tasks they can process concurrently, not how many tasks they can add to the pool. And the creation of the pool is quite cheap (all we do when initializing the pool is spawn limited number of goroutines), so each time user wants to do a group of tasks, they can just create a new pool. This way we can also get rid of the all the errors (ErrPoolClosed & ErrQueueFull) and make the usage of the pool much more straightforward.
- Metrics like processing time, task count, error count, etc.

## Contributing & Development

To run benchmarks and tests, make sure you have `Docker` and `make` installed.

A `.env.example` file is provided, `cp .env.example .env` is needed for the first time.

This project has a entrypoint at `Makefile` for running benchmarks and tests. You can run `make` to see all available commands.

```bash
> make
test       Run tests in Docker
coverage   Generate and open test coverage report
bench      Run benchmarks in Docker
shell      Start a shell in the Docker container
example    Run the basic example
```
