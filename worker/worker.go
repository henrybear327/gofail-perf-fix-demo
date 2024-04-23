package worker

import (
	"log"
	"sync"
)

func WorkerString(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("WorkerString in")
	defer log.Println("WorkerString out")

	// gofail: var SomeFuncString string
	// log.Println("SomeFuncString", SomeFuncString)
}
