package xenia

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/shelf/internal/platform/db"
	"github.com/coralproject/shelf/internal/xenia/mask"
	"gopkg.in/mgo.v2/bson"
)

// processMasks reviews the document for fields that are defined to have
// their values masked.
func processMasks(context interface{}, db *db.DB, collection string, results []bson.M) error {
	masks, err := mask.GetByCollection(context, db, collection)
	if err != nil {

		// If there are no masks to process then great.
		return nil
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

		// We have another JSON document.
		case map[string]interface{}:
			matchMaskField(context, masks, fldVal)

		// We have another BSON document.
		case bson.M:
			matchMaskField(context, masks, fldVal)

		// We have an array of JSON documents.
		case []map[string]interface{}:
			for _, subDoc := range fldVal {
				matchMaskField(context, masks, subDoc)
			}

		// We have an array of BSON documents.
		case []bson.M:
			for _, subDoc := range fldVal {
				matchMaskField(context, masks, subDoc)
			}

		// We have something we can mask.
		default:

			// If this is a field we need to mask?
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

	// Handle the remove mask for all fields.
	if msk.Type == mask.MaskRemove {
		delete(doc, key)
		return nil
	}

	// Handle fields that are not strings.
	switch doc[key].(type) {
	case int, int8, int16, int32, int64:
		doc[key] = 0
		return nil

	case float32, float64:
		doc[key] = 0.00
		return nil

	case string:
		// We need to process strings special.

	default:
		return errors.New("Invalid masking field type")
	}

	// Handle string based fields.
	switch msk.Type[0:3] {
	case mask.MaskAll:
		doc[key] = "******"
		return nil

	case mask.MaskEmail[0:3]:
		v := doc[key].(string)
		i := strings.IndexByte(v, '@')
		if i == -1 {
			return errors.New("Invalid email value")
		}

		doc[key] = "******" + v[i:]

		return nil

	case mask.MaskLeft[0:3]:

		// A left mask defaults to 4 characters to be masked. The user can
		// provide more or less by specifing size, left8. This would use 8
		// instead of 4. If there are less than specified all will be masked.
		chrs := 4
		if msk.Type != mask.MaskLeft {
			var err error
			chrs, err = strconv.Atoi(msk.Type[4:])
			if err != nil {
				log.Error(context, "applyMask", err, "Converting left size")
				return err
			}
		}

		v := doc[key].(string)
		l := len(v)
		if l < chrs {
			chrs = l
		}

		doc[key] = strings.Replace(v, v[:chrs], strings.Repeat("*", chrs), 1)
		return nil

	case mask.MaskRight[0:3]:

		// A right mask defaults to 4 characters to be masked. The user can
		// provide more or less by specifing size, right8. This would use 8
		// instead of 4. If there are less than specified all will be masked.
		chrs := 4
		if msk.Type != mask.MaskRight {
			var err error
			chrs, err = strconv.Atoi(msk.Type[5:])
			if err != nil {
				log.Error(context, "applyMask", err, "Converting right size")
				return err
			}
		}

		v := doc[key].(string)
		l := len(v)
		if l < chrs {
			chrs = l
		}

		doc[key] = strings.Replace(v, v[l-chrs:], strings.Repeat("*", chrs), 1)
		return nil

	default:
		return errors.New("Invalid masking type")
	}
}
