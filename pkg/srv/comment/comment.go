// TODO: add package description
package comment

import "time"

// Action denotes an action taken by someone/something on someone/something.
type Action struct {
	Type   string    `json:"type" bson:"type"`
	UserID string    `json:"userId" bson:"userId"`
	Value  string    `json:"value" bson:"value"`
	Date   time.Time `json:"date" bson:"date"`
}

// Note denotes a note by a user in the system.
type Note struct {
	UserID string    `json:"userId" bson:"userId"`
	Body   string    `json:"body" bson:"body"`
	Date   time.Time `json:"date" bson:"date"`
}

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID           string    `json:"id" bson:"_id"`
	Body         string    `json:"body" bson:"body"`
	ParentID     string    `json:"parentId" bson:"parentId"`
	AssetID      string    `json:"assetId" bson:"assetId"`
	Status       string    `json:"status" bson:"status"`
	CreatedDate  time.Time `json:"createdDate" bson:"createdDate"`
	UpdatedDate  time.Time `json:"updatedDate" bson:"updatedDate"`
	ApprovedDate time.Time `json:"approvedDate" bson:"approvedDate"`
	Actions      []Action  `json:"actions" bson:"actions"`
	Notes        []Note    `json:"notes" bson:"notes"`
}
