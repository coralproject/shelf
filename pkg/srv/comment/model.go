package comment

import (
	"fmt"
	"time"

	"gopkg.in/bluesuncorp/validator.v6"
)

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	config := validator.Config{
		TagName:         "validate",
		ValidationFuncs: validator.BakedInValidators,
	}

	validate = validator.New(config)
}

//==============================================================================

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
	CommentID    string    `json:"comment_id" bson:"comment_id"`
	ParentID     string    `json:"parent_id" bson:"parent_d"`
	AssetID      string    `json:"asset_id" bson:"asset_id"`
	Path         string    `json:"path" bson:"path"`
	Body         string    `json:"body" bson:"body"`
	Status       string    `json:"status" bson:"status"`
	DateApproved time.Time `json:"date_approved" bson:"date_approved"`
	Actions      []Action  `json:"actions" bson:"actions"`
	Notes        []Note    `json:"notes" bson:"notes"`
	DateModified time.Time `json:"date_modified" bson:"date_modified"`
	DateCreated  time.Time `json:"date_created" bson:"date_created"`
}

// Validate performs validation on a Comment value before it is processed.
func (com *Comment) Validate() error {
	errs := validate.Struct(com)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	AssetID    string     `json:"asset_id" bson:"asset_id"`
	SourceID   string     `json:"src_id" bson:"src_id"`
	URL        string     `json:"url" bson:"url"`
	Taxonomies []Taxonomy `json:"taxonomies" bson:"taxonomies"`
}

//==============================================================================

// User denotes a user in the system.
type User struct {
	UserID      string    `json:"user_id" bson:"user_id"`
	UserName    string    `json:"user_name" bson:"user_name"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	LastLogin   time.Time `json:"last_login" bson:"last_login"`
	MemberSince time.Time `json:"member_since" bson:"member_since"`
	TrustScore  float64   `json:"trust_score" bson:"trust_score"`
}
