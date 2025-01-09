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
In this approach I decided to use a parallelized Sieve of Eratosthenes.  
This program runs in  
$$O(N*\log({log{N}}))$$  
This approach is significantly faster than my brute force attempt which took previously  
$$O(N*\sqrt{N})$$  

**Thought Process**  
To parallelize Sieve of Eratosthenes I decided to batch the work evenly among all eight threads. Each thread was responsible for preforming calculations based off the array state.  

One major improvement was to remove all synchronization techniques when performing work.  

 Previously I used a global mutex lock with brute force calculations preformed on each thread (Coarse-grained synchronization attempt). This led to a run time of around 53 seconds to a minute on my machine.  

On my next attempt I tried to parallelized Sieve of Eratosthenes algorithm, but with one major pitfall (leaving the algorithm process sequential while only distributing the product calculations). This led to the runtimes betting effectively the same as the single threaded version. This was due to the unnecessary synchronization cost.  

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

Once finished the threads will proceed to wait until all other threads finish. From there they will begin to go over their sections of the array again and proceed to count and sum the primes they encounter in local variables. After this process has completed they will atomically add to a shared variable to synchronize results.  

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

The final change to the algorithm to parallelize would be to count the top 10 primes. This is a simple addition to the algorithm: checking if it is the last thread (meaning the array is in its final state) then count primes.

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

The new solution on my machine runs on average of 800ms with the lowest being around 700m;while the single threaded implementation runs at around 1.3 seconds. Therefore providing roughly a 50% improvement in runtime.   


### Proof
**Proving Workload**  
$$Thread_{work}(T) \approx 14.3\%$$  
$$Where \sum_{T=1}^8{Thread_{work}(T)} \approx 100 \%$$  
Thus, displaying an even work load among threads.  

## Dining Philosophers 


## 1.4 Solution 


## 1.5 Solution 


## 1.6 Solution 

