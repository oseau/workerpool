package main

import (
	"log"

	"github.com/oseau/workerpool"
)

func main() {
	log.Println("Example Basic Worker Pool starting")
	// Create a new pool with 4 workers and 8 queue size
	pool := workerpool.New(
		workerpool.WithPoolSize(4),
		workerpool.WithQueueSize(8),
	)

	// gracefully stop the pool when the program exits
	defer pool.Stop()

	// Add some send mail tasks
	log.Println("Adding send mail tasks")
	for i := 0; i < 5; i++ {
		if err := pool.Add(&TaskSendMail{}); err != nil {
			log.Printf("Failed to add send mail task %d: %v\n", i, err)
		}
	}
	log.Println("All send mail tasks submitted")

	pool.Wait() // Wait for tasks to complete
	log.Println("All send mail tasks completed")

	// We can add more tasks if needed, add some send sms tasks
	log.Println("Adding more tasks, send sms for demo purposes")
	for i := 0; i < 8; i++ {
		if err := pool.Add(&TaskSendSMS{}); err != nil {
			log.Printf("Failed to add send sms task %d: %v\n", i, err)
		}
	}
	log.Println("All send sms tasks submitted")

	pool.Wait() // Wait for tasks to complete
	log.Println("All send sms tasks completed")
}
