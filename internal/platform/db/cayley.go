package db

import (
	"errors"

	"github.com/cayleygraph/cayley"
	cayleyshelf "github.com/coralproject/shelf/internal/platform/db/cayley"
)

// ErrGraphHandle is returned when a graph handle is not initialized.
var ErrGraphHandle = errors.New("Graph handle not initialized.")

//==============================================================================
// Methods for the DB struct type related to Cayley.

// OpenCayley opens a connection to Cayley and adds that support to the
// database value.
func (db *DB) OpenCayley(context interface{}, mongoURL string) error {
	store, err := cayleyshelf.New(mongoURL)
	if err != nil {
		return err
	}

	db.graphHandle = store
	return nil
}

// GraphHandle returns the Cayley graph handle for graph interactions.
func (db *DB) GraphHandle(context interface{}) (*cayley.Handle, error) {
	if db.graphHandle != nil {
		return db.graphHandle, nil
	}

	return nil, ErrGraphHandle
}

// CloseCayley closes a graph handle value.
func (db *DB) CloseCayley(context interface{}) {
	db.graphHandle.Close()
}
