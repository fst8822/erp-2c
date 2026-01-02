package model

type Supplier struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Dept    int64  `json:"dept,omitempty"`
}

type SupplierDB struct {
	Name    string `db:"name"`
	Email   string `db:"email"`
	Address string `db:"address"`
	Dept    int64  `db:"dept"`
}
