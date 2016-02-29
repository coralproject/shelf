package exec

import (
	"strconv"

	"github.com/coralproject/xenia/pkg/mask"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2/bson"
)

// processMasks reviews the document for fields that are defined to have
// their values masked.
func processMasks(context interface{}, db *db.DB, collection string, results []bson.M) error {
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
func matchMaskField(context interface{}, masks map[string]mask.Mask, doc map[string]interface{}) error {

	// masks: Contains the map of fields that need masking.
	// doc  : The document to check fields against to apply masking.

	for key, value := range doc {

		// What type of value does this field have.
		switch fldVal := value.(type) {

		// We have another document.
		case map[string]interface{}:
			matchMaskField(context, masks, fldVal)

		// We have an array of documents.
		case []map[string]interface{}:
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
	switch msk.Type[0:3] {
	case mask.MaskRemove[0:3]:
		delete(doc, key)

	case mask.MaskAll:
		switch doc[key].(type) {
		case string:
			doc[key] = "******"
		case int, int8, int16, int32, int64:
			doc[key] = 0
		case float32, float64:
			doc[key] = 0.00
		}

	case mask.MaskEmail[0:3]:

	case mask.MaskLeft[0:3]:

		// A left mask defaults to 4 characters to be masked. The user can
		// provide more or less by specifing size, left8. This would use 8
		// instead of 4. If there are less than specified all will be masked.
		chrs := 4
		if msk.Type != "left" {
			var err error
			chrs, err = strconv.Atoi(msk.Type[5:])
			if err != nil {
				log.Error(context, "applyMask", err, "Converting left size")
				return err
			}
		}

	case mask.MaskRight[0:3]:
	}

	return nil
}
