// Package sponge provides support for item importing.
package sponge

import (
	"reflect"

	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/sponge/item"
	"github.com/coralproject/shelf/internal/wire"
)

// Import imports an item into the items collections and into the graph database.
func Import(context interface{}, db *db.DB, graph *cayley.Handle, itm *item.Item) error {
	log.Dev(context, "Import", "Started : ID[%s]", itm.ID)

	// If the item exists and is different than the provided item,
	// we need to remove existing relationships.
	if itm.ID != "" {

		// Get the item if it exists.
		itmOrig, err := item.GetByID(context, db, itm.ID)
		if err != nil {
			if err != item.ErrNotFound {
				log.Error(context, "Import", err, "Completed")
				return err
			}
		}

		// If the item exits, we need to determine if the item is the
		// same as the original.  If not, we need to update the
		// relationships.
		if itmOrig.ID != "" {

			// If the item is identical, we don't have to do anything.
			if reflect.DeepEqual(itmOrig, itm) {
				log.Dev(context, "Import", "Completed")
				return nil
			}

			// If the item is not identical, remove the stale relationships by
			// preparing an item map.
			itmMap := map[string]interface{}{
				"item_id": itmOrig.ID,
				"type":    itmOrig.Type,
				"version": itmOrig.Version,
				"data":    itmOrig.Data,
			}

			// Remove the corresponding relationships from the graph.
			if err := wire.RemoveFromGraph(context, db, graph, itmMap); err != nil {
				log.Error(context, "Import", err, "Completed")
				return err
			}
		}
	}

	// Add the item to the items collection.
	if err := item.Upsert(context, db, itm); err != nil {
		log.Error(context, "Import", err, "Completed")
		return err
	}

	// Prepare the generic item data map.
	itmMap := map[string]interface{}{
		"item_id": itm.ID,
		"type":    itm.Type,
		"version": itm.Version,
		"data":    itm.Data,
	}

	// Infer relationships and add them to the graph.
	if err := wire.AddToGraph(context, db, graph, itmMap); err != nil {
		log.Error(context, "Import", err, "Completed")
		return err
	}

	log.Dev(context, "Import", "Completed")
	return nil
}

// Remove removes an item into the items collection and remove any
// corresponding quads from the graph database.
func Remove(context interface{}, db *db.DB, graph *cayley.Handle, itemID string) error {
	log.Dev(context, "Remove", "Started : ID[%s]", itemID)

	// Get the item from the items collection.
	items, err := item.GetByIDs(context, db, []string{itemID})
	if err != nil {
		log.Error(context, "Remove", err, "Completed")
		return err
	}

	// Prepare the item map data.
	itmMap := map[string]interface{}{
		"item_id": items[0].ID,
		"type":    items[0].Type,
		"version": items[0].Version,
		"data":    items[0].Data,
	}

	// Remove the corresponding relationships from the graph.
	if err := wire.RemoveFromGraph(context, db, graph, itmMap); err != nil {
		log.Error(context, "Remove", err, "Completed")
		return err
	}

	// Delete the item.
	if err := item.Delete(context, db, itemID); err != nil {
		log.Error(context, "Remove", err, "Completed")
		return err
	}

	log.Dev(context, "Remove", "Completed")
	return nil
}
