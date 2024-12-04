package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

const THREAD_COUNT = 8

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// func (h IntHeap) Swap(i,j,in

type Result struct {
	mtx        sync.Mutex
	primeSum   int64
	primeTotal int
	primeNext  int
	heap       IntHeap
}

func (r *Result) inc(prime int) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.primeTotal >= 10 {
		heap.Pop(&r.heap)
	}
	heap.Push(&r.heap, prime)

	r.primeSum += int64(prime)
	r.primeTotal += 1
}

func (r *Result) nextPrime() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	next := r.primeNext
	r.primeNext++

	return next
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func main() {
	var MAX = int(math.Pow10(8))
	r := Result{primeSum: 0, primeTotal: 0, primeNext: 1}

	var wg sync.WaitGroup

	doTest := func() {
		jobCnt := 0
		for {
			i := r.nextPrime()
			if i > MAX {
				break
			}
			jobCnt++

			if isPrime(i) {
				r.inc(i)
			}
		}
		fmt.Println("Jobs from thread", jobCnt)
		wg.Done()
	}

	wg.Add(THREAD_COUNT)
	start := time.Now()
	for i := 0; i < THREAD_COUNT; i++ {
		go doTest()
	}
	wg.Wait()
	fmt.Printf("%v %v %v\n", time.Since(start), r.primeTotal, r.primeSum)
	sort.Ints(r.heap)
	fmt.Println(r.heap)
}
