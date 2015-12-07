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

// Action denotes an action taken by an actor (User) on someone/something.
//   TargetType and Target id may be zero value if data is a subdocument of the Target
//   UserID may be zero value if the data is a subdocument of the actor
type Action struct {
	UserID     string    `json:"user_id" bson:"user_id" validate:"required"`
	Type       string    `json:"type" bson:"type"`
	Value      string    `json:"value" bson:"value"`
	TargetType string    `json:"target_type" bson:"target_type"` // eg: comment, "" for actions existing within target documents
	TargetId   string    `json:"target_type" bson:"target_type"` // eg: 23423
	Date       time.Time `json:"date" bson:"date"`
}

// Note denotes a note by a user in the system.
type Note struct {
	UserID string    `json:"user_id" bson:"user_id"`
	Body   string    `json:"body" bson:"body" validate:"required"`
	Date   time.Time `json:"date" bson:"date"`
}

// Comment denotes a comment by a user in the system.
type Comment struct {
	CommentID    string                 `json:"comment_id" bson:"comment_id"`
	UserId       string                 `json:"user_id" bson:"user_id" validate:"required"`
	ParentID     string                 `json:"parent_id" bson:"parent_d"`
	AssetID      string                 `json:"asset_id" bson:"asset_id"`
	Children     []string               `json:"children" bson:"children"` // experimental
	Path         string                 `json:"path" bson:"path"`
	Body         string                 `json:"body" bson:"body" validate:"required"`
	Status       string                 `json:"status" bson:"status"`
	DateCreated  time.Time              `json:"date_created" bson:"date_created"`
	DateUpdated  time.Time              `json:"date_updated" bson:"date_updated"`
	DateApproved time.Time              `json:"date_approved" bson:"date_approved"`
	Actions      []Action               `json:"actions" bson:"actions"`
	ActionCounts map[string]int         `json:"actionCounts" bson:"actionCounts"`
	Notes        []Note                 `json:"notes" bson:"notes"`
	Stats        map[string]interface{} `json:"stats" bson:"stats"`
	Source       map[string]interface{} `json:"source" bson:"source"` // source document if imported
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
	URL        string     `json:"url" bson:"url" validate:"url"`
	Taxonomies []Taxonomy `json:"taxonomies" bson:"taxonomies"`
}

// Validate performs validation on an Asset value before it is processed.
func (a *Asset) Validate() error {
	errs := validate.Struct(a)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

//==============================================================================

// User denotes a user in the system.
type User struct {
	UserID      string                 `json:"user_id" bson:"user_id"`
	UserName    string                 `json:"user_name" bson:"user_name" validate:"required"`
	Avatar      string                 `json:"avatar" bson:"avatar" validate:"omitempty,url"`
	LastLogin   time.Time              `json:"last_login" bson:"last_login"`
	MemberSince time.Time              `json:"member_since" bson:"member_since"`
	ActionsBy   []Action               `json:"actions_by" bson:"actions_by"`
	ActionsOn   []Action               `json:"actions_on" bson:"actions_on"`
	Notes       []Note                 `json:"notes" bson:"notes"`
	Stats       map[string]interface{} `json:"stats" bson:"stats"`
	Source      map[string]interface{} `json:"source" bson:"source"` // source document if imported
}

// Validate performs validation on a User value before it is processed.
func (u *User) Validate() error {
	errs := validate.Struct(u)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}

	return nil
}
