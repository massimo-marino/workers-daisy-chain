// daisy chain of workers
package workersDaisyChain

// generic data to be transmitted between workers
type dataEnvelope interface{}

// the type of the worker function
type workerFun func(wid uint64, inch chan dataEnvelope, outch chan dataEnvelope)

type workerConfig struct {
	// the number of workers
	numOfWorkers uint64
	// the function every worker runs as a goroutine
	wf workerFun
}

// Start the daisy chain of workers passing the number of workers, a workerFun,
// and the data to be sent to the first worker in the chain.
// Return the final data received from the last worker in the chain
func StartDaisyChainOfWorkers(numOfWorkers uint64, worker workerFun, d dataEnvelope) dataEnvelope {
	// set up the worker config structure
	var wConfig workerConfig

	wConfig.numOfWorkers = numOfWorkers
	wConfig.wf = worker

	// channel used to send data to the first worker
	ch0 := make(chan dataEnvelope)
	defer close(ch0)

	// channel used by the last worker to send data to the caller that
	// triggered the daisy chain
	var chf chan dataEnvelope

	var preCh chan dataEnvelope
	var postCh chan dataEnvelope

	preCh = ch0

	// start the daisy chain of workers
	for wid := uint64(1); wid <= wConfig.numOfWorkers; wid++ {
		postCh = make(chan dataEnvelope)
		go wConfig.wf(wid, preCh, postCh)
		preCh = postCh
	}
	// chf is the channel used by the last worker to send data to the caller
	// that triggered the daisy chain
	chf = postCh

	// send the data to the first worker
	ch0 <- d

	// receive and return the final data from the last worker
	return <-chf
}
