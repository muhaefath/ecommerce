package service

import (
	"context"
	"database/sql"
	"ecommerce/model/entity"
	"ecommerce/model/request"
	"ecommerce/model/response"
	"ecommerce/repository"
	"math"
)

type ecommerceService struct {
	ecommerceRepo repository.EcommerceProvider
}

type EcommerceConfig struct {
	EcommerceRepo repository.EcommerceProvider
}

func NewEcommerceService(config EcommerceConfig) ecommerceService {
	ecommerceProvider := ecommerceService{
		ecommerceRepo: config.EcommerceRepo,
	}

	return ecommerceProvider
}

func (e *ecommerceService) GetProductList(ctx context.Context, payload request.FilterProduct) (response.GetProductListResponse, error) {
	var resp response.GetProductListResponse

	products, err := e.ecommerceRepo.GetProductList(ctx, payload)
	if err != nil {
		return resp, err
	}

	resp.Data = products
	return resp, nil
}

func (e *ecommerceService) CreateProduct(ctx context.Context, request request.UpsertProduct) (err error) {
	product := entity.Product{
		UserID:      request.UserID,
		Sku:         request.Sku,
		Title:       request.Title,
		Category:    request.Category,
		Description: request.Description,
		Etalase:     request.Etalase,
		Price:       request.Price,
		Weight:      request.Weight,
	}

	productID, err := e.ecommerceRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	for _, v := range request.ProductImages {
		productImage := entity.ProductImage{
			ProductID: productID,
			ImageUrl:  v.ImageUrl,
		}

		if v.ShortDescription != "" {
			productImage.ShortDescription = sql.NullString{
				Valid:  true,
				String: v.ShortDescription,
			}
		}

		err = e.ecommerceRepo.CreateProductImages(ctx, productImage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ecommerceService) UpdateProduct(ctx context.Context, id int64, request request.UpsertProduct) (err error) {
	product, err := e.ecommerceRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	productRequest := entity.Product{
		ID:          id,
		UserID:      request.UserID,
		Sku:         request.Sku,
		Title:       request.Title,
		Category:    request.Category,
		Description: request.Description,
		Etalase:     request.Etalase,
		Price:       request.Price,
		Weight:      request.Weight,
		Rating:      product.Rating,
	}

	err = e.ecommerceRepo.UpdateProduct(ctx, productRequest)
	if err != nil {
		return err
	}

	err = e.ecommerceRepo.DeleteProductImagesByID(ctx, id)
	if err != nil {
		return err
	}

	for _, v := range request.ProductImages {
		productImage := entity.ProductImage{
			ProductID: id,
			ImageUrl:  v.ImageUrl,
		}

		if v.ShortDescription != "" {
			productImage.ShortDescription = sql.NullString{
				Valid:  true,
				String: v.ShortDescription,
			}
		}

		err = e.ecommerceRepo.CreateProductImages(ctx, productImage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ecommerceService) GetProductByID(ctx context.Context, id int64) (response.GetProductDetailResponse, error) {
	var resp response.GetProductDetailResponse

	products, err := e.ecommerceRepo.GetProductByID(ctx, id)
	if err != nil {
		return resp, err
	}

	productImages, err := e.ecommerceRepo.GetProductImagesByProductID(ctx, id)
	if err != nil {
		return resp, err
	}

	productReview, err := e.ecommerceRepo.GetProductReviewByProductID(ctx, id)
	if err != nil {
		return resp, err
	}

	resp.Data.Product = products
	resp.Data.ProductImages = productImages
	resp.Data.Review = productReview

	return resp, nil
}

func (e *ecommerceService) CreateProductReview(ctx context.Context, request request.UpsertProductReview) (err error) {

	productReview := entity.ProductReview{
		ProductID: request.ProductID,
		Rating:    request.Rating,
	}

	if request.Comment != "" {
		productReview.Comment = sql.NullString{
			Valid:  true,
			String: request.Comment,
		}
	}

	err = e.ecommerceRepo.CreateProductReview(ctx, productReview)
	if err != nil {
		return err
	}

	product, err := e.ecommerceRepo.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return err
	}

	productReviews, err := e.ecommerceRepo.GetProductReviewByProductID(ctx, request.ProductID)
	if err != nil {
		return err
	}

	totalRating := 0
	for _, v := range productReviews {
		totalRating += v.Rating
	}

	rating := float64(totalRating) / float64(len(productReviews))
	product.Rating = math.Round(rating*10) / 10

	err = e.ecommerceRepo.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
