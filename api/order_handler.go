package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertOrderParams struct {
	CustomerID      string                  `json:"customerId"`
	OrderDate       string                  `json:"orderDate"`
	DeliveryDate    string                  `json:"deliveryDate"`
	PaymentDate     string                  `json:"paymentDate"`
	TotalAmount     float64                 `json:"totalAmount"`
	Status          string                  `json:"status"`
	ShippingAddress string                  `json:"shippingAddress"`
	OrderItems      []InsertOrderItemParams `json:"orderItems"`
}

type InsertOrderItemParams struct {
	Product    InsertOrderProductParams `json:"product"`
	Quantity   int                      `json:"quantity"`
	TotalPrice float64                  `json:"totalPrice"`
}

type InsertOrderProductParams struct {
	ID        string  `json:"id"`
	SKU       string  `json:"sku"`
	UnitPrice float64 `json:"unitPrice"`
}

func (p InsertOrderParams) validate() error {
	return nil
}

type UpdateOrderParams struct {
	CustomerID      string                  `json:"customerId"`
	OrderDate       string                  `json:"orderDate"`
	DeliveryDate    string                  `json:"deliveryDate"`
	PaymentDate     string                  `json:"paymentDate"`
	TotalAmount     float64                 `json:"totalAmount"`
	Status          string                  `json:"status"`
	ShippingAddress string                  `json:"shippingAddress"`
	OrderItems      []UpdateOrderItemParams `json:"orderItems"`
}

type UpdateOrderItemParams struct {
	Product    UpdateOrderProductParams `json:"product"`
	Quantity   int                      `json:"quantity"`
	TotalPrice float64                  `json:"totalPrice"`
}

type UpdateOrderProductParams struct {
	ID        string  `json:"id"`
	SKU       string  `json:"sku"`
	UnitPrice float64 `json:"unitPrice"`
}

func (p UpdateOrderParams) validate() error {
	return nil
}

type OrderHandler struct {
	store *db.Store
}

func NewOrderHandler(store *db.Store) *OrderHandler {
	return &OrderHandler{
		store: store,
	}
}

// HandleGetOrders retrieves a list of orders based on query parameters.
//
// @Summary Get orders
// @Description Retrieves a list of orders based on query parameters.
// @Tags Order
// @Param id query string false "Order ID"
// @Param customerId query string false "Customer ID"
// @Param status query string false "Order status"
// @Produce json
// @Success 200 {array} types.Order
// @Router /order [get]
func (h *OrderHandler) HandleGetOrders(c *fiber.Ctx) error {
	id := c.Query("id")
	customerID := c.Query("customerId")
	status := c.Query("status")

	filter := bson.M{}

	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid order ID",
			})
		}
		filter["_id"] = objID
	}

	if customerID != "" {
		objID, err := primitive.ObjectIDFromHex(customerID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid customer ID",
			})
		}
		filter["customerId"] = objID
	}

	if status != "" {
		filter["status"] = status
	}

	orders, err := h.store.Order.GetOrders(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(orders)
}

