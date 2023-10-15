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

type InsertProductParams struct {
	Name     string    `json:"name"`
	Material string    `json:"material"`
	Color    string    `json:"color"`
	Size     string    `json:"size"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
	Date     time.Time `json:"date"`
	Remark   string    `json:"remark"`
}

func (p InsertProductParams) validate() error {
	// Add validation logic if needed
	return nil
}

type ProductHandler struct {
	store *db.Store
}

func NewProductHandler(store *db.Store) *ProductHandler {
	return &ProductHandler{
		store: store,
	}
}

// HandleGetProducts retrieves a list of products based on query parameters.
//
// @Summary Get products
// @Description Retrieves a list of products based on query parameters.
// @Tags Product
// @Param id query string false "Product ID"
// @Param name query string false "Product name"
// @Param material query string false "Material"
// @Param color query string false "Color"
// @Param size query string false "Size"
// @Param quantity query string false "Quantity"
// @Param price query string false "Price"
// @Param date query string false "Date (format: YYYY-MM-DD)"
// @Param remark query string false "Remark"
// @Produce json
// @Success 200 {array} types.Product
// @Router /product [get]
func (h *ProductHandler) HandleGetProducts(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")
	material := c.Query("material")
	color := c.Query("color")
	size := c.Query("size")
	quantity := c.Query("quantity")
	price := c.Query("price")
	date := c.Query("date")
	remark := c.Query("remark")

	filter := bson.M{}

	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}
		filter["_id"] = objID
	}
	if name != "" {
		filter["name"] = name
	}
	if material != "" {
		filter["material"] = material
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
	if price != "" {
		price, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid price value",
			})
		}
		filter["price"] = price
	}
	if date != "" {
		filter["date"] = date
	}
	if remark != "" {
		filter["remark"] = remark
	}

	products, err := h.store.Product.GetProducts(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No Matches data found",
		})
	}

	return c.JSON(products)
}

// HandleInsertProduct inserts a new product.
//
// @Summary Insert product
// @Description Inserts a new product.
// @Tags Product
// @Accept json
// @Produce json
// @Param product body InsertProductParams true "Product information"
// @Success 200 {object} types.Product
// @Router /product [post]
func (h *ProductHandler) HandleInsertProduct(c *fiber.Ctx) error {
	var params InsertProductParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	product := types.Product{
		Name:     params.Name,
		Material: params.Material,
		Color:    params.Color,
		Size:     params.Size,
		Quantity: params.Quantity,
		Price:    params.Price,
		Date:     params.Date,
		Remark:   params.Remark,
	}

	inserted, err := h.store.Product.InsertProduct(c.Context(), &product)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

// HandleDeleteProduct deletes a product by ID.
//
// @Summary Delete product
// @Description Deletes a product by ID.
// @Tags Product
// @Param id path string true "Product ID"
// @Produce json
// @Success 200 {object} fiber.Map
// @Router /product/{id} [delete]
func (h *ProductHandler) HandleDeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	err = h.store.Product.DeleteProduct(c.Context(), objID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
