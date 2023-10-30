package main

import (
	"ecommerce/httpservice"
	"ecommerce/internal"
	"ecommerce/repository/postgre"
	"ecommerce/service"

	support "ecommerce/utils/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db := internal.NewDatabases(internal.InitConfig(), support.NewLogger())
	ecommerceRepo := postgre.NewEcommerce(db["main"])

	ecommerceService := service.NewEcommerceService(
		service.EcommerceConfig{
			EcommerceRepo: ecommerceRepo,
		},
	)
	httpService := httpservice.NewHandler(httpservice.HandlerConfig{
		EcommerceSrv: &ecommerceService,
	})

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api") // /api

	productApi := api.Group("/product") // /api

	productApi.Get("/list", httpService.GetProductList)
	productApi.Post("/", httpService.CreateProduct)
	productApi.Put("/:product_id", httpService.UpdateProduct)
	productApi.Post("/review", httpService.CreateProductReview)
	productApi.Get("/:product_id", httpService.GetDetailProduct)

	app.Listen(":3000")
}
