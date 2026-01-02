package model

type Buyer struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`
	Dept    int64  `json:"dept,omitempty"`
}

type BuyerDB struct {
	Name    string `db:"name"`
	Email   string `db:"email"`
	Address string `db:"address"`
	Dept    int64  `db:"dept"`
}
