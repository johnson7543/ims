package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessingItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Price      float64            `bson:"price" json:"price"`
	WorkerID   primitive.ObjectID `bson:"workerId" json:"workerId"`
	WorkerName string             `bson:"workerName" json:"workerName"`
	StartDate  time.Time          `bson:"startDate" json:"startDate"`
	EndDate    time.Time          `bson:"endDate" json:"endDate"`
	SKU        string             `bson:"sku" json:"sku"`
	Remarks    string             `bson:"remarks" json:"remarks"`
}
