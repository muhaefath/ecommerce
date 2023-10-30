package entity

import (
	"database/sql"
	"time"
)

type ProductReview struct {
	ID        int64          `db:"id"`
	ProductID int64          `db:"product_id"`
	Rating    int            `db:"rating"`
	Comment   sql.NullString `db:"comment"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
