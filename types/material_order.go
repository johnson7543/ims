package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaterialOrder struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	SellerID           string              `bson:"sellerId" json:"sellerId"`
	SellerName         string              `bson:"sellerName" json:"sellerName"`
	OrderDate          time.Time           `bson:"orderDate" json:"orderDate"`
	DeliveryDate       time.Time           `bson:"deliveryDate" json:"deliveryDate"`
	PaymentDate        time.Time           `bson:"paymentDate" json:"paymentDate"`
	TotalAmount        float64             `bson:"totalAmount" json:"totalAmount"`
	Status             string              `bson:"status" json:"status"`
	MaterialOrderItems []MaterialOrderItem `bson:"materialOrderItems" json:"materialOrderItems"`
}

type MaterialOrderItem struct {
	Material   MaterialOrderMaterial `bson:"material" json:"material"`
	Quantity   int                   `bson:"quantity" json:"quantity"`
	TotalPrice float64               `bson:"totalPrice" json:"totalPrice"`
}

type MaterialOrderMaterial struct {
	MaterialID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Color      string             `bson:"color" json:"color"`
	Size       string             `bson:"size" json:"size"`
	Remarks    string             `bson:"remarks" json:"remarks"`
}
