package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

const PHILO = 5
const SIM_DUR = 5 * time.Second
const LOCK_LEN = 5

var forks [PHILO]QueueLock

type QueueLock struct {
	tail  atomic.Int32
	head  atomic.Int32
	flags [LOCK_LEN]bool
}

func queueLockCreate() QueueLock {
	return QueueLock{
		flags: [LOCK_LEN]bool{true, false},
	}
}

func (l *QueueLock) Lock(id int) {
	idx := (l.tail.Add(1) - 1) % LOCK_LEN
	// fmt.Println(id, "got", idx)
	for !l.flags[idx] {
	}
}

func (l *QueueLock) Unlock(id int) {
	idx := l.head.Add(1)
	// fmt.Println(id, "at", idx-1, "releasing")
	l.flags[(idx-1)%LOCK_LEN] = false
	l.flags[idx%LOCK_LEN] = true
}

func philosopher(id int) {
	for {
		idxLeft := (id + PHILO - 1) % PHILO
		idxRight := id //(id + 1) % PHILO

		fork1 := min(idxLeft, idxRight)
		fork2 := max(idxLeft, idxRight)

		fmt.Println(id, "is hungry taking", fork1, fork2)
		forks[fork1].Lock(id)
		forks[fork2].Lock(id)

		fmt.Println(id, "is eating")
		time.Sleep(time.Millisecond * time.Duration((500 + rand.Intn(2000))))

		fmt.Println(id, "is thinking release", fork1, fork2)
		forks[fork1].Unlock(id)
		forks[fork2].Unlock(id)

		time.Sleep(time.Millisecond * time.Duration((500 + rand.Intn(2000))))
	}
}

func main() {
	for i := range PHILO {
		forks[i] = queueLockCreate()
	}

	for i := range PHILO {
		go philosopher(i)
	}

	time.Sleep(SIM_DUR)
}
