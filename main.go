package main

import (
	"sync"

	worker "example.com/m/v2/worker"
	gofail "go.etcd.io/gofail/runtime"
)

// cd worker && gofail enable && cd .. && go build . && cd worker && gofail disable && cd .. && ./m
func main() {
	var wg sync.WaitGroup

	// before enabling failpoint
	wg.Add(1)
	go worker.WorkerString(&wg)
	wg.Wait()

	// after enabling failpoint
	gofail.Enable("SomeFuncString", `return("Hello_world")`)

	wg.Add(1)
	go worker.WorkerString(&wg)
	wg.Wait()

	gofail.Disable("SomeFuncString")
}
