# JAKE - Assignment 1 Submission_V2

## Table of Contents 
1. How to run 
2. Prime Counter Explanation and Proof
3. Dining Philosophers Explanation
4. 1.4 - 1.6 Textbook Solutions

## How to Run
Required Go Version >= 1.23.3  

**Prime Finder**
```
go run prime-finder.go
```

**Dining Philosophers**
```
go run dining-philosophers.go
```
## Prime Counter 
**Runtime**  
In this approach, I decided to parallelize the Sieve of Eratosthenes algorithm.  
This program runs in  
$$O(N*\log({log{N}}))$$  
which is significantly faster than my brute force attempt running in  
$$O(N*\sqrt{N})$$  

**Thought Process**  
To parallelize Sieve of Eratosthenes I decided to batch the work evenly among all eight threads. Each thread was responsible for preforming calculations based off the array state.  

One major improvement was to remove all synchronization techniques when performing work.  

 Previously I used a global mutex lock with brute force calculations preformed on each thread (Coarse-grained synchronization attempt). This led to a run time of around 53 seconds to a minute+, on my machine.  

On my next attempt I tried to parallelized the Sieve of Eratosthenes algorithm, but with one major pitfall: leaving the algorithm process sequential, while only distributing the product calculations. This led to the runtimes betting effectively the same as the single threaded version, due to the unnecessary synchronization cost.  

While the last attempt was significantly better than the brute force (due to the more optimized algorithm) I still failed to utilize the power of parallelism. To fix this one key observation had to be made "a little bit of redundant work will take far less time to execute than the barriers" my last solution used.  

This brings me to the current implementation which maximizes the parallelism of this algorithm.  

**Program Breakdown**  
For this solution I decided to batch the threads and allow them to work on individual sections of the array, providing a lock-free algorithm.  
```
thread():
    for batch:
        for ele in batch:
            doPrimeCalc(ele) #do prime executes the sieve of eratosthenes algorithm
```

Once finished, the thread will proceed to wait until all other threads complete their work. From there they will begin to go over their sections of the array again and proceed to count and sum the primes they encounter in local variables. After this process has completed they will atomically add to a shared variable to synchronize results.  

```
thread_with_sum_cnt():
    for batch:
        for ele in batch
            doPrimeCalc(ele)
    
    spin(until other threads finish doPrimeCalc operations)

    for batch:
        for ele in batch:
            if prime: 
                localSum+=ele; localCnt+=1

    atomicAdd(globalSum,localSum)
    atomicAdd(globalCnt,localCnt)
```  

This provided a better runtime then doing these calculations on the main thread once the work threads have finished (like in previous solutions). This allows all parts of this program to be parallelized.  

The last thing that needs to be done is parallelizing the Top 10 Primes Fetch. This is a simple addition to the algorithm: checking if the current thread is the last to finish (meaning the array is in its final state) then grab primes.

```
thread_with_top_sum_cnt():
    for batch:
        for ele in batch:
            doPrimeCalc(ele)

    spin()
    if lastThread:
        grabTopPrimes()

    for batch:
        for ele in batch:
            if prime:
                localSum+=ele; localCnt+=1

    atomicAdd(globalSum,localSum)
    atomicAdd(globalCnt,localCnt)
```

**Conclusion**   
Thus, the algorithm is fully parallelized. The only optimization would be to prevent the spin before the sum and count, but I don't think that is possible without jeopardizing the data's integrity.  

The new solution, on my machine, runs on average of 800ms with the lowest being around 700ms; while the single threaded implementation runs at around 1.3 seconds. Therefore providing roughly a 50% improvement in runtime.   


### Proof
**Proving Workload**  
$$Thread_{work}(T) = 13.75\%$$  
$$Where \sum_{T=1}^8{Thread_{work}(T)} = 100 \%$$  
Thus, displaying an even work load among threads.  

## Dining Philosophers 

**Solution**  
To ensure the program was **dead-lock free** I used the resource hierarchy approach. This makes philosophers pick up the lowest numbered fork first. This breaks the cyclic pattern posed in the problem; thus, eliminating the possibility of a deadlock. 

```
philosopher(id):
    left = (id+PHILO_NUM-1) % PHILO_NUM
    right = id #this is due to the layout of the lock array (see code for why)

    firstFork = min(left,right)
    secondFork = max(left,right)

    forks[firstFork].grab()
    forks[secondFork].grab()
```  

To provide **starvation-freedom** I implemented a queue lock, modeled after chapter 7 array based queue, at the fork level.  

Thus when a philosopher attempts to access a fork it will then sit in a queue lock waiting for its turn. This prevents one thread from starving out the others and allows all threads to eventually eat.  

One important note is to not allow philosophers to queue the second fork if it doesn't have access to the first fork. This is because it brings back the cyclic nature of the problem, which invalidates our resource hierarchy.  

