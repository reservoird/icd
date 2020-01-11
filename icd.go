// Package icd provides the interface which plugins must implement
// in order to be used within the reservoird framework. Each plugin
// must provide a 'New' function with the following function signature:
//
//	Queues:
//		New(cfg string) (icd.Queue, error)
//
//	Ingesters:
//		New(cfg string) (icd.Ingester, error)
//
//	Digesters:
//		New(cfg string) (icd.Digester, error)
//
//	Expellers:
//		New(cfg string) (icd.Expeller, error)
//
// Reservoird will not start plugins without the New function as
// defined.
package icd

import (
	"sync"
)

// MonitorControl contain what is needed to monitor and control of reservoird threads
type MonitorControl struct {
	// The channel to send statistics messages
	StatsChan chan interface{}
	// The channel to receive the clear message to clear statistics
	ClearChan chan struct{}
	// The channel to report error messages
	ErrorChan chan error
	// The channel to receive the done message and initiate a graceful shutdown
	DoneChan chan struct{}
	// Call 'defer WaitGroup.Done()' on function start. Reservoird
	// uses this variable to wait for all threads to stop before exiting
	WaitGroup *sync.WaitGroup
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

	// Monitor provides monitoring of queue
	Monitor(
		// Provides the monitor and control
		mc *MonitorControl,
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
		snd Queue,
		// Provies monitor and control
		mc *MonitorControl,
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
		// The queue which data is received from
		rcv Queue,
		// The queue which data is forwarded through
		snd Queue,
		// Provides monitor and control
		mc *MonitorControl,
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
		// The queue(s) which data is received from
		rcv []Queue,
		// Provides monitor and control
		mc *MonitorControl,
	)
}
