package workersDaisyChain

import (
	"fmt"
	"strconv"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////
// Grains problem
//
// Calculate the number of grains of wheat on a chessboard given that the number
// on each square doubles.
// There once was a wise servant who saved the life of a prince.
// The king promised to pay whatever the servant could dream up.
// Knowing that the king loved chess, the servant told the king he would like to
// have grains of wheat. One grain on the first square of a chess board. Two grains
// on the next. Four on the third, and so on.
// There are 64 squares on a chessboard.
// Write code that shows:
//   - how many grains were on each square, and
//   - the total number of grains

// data structure exchanged between workers
type grainsData struct {
	local uint64
	total uint64
	msg   string
}

// the worker function run as a goroutine
func grainWorker(wid uint64, inch chan dataEnvelope, outch chan dataEnvelope) {
	defer close(outch)
	var inDataTyped grainsData

	for inData := range inch {
		inDataTyped = inData.(grainsData)

		// prepare the data for the next worker
		outData := grainsData{inDataTyped.local << 1, inDataTyped.total + (inDataTyped.local << 1), inDataTyped.msg + " " + strconv.FormatUint(wid, 10)}

		fmt.Println(wid, "- Received:", inDataTyped.local, inDataTyped.total, "- Sent:", outData)

		// send the data to the next worker
		outch <- outData
	}
}

func TestGrains(t *testing.T) {
	// We cannot use 64 with uint64 because of overflow
	numOfWorkers := uint64(63)
	// put 1 grain on the first square
	d := grainsData{1, 1, "main"}
	r := StartDaisyChainOfWorkers(numOfWorkers, grainWorker, d)
	expected := uint64(1 << numOfWorkers)

	t.Log("Result:", r.(grainsData).local, "Expected:", expected)

	if r.(grainsData).local != expected {
		t.Fatal("Final value is not the expected one")
	}
}

///////////////////////////////////////////////////////////////////////////////
// Fibonacci Numbers
//
// Given n >= 1 compute the fibonacci numbers fib(1), fib(2), ..., fib(n)
// and print all of them.
// Fibonacci numbers are recursively defined as follows:
//
// fib(0) = 0
// fib(1) = 1
// fib(n) = fib(n-1) + fib(n-2), n >= 2
//
// There are n-1 workers named fibw-1. ..., fibw-n-1
// Every fibw-i, i in [1,n-1] receives in input a pair [fib(i-1),fib(i)], writes
// fib(i), and sends to the next worker fibw-i+1 the pair [fib(i), fib(i-1) + fib(i)]
// The worker fibw-1 receives the pair [fib(0), fib(1)]
// The worker fibw-2 receives the pair [fib(1), fib(2)]
// And so on.
// The last worker fibw-n-1 will send the pair [fib(n-1), fib(n)]

// data structure exchanged between workers
type fibData struct {
	fibp uint64
	fibc uint64
}

// the worker function run as a goroutine
func fibWorker(wid uint64, inch chan dataEnvelope, outch chan dataEnvelope) {
	defer close(outch)
	var inDataTyped fibData

	for inData := range inch {
		inDataTyped = inData.(fibData)

		// prepare the data for the next worker
		outData := fibData{inDataTyped.fibc, inDataTyped.fibp + inDataTyped.fibc}

		fmt.Printf("fibw-%d - fib(%v) = %v\n", wid, wid, inDataTyped.fibc)

		// send the data to the next worker
		outch <- outData
	}
}

func TestFib(t *testing.T) {
	// We can compute up to fib(93) when using uint64's
	numOfWorkers := uint64(93)
	// prepare [fib(0), fib(1)]
	d := fibData{0, 1}
	r := StartDaisyChainOfWorkers(numOfWorkers-1, fibWorker, d)
	fib_93Expected := uint64(12200160415121876738)

	fmt.Printf("fibw-%d - fib(%d) = %v\n", numOfWorkers, numOfWorkers, r.(fibData).fibc)

	t.Log("Result:", r.(fibData).fibc, "Expected:", fib_93Expected)

	if r.(fibData).fibc != fib_93Expected {
		t.Fatal("Final value is not the expected one")
	}
}
