// Itemtypes describe the properties and relationships of Items
package item

//==============================================================================

// RelType holds the config for describing an Item Type's rels
//  We can find related item(s) by querying for the item type and field
type RelType struct {
	Name     string `bson:"name" json:"name"`   // Name of relationship
	Type     string `bson:"type" json:"type"`   // Item Type of target
	Field    string `bson:"field" json:"field"` // field containing foreign key
	Required bool   `bson:"required" json:"required"`
}

//==============================================================================

// ItemType contains all we need to know in order to handle an Item
type Type struct {
	Name    string    `bson:"name" json:"name"`
	IdField string    `bson:"id" json:"id"` // the primary key of this item type
	Rels    []RelType `bson:"rels" json:"rels"`
}

//==============================================================================

// Types is a slice of all the active Item Types
var Types = make(map[string]Type)

//==============================================================================

// Register a Type
//  if a type of the same name was already registered, it will
//  be overwritten by the new type
func RegisterType(t Type) {

	Types[t.Name] = t

}

//==============================================================================

// Unregister a Type
//  returns true if the type is unregistered, false if not found
func UnregisterType(t Type) bool {

	// how is this done, no internetz here
	//return Types.Delete(t.Name)

	return false

}

//==============================================================================

// isTypeRegistered takes a type name and finds if it's in the Types
func isRegistered(n string) bool {

	for _, t := range Types {
		if n == t.Name {
			return true
		}
	}

	return false
}

//==============================================================================

// getRelsByType gets the rels for a given type
func getRelsByType(t string) (*[]RelType, error) {

	rts := Types[t].Rels

	return &rts, nil

}
