package main

import (
	"fmt"
	"sync"
	"time"
)

const PHILO_NUM = 5
const SIM_DURATION = 10 * time.Second

type Resource struct {
	mtx            sync.Mutex
	chopsticksFree [PHILO_NUM]bool
}

func philoSimulate(pos int, res *Resource) {
	fmt.Println("Philosopher:", pos)
	for {
		res.mtx.Lock()
		chopLeft, chopRight := pos, (pos+1)%PHILO_NUM

		// can we take both chopsticks if not we don't take any to prevent deadlock
		if !res.chopsticksFree[chopLeft] || !res.chopsticksFree[chopRight] {
			res.mtx.Unlock()
			continue
		}

		//take chopsticks
		res.chopsticksFree[chopLeft] = false
		res.chopsticksFree[chopRight] = false

		res.mtx.Unlock()
		fmt.Printf("Philosopher %d is eating\n", pos)

		time.Sleep(time.Second)

		// release chopsticks
		res.mtx.Lock()
		res.chopsticksFree[chopLeft] = true
		res.chopsticksFree[chopRight] = true
		res.mtx.Unlock()

		// give opprounity for other threads to take
		time.Sleep(25 * time.Millisecond)

	}
}

func main() {
	res := Resource{}
	for i := 0; i < PHILO_NUM; i++ {
		res.chopsticksFree[i] = true
	}

	for i := 0; i < PHILO_NUM; i++ {
		go philoSimulate(i, &res)
	}

	time.Sleep(SIM_DURATION)
}
