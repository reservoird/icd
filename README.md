# icd

This is the ICD (interface control document) for plugins to be successfully
used within the reservoird system. For reservoird architecture please see
https://github.com/reservoird/reservoird

There are 4 plugin types that the reservoird system supports:

- Queue: Queues provides a communication path between plugin types
- Ingester: Ingests data from a source and forwards to a digester or expeller for further processing
- Digester: Digests data from an ingester and forwards to another digester or an expeller for further processing
- Expeller: Expels data from an ingester or digester external to reservoird.

The godoc provides details of each function required and a description of what function the job should perform.

Example:

We want ingester to read data from stdin. (See https://github.com/reservoird/stdin for actual implementation)

```
package main // plugins require main package name

import (
    // any dependancies
)

func New(cfg string) (icd.Ingester, error) {
    // Set up and configure stdin ingester
}

// long running function with reads from stdin and adds the result to the Queue.
func Ingest(sendQueue Queue, doneChan <-chan struct{}, wg *sync.WaitGroup) {
    // first line of the function must be wg.Done() as reservoird waits for all threads to stop
    // before exiting.
    //
    // reads from stdin and writes to queue
    //
    // non-blocking listen on ingest clearChan each loop, if received clears statistics
    // FLOW: reservoird => monitor => ingest
    //
    // non-blocking send to ingest statsChan which provides the latest snapshot of statistics
    // FLOW: ingest => monitor => reservoird
    //
    // non-blocking listen on reservoird doneChan each loop to see if reservoird is shutting down gracefully
    // FLOW: reservoird => ingest
    // NOTE: only senders should close the queue
}

// long running function which provides statistics and a means to clear statistics
func Monitor(statsChan chan<- string, clearChan <-chan struct{}, doneChan <-chan struct{}, wg *sync.WaitGroup) {
    // first line of the function must be wg.Done() as reseroird waits for all threads to stop
    // before exiting.
    //
    // get stats from ingest thread
    //
    // marshal statistics
    //
    // non-blocking sent statistics to reservoird statsChan
    // FLOW: ingest => monitor => reservoird
    //
    // non-blocking read of reservoird clearChan, if received clear stats from ingest thread
    // FLOW: reservoird => monitor => ingest
    //
    // non-blocking listen on reservoird doneChan each loop to see if reservoird is shutting down gracefully
    // FLOW: reservoird => monitor
}
```
