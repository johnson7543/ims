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
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertMaterialOrderParams struct {
	ID                 string  `json:"id,omitempty"`
	SellerID           string  `json:"sellerID"`
	OrderDate          string  `json:"orderDate"`
	DeliveryDate       string  `json:"deliveryDate"`
	PaymentDate        string  `json:"paymentDate"`
	TotalAmount        float64 `json:"totalAmount"`
	Status             string
	MaterialOrderItems []InsertMaterialOrderItemParams `json:"materialOrderItems"`
}

type InsertMaterialOrderItemParams struct {
	Material   InsertMaterialOrderMaterialParams `json:"material"`
	Quantity   int                               `json:"quantity"`
	TotalPrice float64                           `json:"totalPrice"`
}

type InsertMaterialOrderMaterialParams struct {
	MaterialID string  `json:"id,omitempty"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Color      string  `json:"color"`
	Size       string  `json:"size"`
	Quantity   int     `json:"quantity"`
	Remarks    string  `json:"remarks"`
}

func (p InsertMaterialOrderParams) validate() error {
	return nil
}

type UpdateMaterialOrderParams struct {
	ID           string  `json:"id,omitempty"`
	SellerID     string  `json:"sellerID"`
	OrderDate    string  `json:"orderDate"`
	DeliveryDate string  `json:"deliveryDate"`
	PaymentDate  string  `json:"paymentDate"`
	TotalAmount  float64 `json:"totalAmount"`
	Status       string
}

func (p UpdateMaterialOrderParams) validate() error {
	return nil
}

type MaterialOrderHandler struct {
	store *db.Store
}

func NewMaterialOrderHandler(store *db.Store) *MaterialOrderHandler {
	return &MaterialOrderHandler{
		store: store,
	}
}

// HandleGetMaterialOrders retrieves a list of material orders based on query parameters.
//
// @Summary Get material orders
// @Description Retrieves a list of material orders based on query parameters.
// @Tags MaterialOrder
// @Param id query string false "Material Order ID"
// @Param sellerId query string false "Seller ID"
// @Param status query string false "Material Order status"
// @Produce json
// @Success 200 {array} types.MaterialOrder
// @Router /materialOrder [get]
func (h *MaterialOrderHandler) HandleGetMaterialOrders(c *fiber.Ctx) error {
	id := c.Query("id")
	sellerID := c.Query("sellerId")
	status := c.Query("status")

	filter := bson.M{}

	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid material order ID",
			})
		}
		filter["_id"] = objID
	}

	if sellerID != "" {
		objID, err := primitive.ObjectIDFromHex(sellerID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid seller ID",
			})
		}
		filter["sellerId"] = objID
	}

	if status != "" {
		filter["status"] = status
	}

	materialOrders, err := h.store.MaterialOrder.GetMaterialOrders(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(materialOrders) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No matching data found",
		})
	}

	return c.JSON(materialOrders)
}

// HandleInsertMaterialOrder inserts a new material order.
//
// @Summary Insert material order
// @Description Inserts a new material order.
// @Tags MaterialOrder
// @Accept json
// @Produce json
// @Param materialOrder body InsertMaterialOrderParams true "Material Order information"
// @Success 200 {object} types.MaterialOrder
// @Router /materialOrder [post]
func (h *MaterialOrderHandler) HandleInsertMaterialOrder(c *fiber.Ctx) error {
	var params InsertMaterialOrderParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	orderDateParsed, err := time.Parse(time.RFC3339Nano, params.OrderDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order date format",
		})
	}

	materialOrderItems := make([]types.MaterialOrderItem, len(params.MaterialOrderItems))
	for i, item := range params.MaterialOrderItems {
		materialID, err := primitive.ObjectIDFromHex(item.Material.MaterialID)
		if err != nil {
			return err
		}
		material := types.MaterialOrderMaterial{
			MaterialID: materialID,
			Name:       item.Material.Name,
			Price:      item.Material.Price,
			Color:      item.Material.Color,
			Size:       item.Material.Size,
			Remarks:    item.Material.Remarks,
		}

		materialOrderItems[i] = types.MaterialOrderItem{
			Material:   material,
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
		}

		// update materail information
		if strings.ToLower(params.Status) == "completed" {
			m, err := h.store.Material.GetMaterial(c.Context(), materialID)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					return err
				}

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Material doesn't exist, please create the material first.",
				})
			} else {
				material := *m // make a copy
				material.Name = item.Material.Name
				material.Color = item.Material.Color
				material.Size = item.Material.Size
				material.Quantity += item.Material.Quantity
				material.Remarks = item.Material.Remarks

				priceHistoryEntry := types.PriceHistoryEntry{
					Price:     item.Material.Price,
					UpdatedAt: orderDateParsed,
				}

				material.PriceHistory = append(material.PriceHistory, priceHistoryEntry)

				_, err = h.store.Material.UpdateMaterial(c.Context(), materialID, &material)
				if err != nil {
					return err
				}
			}

		}

	}

	materialOrder := types.MaterialOrder{
		SellerID:           params.SellerID,
		OrderDate:          orderDateParsed,
		TotalAmount:        params.TotalAmount,
		Status:             params.Status,
		MaterialOrderItems: materialOrderItems,
	}

	if params.DeliveryDate != "" {
		deliveryDateParsed, err := time.Parse(time.RFC3339Nano, params.DeliveryDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid delivery date format",
			})
		}
		materialOrder.DeliveryDate = deliveryDateParsed
	}

	if params.PaymentDate != "" {
		paymentDateParsed, err := time.Parse(time.RFC3339Nano, params.PaymentDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid payment date format",
			})
		}
		materialOrder.PaymentDate = paymentDateParsed
	}

	inserted, err := h.store.MaterialOrder.InsertMaterialOrder(c.Context(), &materialOrder)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

// HandleUpdateMaterialOrder updates an existing material order in the system.
// @Summary Update material order
// @Description Update an existing material order in the system.
// @Tags MaterialOrder
// @Accept json
// @Produce json
// @Param id path string true "Material Order ID"
// @Param body body UpdateMaterialOrderParams true "Updated material order details"
// @Success 200 {object} fiber.Map
// @Router /materialOrder/{id} [patch]
func (h *MaterialOrderHandler) HandleUpdateMaterialOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	materialOrderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateMaterialOrderParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	orderDateParsed, err := time.Parse(time.RFC3339Nano, params.OrderDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order date format",
		})
	}

	mo, err := h.store.MaterialOrder.GetMaterialOrder(c.Context(), materialOrderID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Material order doesn't exist.",
		})
	}

	// update materail information if the status changed into completed status
	if strings.ToLower(mo.Status) != "completed" && strings.ToLower(params.Status) == "completed" {
		// update all material items in the material order
		for _, item := range mo.MaterialOrderItems {
			updatedCount, err := h.store.Material.IncreaseMaterialQuantity(c.Context(), item.Material.MaterialID, item.Quantity)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to increase material %s by %d, %s", item.Material.MaterialID, item.Quantity, err.Error()),
				})
			}

			if updatedCount == 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to increase material %s by %d", item.Material.MaterialID, item.Quantity),
				})
			}
		}
	}

	// Decrease material amount if the status changed into cancelled status
	if strings.ToLower(mo.Status) == "completed" && strings.ToLower(params.Status) == "cancelled" {
		for _, item := range mo.MaterialOrderItems {
			updatedCount, err := h.store.Material.DecreaseMaterialQuantity(c.Context(), item.Material.MaterialID, item.Quantity)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to decrease material %s by %d, %s", item.Material.MaterialID, item.Quantity, err.Error()),
				})
			}

			if updatedCount == 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to decrease material %s by %d", item.Material.MaterialID, item.Quantity),
				})
			}
		}
	}

	updatedMaterialOrder := types.MaterialOrder{
		SellerID:    params.SellerID,
		OrderDate:   orderDateParsed,
		TotalAmount: params.TotalAmount,
		Status:      params.Status,
	}

	if params.DeliveryDate != "" {
		deliveryDateParsed, err := time.Parse(time.RFC3339Nano, params.DeliveryDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid delivery date format",
			})
		}
		updatedMaterialOrder.DeliveryDate = deliveryDateParsed
	}

	if params.PaymentDate != "" {
		paymentDateParsed, err := time.Parse(time.RFC3339Nano, params.PaymentDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid payment date format",
			})
		}
		updatedMaterialOrder.PaymentDate = paymentDateParsed
	}

	updateCount, err := h.store.MaterialOrder.UpdateMaterialOrder(c.Context(), materialOrderID, &updatedMaterialOrder)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material Order not found or not updated",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material Order updated successfully",
	})
}

// HandleDeleteMaterialOrder deletes a material order by ID.
//
// @Summary Delete material order
// @Description Deletes a material order by ID.
// @Tags MaterialOrder
// @Param id path string true "Material Order ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /materialOrder/{id} [delete]
func (h *MaterialOrderHandler) HandleDeleteMaterialOrder(c *fiber.Ctx) error {
	materialOrderID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(materialOrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid material order ID",
		})
	}

	deleteCount, err := h.store.MaterialOrder.DeleteMaterialOrder(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material Order not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material Order deleted successfully",
	})
}
