package main

import (
	"fmt"
	"sync"
	"time"
)

const PRIME_MAX = 100000000 // 10^8
// const PRIME_MAX = 10000
const THREAD_CNT = 8
const THREAD_WORKER_CNT = THREAD_CNT - 1
const BATCH_SIZE = 5

var isPrime = [PRIME_MAX + 1]bool{}
var primeGlobal int64 = 0
var jobWaiter sync.WaitGroup

func sieve() {
	var p int64
	for p = 2; p*p <= PRIME_MAX; p++ {
		if isPrime[p] {
			// work
			primeGlobal = p
			jobWaiter.Add(THREAD_WORKER_CNT)
			jobWaiter.Wait()
		}
	}

	primeGlobal = PRIME_MAX + 1 // set the quit flag
}

// take some bound and calcualte for them
func worker(id int64) {
	for primeGlobal <= PRIME_MAX {
		primeLocal := primeGlobal

		start := primeLocal * primeLocal
		stagger := id * primeLocal

		for i := start + stagger; i <= PRIME_MAX; i += THREAD_WORKER_CNT * primeLocal {
			// fmt.Printf("worker: %v evluating: %v\n", id, i)
			isPrime[i] = false
		}

		jobWaiter.Done()
		for primeLocal == primeGlobal {
		}
	}
	fmt.Println("thread killed", id)
}

func main() {
	for i := 0; i <= PRIME_MAX; i++ {
		isPrime[i] = true
	}

	start := time.Now()
	for i := 0; i < THREAD_WORKER_CNT; i++ {
		go worker(int64(i))
	}
	sieve()
	totalTime := time.Since(start)

	total := 0
	for i := 2; i <= PRIME_MAX; i++ {
		if isPrime[i] {
			// fmt.Printf("%v ", i)
			total += 1
		}
	}
	fmt.Println("Hello world!", total, totalTime)
}
