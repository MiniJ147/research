package main

import (
	"fmt"
	"sync"
	"time"
)

const PRIME_MAX = 100000000 // 10^8
const THREAD_CNT = 8
const THREAD_WORKER_CNT = THREAD_CNT - 1
const THREAD_QUIT_FLAG = -1

var isPrime = [PRIME_MAX + 1]bool{}
var primeGlobal int64 = 0
var jobWaiter sync.WaitGroup

// sequential sieve of eratosthenes alogrithim
func sieve() {
	var p int64
	for p = 2; p*p <= PRIME_MAX; p++ {
		if isPrime[p] {
			primeGlobal = p

			jobWaiter.Add(THREAD_WORKER_CNT)
			jobWaiter.Wait()
		}
	}

	primeGlobal = THREAD_QUIT_FLAG // set the quit flag
}

// take some bound and calcualte for them
func worker(id int64) {
	jobCnt := 0

	for primeGlobal != THREAD_QUIT_FLAG {
		primeLocal := primeGlobal

		start := primeLocal * primeLocal
		stagger := id * primeLocal // stagger for batch work

		for i := start + stagger; i <= PRIME_MAX; i += THREAD_WORKER_CNT * primeLocal {
			isPrime[i] = false
			jobCnt++
		}

		jobWaiter.Done()

		// spin until new prime is found
		for primeLocal == primeGlobal {
		}
	}
	fmt.Printf("thread %v killed and completed %v jobs\n", id, jobCnt)
}

func printTop10Primes() {
	list := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	idx := 9
	for i := PRIME_MAX; i > 0 && idx >= 0; i-- {
		if isPrime[i] {
			list[idx] = i
			idx--
		}
	}

	fmt.Println(list)
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

	primeTotal := 0
	var primeSum int64 = 0

	for i := 2; i <= PRIME_MAX; i++ {
		if isPrime[i] {
			primeSum += int64(i)
			primeTotal += 1
		}
	}

	fmt.Printf("Total Primes: %v | Prime Sum: %v | Total Time: %v | Top 10 Primes: ", primeTotal, primeSum, totalTime)
	printTop10Primes()
}
