package api

import (
	"fmt"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertWorkerParams struct {
	Company     string `json:"company"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	TaxIdNumber string `json:"taxIdNumber"`
}

func (p InsertWorkerParams) validate() error {
	return nil
}

type UpdateWorkerParams struct {
	Company     string `json:"company" form:"company"`
	Name        string `json:"name" form:"name"`
	Phone       string `json:"phone" form:"phone"`
	Address     string `json:"address" form:"address"`
	TaxIdNumber string `json:"taxIdNumber" form:"taxIdNumber"`
}

func (p *UpdateWorkerParams) validate() error {
	return nil
}

type WorkerHandler struct {
	store *db.Store
}

func NewWorkerHandler(store *db.Store) *WorkerHandler {
	return &WorkerHandler{
		store: store,
	}
}

// HandleGetWorkers retrieves a list of workers based on query parameters.
//
// @Summary Get workers
// @Description Retrieves a list of workers based on query parameters.
// @Tags Worker
// @Param id query string false "Worker ID"
// @Param company query string false "Company name"
// @Param name query string false "Worker name"
// @Param phone query string false "Phone number"
// @Param address query string false "Address"
// @Param taxIdNumber query string false "Tax ID number"
// @Produce json
// @Success 200 {array} types.Worker
// @Router /worker [get]
func (h *WorkerHandler) HandleGetWorkers(c *fiber.Ctx) error {
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
				"error": "Invalid worker ID",
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

	workers, err := h.store.Worker.GetWorkers(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(workers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(workers)
}

// HandleInsertWorker inserts a new worker.
//
// @Summary Insert worker
// @Description Inserts a new worker.
// @Tags Worker
// @Accept json
// @Produce json
// @Param worker body InsertWorkerParams true "Worker information"
// @Success 200 {object} types.Worker
// @Router /worker [post]
func (h *WorkerHandler) HandleInsertWorker(c *fiber.Ctx) error {
	var params InsertWorkerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	worker := types.Worker{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	inserted, err := h.store.Worker.InsertWorker(c.Context(), &worker)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Worker inserted successfully, ID: %s, Name: %s", inserted.ID.Hex(), inserted.Name),
	})
}

// HandleUpdateWorker updates an existing worker in the system.
// @Summary Update worker
// @Description Update an existing worker in the system.
// @Tags Worker
// @Accept json
// @Produce json
// @Param id path string true "Worker ID"
// @Param body body UpdateWorkerParams true "Updated worker details"
// @Success 200 {object} fiber.Map
// @Router /worker/{id} [patch]
func (h *WorkerHandler) HandleUpdateWorker(c *fiber.Ctx) error {
	id := c.Params("id")
	workerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateWorkerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	updatedWorker := types.Worker{
		Company:     params.Company,
		Name:        params.Name,
		Phone:       params.Phone,
		Address:     params.Address,
		TaxIdNumber: params.TaxIdNumber,
	}

	updateCount, err := h.store.Worker.UpdateWorker(c.Context(), workerID, &updatedWorker)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Worker not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Worker updated successfully",
	})
}

// HandleDeleteWorker deletes a worker by ID.
//
// @Summary Delete worker
// @Description Deletes a worker by ID.
// @Tags Worker
// @Param id path string true "Worker ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /worker/{id} [delete]
func (h *WorkerHandler) HandleDeleteWorker(c *fiber.Ctx) error {
	workerID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(workerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid worker ID",
		})
	}

	deleteCount, err := h.store.Worker.DeleteWorker(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Worker deleted successfully",
	})
}
