package request

type UpsertProduct struct {
	UserID        int64   `json:"user_id"`
	Sku           string  `json:"sku"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	Etalase       string  `json:"etalase"`
	Weight        float64 `json:"weight"`
	Price         int64   `json:"price"`
	ProductImages []struct {
		ImageUrl         string `json:"image_url"`
		ShortDescription string `json:"short_description"`
	} `json:"product_images"`
}

type UpsertProductReview struct {
	ProductID int64  `json:"product_id"`
	Comment   string `json:"comment"`
	Rating    int    `json:"rating"`
}

type FilterProduct struct {
	Search string `json:"search"`
	Sort   string `json:"sort"`
	IsAsc  bool   `json:"is_asc"`
}
