package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PreviousOrgs struct{
   Name string `bson:"name" json:"name"`
   Duration string `bson:"duration" json:"duration"`
}

type Employee struct{
	ID *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name *string `bson:"name,omitempty" json:"name,omitempty"`
	Email *string `bson:"email,omitempty" json:"email,omitempty"`
	Position *string `bson:"position,omitempty" json:"position,omitempty"`
	Salary *float64 `bson:"salary,omitempty" json:"salary,omitempty"`
	Joining *time.Time `bson:"joining,omitempty" json:"joining,omitempty"`
	PreviousOrgs *[]PreviousOrgs `bson:"previousOrgs,omitempty" json:"previousOrgs,omitempty"`
}