package response

import "ecommerce/model/entity"

type BaseResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ProductDetail struct {
	Product       entity.Product         `json:"product"`
	ProductImages []entity.ProductImage  `json:"product_images"`
	Review        []entity.ProductReview `json:"review"`
}

type GetProductListResponse struct {
	Data []entity.Product `json:"data"`
	BaseResponse
}

type GetProductDetailResponse struct {
	Data ProductDetail `json:"data"`
	BaseResponse
}
