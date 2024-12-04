# JAKE - Assignment 1 Submission

## How to Run
Required Go Version 1.23.3  

**Prime Finder**
```
go run prime-finder.go
```

**Dining Philosophers**
```
go run dining-philosophers.go
```

## Prime Counter 

**Program breakdown**  
To calculate primes I choose to use an optimized brute force approach where running each findPrime calculation takes 
$$O(N\sqrt{N}); N=10^8$$

At program start 8 threads will spawn where there is a global mutex shared among the threads. This mutex is in charge of keeping track of state and providing the next prime for the threads. The use of the mutex is important here so we can prevent race conditions.  

Each thread will get a unique number that is given by the sequential increment Mutex for each request. This thread will be required to check if this given number is in fact a prime number. If it is it will go into the Mutex and update its state.  

To keep track of the 10 highest primes I use a max heap which stores a max of 10 elements at a time. This allows for constant time tracking of the most recent 10 highest primes.

### Proof

$$Thread_{work}(i) \approx 12.5\%$$  
$$Where \sum_{i=1}^8{Thread_{work}(i)} = 100\%$$  
$$Time_{avg} \approx 75sec$$
$$Where\;the\;work\;of,\\ Thread_{single} \approx Time_{avg} * Thread_{total} = 75 * 8 = 600sec$$  
**Thus, we can see an 8x increase in performance and an even distribution of workload among threads** 


### Potential Improvements  

One major bottleneck is the prime counter algorithm. This brute force approach leads to a lot of repeated work and research into a more effective algorithm would be beneficial. One algorithm I dived into was Sieve of Eratosthenes, however this requires a shared memory pool among threads. Initially this implementation seemed overly complex, but upon reflection it seems rather trivial. Thus, if a more efficient solution is required it would be beneficial to implement this algorithm with the added cost of complexity and memory usage.

**Prime-Finder Runtime Improvement**   
$$O_{bruteforce}(N\sqrt{N}) > O_{optimized}(N\log{\log{N}})$$

## Dining Philosophers

**Solution**  
To ensure that this program was **dead-lock free** I forced philosophers to pick up both chopsticks at a time. (prevents from holding one chopstick)

```
Algorithm (chopstickLogic)

if chopstick_left == FREE && chopstick_right == FREE
    PICK UP BOTH
else
    WAIT
```

To prevent race conditions I made sure that I locked the global mutex which represented the chopsticks state. This would ensure that no two threads would attempt to take both chopsticks at the same time, given that multiple threads can have this given state.

```
Algorithm preventing race conditions

mutex.lock()
chopstickLogic()
mutex.unlock()
```

**Issues with starvation freedom**  
My initial solution for starvation freedom was to keep a priority queue of wait times per thread. Thus, the most starved thread would get access to the resource first. However, I had no idea how to implement it in a concurrent way.  
Upon reflection I realized my implemented solution was not starvation free due to numerous key assumptions being made
1. each thread uses resource for fixed time
2. each thread waits for fix time 

This results in a round robin effect where one set of philosophers get to eat and then the next opposite set gets to eat.

A proper solution would be
```
Proper Algorithm 

mutex: WaitTimes[id] = (waitTime)

// assume mutex calls are correctly made and code handles id check properly via mod
ThreadEat(){
    id = threadId
    currentWait = WaitTimes[id]
    if  currentWait >= WaitTimes[id-1] AND currentWait >= WaitTimes[id+1]
        EAT
    else
        WAIT
}

```

This algorithm provides starvation freedom as it ensures the most starved threads go first. Once a thread request a resource it will check its neighboring threads to ensure that it is the most starved. If it is not the most starve it will ceded to that thread. This is critical in a concurrent environment, as we cannot just simply sequentially go through the most starved to the least. This is because each thread will attempt to access the resource at a different time. Thus, having this conditional check allows this program to be starvation free.