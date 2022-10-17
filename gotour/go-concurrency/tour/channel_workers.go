package main

import (
	"fmt"
	"sync"
	"time"
)

func sqrWorker(wg *sync.WaitGroup, tasks <-chan int, results chan<- int ,instance int) {
	for num := range tasks {
		time.Sleep(time.Millisecond)
		fmt.Printf("[worker %v] Sending result by worker %v\n", instance, instance)
		results <- num * num
	}

	wg.Done()
}

// test3
func test3() {
	fmt.Println("[main] main() started")

	var wg sync.WaitGroup

	tasks := make(chan int, 10)
	results := make(chan int, 10)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go sqrWorker(&wg, tasks, results, i)
	}

	for i := 0; i < 5; i++ {
		tasks <- i * 2 // non-blocking as buffer capacity is 10
	}

	fmt.Println("[main] Wrote 5 tasks")

	close(tasks)

	wg.Wait()

	for i := 0; i < 5; i++ {
		result := <-results // non-blocking because buffer is non-empty
		fmt.Println("[main] Result", i, ":", result)
	}	

	fmt.Println("[main] main() stopped")
}
