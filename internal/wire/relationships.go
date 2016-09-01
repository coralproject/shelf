package wire

import (
	"fmt"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/coralproject/shelf/internal/wire/pattern"
	validator "gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

// validate is used to perform field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

// QuadParams contains information needed to add/remove relationships
// to/from the cayley graph.
type QuadParams struct {
	Subject   string `validate:"required,min=2"`
	Predicate string `validate:"required,min=2"`
	Object    string `validate:"required,min=2"`
}

// Validate checks the QuadParams value for consistency.
func (q *QuadParams) Validate() error {
	if err := validate.Struct(q); err != nil {
		return err
	}
	return nil
}

// AddToGraph adds relationships as quads into the cayley graph.
func AddToGraph(context interface{}, store *cayley.Handle, quadParams []QuadParams) error {
	log.Dev(context, "AddToGraph", "Started : %d Relationships", len(quadParams))

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
		log.Error(context, "AddToGraph", err, "Completed")
		return err
	}

	log.Dev(context, "AddToGraph", "Completed")
	return nil
}

// RemoveFromGraph removes relationship quads from the cayley graph.
func RemoveFromGraph(context interface{}, store *cayley.Handle, quadParams []QuadParams) error {
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

// InferRelationships infers realtionships based on patterns corresponding to
// types of items.
func InferRelationships(context interface{}, mgoDB *db.DB, items []bson.M) ([]QuadParams, error) {
	log.Dev(context, "InferRelationships", "Started : %d Items", len(items))

	// Infer relationships for each provided item.
	var quadParams []QuadParams
	for _, item := range items {

		// Get the relevant pattern.
		itemType := item["type"]
		p, err := pattern.GetByType(context, mgoDB, itemType.(string))
		if err != nil {
			continue
		}

		// Loop over inferences in the pattern.
		data := item["data"]
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("Could not parse item data")
			log.Error(context, "InferRelationships", err, "Completed")
			return nil, err
		}
		for _, inference := range p.Relationships {

			// Check for the relevant field in the item.
			if val, ok := dataMap[inference.RelIDField].(string); ok {

				// Add the relationship parameters.
				switch inference.Direction {
				case inString:
					quad := QuadParams{
						Subject:   val,
						Predicate: inference.Predicate,
						Object:    item["item_id"].(string),
					}
					quadParams = append(quadParams, quad)
				case outString:
					quad := QuadParams{
						Subject:   item["item_id"].(string),
						Predicate: inference.Predicate,
						Object:    val,
					}
					quadParams = append(quadParams, quad)
				}
			}
		}
	}

	log.Dev(context, "InferRelationships", "Completed")
	return quadParams, nil
}
