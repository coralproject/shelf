package wire_test

// getPosViews returns the table for the testing.
func getNegViews() []execView {
	return []execView{
		nameFail(),
		invalidStartType(),
		invalidRelationship(),
	}
}

// nameFail tries to execute a view that doesn't exist.
func nameFail() execView {
	return execView{
		fail:       true,
		viewName:   "this view name does not exist",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     0,
		collection: "",
		results:    []string{},
	}
}

// invalidStartType tries to execute a view with an invalid start type.
func invalidStartType() execView {
	return execView{
		fail:       true,
		viewName:   "comments from authors flagged by a user",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     0,
		collection: "",
		results:    []string{},
	}
}

// invalidRelationship tries to execute a view with an invalid relationship.
func invalidRelationship() execView {
	return execView{
		fail:       true,
		viewName:   "has invalid starting relationship",
		itemKey:    "80aa936a-f618-4234-a7be-df59a14cf8de",
		number:     0,
		collection: "",
		results:    []string{},
	}
}
