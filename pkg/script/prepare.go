package script

// prepareForInsert walks the document preprocessing keys for insert.
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

		// Replace any key we find that needs to be changed.
		if key[0] == '$' {
			delete(commands, key)
			commands["_$"+key[1:]] = value
		}
	}
}

// prepareForUse walks the document preprocessing keys for use.
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

		// Replace any key we find that needs to be changed.
		if key[0:2] == "_$" {
			delete(commands, key)
			commands[key[1:]] = value
		}
	}
}