Why?  
Say all threads besides one win the lower number fork. The one thread which didn't win the lower fork goes and enqueues themselves into the higher fork and wins the queue (getting the lock). All threads would then be in a cycle thus a deadlock. So, we must ensure that we only enqueue one fork at a time (in order).  

This also allows other threads to work by preventing threads from holding the lock when they aren't able to use the resource (waiting for the lower fork, but holding the upper fork).  

**Conclusion**  
Since, the queue locks are implemented at the fork level utilizing the resource hierarchy approach we provide dead-lock freedom, starvation-freedom, and effective parallelism (no one global lock)

## 1.4 Solution 
**Light Off**  
The solution to this problem involves marking one of prisoners as a (consumer thread) and the rest as producer threads.

When the producer thread enters the room they should flip the light switch on {if it is off and they haven't turned the light on before}.

When a consumer enters a room if the light switch is on they should flip it off and add to their mental count. If it equals P-1 they should declare free.  

This works because this guarantees that the consumer will count only when new producers enter, since there is no repeated work. Additionally, we are guaranteed to go at least N times, so since there is no limit, we are guaranteed to eventually visit the rooms at least once.

**Light Random**  
The solution is to keep the same producer and consumer thread ratio.  
From there the Producer should  
1. Turn on the light if {it was off and he has turned it on at most once before}.
2. Leave the light on {if light is already on}  

The consumer should:  
1. Turn off the light and count, then keep track of how many times he has turned off the light until he reaches 2*(P-1).

It needs to be 2*(P-1) to account for the random state. If the light is on it will result in a false positive, thus we would need to go through the cycle twice. This will allow us to fliter out the false positive (if state was on) or count the prisoners twice (if state was off).  


## 1.5 Solution 
The solution involves encoding the information about the color order. To achieve this the first prisoner will have to sacrifice themselves since there is no way they could know their color.
The first prisoner should count the number of blue hats (b) and number of red hats (r). From there if R%2==1 say red else say blue.  
From there each prisoner should keep a mental note of the result (A=1 if red else A=0). Everytime the prisoner hears a red hat change A=((A+1)%2).  
Once they get the chance to answer they should count the number of red hats in front of them and respond red if (red-hats-in-front % 2 != A) else blue.  
This works because A keeps the current state of red hats in the list, so if removing from the list breaks that current state then you have a red hat on.  
This allows for an easy encoding algorithm for the prisoners to follow only having to flip a binary bit representing odd or even number of red hats.

## 1.6 Solution 
The question gives us the following information  
$$Parallel=85\%$$  
$$Sequential=15\%$$  
$$MemoryAccess=20\%$$  
$$CacheMiss=\frac{N}{N+10}$$  

Sequential  
$$Freq_{mem}=0.20$$  
$$Cost_{mem}=1$$  
  
Parallel:
$$Freq_{mem}=0.20$$  
$$Cost_{mem}(N)=3N+11$$  
  
With some other given information we can reduce the following  
Both sequential and parallel cost:  
$$Avg = Freq_{mem}*Cost_{mem} + Freq_{other}*Cost_{other}$$  

For Sequential  
$$S_{avg}= 0.20 * 14 + 0.8 * 1 = 3.6$$  
$$S_{total} = Freq*Avg$$  
$$S_{total} = 0.15 * 3.6=0.54$$  

For parallel    
$$P_{avg}(N) = 0.20 * (3N+11) + 0.8 * 1$$  
$$P_{avg}(N) = 0.6N+2.2$$  
$$P_{total}(N) = \frac{P_{avg}(N) * CacheMiss}{N}$$  
$$P_{total}(N) = 0.85*(0.6N+2.2)*(1+\frac{N}{N+10})*\frac{1}{N}$$  
  
For total Cost  
$$Total(N) = S_{total} + P_{total}(N)$$  
$$T(N) = 0.54 + 0.85*(0.6N+2.2)*(1+\frac{N}{N+10})*\frac{1}{N}$$  

Now that we have our T(N) we can analyze this function to find the optimal amount of cores.  
We can do this via two methods  
**Method 1:**  
Find the critical points.  
$${\frac{dT(N)}{dN}}=0$$  
Then take the critical points to find the minimum.  

**Method 2:**  
Plugging T(N) into a graphic calculator and visually analyzing it.  

**Method of choice:**  
I chose method 2, since analyzing T(N) can be algebraically complicated and not lead to the most optimal amount of cores.
This is because it provides the abs minimum and extra work would be required to find the most effective amount. Thus, the simpler option is method 2.  

When plugging T(N) into a graphic calculator we can see the optimal range of cores is between 4 and 8. After 8 cores the performance increase is so marginal that it would be a waste of resources.