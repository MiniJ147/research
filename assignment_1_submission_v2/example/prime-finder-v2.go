package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const PRIME_MAX = 100000000 // 10^8
const THREAD_CNT = 8
const BATCH_SIZE = 10
const BATCH_SIZE_CNT = 50

var nonPrime = [PRIME_MAX + 1]bool{true, true}
var totalPrimes int64 = 0
var primeSum int64 = 0
var workDone int32 = 0
var topPrimes = [10]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}

func grabTopPrimes() {
	j := 9
	for i := PRIME_MAX; i >= 0 && j >= 0; i-- {
		if !nonPrime[i] {
			topPrimes[j] = i
			j--
		}
	}
}

func work(id int, wg *sync.WaitGroup) {
	jobCnt := 0
	var sum int64 = 0

	// run through your batches and caclualate primes
	for i := (BATCH_SIZE * id); i <= PRIME_MAX; i += (BATCH_SIZE * THREAD_CNT) {
		for j := 0; j <= BATCH_SIZE && j+i <= PRIME_MAX; j++ {
			k := i + j
			if !nonPrime[k] {
				for z := k * k; z <= PRIME_MAX; z += k {
					nonPrime[z] = true
				}
			}

			jobCnt++
		}
	}

	// count primes

	// if it is the last thread to be modifying the nonPrime Array should begin the grab
	if atomic.AddInt32(&workDone, 1) == THREAD_CNT {
		grabTopPrimes()
	}

	// spin until all modifying threads are done
	for workDone != THREAD_CNT {

	}

	//rerun over the batch and find the primes
	var cnt int64 = 0
	for i := (BATCH_SIZE_CNT * id); i <= PRIME_MAX; i += (BATCH_SIZE_CNT * THREAD_CNT) {
		for j := 0; j <= BATCH_SIZE_CNT && j+i <= PRIME_MAX; j++ {
			if !nonPrime[i+j] {
				cnt++
				sum += int64(i + j)
			}
		}
	}

	fmt.Printf("Thread %d did %d jobs\n", id, jobCnt)
	atomic.AddInt64(&totalPrimes, cnt)
	atomic.AddInt64(&primeSum, sum)
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(THREAD_CNT)
	start := time.Now()
	for i := 0; i < THREAD_CNT; i++ {
		go work(i, &wg)
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Printf("Total Time: %v, total primes: %v, prime sum: %v, top primes %v\n", end, totalPrimes, primeSum, topPrimes)
}
