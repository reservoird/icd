// Package icd provides the interface which plugins must implement
// in order to be used within the reservoird framework. Each plugin
// must provide a 'New' function with the following function signature:
//
//   New(cfg string, stats chan<- string) (icd.<interface>, error)
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
	Name() string               // Name of the queue
	Put(interface{}) error      // Put an item into the queue
	Get() (interface{}, error)  // Get an item from the queue
	Peek() (interface{}, error) // Peek at an item in the queue
	Len() int                   // Len returns the length of the queue
	Cap() int                   // Cap returns the capacity of the queue
	Clear()                     // Clear zeros out the queue
	Close() error               // Close closes the queue, can no longer be used
	Closed() bool               // Closed provides the state of the queue, open or closed
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
}
