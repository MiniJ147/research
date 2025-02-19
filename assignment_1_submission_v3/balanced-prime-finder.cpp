#include <atomic>
#include <iostream>
#include <thread>
#include <cmath>
#include <vector>

const int FLAG_KILL = -1;
const int FLAG_IDLE = 0;


const int GROUP_SIZE = 4;
const int WORKERS = 3;
const int TRACERS = 2;
const int THREAD_NUM = TRACERS + TRACERS * WORKERS;

const int STAGE_COUNT = TRACERS;

const int PRIME_MAX = 100000000;

bool nonPrime[PRIME_MAX+1] = {true, true};
int BATCH_SIZE = 3;

int topPrimes[10];
int jobs[TRACERS];
long works[THREAD_NUM];

std::atomic<int> signals[TRACERS];
std::atomic<int> stage;
std::atomic<int> primesFound = 0; 
std::atomic<long long> primeSum = 0;

// count over primes for count and sum values
void count(int id, int& primes, long long& sum){
    for(int i=BATCH_SIZE * id; i <= PRIME_MAX; i += BATCH_SIZE * THREAD_NUM){
        for(int j=0; j<BATCH_SIZE && i+j <= PRIME_MAX; j++){
            if(!nonPrime[i+j]){
                primes++;
                sum += i+j;
            }
        }
    }
}


void worker(int id, int owner, int offset){
    int work = 0;
    int primes = 0;
    long long sum = 0;

    // printf("worker %d\n",id);
    while(jobs[owner]!=FLAG_KILL){
        int local = jobs[owner];

        // not ready to count
        if(local==FLAG_IDLE || local==FLAG_KILL){
            continue;
        }

        // grab our starting point
        unsigned long long curr = local * local;
        curr += local * offset; 

        // count over our batch
        for(;curr<=PRIME_MAX; curr += local * GROUP_SIZE){
            work++;
            nonPrime[curr] = true;
        }

        // signal we are done
        signals[owner]++;
        
        // wait for new job or death signal
        while(local==jobs[owner] && local!=FLAG_KILL);
    }

    // wait for counting stage 
    while(stage!=STAGE_COUNT);
    count(id,primes,sum); 

    // printf("Thread: %d did %d work\n",id, work); 
    works[id] = work;
    primesFound += primes;
    primeSum += sum;

}

void trace(int id){
    int work = 0;
    int primes = 0;
    long long sum = 0;
    int upperBound = sqrt(PRIME_MAX);

    for(int i=BATCH_SIZE * id; i <= upperBound; i += BATCH_SIZE * TRACERS){
        for(int j=0; j<BATCH_SIZE && i+j <= PRIME_MAX; j++){
            int p = i+j;

            if(nonPrime[p]){
                continue;
            }

            jobs[id] = p;
            for(unsigned long long k=p*p;  k<=PRIME_MAX; k+=(p*GROUP_SIZE)){
                work++;
                nonPrime[k] = true;
            }

            // wait for workers to finish to avoid unmarked areas
            while(signals[id]<WORKERS);
            signals[id] = 0;
       }
    }

    // telling workers to stop counting
    jobs[id] = FLAG_KILL;

    // upgrading stage and if we are last thread we will grab top primes
    int prev = stage.fetch_add(1);
    int cnt = 9;
    for(int i=PRIME_MAX; i > 0 && prev == TRACERS-1 && cnt >= 0; i--){
        if(!nonPrime[i]){
            topPrimes[cnt] = i;
            cnt--;
        }
    }
    while(stage!=STAGE_COUNT);
    count(id,primes,sum); 

    // pushing values
    works[id] = work;
    primesFound += primes;
    primeSum += sum;
}



int main(){
    std::vector<std::thread> threads;

    auto start = std::chrono::high_resolution_clock::now();

    // spawning threads
    for(int i=0; i<TRACERS; i++){
        for(int j=0; j<WORKERS;j++){
            int id = (TRACERS * i)+j+(1*i)+TRACERS; // dont worry just weird math to get to cycle ids 
            threads.push_back(std::thread(worker,id,i,j+1));
        }
    }

    
    for(int i=0; i<TRACERS; i++){
        threads.push_back(std::thread(trace,i));
    }

    // waiting for threads
    for(int i=0; i<threads.size(); i++){
        threads[i].join();
    }
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);

    //printing information 
    printf("\nTime Taken: %ld ms Primes Found: %d Prime Sum: %lld Top Primes: ",duration.count(), primesFound.load(), primeSum.load());
    for(int i=0; i<10;i++){
        printf("%d ",topPrimes[i]);
    }
    printf("\n\n");

    // getting workload
    long workdone = 0;
    for(int i=0; i<THREAD_NUM; i++){
        workdone += works[i];
    }
    for(int i=0; i<THREAD_NUM; i++){
        printf("Thread: did %d %ld work / %.2f%%\n",i,works[i], (float)works[i]/workdone);
    }
    printf("total work: %ld\n\n",workdone);

    return 0;
}
