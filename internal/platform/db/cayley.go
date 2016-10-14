package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cayleygraph/cayley"
	cayleymongo "github.com/cayleygraph/cayley/graph/mongo"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
)

// ErrGraphHandle is returned when a graph handle is not initialized.
var ErrGraphHandle = errors.New("Graph handle not initialized.")

//==============================================================================

// cayleyDB maintains a master handle for a cayley database.
type cayleyDB struct {
	qs *cayley.Handle
}

// masterCayley manages a set of different cayley master handles.
var masterCayley = struct {
	sync.RWMutex
	qs map[string]cayleyDB
}{
	qs: make(map[string]cayleyDB),
}

// RegMasterHandle adds a new master handle to the set. If no url is provided,
// it will default to localhost:27017.
func RegMasterHandle(context interface{}, name string, url string) error {
	masterCayley.Lock()
	defer masterCayley.Unlock()

	if _, exists := masterCayley.qs[name]; exists {
		return errors.New("Master session already exists")
	}

	store, err := cayleyshelf.New(url)
	if err != nil {
		return err
	}

	masterCayley.qs[name] = cayleyDB{
		qs: store,
	}

	return nil
}

//==============================================================================

// NewCayley adds support to a DB value for cayley based on a registered
// master cayley handle.
func (db *DB) NewCayley(context interface{}, name string) error {
	var exists bool
	var handle cayleyDB
	masterCayley.RLock()
	{
		handle, exists = masterCayley.qs[name]
	}
	masterCayley.RUnlock()

	if !exists {
		return fmt.Errorf("Master sesssion %q does not exist", name)
	}

	mongoStore, ok := handle.qs.QuadStore.(*cayleymongo.QuadStore)
	if !ok {
		return fmt.Errorf("Cayley quadstore not using the mongo backend")
	}

	store := mongoStore.Copy()
	handleOut := cayley.Handle{
		QuadStore:  store,
		QuadWriter: handle.qs.QuadWriter,
	}

	db.graphHandle = &handleOut

	return nil
}

//==============================================================================
// Methods for the DB struct type related to Cayley.

// GraphHandle returns the Cayley graph handle for graph interactions.
func (db *DB) GraphHandle(context interface{}) (*cayley.Handle, error) {
	if db.graphHandle != nil {
		return db.graphHandle, nil
	}

	return nil, ErrGraphHandle
}

// CloseCayley closes a graph handle value.
func (db *DB) CloseCayley(context interface{}) error {
	mongoStore, ok := db.graphHandle.QuadStore.(*cayleymongo.QuadStore)
	if !ok {
		return fmt.Errorf("Cayley quadstore not using the mongo backend")
	}

	mongoStore.Release()

	return nil
}
