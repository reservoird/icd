# icd

This is the ICD (interface control document) for plugins to be successfully
used within the reservoird system. For reservoird architecture please see
github.com/reservoird/reservoird

There are 4 plugin types that the reservoird system supports:

- Queue: Queues provides a communication path between plugin types
- Ingester: Ingests data from a source and forwards to a digester or expeller for further processin
- Digester: Digests data from an ingester and forwards to another digester or an expeller for further processing
- Expeller: Expels data from an ingester or digester external to reservoird.

The godoc provides details of each function required and a description of what function the job should perform.

Example:

We want ingester to read data from stdin. (See github.com/reservoird/stdin for actual implementation)

```
package main // plugins require main package name

import (
    // any dependancies
)

func New(cfg string, stats chan <-string) (icd.Ingester, error) {
    // Set up and configure stdin ingester
}

func Ingest(queue Queue, doneChan <-chan struct{}, wg *sync.WaitGroup) {
    // long running function with reads from stdin and adds the
    // result to the Queue.
    //
    // first line of the function must be wg.Done() as reservoird waits for all threads to stop
    // before exiting.
    //
    // non-blocking listen on doneChan each loop to see if reservoird is shutting down gracefully
    //
    // non-blocking send on stats channel each loop provides ingester statistics information
}
```