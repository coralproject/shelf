//go:generate go-bindata -pkg queries -o assets.go sets/

package queries

import "fmt"

// Load unmarshals the specified fixture into the provided
// data value.
func Load(name string) ([]byte, error) {

	// Load the fixtures bytes into the byte slice.
	return Asset(fmt.Sprintf("sets/%s.json", name))
}
