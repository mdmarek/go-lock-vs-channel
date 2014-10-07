package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

const Iterations = 1000000

type SharedState struct {
	X     uint64
	L     sync.Mutex
	W     sync.WaitGroup
	Stats map[string]uint64
}

func process(name string, ss *SharedState) {
	ss.L.Lock()
	defer ss.L.Unlock()
	ss.X = ss.X + 1
	ss.Stats[name] = ss.Stats[name] + 1
}

func processLoop(name string, ss *SharedState) {
	for ss.X < Iterations {
		process(name, ss)
	}
	ss.W.Done()
}

func main() {
	var threads = flag.Int("threads", 2, "number of concurrent threads of execution")
	flag.Parse()

	ss := &SharedState{0, sync.Mutex{}, sync.WaitGroup{}, make(map[string]uint64)}
	ss.W.Add(*threads)

	t_start := time.Now()

	for i := 0; i < *threads; i++ {
		go processLoop(fmt.Sprintf("%02d", i), ss)
	}

	ss.W.Wait()

	t_finish := time.Now()

	duration := t_finish.Sub(t_start)
	fmt.Printf("Total iterations: %v; total time: %v, %vns/op\n", Iterations, duration, duration.Nanoseconds()/Iterations)
	for thread, runs := range ss.Stats {
		fmt.Printf("    Thread %v ran: %v times\n", thread, runs)
	}
}
