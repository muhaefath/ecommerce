package repository

import (
	"context"
	"database/sql"
	"ecommerce/model/entity"
	"ecommerce/model/request"
)

type QueryProvider interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type EcommerceProvider interface {
	GetProductList(ctx context.Context, payload request.FilterProduct) (response []entity.Product, err error)
	CreateProduct(ctx context.Context, request entity.Product) (id int64, err error)
	UpdateProduct(ctx context.Context, request entity.Product) (err error)
	GetProductByID(ctx context.Context, id int64) (response entity.Product, err error)
	GetProductImagesByProductID(ctx context.Context, id int64) (response []entity.ProductImage, err error)
	GetProductReviewByProductID(ctx context.Context, id int64) (response []entity.ProductReview, err error)
	CreateProductReview(ctx context.Context, payload entity.ProductReview) (err error)
	CreateProductImages(ctx context.Context, payload entity.ProductImage) (err error)
	DeleteProductImagesByID(ctx context.Context, productID int64) (err error)
}
