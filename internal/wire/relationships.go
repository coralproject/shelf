package wire

import (
	"errors"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/wire/pattern"
	validator "gopkg.in/bluesuncorp/validator.v8"
)

var (
	// validate is used to perform field validation.
	validate *validator.Validate

	// ErrItemType is used in item parsing.
	ErrItemType = errors.New("Could not parse item type")

	// ErrItemData is used in item parsing.
	ErrItemData = errors.New("Could not parse item data")

	// ErrItemID is used in item parsing.
	ErrItemID = errors.New("Could not parse item ID")
)

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

// QuadParam contains information needed to add/remove relationships
// to/from the cayley graph.
type QuadParam struct {
	Subject   string `validate:"required,min=2"`
	Predicate string `validate:"required,min=2"`
	Object    string `validate:"required,min=2"`
}

// Validate checks the QuadParams value for consistency.
func (q *QuadParam) Validate() error {
	if err := validate.Struct(q); err != nil {
		return err
	}
	return nil
}

// AddToGraph adds relationships as quads into the cayley graph.
func AddToGraph(context interface{}, db *db.DB, store *cayley.Handle, item map[string]interface{}) error {
	log.Dev(context, "AddToGraph", "Started : %v", item)

	// Infer the relationships in the item.
	quadParams, err := inferRelationships(context, db, item)
	if err != nil {
		log.Error(context, "AddToGraph", err, "Completed")
		return err
	}

	// Convert the given parameters into cayley quads.
	tx := cayley.NewTransaction()
	for _, params := range quadParams {

		// Validate the parameters.
		if err := params.Validate(); err != nil {
			log.Error(context, "AddToGraph", err, "Completed")
			return err
		}

		// Form the cayley quad.
		quad := quad.Make(params.Subject, params.Predicate, params.Object, "")
		tx.AddQuad(quad)
	}

	// Apply the transaction.
	if err := store.ApplyTransaction(tx); err != nil {
		if !graph.IsQuadExist(err) {
			log.Error(context, "AddToGraph", err, "Completed")
			return err
		}
	}

	log.Dev(context, "AddToGraph", "Completed")
	return nil
}

// RemoveFromGraph removes relationship quads from the cayley graph.
func RemoveFromGraph(context interface{}, store *cayley.Handle, quadParams []QuadParam) error {
	log.Dev(context, "RemoveFromGraph", "Started : %d Relationships", len(quadParams))

	// Convert the given parameters into cayley quads.
	tx := cayley.NewTransaction()
	for _, params := range quadParams {

		// Validate the parameters.
		if err := params.Validate(); err != nil {
			log.Error(context, "RemoveFromGraph", err, "Completed")
			return err
		}

		// Form the cayley quad.
		quad := quad.Make(params.Subject, params.Predicate, params.Object, "")
		tx.RemoveQuad(quad)
	}

	// Apply the transaction.
	if err := store.ApplyTransaction(tx); err != nil {
		log.Error(context, "RemoveFromGraph", err, "Completed")
		return err
	}

	log.Dev(context, "RemoveFromGraph", "Completed")
	return nil
}

// inferRelationships infers realtionships based on patterns corresponding to
// a type of item.
func inferRelationships(context interface{}, db *db.DB, itemIn map[string]interface{}) ([]QuadParam, error) {

	// Parse the item.
	item, err := itemParse(itemIn)
	if err != nil {
		return nil, err
	}

	// Get the relevant pattern.
	p, err := pattern.GetByType(context, db, item.itemType)
	if err != nil {
		return nil, err
	}

	// Loop over inferences in the pattern.
	var qps []QuadParam
	for _, inf := range p.Inferences {

		// Check for the relevant field in the item.
		if relID, ok := item.itemData[inf.RelIDField]; ok {

			// Add the relationship parameters.
			switch inf.Direction {
			case inString:
				qp := QuadParam{
					Subject:   relID,
					Predicate: inf.Predicate,
					Object:    item.itemID,
				}
				qps = append(qps, qp)
			case outString:
				qp := QuadParam{
					Subject:   item.itemID,
					Predicate: inf.Predicate,
					Object:    relID,
				}
				qps = append(qps, qp)
			}
		}
	}

	return qps, nil
}

// parsedItem contains the structure of the item.
type parsedItem struct {
	itemID   string
	itemType string
	itemData map[string]string
}

// itemParse parses a general map[string]interface{} into a parsedItem value,
// validating the fields required for relationship inference.
func itemParse(itemIn map[string]interface{}) (parsedItem, error) {

	// Validate and extract the item type.
	val, ok := itemIn["type"]
	if !ok {
		return parsedItem{}, ErrItemType
	}

	itemType, ok := val.(string)
	if !ok {
		return parsedItem{}, ErrItemType
	}

	if itemType == "" {
		return parsedItem{}, ErrItemType
	}

	// Validate and extract the item ID.
	val, ok = itemIn["item_id"]
	if !ok {
		return parsedItem{}, ErrItemID
	}

	itemID, ok := val.(string)
	if !ok {
		return parsedItem{}, ErrItemType
	}

	if itemID == "" {
		return parsedItem{}, ErrItemType
	}

	// Validate and extract the item data.
	val, ok = itemIn["data"]
	if !ok {
		return parsedItem{}, ErrItemData
	}

	dataMap, ok := val.(map[string]interface{})
	if !ok {
		return parsedItem{}, ErrItemData
	}

	itemData := make(map[string]string)
	for k, v := range dataMap {

		vString, ok := v.(string)
		if !ok {
			continue
		}
		itemData[k] = vString
	}

	if len(itemData) == 0 {
		return parsedItem{}, ErrItemData
	}

	// Create and return the parsed item value.
	itemOut := parsedItem{
		itemID:   itemID,
		itemType: itemType,
		itemData: itemData,
	}
	return itemOut, nil
}
