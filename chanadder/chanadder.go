package main

import (
	"flag"
	"fmt"
	"time"
)

const Iterations = 1000000

type SharedState struct {
	X     uint64
	Stats map[string]uint64
}

func process(name string, inprocess chan *SharedState, processed chan<- *SharedState) {
	for ss := range inprocess {
		if ss.X >= Iterations {
			processed <- ss
			return
		}
		ss.X = ss.X + 1
		ss.Stats[name] = ss.Stats[name] + 1
		inprocess <- ss
	}
}

func main() {
	var threads = flag.Int("threads", 2, "number of concurrent threads of execution")
	flag.Parse()

	inprocess := make(chan *SharedState)
	processed := make(chan *SharedState)

	go func() {
		inprocess <- &SharedState{0, make(map[string]uint64)}
	}()

	t_start := time.Now()

	for i := 0; i < *threads; i++ {
		go process(fmt.Sprintf("%02d", i), inprocess, processed)
	}

	ss := <-processed

	t_finish := time.Now()

	duration := t_finish.Sub(t_start)
	fmt.Printf("Total iterations: %v; total time: %v, %vns/op\n", Iterations, duration, duration.Nanoseconds()/Iterations)
	for thread, runs := range ss.Stats {
		fmt.Printf("    Thread %v ran: %v times\n", thread, runs)
	}
}
