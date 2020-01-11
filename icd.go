// Package icd provides the interface which plugins must implement
// in order to be used within the reservoird framework. Each plugin
// must provide a 'New' function with the following function signature:
//
//	Queues:
//		New(cfg string, monitor *icd.Monitor) (icd.Queue, error)
//
//	Ingesters:
//		New(cfg string, flow *icd.Flow, monitor *icd.Monitor) (icd.Ingester, error)
//
//	Digesters:
//		New(cfg string, flow *icd.Flow, monitor *icd.Monitor) (icd.Digester, error)
//
//	Expellers:
//		New(cfg string, flow *icd.Flow, monitor *icd.Monitor) (icd.Expeller, error)
//
// Reservoird will not start plugins without the New function as
// defined.
package icd

import (
	"sync"
)

// Flow provides channels and control for the flow threads
type Flow struct {
	// The channel to receive the done message and initiate a graceful shutdown
	doneChan chan struct{}
	// Call 'defer waitGroup.Done()' on flow function start. Reservoird
	// uses this variable to wait for all threads to stop before exiting
	wg *sync.WaitGroup
}

// Monitor provides channels and control for the monitor threads
type Monitor struct {
	// The channel to send statistics messages
	statsChan chan string
	// The channel to receive the clear message to clear statistics
	clearChan chan struct{}
	// The channel to report error messages
	errorChan chan error
	// The channel to receive the done message and initiate a graceful shutdown
	doneChan chan struct{}
	// Call 'defer waitGroup.Done()' on monitor function start. Reservoird
	// uses this variable to wait for all threads to stop before exiting
	wg *sync.WaitGroup
}

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
	Monitor()
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
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Ingest.
	Monitor()
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
		// The queue which data is received from
		recvQueue Queue,
		// The queue which data is forwarded through
		sendQueue Queue,
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Digest.
	Monitor()
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
		// The queue(s) which data is received from
		recvQueues []Queue,
	)

	// Monitor provides a method for sending statistics messages
	// and for receiving the clear statistisics message and the done message.
	//
	// NOTE: monitor runs in a separate thread from Expel.
	Monitor()
}
