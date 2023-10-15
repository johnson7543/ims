package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Material struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Color    string             `bson:"color" json:"color"`
	Size     string             `bson:"size" json:"size"`
	Quantity string             `bson:"quantity" json:"quantity"`
	Remarks  string             `bson:"remarks" json:"remarks"`
}
