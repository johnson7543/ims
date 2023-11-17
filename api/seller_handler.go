package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertSellerParams struct {
	Company     string `json:"company"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	TaxIdNumber string `json:"taxIdNumber"`
}

func (p InsertSellerParams) validate() error {
	return nil
}

type UpdateSellerParams struct {
	Company     string `json:"company" form:"company"`
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	Address     string `json:"address" form:"address"`
	TaxIdNumber string `json:"taxIdNumber" form:"taxIdNumber"`
}

func (p *UpdateSellerParams) validate() error {
	return nil
}

type SellerHandler struct {
	store *db.Store
}

func NewSellerHandler(store *db.Store) *SellerHandler {
	return &SellerHandler{
		store: store,
	}
}

// HandleGetSellers retrieves a list of sellers based on query parameters.
//
// @Summary Get sellers
// @Description Retrieves a list of sellers based on query parameters.
// @Tags Seller
// @Param id query string false "Seller ID"
// @Param company query string false "Company name"
// @Param name query string false "Seller name"
// @Param phone query string false "Phone number"
// @Param address query string false "Address"
// @Param taxIdNumber query string false "Tax ID number"
// @Produce json
// @Success 200 {array} types.Seller
// @Router /seller [get]
func (h *SellerHandler) HandleGetSellers(c *fiber.Ctx) error {
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
				"error": "Invalid seller ID",
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

	sellers, err := h.store.Seller.GetSellers(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(sellers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(sellers)
}

// HandleInsertSeller inserts a new seller.
//
// @Summary Insert seller
// @Description Inserts a new seller.
// @Tags Seller
// @Accept json
// @Produce json
// @Param seller body InsertSellerParams true "Seller information"
// @Success 200 {object} types.Seller
// @Router /seller [post]
func (h *SellerHandler) HandleInsertSeller(c *fiber.Ctx) error {
	var params InsertSellerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	seller := types.Seller{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	inserted, err := h.store.Seller.InsertSeller(c.Context(), &seller)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

// HandleUpdateSeller updates an existing seller in the system.
// @Summary Update seller
// @Description Update an existing seller in the system.
// @Tags Seller
// @Accept json
// @Produce json
// @Param id path string true "Seller ID"
// @Param body body UpdateSellerParams true "Updated seller details"
// @Success 200 {object} fiber.Map
// @Router /seller/{id} [patch]
func (h *SellerHandler) HandleUpdateSeller(c *fiber.Ctx) error {
	id := c.Params("id")
	sellerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateSellerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	updatedSeller := types.Seller{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	updateCount, err := h.store.Seller.UpdateSeller(c.Context(), sellerID, &updatedSeller)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Seller not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Seller updated successfully",
	})
}

// HandleDeleteSeller deletes a seller by ID.
//
// @Summary Delete seller
// @Description Deletes a seller by ID.
// @Tags Seller
// @Param id path string true "Seller ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /seller/{id} [delete]
func (h *SellerHandler) HandleDeleteSeller(c *fiber.Ctx) error {
	sellerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(sellerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid seller ID",
		})
	}

	deleteCount, err := h.store.Seller.DeleteSeller(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Seller not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Seller deleted successfully",
	})
}
