package wire

import (
	"fmt"

	"github.com/ardanlabs/kit/db"
	"github.com/coralproject/shelf/internal/wire/relationship"
	"github.com/coralproject/shelf/internal/wire/view"
)

// validateStartType verifies the start type of a view path.
func validateStartType(context interface{}, db *db.DB, v *view.View) error {

	// Extract the first level relationship predicate.
	var firstRel string
	var firstDir string
	for _, segment := range v.Path {
		if segment.Level == 1 {
			firstRel = segment.Predicate
			firstDir = segment.Direction
		}
	}

	// Get the relationship metadata.
	rel, err := relationship.GetByPredicate(context, db, firstRel)
	if err != nil {
		return err
	}

	// Get the relevant item types based on the direction of the
	// first relationship in the path.
	var itemTypes []string
	switch firstDir {
	case "out":
		itemTypes = rel.SubjectTypes
	case "in":
		itemTypes = rel.ObjectTypes
	}

	// Validate the starting type provided in the view.
	for _, itemType := range itemTypes {
		if itemType == v.StartType {
			return nil
		}
	}

	return fmt.Errorf("Start type %s does not match relationship subject types %v", v.StartType, itemTypes)
}
