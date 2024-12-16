package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const PHILO_NUM = 5
const SIM_DURATION = 10 * time.Second

type Queue struct {
	data [PHILO_NUM]int
	head int
	size int
}

func (q *Queue) add(val int) {
	if q.size == PHILO_NUM {
		panic("failed add: cannot add to already full queue")
	}

	q.data[(q.head+q.size)%PHILO_NUM] = val
	q.size++
}

func (q *Queue) pop() int {
	if q.size <= 0 {
		panic("failed pop: size is negative or empty")
	}

	val := q.data[q.head]
	q.data[q.head] = 0
	q.head = (q.head + 1) % PHILO_NUM
	q.size--

	return val
}

func (q Queue) index(idx int) int {
	if idx < 0 || idx >= PHILO_NUM {
		panic("index out of bounds")
	}

	return q.data[(q.head+idx)%PHILO_NUM]
}

func (q Queue) find(target int) int {
	for i := 0; i < q.size; i++ {
		res := q.index(i)
		if res == target {
			return i
		}
	}
	return -1
}

// is a in front of b
func (q Queue) isInFront(a int, b int) bool {
	for i := 0; i < q.size; i++ {
		res := q.index(i)
		if res == a {
			return true
		}
		if res == b {
			return false
		}
	}
	return false
}

type Arbitrator struct {
	queue   Queue
	isTaken [PHILO_NUM]bool
	mtx     sync.Mutex
}

func (a *Arbitrator) canEat(id int) bool {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	left := (id + PHILO_NUM - 1) % PHILO_NUM
	right := (id + 1) % PHILO_NUM

	inQueue := a.queue.find(id) != -1
	if a.queue.isInFront(left, id) || a.queue.isInFront(right, id) {
		if !inQueue {
			a.queue.add(id)
		}
		return false
	}

	if a.isTaken[left] || a.isTaken[right] {
		if !inQueue {
			a.queue.add(id)
		}
		return false
	}

	if inQueue && a.queue.find(id) != 0 {
		return false
	}

	if inQueue && a.queue.find(id) == 0 {
		a.queue.pop()
	}

	a.isTaken[id] = true
	return true
}

func (a *Arbitrator) release(id int) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.isTaken[id] = false
}

const (
	STATE_THINKING = "THINKING"
	STATE_HUNGRY   = "HUNGRY"
	STATE_EATING   = "EATING"
)

func main() {
	arb := Arbitrator{}
	for i := range PHILO_NUM {
		go func(id int) {
			state := STATE_THINKING
			for {
				time.Sleep(1000 * time.Millisecond)
				roll := rand.Intn(100)

				if state == STATE_THINKING && roll >= 50 {
					state = STATE_HUNGRY
				}

				if state == STATE_HUNGRY && arb.canEat(id) {
					state = STATE_EATING
				}

				if state == STATE_EATING && roll < 30 {
					arb.release(id)
					state = STATE_THINKING
				}

				fmt.Printf("Philospher: %v is %v\n", id, state)
			}
		}(i)
	}

	time.Sleep(SIM_DURATION)
}
