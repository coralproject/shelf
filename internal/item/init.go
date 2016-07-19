// initialize the item system, including loading type, rel and storage config
package item

import (
	"encoding/json"
	"os"
)

var path, typesFile string

// defaults sets the default sources for item type information
func defaults() {
	// load default path, may be overwritten by plugin architecture
	path = os.Getenv("GOPATH") + "/src/github.com/coralproject/xenia/cmd/sponged/config/"

	typesFile = "types.json"

}

// RegisterTypes reads in a .json file of item types and rel data, validates it,
//  then registeres into the Type system. Types follow Upsert principles, meaning
//  that exiting types will remain or be overwritten
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

// Initialize resets the state of the item type, rels and storage
func Initialize() error {

	// Todo, incorporate plugin system
	defaults()

	err := RegisterTypes()
	if err != nil {
		return err
	}

	// TODO: initialize any custom storage drivers?

	return nil
}
