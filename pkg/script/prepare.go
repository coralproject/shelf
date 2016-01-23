package script

import (
	"strings"
)

// prepareForInsert walks the document preprocessing keys for insert.
//
// MongoDB will not let us save field names with '$' in the beginning or
// using dot (name.name) notation. We need to change that out to save.
func prepareForInsert(commands map[string]interface{}) {
	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			prepareForInsert(doc)

		// We have an array of values.
		case []interface{}:

			// Iterate over the array of values.
			for _, subDoc := range doc {

				// I only care about documents because we are looking for keys.
				if cmd, ok := subDoc.(map[string]interface{}); ok {
					prepareForInsert(cmd)
				}
			}
		}

		if key[0] == '$' {

			// Replace any key we find starts with $.
			delete(commands, key)
			commands["_"+key] = value

		} else {

			// Replace any key we find that has dot notation.
			if idx := strings.Index(key, "."); idx != -1 {
				delete(commands, key)
				commands[key[0:idx]+"*"+key[idx+1:]] = value
			}

		}
	}
}

// prepareForUse walks the document preprocessing keys for use.
//
// MongoDB will not let us save field names with '$' in the beginning or
// using dot (name.name) notation. We need to change that out to save. But
// when we get the document back, we need to replace things back.
func prepareForUse(commands map[string]interface{}) {
	for key, value := range commands {

		// Test for the type of value we have.
		switch doc := value.(type) {

		// We have another document.
		case map[string]interface{}:
			prepareForUse(doc)

		// We have an array of values.
		case []interface{}:

			// Iterate over the array of values.
			for _, subDoc := range doc {

				// I only care about documents because we are looking for keys.
				if cmd, ok := subDoc.(map[string]interface{}); ok {
					prepareForUse(cmd)
				}
			}
		}

		if key[0:2] == "_$" {

			// Replace any key we find starts with _$.
			delete(commands, key)
			commands[key[1:]] = value

		} else {

			// Replace any key we find that has *.
			if idx := strings.Index(key, "*"); idx != -1 {
				delete(commands, key)
				commands[key[0:idx]+"."+key[idx+1:]] = value
			}

		}
	}
}
