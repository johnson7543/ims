package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceHistoryEntry struct {
	Price     float64   `bson:"price" json:"price"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Material struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string              `bson:"name" json:"name"`
	Color        string              `bson:"color" json:"color"`
	Size         string              `bson:"size" json:"size"`
	Quantity     string              `bson:"quantity" json:"quantity"`
	Remarks      string              `bson:"remarks" json:"remarks"`
	PriceHistory []PriceHistoryEntry `bson:"price_history" json:"price_history"`
}
