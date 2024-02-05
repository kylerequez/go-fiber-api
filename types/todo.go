package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title string             `bson:"title,omitempty" json:"title"`
	Body  string             `bson:"body,omitempty" json:"body"`
}
