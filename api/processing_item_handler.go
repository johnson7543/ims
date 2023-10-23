package api

import (
	"strconv"
	"time"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertProcessingItemParams struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	WorkerID  string  `json:"workerId"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	SKU       string  `json:"sku"`
	Remarks   string  `json:"remarks"`
}

func (p InsertProcessingItemParams) validate() error {
	// Add validation logic if needed
	return nil
}

type UpdateProcessingItemParams struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	WorkerID  string  `json:"workerId"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	SKU       string  `json:"sku"`
	Remarks   string  `json:"remarks"`
}

func (p *UpdateProcessingItemParams) validate() error {
	// You can add validation logic here if needed.
	// For example, validate that the date formats are correct.
	return nil
}

type ProcessingItemHandler struct {
	store *db.Store
}

func NewProcessingItemHandler(store *db.Store) *ProcessingItemHandler {
	return &ProcessingItemHandler{
		store: store,
	}
}

// HandleGetProcessingItems retrieves a list of processing items based on query parameters.
//
// @Summary Get processing items
// @Description Retrieves a list of processing items based on query parameters.
// @Tags Processing Item
// @Param id query string false "Processing item ID"
// @Param name query string false "Processing item name"
// @Param quantity query string false "Quantity"
// @Param price query string false "Price"
// @Param workerId query string false "Worker ID"
// @Param startDate query string false "Start date (format: YYYY-MM-DD)"
// @Param endDate query string false "End date (format: YYYY-MM-DD)"
// @Param sku query string false "Product ID"
// @Param remarks query string false "Remarks"
// @Produce json
// @Success 200 {array} types.ProcessingItem
// @Router /processingItem [get]
func (h *ProcessingItemHandler) HandleGetProcessingItems(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")
	quantity := c.Query("quantity")
	price := c.Query("price")
	workerID := c.Query("workerId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	productID := c.Query("sku")
	remarks := c.Query("remarks")

	filter := bson.M{}

	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid processing item ID",
			})
		}
		filter["_id"] = objID
	}
	if name != "" {
		filter["name"] = name
	}
	if quantity != "" {
		quantity, err := strconv.Atoi(quantity)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid quantity value",
			})
		}
		filter["quantity"] = quantity
	}
	if price != "" {
		price, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid price value",
			})
		}
		filter["price"] = price
	}
	if workerID != "" {
		objID, err := primitive.ObjectIDFromHex(workerID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid worker ID",
			})
		}
		filter["workerId"] = objID
	}
	if startDate != "" {
		startDateParsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format",
			})
		}
		filter["startDate"] = startDateParsed
		if endDate != "" {
			endDateParsed, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid end date format",
				})
			}
			if startDateParsed.After(endDateParsed) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Start date cannot be after end date",
				})
			}
			filter["endDate"] = endDateParsed
		}
	}
	if productID != "" {
		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}
		filter["sku"] = objID
	}
	if remarks != "" {
		filter["remarks"] = remarks
	}

	processingItems, err := h.store.ProcessingItem.GetProcessingItems(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(processingItems) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(processingItems)
}

// HandleInsertProcessingItem inserts a new processing item.
//
// @Summary Insert processing item
// @Description Inserts a new processing item.
// @Tags Processing Item
// @Accept json
// @Produce json
// @Param processingItem body InsertProcessingItemParams true "Processing item information"
// @Success 200 {object} types.ProcessingItem
// @Router /processingItem [post]
func (h *ProcessingItemHandler) HandleInsertProcessingItem(c *fiber.Ctx) error {
	var params InsertProcessingItemParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	workerID, err := primitive.ObjectIDFromHex(params.WorkerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid worker ID",
		})
	}

	processingItem := types.ProcessingItem{
		Name:     params.Name,
		Quantity: params.Quantity,
		Price:    params.Price,
		WorkerID: workerID,
		Remarks:  params.Remarks,
	}

	if params.StartDate != "" {
		startDateParsed, err := time.Parse("2006-01-02", params.StartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format",
			})
		}
		processingItem.StartDate = startDateParsed
	}

	if params.EndDate != "" {
		endDateParsed, err := time.Parse("2006-01-02", params.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end date format",
			})
		}
		processingItem.EndDate = endDateParsed
	}

	if params.SKU != "" {
		productID, err := primitive.ObjectIDFromHex(params.SKU)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}
		processingItem.SKU = productID
	}

	inserted, err := h.store.ProcessingItem.InsertProcessingItem(c.Context(), &processingItem)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

// HandleUpdateProcessingItem updates an existing processing item in the system.
//
// @Summary Update processing item
// @Description Update an existing processing item in the system.
// @Tags Processing Item
// @Accept json
// @Produce json
// @Param id path string true "Processing Item ID"
// @Param body body UpdateProcessingItemParams true "Updated processing item details"
// @Success 200 {object} fiber.Map
// @Router /processingItem/{id} [patch]
func (h *ProcessingItemHandler) HandleUpdateProcessingItem(c *fiber.Ctx) error {
	id := c.Params("id")
	processingItemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Processing Item ID",
		})
	}

	var params UpdateProcessingItemParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	workerID, err := primitive.ObjectIDFromHex(params.WorkerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid worker ID",
		})
	}

	updatedProcessingItem := types.ProcessingItem{
		Name:     params.Name,
		Quantity: params.Quantity,
		Price:    params.Price,
		WorkerID: workerID,
		Remarks:  params.Remarks,
	}

	if params.StartDate != "" {
		startDateParsed, err := time.Parse("2006-01-02", params.StartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format",
			})
		}
		updatedProcessingItem.StartDate = startDateParsed
	}

	if params.EndDate != "" {
		endDateParsed, err := time.Parse("2006-01-02", params.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end date format",
			})
		}
		updatedProcessingItem.EndDate = endDateParsed
	}

	if params.SKU != "" {
		productID, err := primitive.ObjectIDFromHex(params.SKU)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}
		updatedProcessingItem.SKU = productID
	}

	updateCount, err := h.store.ProcessingItem.UpdateProcessingItem(c.Context(), processingItemID, &updatedProcessingItem)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Processing Item not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Processing Item updated successfully",
	})
}

// HandleDeleteProcessingItem deletes a processing item by ID.
//
// @Summary Delete processing item
// @Description Deletes a processing item by ID.
// @Tags Processing Item
// @Param id path string true "Processing item ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /processingItem/{id} [delete]
func (h *ProcessingItemHandler) HandleDeleteProcessingItem(c *fiber.Ctx) error {
	processingItemID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(processingItemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid processing item ID",
		})
	}

	deleteCount, err := h.store.ProcessingItem.DeleteProcessingItem(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Processing item deleted successfully",
	})
}
