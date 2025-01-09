package main

import (
	"fmt"
	"math/rand"
	// "math/rand"
	"sync/atomic"
	"time"
)

const (
	PHILO     = 4 
	LOCK_SIZE = 2
	SIM_DUR   = time.Second * 10

	STATE_THINKING = "THINKING"
	STATE_HUNGRY   = "HUNGRY"
	STATE_EATING   = "EATING"
)

var forks [PHILO]QueueLock

type QueueLock struct {
	head atomic.Int32
	tail atomic.Int32
	flag [LOCK_SIZE]bool
}

func (l *QueueLock) lock() {
	slot := (l.tail.Add(1) - 1) % LOCK_SIZE
	for !l.flag[slot] {

	}
}

func (l *QueueLock) UnLock() {
	slot := l.head.Load()
	l.flag[slot%LOCK_SIZE] = false
	l.flag[(slot+1)%LOCK_SIZE] = true
	l.head.Add(1)
}

func newLock() QueueLock {
	return QueueLock{
		flag: [2]bool{true, false},
	}
}

func main() {
	for i := 0; i < PHILO; i++ {
		forks[i] = newLock()
	}

	for i := 0; i < PHILO; i++ {
		go func() {
			state := STATE_THINKING
			for {
				time.Sleep(time.Second)
				roll := rand.Intn(100)

				if state == STATE_THINKING && roll >= 50 {
					fmt.Printf("Philospher: %v is now hungry\n", i)
					left := (i + PHILO - 1) % PHILO
					right := (i + 1) % PHILO

					forks[left].lock()
					forks[right].lock()
					fmt.Printf("Philospher: %v is eating\n", i)
					time.Sleep(time.Millisecond * time.Duration((500 + rand.Intn(1500))))
					fmt.Printf("Philospher: %v stopped eating\n", i)
					forks[left].UnLock()
					forks[right].UnLock()
				}
				// fmt.Printf("Philospher: %v is thinking\n", i)
			}
		}()
	}
	// f := func(id int) {
	// 	fmt.Println(id, "locking")
	// 	forks[0].lock()
	// 	fmt.Println(id, "in lock")
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println(id, "unlocking")
	// 	forks[0].UnLock()
	// }
	//
	// go f(1)
	// go f(2)
	time.Sleep(SIM_DUR)
}
