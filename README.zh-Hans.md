# workerpool

[![en](https://img.shields.io/badge/lang-en-blue.svg)](https://github.com/oseau/workerpool/blob/main/README.md)
[![zh-Hans](https://img.shields.io/badge/lang-zh--Hans-blue.svg)](https://github.com/oseau/workerpool/blob/main/README.zh-Hans.md)

[![Go](https://github.com/oseau/workerpool/actions/workflows/codecov.yml/badge.svg)](https://github.com/oseau/workerpool/actions/workflows/codecov.yml)
[![codecov](https://codecov.io/gh/oseau/workerpool/branch/main/graph/badge.svg)](https://codecov.io/gh/oseau/workerpool)

workerpool 是一个简单且高效的 Go 语言工作池实现。

## 安装

```bash
go get github.com/oseau/workerpool
```

## 使用

查看 `examples` 目录中的示例了解更多信息。

## 特性

- 零依赖
- 100% 测试覆盖度
- 简单、组合式 API
  - `New()` 用于创建一个工作池，
    - 可选配置可以通过 `WithPoolSize()` 设置池大小，固定 worker 数量可以带来可预测的资源使用
    - 可选配置可以通过 `WithQueueSize()` 设置队列大小
  - `Add()` 用于添加任务，非阻塞，如果工作池已关闭或队列已满，则返回错误
  - `Wait()` 用于等待所有任务完成
  - `Stop()` 用于停止工作池
- 支持 context - 遵循 Go 语言的惯例
- 通过 `Task` 接口实现与用户端的解耦
- 高效的 worker 池实现（benchmark 见下文）
  - 0 分配 - 所有分配都在池创建时发生

## 基准测试

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

本次基准测试运行在 MacBook Air M1 上，8GB RAM，使用 [OrbStack](https://www.orbstack.dev/) 作为容器运行时，限制 CPU 为 100% 和可用内存为 1GB。

由于我们的 `workerpool` 实现功能相对简单，如果与其他支持更多功能的库做对比可能得出的结论不具有参考性。因此我们只与原始 goroutines 进行对比，以展示性能差异。

本 benchmark 旨在为用户提供一个大致的性能参考。我们的 `workerpool` 在小型（100 个任务）、中型（1000 个任务）和高负载（10000 个任务）使用场景下表现良好，并且随着负载增加，性能提升明显。当负载较小时，它与其他实现性能相当，但在高负载时，它明显优于原始 goroutines 和未缓冲 channel。

## 错误处理

`workerpool` 提供了清晰的错误类型，用于常见的场景：

- `ErrPoolClosed`: 当尝试向已关闭的池中添加任务时返回
- `ErrQueueFull`: 当任务队列已满时返回

这些错误由 `Add()` 方法返回，并由调用者处理。我们可能会在未来移除这些错误（见 `可能的改进` 部分）。

## 可能的改进

- 本项目的初始要求是创建一个工作池，监听传入的任务，并以 FIFO 方式处理它们，并在所有任务完成后继续等待新任务，直到手动停止。这提供了一种简单的方法来控制可并发处理的任务总数。如果这个控制不是必须的，我们可以在初始化工作池的时候以参数形式传入任务，这样可以去掉 `Add()` 方法、 `queueSize` 选项以及目前定义的两个错误 `ErrPoolClosed` 和 `ErrQueueFull`。这将大大简化用户的使用，只需要初始化并 `Wait()`。这个改进可能是有意义的，因为用户通常并不关心工作池的内部实现，比如缓冲区的队列长度，而更关心的是可并行执行的任务数。
- 运行信息，如处理时间、已处理任务数量、错误计数等。

## 贡献 & 开发

要运行基准测试和测试，请确保您已安装 `Docker` 和 `make`。

根目录下的 `.env.example` 仅供参考，首次运行时需要 `cp .env.example .env`。

您可以在项目的根目录运行 `make` 查看所有可用命令。

```bash
> make
test       Run tests in Docker
coverage   Generate and open test coverage report
bench      Run benchmarks in Docker
shell      Start a shell in the Docker container
example    Run the basic example
```

