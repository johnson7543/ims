package api

import (
	"fmt"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertCustomerParams struct {
	Company     string `json:"company"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	TaxIdNumber string `json:"taxIdNumber"`
}

func (p InsertCustomerParams) validate() error {
	return nil
}

type UpdateCustomerParams struct {
	Company     string `json:"company" form:"company"`
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	Address     string `json:"address" form:"address"`
	TaxIdNumber string `json:"taxIdNumber" form:"taxIdNumber"`
}

func (p *UpdateCustomerParams) validate() error {
	return nil
}

type CustomerHandler struct {
	store *db.Store
}

func NewCustomerHandler(store *db.Store) *CustomerHandler {
	return &CustomerHandler{
		store: store,
	}
}

// HandleGetCustomers retrieves a list of customers based on query parameters.
//
// @Summary Get customers
// @Description Retrieves a list of customers based on query parameters.
// @Tags Customer
// @Param id query string false "Customer ID"
// @Param company query string false "Company name"
// @Param name query string false "Customer name"
// @Param phone query string false "Phone number"
// @Param address query string false "Address"
// @Param taxIdNumber query string false "Tax ID number"
// @Produce json
// @Success 200 {array} types.Customer
// @Router /customer [get]
func (h *CustomerHandler) HandleGetCustomers(c *fiber.Ctx) error {
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
				"error": "Invalid customer ID",
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

	customers, err := h.store.Customer.GetCustomers(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(customers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(customers)
}

// HandleInsertCustomer inserts a new customer.
//
// @Summary Insert customer
// @Description Inserts a new customer.
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body InsertCustomerParams true "Customer information"
// @Success 200 {object} types.Customer
// @Router /customer [post]
func (h *CustomerHandler) HandleInsertCustomer(c *fiber.Ctx) error {
	var params InsertCustomerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	customer := types.Customer{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	inserted, err := h.store.Customer.InsertCustomer(c.Context(), &customer)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Customer inserted successfully, ID: %s, Name: %s", inserted.ID, inserted.Name),
	})
}

// HandleUpdateCustomer updates an existing customer in the system.
// @Summary Update customer
// @Description Update an existing customer in the system.
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param body body UpdateCustomerParams true "Updated customer details"
// @Success 200 {object} fiber.Map
// @Router /customer/{id} [patch]
func (h *CustomerHandler) HandleUpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateCustomerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	updatedCustomer := types.Customer{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	updateCount, err := h.store.Customer.UpdateCustomer(c.Context(), customerID, &updatedCustomer)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Customer updated successfully",
	})
}

// HandleDeleteCustomer deletes a customer by ID.
//
// @Summary Delete customer
// @Description Deletes a customer by ID.
// @Tags Customer
// @Param id path string true "Customer ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /customer/{id} [delete]
func (h *CustomerHandler) HandleDeleteCustomer(c *fiber.Ctx) error {
	customerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	deleteCount, err := h.store.Customer.DeleteCustomer(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Customer deleted successfully",
	})
}
