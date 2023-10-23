package api

import (
	"strconv"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertMaterialParams struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Size     string `json:"size"`
	Quantity string `json:"quantity"`
	Remarks  string `json:"remarks"`
}

func (p InsertMaterialParams) validate() error {

	return nil
}

type UpdateMaterialParams struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Size     string `json:"size"`
	Quantity string `json:"quantity"`
	Remarks  string `json:"remarks"`
}

func (p UpdateMaterialParams) validate() error {

	return nil
}

type MaterialHandler struct {
	store *db.Store
}

func NewMaterialHandler(store *db.Store) *MaterialHandler {
	return &MaterialHandler{
		store: store,
	}
}

// HandleGetMaterials retrieves a list of materials based on the provided filters.
// @Summary Get materials
// @Description Get a list of materials based on the provided filters.
// @Tags Material
// @Produce json
// @Param id query string false "Material ID (optional)"
// @Param name query string false "Material name (optional)"
// @Param color query string false "Material color (optional)"
// @Param size query string false "Material size (optional)"
// @Param quantity query string false "Material quantity (optional)"
// @Param remarks query string false "Material remarks (optional)"
// @Success 200 {array} types.Material
// @Router /material [get]
func (h *MaterialHandler) HandleGetMaterials(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")
	color := c.Query("color")
	size := c.Query("size")
	quantity := c.Query("quantity")
	remarks := c.Query("remarks")

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
	if name != "" {
		filter["name"] = name
	}
	if color != "" {
		filter["color"] = color
	}
	if size != "" {
		filter["size"] = size
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
	if remarks != "" {
		filter["remarks"] = remarks
	}

	materials, err := h.store.Material.GetMaterials(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(materials) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(materials)
}

// HandleInsertMaterial inserts a new material into the system.
// @Summary Insert material
// @Description Insert a new material into the system.
// @Tags Material
// @Accept json
// @Produce json
// @Param body body InsertMaterialParams true "Material details"
// @Success 200 {object} types.Material
// @Router /material [post]
func (h *MaterialHandler) HandleInsertMaterial(c *fiber.Ctx) error {
	var params InsertMaterialParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	material := types.Material{
		Name:     params.Name,
		Color:    params.Color,
		Size:     params.Size,
		Quantity: params.Quantity,
		Remarks:  params.Remarks,
	}

	inserted, err := h.store.Material.InsertMaterial(c.Context(), &material)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

// HandleUpdateMaterial updates an existing material in the system.
// @Summary Update material
// @Description Update an existing material in the system.
// @Tags Material
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param body body UpdateMaterialParams true "Updated material details"
// @Success 200 {object} fiber.Map
// @Router /material/{id} [patch]
func (h *MaterialHandler) HandleUpdateMaterial(c *fiber.Ctx) error {
	id := c.Params("id")
	materialID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var params UpdateMaterialParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	updatedMaterial := types.Material{
		Name:     params.Name,
		Color:    params.Color,
		Size:     params.Size,
		Quantity: params.Quantity,
		Remarks:  params.Remarks,
	}

	updateCount, err := h.store.Material.UpdateMaterial(c.Context(), materialID, &updatedMaterial)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Material updaated successfully",
	})
}

// HandleDeleteMaterial deletes a material from the system.
// @Summary Delete material
// @Description Delete a material from the system.
// @Tags Material
// @Param id path string true "Material ID"
// @Success 200 {object} fiber.Map
// @Router /material/{id} [delete]
func (h *MaterialHandler) HandleDeleteMaterial(c *fiber.Ctx) error {
	materialID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(materialID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid material ID",
		})
	}

	deleteCount, err := h.store.Material.DeleteMaterial(c.Context(), objID)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Material not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Material deleted successfully",
	})
}

// HandleGetMaterialColors retrieves a list of unique material colors.
// @Summary Get material colors
// @Description Get a list of unique material colors.
// @Tags Material
// @Produce json
// @Success 200 {array} string
// @Router /material/colors [get]
func (h *MaterialHandler) HandleGetMaterialColors(c *fiber.Ctx) error {
	colors, err := h.store.Material.GetMaterialColors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve material colors",
		})
	}

	return c.JSON(colors)
}

// HandleGetMaterialSizes retrieves a list of unique material sizes.
// @Summary Get material sizes
// @Description Get a list of unique material sizes.
// @Tags Material
// @Produce json
// @Success 200 {array} string
// @Router /material/sizes [get]
func (h *MaterialHandler) HandleGetMaterialSizes(c *fiber.Ctx) error {
	sizes, err := h.store.Material.GetMaterialSizes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve material sizes",
		})
	}

	return c.JSON(sizes)
}
