package model

type ProductToSave struct {
	ProductName  string `json:"product-name"`
	ProductGroup string `json:"product-group"`
	Image        byte   `json:"image"`
	Stock        int64  `json:"stock"`
	Price        int64  `json:"price"`
}

type ProductUpdate struct {
	ProductName  string `json:"product-name"`
	ProductGroup string `json:"product-group"`
	Image        byte   `json:"image"`
}

type ProductStockUpdate struct {
	Id          int64  `json:"id"`
	ProductName string `json:"product-name"`
	Stock       int64  `json:"stock"`
}

type ProductPriceUpdate struct {
	Id          int64  `json:"id"`
	ProductName string `json:"product-name"`
	Price       int64  `json:"price"`
}

type ProductDomain struct {
	Id           int64
	ProductName  string
	ProductGroup string
	Image        byte
	Stock        int64
	Price        int64
}

type ProductDB struct {
	Id           int64  `db:"id"`
	ProductName  string `db:"product_name"`
	ProductGroup string `db:"product_group"`
	Image        byte   `db:"image"`
	Stock        int64  `db:"stock"`
	Price        int64  `db:"price"`
}
