// Package data handles the type infomation needed to turn data into items
package data

import (
	"encoding/json"
	"os"
)

//==============================================================================

// Type contains all we need to know in order to handle data
type Type struct {
	Name    string `bson:"name" json:"name"`
	IDField string `bson:"source_id" json:"source_id"` // field containing primary key
	Version string `bson:"version" json:"version"`     // current version
}

//==============================================================================

// Types is a slice of all the active Data Types
var Types = make(map[string]Type)

//==============================================================================

// RegisterType sets a type for use in the platform. If a type of the same name
// was already registered, it will be overwritten by the new type.
func RegisterType(t Type) {

	Types[t.Name] = t

}

//==============================================================================

// UnregisterType removes a type from use in the platform. Returns true if the
// type is unregistered, false if not found. Removing a type will disable itemization
// of that data type and has no effect on items already in the platform.
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

// TODO: integrate this type loading system with Coral Platform config loading

var path, typesFile string

// defaults sets the default sources for item type information
func defaults() {
	// load default path, may be overwritten by plugin architecture
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/cmd/sponged/config/"

	typesFile = "types.json"

}

// RegisterTypes reads in a .json file of item types and rel data, validates it,
// then registeres into the Type system. Types follow Upsert principles, meaning
// that exiting types will remain or be overwritten.
func RegisterTypes() error {

	file, err := os.Open(path + typesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// grab the item type fixture file
	var types []Type
	err = json.NewDecoder(file).Decode(&types)
	if err != nil {
		return err
	}

	// retister the item types
	for _, t := range types {
		RegisterType(t)
	}

	return nil
}

// Initialize resets the state of the item type.
func Initialize() error {

	// Todo, incorporate plugin system
	defaults()

	err := RegisterTypes()
	if err != nil {
		return err
	}

	return nil
}
