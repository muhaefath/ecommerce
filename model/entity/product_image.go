package entity

import (
	"database/sql"
	"time"
)

type ProductImage struct {
	ID               int64          `db:"id"`
	ProductID        int64          `db:"product_id"`
	ImageUrl         string         `db:"image_url"`
	ShortDescription sql.NullString `db:"short_description"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}
