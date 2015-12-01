package comment

import (
	"time"
)

// Action denotes an action taken by someone/something on someone/something.
type Action struct {
	Type   string    `json:"type" bson:"type"`
	UserID string    `json:"user_id" bson:"user_id"`
	Value  string    `json:"value" bson:"value"`
	Date   time.Time `json:"date" bson:"date"`
}

// Note denotes a note by a user in the system.
type Note struct {
	UserID string    `json:"user_id" bson:"user_id"`
	Body   string    `json:"body" bson:"body"`
	Date   time.Time `json:"date" bson:"date"`
}

// Comment denotes a comment by a user in the system.
type Comment struct {
	ID           string    `json:"id" bson:"_id"`
	ParentID     string    `json:"parent_id" bson:"parent_d"`
	AssetID      string    `json:"asset_id" bson:"asset_id"`
	Path         string    `json:"path" bson:"path"`
	Body         string    `json:"body" bson:"body"`
	Status       string    `json:"status" bson:"status"`
	DateCreated  time.Time `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time `json:"date_approved" bson:"date_approved"`
	Actions      []Action  `json:"actions" bson:"actions"`
	Notes        []Note    `json:"notes" bson:"notes"`
}
