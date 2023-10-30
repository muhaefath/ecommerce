package postgre

import (
	"context"
	"ecommerce/model/entity"
	"ecommerce/model/request"
	"ecommerce/repository"
	sdkSql "ecommerce/utils/sql"
	"fmt"
)

type ecommerceRepo struct {
	baseRepo
}

// NewPromotion is function to initialize promotion repository logic.
func NewEcommerce(db sdkSql.DBer) repository.EcommerceProvider {
	return &ecommerceRepo{
		baseRepo: baseRepo{db: db},
	}
}

func (e *ecommerceRepo) GetProductList(ctx context.Context, payload request.FilterProduct) (response []entity.Product, err error) {
	var products []entity.Product
	sort := "DESC"
	if payload.IsAsc {
		sort = "ASC"
	}

	selectQuery := `
		SELECT
			*
		FROM
			products
		WHERE 
			sku ilike '%' || $1 || '%'
		OR 
			category ilike '%' || $1 || '%'
		OR 
			etalase ilike '%' || $1 || '%'
		OR 
			title ilike '%' || $1 || '%'
	`
	selectQuery += fmt.Sprintf(" ORDER BY %s %s", payload.Sort, sort)

	err = e.db.SelectContext(ctx, &products, selectQuery, payload.Search)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (e *ecommerceRepo) CreateProduct(ctx context.Context, payload entity.Product) (id int64, err error) {
	var lastInsertId int64
	err = e.db.GetContext(ctx, &lastInsertId,
		`INSERT INTO 
		products ( sku, title, description, category, etalase, weight, price, user_id) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`, payload.Sku, payload.Title, payload.Description, payload.Category, payload.Etalase,
		payload.Weight, payload.Price, payload.UserID)

	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (e *ecommerceRepo) UpdateProduct(ctx context.Context, payload entity.Product) (err error) {
	_, err = e.db.ExecContext(ctx,
		`UPDATE
		products
	SET
		sku=$1,
		title=$2,
		description=$3,
		category=$4,
		etalase=$5,
		weight=$6,
		price=$7,
		rating=$8
	WHERE
		id=$9`, payload.Sku, payload.Title, payload.Description, payload.Category, payload.Etalase,
		payload.Weight, payload.Price, payload.Rating, payload.ID)

	if err != nil {
		return err
	}

	return nil
}

func (e *ecommerceRepo) GetProductByID(ctx context.Context, id int64) (response entity.Product, err error) {
	var products entity.Product

	selectQuery := `
		SELECT
			*
		FROM
			products
		WHERE
			id = $1
	`
	err = e.db.GetContext(ctx, &products, selectQuery, id)
	if err != nil {
		return entity.Product{}, err
	}

	return products, nil
}

func (e *ecommerceRepo) GetProductImagesByProductID(ctx context.Context, id int64) (response []entity.ProductImage, err error) {
	var productImage []entity.ProductImage

	selectQuery := `
		SELECT
			*
		FROM
			product_images
		WHERE
			product_id = $1
	`
	err = e.db.SelectContext(ctx, &productImage, selectQuery, id)
	if err != nil {
		return []entity.ProductImage{}, err
	}

	return productImage, nil
}

func (e *ecommerceRepo) GetProductReviewByProductID(ctx context.Context, id int64) (response []entity.ProductReview, err error) {
	var productReviews []entity.ProductReview

	selectQuery := `
		SELECT
			*
		FROM
			product_reviews
		WHERE
			product_id = $1
	`
	err = e.db.SelectContext(ctx, &productReviews, selectQuery, id)
	if err != nil {
		return []entity.ProductReview{}, err
	}

	return productReviews, nil
}

func (e *ecommerceRepo) CreateProductReview(ctx context.Context, payload entity.ProductReview) (err error) {
	_, err = e.db.ExecContext(ctx,
		`INSERT INTO 
			product_reviews ( product_id, comment, rating) 
		VALUES 
			($1, $2, $3)`, payload.ProductID, payload.Comment, payload.Rating)
	if err != nil {
		return err
	}

	return nil
}

func (e *ecommerceRepo) CreateProductImages(ctx context.Context, payload entity.ProductImage) (err error) {
	_, err = e.db.ExecContext(ctx,
		`INSERT INTO 
			product_images ( product_id, image_url, short_description) 
		VALUES 
			($1, $2, $3)`, payload.ProductID, payload.ImageUrl, payload.ShortDescription)
	if err != nil {
		return err
	}

	return nil
}

func (e *ecommerceRepo) DeleteProductImagesByID(ctx context.Context, productID int64) (err error) {
	query := `
	DELETE FROM 
		product_images
	WHERE
		product_id = $1`

	_, err = e.DB().ExecContext(ctx, query, productID)

	return nil
}
