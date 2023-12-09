package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Buyer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Company     string             `bson:"company" json:"company"`
	Name        string             `bson:"name" json:"name"`
	Phone       string             `bson:"phone" json:"phone"`
	Address     string             `bson:"address" json:"address"`
	TaxIdNumber string             `bson:"taxIdNumber" json:"taxIdNumber"`
}
