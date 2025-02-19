/*
 * C++ version of the orginal go implementation but with a changed batch size
*/


#include <atomic>
#include <cmath>
#include <chrono>
#include <iostream>
#include <thread>
#include <vector>

const int PRIME_MAX = 100000000;
const int THREAD_NUM = 8;
int BATCH_SIZE = 3;

std::atomic<int> workers = 0;
std::atomic<int> primesFound = 0;
std::atomic<long long> workdone = 0; 
std::atomic<long long> primeSum = 0;

bool nonPrime[PRIME_MAX+1] = {true, true};
long threadsWork[THREAD_NUM];
int topPrimes[10];


void work(int id){
    int work = 0;
    int upperBound = sqrt(PRIME_MAX);
    int primes = 0;
    long long sum = 0;
    int local_batch = BATCH_SIZE;
    int bcn = 0;
    int workid = id;

    for(int i=BATCH_SIZE * workid; i <= upperBound; i += BATCH_SIZE * THREAD_NUM){
        for(int j=0; j<BATCH_SIZE && i+j <= PRIME_MAX; j++){
            int p = i+j;
            // work++;

            if(nonPrime[p]){
                continue;
            }

            for(unsigned long long k=p*p;  k<=PRIME_MAX; k+=p){
                work++;
                nonPrime[k] = true;

            }
        }
    }


    int prev = workers.fetch_add(1);
    int cnt = 9;
    for(int i=PRIME_MAX; i > 0 && prev == THREAD_NUM-1 && cnt >= 0; i--){
        if(!nonPrime[i]){
            topPrimes[cnt] = i;
            cnt--;
        }
    }
    while(workers != THREAD_NUM){}


    for(int i=BATCH_SIZE * id; i <= PRIME_MAX; i += BATCH_SIZE * THREAD_NUM){
        for(int j=0; j<BATCH_SIZE && i+j <= PRIME_MAX; j++){
            if(!nonPrime[i+j]){
                primes++;
                sum += i+j;
            }
        }
    }
     
    // printf("%d %d %d\n",id,work,primes);
    primesFound += primes;
    primeSum += sum;
    workdone += work;
    threadsWork[id] = work;
}

int main(){
    std::vector<std::thread> threads; 
    auto start = std::chrono::high_resolution_clock::now();

    // spawning theads
    for(int i=0; i<THREAD_NUM; i++){
        threads.push_back(std::thread(work,i));
    }

    // waiting on threads
    for(int i=0; i<THREAD_NUM; i++){
        threads[i].join();
    }
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);

    // printing information
    printf("Time Taken: %ld ms Primes Found: %d Prime Sum: %lld Top Primes: ",duration.count(), primesFound.load(), primeSum.load());
    for(int i=0; i<10;i++){
        printf("%d ",topPrimes[i]);
    }
    printf("\n");
    for(int i=0; i<THREAD_NUM; i++){
        printf("%d %ld %.2f\n",i,threadsWork[i], (float)threadsWork[i]/workdone.load());
    }
    printf("%lld\n",workdone.load());
    return 0;
}
