package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerID      primitive.ObjectID `bson:"customerId" json:"customerId"`
	CustomerName    string             `bson:"customerName" json:"customerName"`
	OrderDate       time.Time          `bson:"orderDate" json:"orderDate"`
	DeliveryDate    time.Time          `bson:"deliveryDate" json:"deliveryDate"`
	PaymentDate     time.Time          `bson:"paymentDate" json:"paymentDate"`
	TotalAmount     float64            `bson:"totalAmount" json:"totalAmount"`
	Status          string             `bson:"status" json:"status"`
	ShippingAddress string             `bson:"shippingAddress" json:"shippingAddress"`
	OrderItems      []OrderItem        `bson:"orderItems" json:"orderItems"`
}

// OrderItem represents an item within a customer order.
type OrderItem struct {
	Product    OrderProduct `bson:"product" json:"product"`
	Quantity   int          `bson:"quantity" json:"quantity"`
	TotalPrice float64      `bson:"totalPrice" json:"totalPrice"`
}

// OrderProduct represents a product associated with an order item.
type OrderProduct struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SKU       string             `bson:"sku,omitempty" json:"sku,omitempty"`
	Name      string             `bson:"name" json:"name"`
	UnitPrice float64            `bson:"unitPrice" json:"unitPrice"`
}
