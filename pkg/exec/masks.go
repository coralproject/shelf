package exec

import (
	"github.com/coralproject/xenia/pkg/mask"

	"github.com/ardanlabs/kit/db"
	"gopkg.in/mgo.v2/bson"
)

// ProcessMasks reviews the document for fields that are defined to have
// their values masked. This function is exported so it can be tested.
func ProcessMasks(context interface{}, db *db.DB, collection string, results []bson.M) error {
	masks, err := mask.GetByCollection(context, db, collection)
	if err != nil {
		return err
	}

	for _, doc := range results {
		if err := matchMaskField(context, masks, doc); err != nil {
			return err
		}
	}

	return nil
}

// matchMaskField checks the specificed document against the masks and updated any
// field values that match based on the configured masking operation.
func matchMaskField(context interface{}, masks map[string]mask.Mask, doc bson.M) error {

	// masks: Contains the map of fields that need masking.
	// doc  : The document to check fields against to apply masking.

	for key, value := range doc {

		// What type of value does this field have.
		switch fldVal := value.(type) {

		// We have another document.
		case bson.M:
			matchMaskField(context, masks, fldVal)

		// We have an array of documents.
		case []bson.M:
			for _, subDoc := range fldVal {
				matchMaskField(context, masks, subDoc)
			}

		// We have something we can mask.
		default:

			// If this a field we need to mask?
			msk, exists := masks[key]
			if !exists {
				continue
			}

			// Apply the mask against the value.
			if err := applyMask(context, msk, doc, key); err != nil {
				return err
			}
		}
	}

	return nil
}

// applyMask performs the specified masking operation.
func applyMask(context interface{}, msk mask.Mask, doc bson.M, key string) error {
	switch msk.Type {
	case mask.MaskRemove:
		delete(doc, key)

	case mask.MaskAll:
		doc[key] = "******"

	case mask.MaskEmail:
	case mask.MaskLeft:
	case mask.MaskRight:
	}

	return nil
}
