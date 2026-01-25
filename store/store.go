package store

import (
	"erp-2c/store/pg"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	UserRepo    UserRepository
	ProductRepo ProductRepository
	Delivery    DeliveryRepository
}

func NewStore(db *sqlx.DB) *Store {

	return &Store{
		UserRepo:    pg.NewUserRepository(db),
		ProductRepo: pg.NewProductRepository(db),
		Delivery:    pg.NewDeliveryRepository(db),
	}
}
