package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Process(ch chan int, i int) {
	time.Sleep(time.Second)
	for {
		select {
		case <-ch:
			fmt.Printf("Process %d return\n", i)
			return
		}
	}
}

func TestChannel(t *testing.T) {
	channels := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		channels[i] = make(chan int)
		go Process(channels[i], i)
	}

	for _, ch := range channels {
		ch <- 1
		// fmt.Println("Routine ", i, " quit!")
	}
}

func DoWork(wg *sync.WaitGroup, i int) {
	time.Sleep(1 * time.Second)
	fmt.Println("Goroutine ", i, " finished!")
	wg.Done()
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go DoWork(&wg, i)
	}
	wg.Wait()
	fmt.Println("All goroutines finished...")
}
