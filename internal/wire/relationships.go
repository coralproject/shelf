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

	// Verify the starting type.
	verify := false
	switch firstDir {
	case "out":
		for _, itemType := range rel.SubjectTypes {
			if itemType == v.StartType {
				verify = true
			}
		}
	case "in":
		for _, itemType := range rel.ObjectTypes {
			if itemType == v.StartType {
				verify = true
			}
		}
	}

	if !verify {
		return fmt.Errorf("Start type %s does not match relationship subject types %v", v.StartType, rel.SubjectTypes)
	}

	return nil
}