// HandleInsertOrder inserts a new order.
//
// @Summary Insert order
// @Description Inserts a new order.
// @Tags Order
// @Accept json
// @Produce json
// @Param order body InsertOrderParams true "Order information"
// @Success 200 {object} types.Order
// @Router /order [post]
func (h *OrderHandler) HandleInsertOrder(c *fiber.Ctx) error {
	var params InsertOrderParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	customerID, err := primitive.ObjectIDFromHex(params.CustomerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	orderDateParsed, err := time.Parse(time.RFC3339Nano, params.OrderDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order date format",
		})
	}

	orderItems := make([]types.OrderItem, len(params.OrderItems))
	for i, item := range params.OrderItems {
		productID, err := primitive.ObjectIDFromHex(item.Product.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		orderItems[i] = types.OrderItem{
			Product: types.OrderProduct{
				ID:        productID,
				SKU:       item.Product.SKU,
				UnitPrice: item.Product.UnitPrice,
			},
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
		}
	}

	order := types.Order{
		CustomerID:      customerID,
		OrderDate:       orderDateParsed,
		TotalAmount:     params.TotalAmount,
		Status:          params.Status,
		ShippingAddress: params.ShippingAddress,
		OrderItems:      orderItems,
	}

	if params.DeliveryDate != "" {
		deliveryDateParsed, err := time.Parse(time.RFC3339Nano, params.DeliveryDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid delivery date format",
			})
		}
		order.DeliveryDate = deliveryDateParsed
	}

	if params.PaymentDate != "" {
		paymentDateParsed, err := time.Parse(time.RFC3339Nano, params.PaymentDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid payment date format",
			})
		}
		order.PaymentDate = paymentDateParsed
	}

	inserted, err := h.store.Order.InsertOrder(c.Context(), &order)
	if err != nil {
		return err
	}

	// Decrease product quantities after successfully inserting the order
	for _, item := range params.OrderItems {
		productID, err := primitive.ObjectIDFromHex(item.Product.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid product ID: %s", item.Product.ID),
			})
		}

		updatedCount, err := h.store.Product.DecreaseProductQuantity(c.Context(), productID, item.Quantity)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to decrease product %s by %d, %s", item.Product.ID, item.Quantity, err.Error()),
			})
		}

		if updatedCount == 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to decrease product %s by %d", item.Product.ID, item.Quantity),
			})
		}

	}

	return c.JSON(inserted)
}

// HandleUpdateOrder updates an existing order in the system.
// @Summary Update order
// @Description Update an existing order in the system.
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param body body UpdateOrderParams true "Updated order details"
// @Success 200 {object} fiber.Map
// @Router /order/{id} [patch]
func (h *OrderHandler) HandleUpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateOrderParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	customerID, err := primitive.ObjectIDFromHex(params.CustomerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	orderDateParsed, err := time.Parse(time.RFC3339Nano, params.OrderDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order date format",
		})
	}

	orderItems := make([]types.OrderItem, len(params.OrderItems))
	for i, item := range params.OrderItems {
		productID, err := primitive.ObjectIDFromHex(item.Product.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		orderItems[i] = types.OrderItem{
			Product: types.OrderProduct{
				ID:        productID,
				SKU:       item.Product.SKU,
				UnitPrice: item.Product.UnitPrice,
			},
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
		}
	}

	updatedOrder := types.Order{
		CustomerID:      customerID,
		OrderDate:       orderDateParsed,
		TotalAmount:     params.TotalAmount,
		Status:          params.Status,
		ShippingAddress: params.ShippingAddress,
		OrderItems:      orderItems,
	}

	if params.DeliveryDate != "" {
		deliveryDateParsed, err := time.Parse(time.RFC3339Nano, params.DeliveryDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid delivery date format",
			})
		}
		updatedOrder.DeliveryDate = deliveryDateParsed
	}

	if params.PaymentDate != "" {
		paymentDateParsed, err := time.Parse(time.RFC3339Nano, params.PaymentDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid payment date format",
			})
		}
		updatedOrder.PaymentDate = paymentDateParsed
	}

	updateCount, err := h.store.Order.UpdateOrder(c.Context(), orderID, &updatedOrder)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	existingOrder, err := h.store.Order.GetOrders(c.Context(), bson.M{"_id": orderID})
	if err != nil {
		return err
	}

	if strings.ToLower(existingOrder[0].Status) != "canceled" && strings.ToLower(params.Status) == "canceled" {
		for _, item := range orderItems {
			updatedCount, err := h.store.Product.IncreaseProductQuantity(c.Context(), item.Product.ID, item.Quantity)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to increase product %s by %d, %s", item.Product.ID, item.Quantity, err.Error()),
				})
			}

			if updatedCount == 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to increase product %s by %d", item.Product.ID, item.Quantity),
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"message": "Order updated successfully",
	})
}

// HandleDeleteOrder deletes an order by ID.
//
// @Summary Delete order
// @Description Deletes an order by ID.
// @Tags Order
// @Param id path string true "Order ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /order/{id} [delete]
func (h *OrderHandler) HandleDeleteOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	deleteCount, err := h.store.Order.DeleteOrder(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order deleted successfully",
	})
}
