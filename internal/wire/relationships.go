package wire

import (
	"github.com/ardanlabs/kit/log"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	validator "gopkg.in/bluesuncorp/validator.v8"
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
