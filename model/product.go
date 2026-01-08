package model

type Product struct {
	ID           int64  `json:"id,omitempty"`
	ProductName  string `json:"product-name" validate:"required"`
	ProductGroup string `json:"product-group,omitempty"`
	Image        byte   `json:"image,omitempty"`
	Stock        int64  `json:"stock,omitempty"`
	Price        int64  `json:"price,omitempty"`
}

type ProductDB struct {
	ID           int64  `db:"id"`
	ProductName  string `db:"product_name"`
	ProductGroup string `db:"product_group"`
	Image        byte   `db:"image"`
	Stock        int64  `db:"stock"`
	Price        int64  `db:"price"`
}
