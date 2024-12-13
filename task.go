package workerpool

// Task represents a unit of work that can be executed by the worker pool.
// Implementations must be safe for concurrent execution.
type Task interface {
	// Execute performs the task and returns an error if the task failed.
	Execute() error
}
