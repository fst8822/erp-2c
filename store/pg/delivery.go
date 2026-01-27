package pg

import (
	"database/sql"
	"erp-2c/lib/types"
	"erp-2c/model"
	"errors"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (d *DeliveryRepository) SaveWithItems(
	tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error) {

	queryOne := `INSERT INTO delivery(recipient, address, status, created_at)  
			  		VALUES ($1, $2, $3, $4) RETURNING id`

	err := tx.QueryRowx(
		queryOne,
		deliveryWithItems.DeliveryDB.Recipient,
		deliveryWithItems.DeliveryDB.Address,
		deliveryWithItems.DeliveryDB.Status,
		deliveryWithItems.DeliveryDB.CreatedAt,
	).Scan(&deliveryWithItems.DeliveryDB.ID)

	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to insert delivery",
			errors.Join(err, types.ErrInspectedSQL))
	}

	queryTwo := `INSERT INTO delivery_items(delivery_id, product_id, quantity, item_price)
			    	VALUES (:delivery_id, :product_id, :quantity, :item_price)`

	for i := range deliveryWithItems.DeliveryItemsDB {
		deliveryWithItems.DeliveryItemsDB[i].DeliveryID = deliveryWithItems.DeliveryDB.ID
	}
	_, err = tx.NamedExec(queryTwo, deliveryWithItems.DeliveryItemsDB)
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to insert delivery items",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return &deliveryWithItems.DeliveryDB, nil
}

func (d *DeliveryRepository) GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (*model.DeliveryWithItemsDB, error) {
	var deliveryWithItems model.DeliveryWithItemsDB

	err := tx.Get(&deliveryWithItems.DeliveryDB, `SELECT * FROM delivery WHERE id = $1`, deliveryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewAppErr("delivery not found", types.ErrNotFound)
		}
		return nil, types.NewAppErr(" inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrInspectedSQL))
	}

	err = tx.Select(
		&deliveryWithItems.DeliveryItemsDB,
		`SELECT * FROM delivery_items where delivery_id = $1`,
		deliveryId)
	if err != nil {
		return nil, types.NewAppErr(" inspected SQL error, failed to get item",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return &deliveryWithItems, nil
}

func (d *DeliveryRepository) GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error) {
	var deliverListDB model.DeliverListDB

	err := tx.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery")
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrNotFound))
	}

	err = tx.Select(&deliverListDB.ItemsDB, "SELECT * FROM delivery_items")
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrNotFound))
	}
	return &deliverListDB, nil
}

func (d *DeliveryRepository) GetWithItemsByStatus(tx *sqlx.Tx, status string) (*model.DeliverListDB, error) {
	var deliverListDB model.DeliverListDB

	err := tx.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery d where d.status = $1", status)
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrNotFound))
	}

	err = tx.Select(&deliverListDB.ItemsDB, `SELECT di.id, di.delivery_id, di.product_id, di.item_price, di.quantity 
												   FROM delivery_items di join delivery d on di.delivery_id = d.id
												   WHERE d.status = $1`, status)
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrNotFound))
	}
	return &deliverListDB, nil
}

func (d *DeliveryRepository) UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func (d *DeliveryRepository) DeleteById(tx *sqlx.Tx, deliveryId int64) error { return nil }
