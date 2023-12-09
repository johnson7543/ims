package api

import (
	"fmt"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertBuyerParams struct {
	Company     string `json:"company"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	TaxIdNumber string `json:"taxIdNumber"`
}

func (p InsertBuyerParams) validate() error {
	return nil
}

type UpdateBuyerParams struct {
	Company     string `json:"company" form:"company"`
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	Address     string `json:"address" form:"address"`
	TaxIdNumber string `json:"taxIdNumber" form:"taxIdNumber"`
}

func (p *UpdateBuyerParams) validate() error {
	return nil
}

type BuyerHandler struct {
	store *db.Store
}

func NewBuyerHandler(store *db.Store) *BuyerHandler {
	return &BuyerHandler{
		store: store,
	}
}

// HandleGetBuyers retrieves a list of buyers based on query parameters.
//
// @Summary Get buyers
// @Description Retrieves a list of buyers based on query parameters.
// @Tags Buyer
// @Param id query string false "Buyer ID"
// @Param company query string false "Company name"
// @Param name query string false "Buyer name"
// @Param phone query string false "Phone number"
// @Param address query string false "Address"
// @Param taxIdNumber query string false "Tax ID number"
// @Produce json
// @Success 200 {array} types.Buyer
// @Router /buyer [get]
func (h *BuyerHandler) HandleGetBuyers(c *fiber.Ctx) error {
	id := c.Query("id")
	company := c.Query("company")
	name := c.Query("name")
	phone := c.Query("phone")
	address := c.Query("address")
	taxIdNumber := c.Query("taxIdNumber")

	filter := bson.M{}

	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid buyer ID",
			})
		}
		filter["_id"] = objID
	}
	if company != "" {
		filter["company"] = company
	}
	if name != "" {
		filter["name"] = name
	}
	if phone != "" {
		filter["phone"] = phone
	}
	if address != "" {
		filter["address"] = address
	}
	if taxIdNumber != "" {
		filter["taxIdNumber"] = taxIdNumber
	}

	buyers, err := h.store.Buyer.GetBuyers(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(buyers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(buyers)
}

// HandleInsertBuyer inserts a new buyer.
//
// @Summary Insert buyer
// @Description Inserts a new buyer.
// @Tags Buyer
// @Accept json
// @Produce json
// @Param buyer body InsertBuyerParams true "Buyer information"
// @Success 200 {object} types.Buyer
// @Router /buyer [post]
func (h *BuyerHandler) HandleInsertBuyer(c *fiber.Ctx) error {
	var params InsertBuyerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	buyer := types.Buyer{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	inserted, err := h.store.Buyer.InsertBuyer(c.Context(), &buyer)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Buyer inserted successfully, ID: %s, Name: %s", inserted.ID.Hex(), inserted.Name),
	})
}

// HandleUpdateBuyer updates an existing buyer in the system.
// @Summary Update buyer
// @Description Update an existing buyer in the system.
// @Tags Buyer
// @Accept json
// @Produce json
// @Param id path string true "Buyer ID"
// @Param body body UpdateBuyerParams true "Updated buyer details"
// @Success 200 {object} fiber.Map
// @Router /buyer/{id} [patch]
func (h *BuyerHandler) HandleUpdateBuyer(c *fiber.Ctx) error {
	id := c.Params("id")
	buyerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateBuyerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	updatedBuyer := types.Buyer{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	updateCount, err := h.store.Buyer.UpdateBuyer(c.Context(), buyerID, &updatedBuyer)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Buyer not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Buyer updated successfully",
	})
}

// HandleDeleteBuyer deletes a buyer by ID.
//
// @Summary Delete buyer
// @Description Deletes a buyer by ID.
// @Tags Buyer
// @Param id path string true "Buyer ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /buyer/{id} [delete]
func (h *BuyerHandler) HandleDeleteBuyer(c *fiber.Ctx) error {
	buyerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(buyerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid buyer ID",
		})
	}

	deleteCount, err := h.store.Buyer.DeleteBuyer(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Buyer not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Buyer deleted successfully",
	})
}
