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

	// Put puts an item into the queue
	Put(interface{}) error

	// Get gets the next item from the queue
	Get() (interface{}, error)

	// Peek peeks at the next item in the queue
	Peek() (interface{}, error)

	// Len returns the length of the queue
	Len() int

	// Cap returns the capacity of the queue, if unbounded return -1
	Cap() int

	// Clears the queue
	Clear()

	// Close closes the queue, no longer usable
	Close() error

	// Closed returns whether or not the queue is closed
	Closed() bool

	// Returns the marshalled stats
	//
	// NOTE: Stats functions run in a different thread than
	// queue access therefore access to stats MUST be thread safe
	Stats() (string, error)

	// Clears stats
	//
	// NOTE: Stats functions run in a different thread than
	// queue access therefore access to stats MUST be thread safe
	ClearStats()
}

// Ingester is the inteface for the reservoird ingester plugin type. This
// plugin type ingests data from a data source and forwards that data through
// the queue for further processing. This is the source points
type Ingester interface {
	// Name provides the name of the ingest plugin
	Name() string

	// Ingest is a long running function which captures and forwards data
	// through the queue for further processing
	//
	// outQueue: The queue which data is forwarded through
	// doneChan: The channel used to gracefully stop the long running function.
	//     If data is present on this channel initiate graceful shutdown
	// waitGroup: Call 'waitGroup.Done()' on function start. Reservoird uses
	//     this variable to wait for all threads to stop before exiting
	// error: Returns and error if there is an issue.
	Ingest(outQueue Queue, doneChan <-chan struct{}, waitGroup *sync.WaitGroup) error

	// Returns marshalled stats
	//
	// NOTE: Stats functions run in a different thread than
	// ingest function therefore access to stats MUST be thread safe
	Stats() (string, error)

	// Clears stats
	//
	// NOTE: Stats functions run in a different thread than
	// ingest function therefore access to stats MUST be thread safe
	ClearStats()
}

// Digester is the inteface for the reservoird digester plugin type. This
// plugin type digests data from an ingester queue and forwards that data through
// the another queue for further processing.
type Digester interface {
	// Name provides the name of the digest plugin
	Name() string

	// Digest is a long running function which captures data from one queue,
	// processes, and then forwards data through another queue for further processing
	//
	// inQueue: The queue which the digester receives data from
	// outQueue: The queue which the digester forwards data through
	// doneChan: The channel used to gracefully stop the long running function.
	//     If data is present on this channel initiate graceful shutdown
	// waitGroup: Call 'waitGroup.Done()' on function start. Reservoird uses
	//     this variable to wait for all threads to stop before exiting
	// error: Returns and error if there is an issue.
	Digest(inQueue Queue, outQueue Queue, doneChan <-chan struct{}, waitGroup *sync.WaitGroup) error

	// Returns marshalled stats
	//
	// NOTE: Stats functions run in a different thread than
	// digest function therefore access to stats MUST be thread safe
	Stats() (string, error)

	// Clears stats
	//
	// NOTE: Stats functions run in a different thread than
	// digest function therefore access to stats MUST be thread safe
	ClearStats()
}

// Expeller is the inteface for the reservoird expeller plugin type. This
// plugin type receives data from a queue and expels the data outside of reservorid.
// This is the termination points.
type Expeller interface {
	// Name provides the name of the expeller plugin
	Name() string

	// Expeller is a long running function which captures data from one queue,
	// processes, and then forwards data through another queue for further processing
	//
	// inQueues: The queues which the expeller receives data from
	// doneChan: The channel used to gracefully stop the long running function.
	//     If data is present on this channel initiate graceful shutdown
	// waitGroup: Call 'waitGroup.Done()' on function start. Reservoird uses
	//     this variable to wait for all threads to stop before exiting
	// error: Returns and error if there is an issue.
	Expel(inQueues []Queue, doneChan <-chan struct{}, waitGroup *sync.WaitGroup) error

	// Returns marshalled stats
	//
	// NOTE: Stats functions run in a different thread than
	// digest function therefore access to stats MUST be thread safe
	Stats() (string, error)

	// Clears stats
	//
	// NOTE: Stats functions run in a different thread than
	// expel function therefore access to stats MUST be thread safe
	ClearStats()
}
