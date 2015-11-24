package comment

import "time"

//User denotes a user in the system
type User struct {
	ID          string    `json:"id" bson:"_id"`
	Email       string    `json:"email" bson:"email"`
	Password    string    `json:"password" bson:"password"`
	DisplayName string    `json:"displayName" bson:"displayName"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	Name        string    `json:"name" bson:"name"`
	Address     string    `json:"address" bson:"address"`
	City        string    `json:"city" bson:"city"`
	State       string    `json:"state" bson:"state"`
	ZipCode     string    `json:"zipCode" bson:"zipCode"`
	Gender      string    `json:"gender" bson:"gender"`
	AgeGroup    string    `json:"ageGroup" bson:"ageGroup"`
	LastLogin   time.Time `json:"lastLogin" bson:"lastLogin"`
	MemberSince time.Time `json:"memberSince" bson:"memberSince"`
	TrustScore  int       `json:"trustScore" bson:"trustScore"`
}
