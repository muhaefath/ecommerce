package httpservice

import (
	"ecommerce/model/request"
	"ecommerce/model/response"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// NewHandler is method to create handler httpservice for event
func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		ecommerceSrv: cfg.EcommerceSrv,
	}
}

// GetProductList is a handler to get product list
// GetProductList godoc
// @Summary      get list of product
// @Description  get list of product
// @Tags         Product
// @Param FilterProduct body request.FilterProduct true "FilterProduct"

// @Success 200 {object} response.GetProductListResponse{}
// @Failure 400 {object} response.Error{}
// @ID v1-GetProductList
// @Router       /product/list   [get]
func (d *Handler) GetProductList(c *fiber.Ctx) error {
	request := request.FilterProduct{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	resp, err := d.ecommerceSrv.GetProductList(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	resp.StatusCode = http.StatusOK
	resp.Message = "success"

	return c.Status(http.StatusOK).JSON(resp)
}

// CreateProduct is a handler to create a product
// CreateProduct godoc
// @Summary      create a product
// @Description  create a product
// @Tags         Product
// @Param UpsertProduct body request.UpsertProduct true "UpsertProduct"
// @Success 200 {object} response.BaseResponse{}
// @Failure 400 {object} response.Error{}
// @ID v1-CreateProduct
// @Router       /product   [post]
func (d *Handler) CreateProduct(c *fiber.Ctx) error {
	request := request.UpsertProduct{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	err := d.ecommerceSrv.CreateProduct(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "success",
	})
}

// UpdateProduct is a handler to update a product
// UpdateProduct godoc
// @Summary      update a product
// @Description  update a product
// @Tags         Product
// @Param 	product_id path  string true "Product ID"
// @Param UpsertProduct body request.UpsertProduct true "UpsertProduct"
// @Success 200 {object} response.BaseResponse{}
// @Failure 400 {object} response.Error{}
// @ID v1-UpdateProduct
// @Router       /product/{product_id}    [put]
func (d *Handler) UpdateProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "product_id can'b be null and should be an integer",
		})
	}

	request := request.UpsertProduct{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	err = d.ecommerceSrv.UpdateProduct(c.Context(), int64(productID), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	})
}

// CreateProductReview is a handler to create a product review
// CreateProductReview godoc
// @Summary      create a product review
// @Description  create a product review
// @Tags         Product
// @Param UpsertProductReview body request.UpsertProductReview true "UpsertProductReview"
// @Success 200 {object} response.BaseResponse{}
// @Failure 400 {object} response.Error{}
// @ID v1-CreateProductReview
// @Router       /product/review   [post]
func (d *Handler) CreateProductReview(c *fiber.Ctx) error {
	request := request.UpsertProductReview{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	err := d.ecommerceSrv.CreateProductReview(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "success",
	})
}

// GetDetailProduct is a handler to get a product
// GetDetailProduct godoc
// @Summary      get a product
// @Description  get a product
// @Tags         Product
// @Param 	product_id path  string true "Product ID"
// @Success 200 {object} response.BaseResponse{}
// @Failure 400 {object} response.Error{}
// @ID v1-GetDetailProduct
// @Router       /product/{product_id}   [get]
func (d *Handler) GetDetailProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "product_id can'b be null and should be an integer",
		})
	}

	resp, err := d.ecommerceSrv.GetProductByID(c.Context(), int64(productID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	resp.StatusCode = http.StatusOK
	resp.Message = "success"

	return c.Status(http.StatusOK).JSON(resp)
}
