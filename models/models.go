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
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string `bson:"name" json:"name"`
	Position string `bson:"position" json:"position"`
	Salary float64 `bson:"salary" json:"salary"`
	Joining time.Time `bson:"joining" json:"joining"`
	PreviousOrgs []PreviousOrgs `bson:"previousOrgs" json:"previousOrgs"`
}