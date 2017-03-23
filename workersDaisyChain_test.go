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
		// only the first worker must do this
		if inDataTyped.local == 0 {
			inDataTyped.local = 1
			inDataTyped.total = 1
		}

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
	d := grainsData{0, 0, "main"}
	r := StartDaisyChainOfWorkers(numOfWorkers, grainWorker, d)
	expected := uint64(1 << numOfWorkers)

	t.Log("Result:", r.(grainsData).local, "Expected:", expected)

	if r.(grainsData).local != expected {
		t.Fatal("Final value is not the expected one")
	}
}

///////////////////////////////////////////////////////////////////////////////
