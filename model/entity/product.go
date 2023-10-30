package entity

import (
	"time"
)

type Product struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	Sku         string    `db:"sku"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Category    string    `db:"category"`
	Etalase     string    `db:"etalase"`
	Weight      float64   `db:"weight"`
	Price       int64     `db:"price"`
	Rating      float64   `db:"rating"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
