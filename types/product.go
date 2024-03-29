package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SKU      string             `bson:"sku" json:"sku"`
	Name     string             `bson:"name" json:"name"`
	Material string             `bson:"material" json:"material"`
	Color    string             `bson:"color" json:"color"`
	Type     string             `bson:"type" json:"type"`
	Size     string             `bson:"size" json:"size"`
	Quantity int                `bson:"quantity" json:"quantity"`
	Price    float64            `bson:"price" json:"price"`
	Date     time.Time          `bson:"date" json:"date"`
	Remark   string             `bson:"remark" json:"remark"`
}
