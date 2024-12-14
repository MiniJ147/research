package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const PRIME_MAX = 100000000 // 10^8
// const PRIME_MAX = 10000
const THREAD_CNT = 8
const THREAD_BATCH_SIZE = 50 

var isPrime = [PRIME_MAX + 1]bool{}
var primeGlobalCnt int32 = 0
var primeGlobalSum int64 = 0

func primePrintTop10() {
	i := PRIME_MAX
	cnt := 0
	primes := []int{}

	for cnt < 10 {
		if isPrime[i] {
			primes = append(primes, i)
			cnt += 1
		}
		i -= 1
	}

	sort.Ints(primes)
	fmt.Println(primes)
}

func primeValidate(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 || !isPrime[n] {
			return false
		}
	}

	for i := n * 2; i <= PRIME_MAX; i += n {
		isPrime[i] = false
	}

	return true
}

func worker(start int, wg *sync.WaitGroup) {
	defer wg.Done()

	primeCurr := start
	var primeLocalCnt int32 = 0
	var primeLocalSum int64 = 0
	jobCnt := 0

	for primeCurr <= PRIME_MAX {
		for i := primeCurr; i < primeCurr+THREAD_BATCH_SIZE && i <= PRIME_MAX; i++ {
			jobCnt++

			if primeValidate(i) {
				// fmt.Printf("thread %d found %d as prime\n", start, i)
				primeLocalCnt += 1
				primeLocalSum += int64(i)
			}
		}

		// fmt.Printf("thread: %d, start: %d, end: %d, primes found: %d\n", start, primeCurr, primeCurr+THREAD_BATCH_SIZE, primeLocalCnt)
		primeCurr += (THREAD_BATCH_SIZE * THREAD_CNT)
	}

	fmt.Printf("thread: %d, found: %d primes, did: %d jobs. Local sum: %d\n", start, primeLocalCnt, jobCnt, primeLocalSum)
	atomic.AddInt32(&primeGlobalCnt, primeLocalCnt)
	atomic.AddInt64(&primeGlobalSum, primeLocalSum)
}

func main() {
	for i := 0; i < PRIME_MAX; i++ {
		isPrime[i] = true
	}

	wg := sync.WaitGroup{}
	wg.Add(THREAD_CNT)

	start := time.Now()

	for i := 0; i < THREAD_CNT; i++ {
		go worker(i*THREAD_BATCH_SIZE, &wg)
	}
	wg.Wait()

	fmt.Printf("Total Time: %v | Prime Count: %v | Prime Sum: %v | greatest 10", time.Since(start), primeGlobalCnt, primeGlobalSum)
	primePrintTop10()
}
