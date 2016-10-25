package wire_test

// getPosViews returns the table for the testing.
func getPosViews() []execView {
	return []execView{
		basic(),
		basicPersist(),
		persistWithBuffer(),
		returnRoot(),
		splitPath(),
		backwards(),
	}
}

// basic starts with a simple view set returning all comment items authored by a user.
func basic() execView {
	return execView{
		fail:       false,
		viewName:   "user comments",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     2,
		collection: "",
		results: []string{
			`"item_id":"WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"`,
			`"item_id":"WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"`,
		},
	}
}

// basicPersist persists a thread view in mongo.
func basicPersist() execView {
	return execView{
		fail:       false,
		viewName:   "thread",
		itemKey:    "c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
		number:     5,
		collection: "testcollection",
		results: []string{
			`"item_id":"WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"`,
			`"item_id":"WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"`,
			`"item_id":"WTEST_d16790f8-13e9-4cb4-b9ef-d82835589660"`,
			`"item_id":"WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"`,
			`"item_id":"WTEST_a63af637-58af-472b-98c7-f5c00743bac6"`,
			`"related":{"author":["WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"]}`,
			`"related":{"author":["WTEST_a63af637-58af-472b-98c7-f5c00743bac6"]}`,
		},
	}
}

// persistWithBuffer persists a thread view in mongo and uses a buffer limit when saving
// items out to mongo.
func persistWithBuffer() execView {
	return execView{
		fail:        false,
		viewName:    "thread",
		itemKey:     "c1b2bbfe-af9f-4903-8777-bd47c4d5b20a",
		number:      5,
		bufferLimit: 2,
		collection:  "testcollection",
		results: []string{
			`"item_id":"WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"`,
			`"item_id":"WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"`,
			`"item_id":"WTEST_d16790f8-13e9-4cb4-b9ef-d82835589660"`,
			`"item_id":"WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"`,
			`"item_id":"WTEST_a63af637-58af-472b-98c7-f5c00743bac6"`,
			`"related":{"author":["WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"]}`,
			`"related":{"author":["WTEST_a63af637-58af-472b-98c7-f5c00743bac6"]}`,
		},
	}
}

// returnRoot executes a simple view set returning all comment items authored by a user,
// as well as the root user item.
func returnRoot() execView {
	return execView{
		fail:       false,
		viewName:   "user comments return root",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     3,
		collection: "",
		results: []string{
			`"item_id":"WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"`,
			`"item_id":"WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"`,
			`"item_id":"WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"`,
			`"related":{"comment":["WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82","WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"]}`,
		},
	}
}

// splitPath executes a view with a split or branched path.
func splitPath() execView {
	return execView{
		fail:       false,
		viewName:   "split_path",
		itemKey:    "a63af637-58af-472b-98c7-f5c00743bac6",
		number:     3,
		collection: "",
		results: []string{
			`"item_id":"WTEST_d16790f8-13e9-4cb4-b9ef-d82835589660"`,
			`"item_id":"WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"`,
			`"item_id":"WTEST_a63af637-58af-472b-98c7-f5c00743bac6"`,
			`"related":{"comment":["WTEST_d16790f8-13e9-4cb4-b9ef-d82835589660"],"flagged_item":["WTEST_80aa936a-f618-4234-a7be-df59a14cf8de"]}`,
		},
	}
}

// backwards executes a view from lower level to high level items.
func backwards() execView {
	return execView{
		fail:       false,
		viewName:   "thread_backwards",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     3,
		collection: "",
		results: []string{
			`"item_id":"WTEST_d1dfa366-d2f7-4a4a-a64f-af89d4c97d82"`,
			`"item_id":"WTEST_6eaaa19f-da7a-4095-bbe3-cee7a7631dd4"`,
			`"item_id":"WTEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a"`,
			`"related":{"asset":["WTEST_c1b2bbfe-af9f-4903-8777-bd47c4d5b20a"]}`,
		},
	}
}
