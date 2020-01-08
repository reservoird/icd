// Package icd provides the interface which plugins must implement
// in order to be used within the reservoird framework. Each plugin
// must provide a 'New' function with the following function signature:
//
//   New(cfg string) (icd.<interface>, error)
//
// Reservoird will not start plugins without the New function as
// defined.
package icd

import (
	"sync"
)

// Queue is the inteface for the reservoird queue plugin type.
// This plugin provides the means for communication between
// the ingester, digester, and expeller reservoird plugin
// types.
type Queue interface {
	// Name provides the name of the queue
	Name() string

	// Put puts an item into the the queue
	Put(interface{}) error

	// Get gets the next item from the queue
	Get() (interface{}, error)

	// Len returns the number of items in the queue
	Len() int

	// Cap returns the maximum number of items the queue can hold,
	// if unbounded return -1
	Cap() int

	// Clears the queue, i.e. Len() = 0
	Clear()

	// Close closes the queue, no longer usable
	Close() error

	// Closed returns whether or not the queue is closed
	Closed() bool

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from queue access functions.
	Monitor(
		// The channel to send statistics messages
		statsChan chan<- string,
		// The channel to receive the clear message to clear statistics
		clearChan <-chan struct{},
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)
}

// Ingester is the inteface for the reservoird ingester plugin type. This
// plugin type ingests data from a data source and forwards that data through
// the queue for further processing.
//
// This is considered the input into reservoird.
type Ingester interface {
	// Name returns the name of the ingest plugin
	Name() string

	// Running returns whether or not ingest is running
	Running() bool

	// Ingest is a long running function which captures and forwards data
	// through the queue for further processing.
	Ingest(
		// The queue which data is forwarded through
		sendQueue Queue,
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Ingest.
	Monitor(
		// The channel to send statistics messages
		statsChan chan<- string,
		// The channel to receive the clear message to clear statistics
		clearChan <-chan struct{},
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)
}

// Digester is the inteface for the reservoird digester plugin type. This
// plugin type digests data from an ingester queue and forwards that
// data through the another queue for further processing.
type Digester interface {
	// Name provides the name of the digest plugin
	Name() string

	// Running returns whether or not digest is running
	Running() bool

	// Digest is a long running function which captures data from one queue,
	// processes the data, then forwards the processed data through
	// another queue for further processing.
	//
	// This is considered the filters, annotations, and transformation
	// within reservoird.
	Digest(
		// recvQueue: The queue which data is received from
		recvQueue Queue,
		// The queue which data is forwarded through
		sendQueue Queue,
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Digest.
	Monitor(
		// The channel to send statistics messages
		statsChan chan<- string,
		// The channel to receive the clear message to clear statistics
		clearChan <-chan struct{},
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)
}

// Expeller is the inteface for the reservoird expeller plugin type. This
// plugin type receives data from a queue and expels the data outside
// of reservorid.
//
// This is considered the output of reservoird.
type Expeller interface {
	// Name provides the name of the expeller plugin
	Name() string

	// Running returns whether or not expel is running
	Running() bool

	// Expeller is a long running function which captures data from one queue,
	// processes, and then forwards data through another queue for
	// further processing.
	Expel(
		// The queues which data is received from
		recvQueues []Queue,
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Expel.
	Monitor(
		// The channel to send statistics messages
		statsChan chan<- string,
		// The channel to receive the clear message to clear statistics
		clearChan <-chan struct{},
		// The channel to receive the done message and initiate a graceful shutdown
		doneChan <-chan struct{},
		// Call 'defer waitGroup.Done()' on function start. Reservoird
		// uses this variable to wait for all threads to stop before exiting
		waitGroup *sync.WaitGroup,
	)
}
