package main

import (
	"fmt"
	"sync"
	"time"
)

const PRIME_MAX = 100000000 // 10^8
const THREAD_CNT = 8
const BATCH_SIZE = 10

var nonPrime = [PRIME_MAX + 1]bool{true, true}

func work(id int, wg *sync.WaitGroup) {
	for i := (BATCH_SIZE * id); i <= PRIME_MAX; i += (BATCH_SIZE * THREAD_CNT) {
		for j := 0; j <= BATCH_SIZE && j+i <= PRIME_MAX; j++ {
			k := i + j
			if !nonPrime[k] {
				for z := k * k; z <= PRIME_MAX; z += k {
					nonPrime[z] = true
				}
			}
		}
	}
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
	fmt.Println(time.Since(start))
}
