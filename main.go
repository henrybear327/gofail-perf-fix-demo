package main

import (
	"log"
	"sync"
	"time"

	worker "example.com/m/v2/worker"
	gofail "go.etcd.io/gofail/runtime"
)

// cd worker && gofail enable && cd .. && go build . && cd worker && gofail disable && cd .. && ./m
func main() {
	// goal: try to see how failpoint is blocking

	{
		// expectation: this part of the code will take about 3s to execute only
		log.Println("Stage 1: Run 3 workers under normal logic")
		var wg sync.WaitGroup
		wg.Add(1)
		go worker.Worker1(&wg)

		wg.Add(1)
		go worker.Worker2(&wg)

		wg.Add(1)
		go worker.Worker3(&wg)

		wg.Wait()
		log.Println("Stage 1: Done")
	}

	{
		// expectation: this part of the code will take about 6s to execute only, if gofail is non-blocking. Otherwise, about 12s
		log.Println("Stage 2: Run 3 workers under gofail logic")
		/*
			ISSUE: the gofail implementation up till commit 93c579a86c46 will be executing the program sequentially

			Due to the execution and enable/disable flows are under the same locking mechanism, only one of the actions can make progress at a given moment
		*/

		var wg sync.WaitGroup
		gofail.Enable("worker1Failpoint", `sleep("3s")`)
		wg.Add(1)
		go worker.Worker1(&wg)
		time.Sleep(10 * time.Millisecond)

		gofail.Enable("worker2Failpoint", `sleep("3s")`)
		wg.Add(1)
		go worker.Worker2(&wg)
		time.Sleep(10 * time.Millisecond)

		gofail.Enable("worker3Failpoint", `sleep("3s")`)
		wg.Add(1)
		go worker.Worker3(&wg)
		time.Sleep(10 * time.Millisecond)

		wg.Wait()
		gofail.Disable("worker1Failpoint")
		gofail.Disable("worker2Failpoint")
		gofail.Disable("worker3Failpoint")
		log.Println("Stage 2: Done")
	}
}
