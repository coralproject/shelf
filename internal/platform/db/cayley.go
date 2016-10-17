package db

import (
	"errors"
	"fmt"

	"github.com/cayleygraph/cayley"
	cayleykit "github.com/coralproject/shelf/internal/platform/db/cayley"
)

// ErrGraphHandle is returned when a graph handle is not initialized.
var ErrGraphHandle = errors.New("Graph handle not initialized.")

//==============================================================================

// NewCayley adds support to a DB value for cayley based on a registered
// master cayley handle.
func (db *DB) NewCayley(context interface{}, name string) error {
	var masterDB mgoDB
	var exists bool
	masterMGO.RLock()
	{
		masterDB, exists = masterMGO.ses[name]
	}
	masterMGO.RUnlock()

	if !exists {
		return fmt.Errorf("Master sesssion %q does not exist", name)
	}

	ses := masterDB.ses.Copy()
	store, err := cayleykit.New("", ses)
	if err != nil {
		return err
	}

	db.graphHandle = store

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
func (db *DB) CloseCayley(context interface{}) {
	db.graphHandle.Close()
}
