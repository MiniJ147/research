# JAKE - Assignment 1 Submission_V2

## Table of Contents 
1. How to run 
2. Prime Counter Explanation and Proof
3. Dining Philosophers Explanation
4. 1.4 - 1.6 Textbook Solutions

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
**program breakdown**  
In this approach I decided to parallelize the Sieve of Eratosthenes algorithm.  
$$O(N*loglogN);N=10^8$$  
Previous I used a brute force attempt with   
$$O(N\sqrt{N})$$  
which we will see provided a massive improvement in efficiency.  

To parallelize Sieve of Eratosthenes I decided to have prime tracking thread called sieve, and 7 other worker threads.  
The job of the worker threads was to evenly split up the workload of calculating the products of whatever prime the algorithm is on.  
The job of the sieve thread is to iterate through the array and track the next available prime.  

The reason for this implementation is that Sieve of Eratosthenes is a sequential algorithm, meaning that the results of k+1 depend on the calculation of k. Therefore, if we parallelize the algorithm itself we would lose the data integrity of the boolean array.  

To get around this I decided to parallelize the multiple calculation part of the algorithm (worker threads) and leave the next prime checker sequential (Sieve thread). This keeps the data integrity of the boolean array since k-1 will always be fully calculated.  

Furthermore, what allowed me to parallelize the multiple calculations is the fact that it does not matter which order it is completed in, meaning, we could batch it among the threads.  

To ensure that the threads were synchronized I used a barrier technique.  

**algorithm run down**
```
KILL = -1
currPrime = 0

Sieve thread():
    for num in 10^8:
        if isPrime[num]:
            currPrime = num

            barrier.add(worker_threads)
            barrier.wait() // wait until all threads finished

    currPrime = KILL

Worker thread():
    while currPrime != KILL:
        localPrime = currPrime

        for num in batch:
            isPrime[num] = false

        barrier.done(1) 
        wait_until(localPrime!=currPrime) // detect new prime 
```

The barrier ensures that the calculation is finished before we can move onto the next prime.  

To detect if we have received a new prime we can compare the global prime (current) with our local prime (previous). If they are different we know we have a new prime to calculate.  

The batch is necessary to ensure each thread is doing unique work (no duplicates). This is done through an interval technique where  
Thread 1: [1,100)  
Thread 2: [100,200)  
and so on... 

With this approach the work is completed in around 1 second on my machine. With the lowest value taking less than a second.  

This is an 76x approvement over my last submission.

### Proof
**Proving Workload**  
$$Thread_{work}(T) \approx 14.3\%$$  
$$Where \sum_{T=1}^8{Thread_{work}(T)} \approx 100 \%$$  
Thus, displaying an even work load among threads.  

**Proving Efficiency**  
$$Thread_{single} \approx Time_{avg} * Thread_{total} = 1 * 7 = 7sec$$

## Dining Philosophers 


## 1.4 Solution 


## 1.5 Solution 


## 1.6 Solution 

