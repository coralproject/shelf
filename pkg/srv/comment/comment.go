package comment

import "time"

type Action struct {
	Type    string 		`json:"type" bson:"type"`
	UserId  string 		`json:"userId" bson:"userId"`
	Value   string 		`json:"value" bson:"value"`
	Date    time.Time 	`json:"date" bson:"date"`
}
type ActionArray []Action


type Note struct {
	UserId  string 		`json:"userId" bson:"userId"`
	Body    string 		`json:"body" bson:"body"`
	Date    time.Time 	`json:"date" bson:"date"`
}
type NoteArray []Note


type Comment struct {
	Id              string 			`json:"id" bson:"_id"`
	Body            string 			`json:"body" bson:"body"`
	ParentId    	string 			`json:"parentId" bson:"parentId"`
	AssetId    	string 			`json:"assetId" bson:"assetId"`
	Status    	string 			`json:"status" bson:"status"`
	CreatedDate    	time.Time 		`json:"createdDate" bson:"createdDate"`
	UpdatedDate    	time.Time 		`json:"updatedDate" bson:"updatedDate"`
	ApprovedDate    time.Time 		`json:"approvedDate" bson:"approvedDate"`
	Actions		ActionArray		`json:"actions" bson:"actions"`
	Notes		NoteArray		`json:"notes" bson:"notes"`
}
type CommentArray []Comment
