package icd

import (
	"sync"
)

// Queue provides interface for queue plugins
type Queue interface {
	Name() string
	Put(interface{}) error
	Get() (interface{}, error)
	Peek() (interface{}, error)
	Len() int
	Cap() int
	Clear()
	Close() error
}

// Ingester provides interface for plugins that ingest (push) data into reservoird
// struct channel and wait group are for graceful shutdown of ingester plugin
type Ingester interface {
	Name() string
	Ingest(Queue, <-chan struct{}, *sync.WaitGroup) error
}

// Digester provides interface for plugins that filter/annotate (push/pop) data within reservoird
// struct channel and wait group are for graceful shutdown of digester plugin
type Digester interface {
	Name() string
	Digest(Queue, Queue, <-chan struct{}, *sync.WaitGroup) error
}

// Expeller provides interface for plugins that expel (pop) data outof reservoird
// struct channel and wait group are for graceful shutdown of expeller plugin
type Expeller interface {
	Name() string
	Expel([]Queue, <-chan struct{}, *sync.WaitGroup) error
}
