package model

type ProductToSave struct {
	ProductName  string `json:"product-name" validate:"required,min=1"`
	ProductGroup string `json:"product-group,omitempty" validate:"omitempty"`
	Image        []byte `json:"image,omitempty" validate:"omitempty"`
	Stock        int64  `json:"stock,omitempty" validate:"omitempty"`
	Price        int64  `json:"price,omitempty" validate:"omitempty,gte=0"`
}

type ProductUpdate struct {
	ProductName  *string `json:"product-name"`
	ProductGroup *string `json:"product-group"`
	Image        *byte   `json:"image"`
	Stock        *int64  `json:"stock"`
	Price        *int64  `json:"price"`
}

type ProductDomain struct {
	Id           int64
	ProductName  string
	ProductGroup string
	Image        []byte
	Stock        int64
	Price        int64
}

type ProductResponse struct {
	Id           int64  `json:"id,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	ProductGroup string `json:"product_group,omitempty"`
	Image        []byte `json:"image,omitempty"`
	Stock        int64  `json:"stock,omitempty"`
	Price        int64  `json:"price,omitempty"`
}
