package comment

import (
	"time"
)

//User denotes a user in the system
type User struct {
	ID          string    `json:"id" bson:"_id"`
	UserName    string    `json:"user_name" bson:"user_name"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	LastLogin   time.Time `json:"last_login" bson:"last_login"`
	MemberSince time.Time `json:"member_since" bson:"member_since"`
	TrustScore  float64   `json:"trust_score" bson:"trust_score"`
}
