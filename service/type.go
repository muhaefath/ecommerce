package service

import (
	"context"
	"ecommerce/model/request"
	"ecommerce/model/response"
)

type EcommerceProvider interface {
	GetProductList(context.Context, request.FilterProduct) (response.GetProductListResponse, error)
	CreateProduct(ctx context.Context, request request.UpsertProduct) (err error)
	UpdateProduct(ctx context.Context, id int64, request request.UpsertProduct) (err error)
	GetProductByID(ctx context.Context, id int64) (response response.GetProductDetailResponse, err error)
	CreateProductReview(ctx context.Context, request request.UpsertProductReview) (err error)
}
