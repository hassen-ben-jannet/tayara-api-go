package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserList struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Email      string             `json:"email,omitempty"`
	Profile_id string             `json:"profile_id,omitempty"`
	Password   string             `json:"password,omitempty"`
}
type Phone struct {
	Number string
}
