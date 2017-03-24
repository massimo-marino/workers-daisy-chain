# Daisy Chain of Workers

**Files:** *workersDaisyChain.go, workersDaisyChain_test.go*

A set of goroutines (workers) is created in a daisy chain, so that every worker W~i~ has an input channel form worker W~i-1~ and an output channel to worker W~i+1~.

The caller that triggers the chain of workers must call `StartDaisyChainOfWorkers()`

```go
// Start the daisy chain of workers passing the number of workers, a workerFun,
// and the data to be sent to the first worker in the chain.
// Return the final data received from the last worker in the chain
func StartDaisyChainOfWorkers(numOfWorkers uint64, worker workerFun, d dataEnvelope) dataEnvelope
```
passing:

- `numOfWorkers` the number of concurrent workers
- `worker` the worker function that every worker must run
- `d` the data that must be sent to the first worker W~1~

The worker function is defined as follows:
```go
type workerFun func(wid uint64, inch chan dataEnvelope, outch chan dataEnvelope)
```
where:

- `wid` is the worker id in `[1,numOfWorkers]`
- `inch` is the input channel from worker W~i-1~
- `outch` is the output channel to worker W~i+1~

In general, the worker function will:

- read data from the input channel
- process the data
- write the processed data to the output channel for the next worker


### The Examples

**File:** *workersDaisyChain_test.go*

#### Grains
Calculate the number of grains of wheat on a chessboard given that the number on each square doubles.
There once was a wise servant who saved the life of a prince.
The king promised to pay whatever the servant could dream up.
Knowing that the king loved chess, the servant told the king he would like to have grains of wheat. One grain on the first square of a chess board. Two grains on the next. Four on the third, and so on.
There are 64 squares on a chessboard.
Write code that shows:

- how many grains were on each square, and
- the total number of grains

#### Fibonacci Numbers
Given `n >= 1` compute the fibonacci numbers `fib(1), fib(2), ..., fib(n)` and print all of them.  
Fibonacci numbers are recursively defined as follows:
``
fib(0) = 0  
fib(1) = 1  
fib(n) = fib(n-1) + fib(n-2), n >= 2  
``

There are n workers named `fibw-1. ..., fibw-n`  
Every `fibw-i`, `i` in `[1,n]` receives in input a pair `[fib(i-1),fib(i)]`, writes
`fib(i)`, and sends to the next worker `fibw-i+1` the pair `[fib(i), fib(i-1) + fib(i)]`  
The worker `fibw-1` receives the pair `[fib(0), fib(1)]`  
The worker `fibw-2` receives the pair `[fib(1), fib(2)]`  
And so on.  
The last worker `fibw-n` will send the pair `[fib(n-1), fib(n)]`  
