package pg

import (
	"erp-2c/lib/types"
	"erp-2c/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (d *DeliveryRepository) SaveDelivery(tx *sqlx.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error) {
	query := `INSERT INTO delivery(recipient, address, status, created_at)  
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := tx.QueryRow(
		query,
		deliveryDB.Recipient,
		deliveryDB.Address,
		deliveryDB.StatusDelivery,
		deliveryDB.CreatedAt,
	).Scan(&deliveryDB.Id)

	if err != nil {
		return nil, types.NewAppErr("inspected SQL error",
			fmt.Errorf("failed to insert delivery %w: %w", err, types.ErrInspectedSQL))
	}

	return &deliveryDB, nil
}

func (d *DeliveryRepository) GetById(tx *sqlx.Tx, deliveryId int64) (*model.DeliveryDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetAll() (*[]model.ProductDomain, error) {
	return nil, nil
}

func (d *DeliveryRepository) GetByStatus(tx *sqlx.Tx, status string) (*model.DeliveryDB, error) {
	return nil, nil
}

func (d *DeliveryRepository) UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func (d *DeliveryRepository) DeleteById(tx *sqlx.Tx, deliveryId int64) error { return nil }

func (d *DeliveryRepository) SaveDeliveryProducts(tx *sqlx.Tx, deliveryProductsDB []model.DeliveryProductDB) error {
	query := `INSERT INTO delivery_product(delivery_id, product_id, quantity, unit_price, total_amount)
			   VALUES (:delivery_id, :product_id, :quantity, :unit_price, :total_amount)`

	_, err := tx.NamedExec(query, deliveryProductsDB)
	if err != nil {
		return types.NewAppErr("inspected SQL error",
			fmt.Errorf("failed to insert delivery %w: %w", err, types.ErrInspectedSQL))
	}
	return nil
}
