package pg

import (
	"database/sql"
	"erp-2c/lib/types"
	"erp-2c/model"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (d *DeliveryRepository) Save(tx *sqlx.Tx, deliveryDB model.DeliveryDB) (*model.DeliveryDB, error) {
	queryOne := `INSERT INTO delivery(recipient, address, status, created_at)  
			  		VALUES ($1, $2, $3, $4) RETURNING id`

	var err error
	if tx == nil {
		err = d.db.QueryRowx(
			queryOne,
			deliveryDB.Recipient,
			deliveryDB.Address,
			deliveryDB.Status,
			deliveryDB.CreatedAt,
		).Scan(&deliveryDB.ID)
	} else {
		err = tx.QueryRowx(
			queryOne,
			deliveryDB.Recipient,
			deliveryDB.Address,
			deliveryDB.Status,
			deliveryDB.CreatedAt,
		).Scan(&deliveryDB.ID)
	}

	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to insert delivery",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return &deliveryDB, err
}

func (d *DeliveryRepository) GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (*model.DeliveryWithItemsDB, error) {
	var deliveryWithItems model.DeliveryWithItemsDB
	queryGet := `SELECT * FROM delivery WHERE id = $1`
	querySelect := `SELECT * FROM delivery_items where delivery_id = $1`

	var err error
	if tx == nil {
		err = d.db.Get(&deliveryWithItems.DeliveryDB, queryGet, deliveryId)

	} else {
		err = tx.Get(&deliveryWithItems.DeliveryDB, queryGet, deliveryId)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewAppErr("delivery not found", types.ErrNotFound)
		}
		return nil, types.NewAppErr(" inspected SQL error, failed to get delivery by id",
			errors.Join(err, types.ErrInspectedSQL))
	}

	if tx == nil {
		err = d.db.Select(
			&deliveryWithItems.DeliveryItemsDB,
			querySelect,
			deliveryId)
	} else {
		err = tx.Select(
			&deliveryWithItems.DeliveryItemsDB,
			querySelect,
			deliveryId)
	}

	if err != nil {
		return nil, types.NewAppErr(" inspected SQL error, failed to get item",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return &deliveryWithItems, nil
}

func (d *DeliveryRepository) GetAllByStatus(tx *sqlx.Tx, status model.DeliveryStatus) ([]model.DeliveryDB, error) {
	var delivers []model.DeliveryDB
	query := `SELECT * FROM delivery WHERE status = $1 FOR UPDATE SKIP LOCKED LIMIT 5`

	var err error
	if tx == nil {
		err = d.db.Select(&delivers, query, status)
	} else {
		err = tx.Select(&delivers, query, status)
	}
	if err != nil {
		return nil, types.NewAppErr(" inspected SQL error, failed to get deliveries by status",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return delivers, nil
}

func (d *DeliveryRepository) GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error) {
	var deliverListDB model.DeliverListDB

	var err error
	if tx == nil {
		err = d.db.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery")
	} else {
		err = tx.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery")
	}
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get deliveries",
			errors.Join(err, types.ErrNotFound))
	}

	if tx == nil {
		err = d.db.Select(&deliverListDB.ItemsDB, "SELECT * FROM delivery_items")
	} else {
		err = tx.Select(&deliverListDB.ItemsDB, "SELECT * FROM delivery_items")
	}
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get delivery",
			errors.Join(err, types.ErrNotFound))
	}
	return &deliverListDB, nil
}

func (d *DeliveryRepository) GetAllWithItemsByStatus(tx *sqlx.Tx, status model.DeliveryStatus) (*model.DeliverListDB, error) {
	var deliverListDB model.DeliverListDB

	var err error
	if tx == nil {
		err = d.db.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery d where d.status = $1", status)
	} else {
		err = tx.Select(&deliverListDB.DeliveriesDB, "SELECT * FROM delivery d where d.status = $1", status)
	}
	if err != nil {
		return nil, types.NewAppErr("inspected SQL error, failed to get deliveries by status",
			errors.Join(err, types.ErrNotFound))
	}

	if tx == nil {
		err = d.db.Select(&deliverListDB.ItemsDB, `SELECT di.id, di.delivery_id, di.product_id, di.item_price, di.quantity 
												   FROM delivery_items di join delivery d on di.delivery_id = d.id
												   WHERE d.status = $1`, status)
	} else {
		err = tx.Select(&deliverListDB.ItemsDB, `SELECT di.id, di.delivery_id, di.product_id, di.item_price, di.quantity 
												   FROM delivery_items di join delivery d on di.delivery_id = d.id
												   WHERE d.status = $1`, status)
	}
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

func (d *DeliveryRepository) UpdateStatusById(tx *sqlx.Tx, id int64, status model.DeliveryStatus) error {

	query := `UPDATE delivery SET status = $1 WHERE id = $2`
	var result sql.Result
	var err error

	if tx == nil {
		result, err = d.db.Exec(query, status, id)
	} else {
		result, err = tx.Exec(query, status, id)
	}

	if err != nil {
		return types.NewAppErr(" inspected SQL error, failed to change deliveries status",
			errors.Join(err, types.ErrInspectedSQL))
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return types.NewAppErr(" inspected SQL error, failed to get RowsAffected",
			errors.Join(err, types.ErrInspectedSQL))
	}
	if rows == 0 {
		return types.NewAppErr(" inspected SQL error, no change statuses deliveries",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return nil
}

func (d *DeliveryRepository) UpdateStatusByIds(tx *sqlx.Tx, groups map[model.DeliveryStatus][]int64) error {
	query := `UPDATE delivery SET status = $1 WHERE id = any ($2)`
	var result sql.Result
	var err error

	if tx == nil {
		for status, ids := range groups {
			result, err = d.db.Exec(query, status, pq.Array(ids))
		}

	} else {
		for status, ids := range groups {
			result, err = tx.Exec(query, status, pq.Array(ids))
		}
	}

	if err != nil {
		return types.NewAppErr(" inspected SQL error, failed to change deliveries status",
			errors.Join(err, types.ErrInspectedSQL))
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return types.NewAppErr(" inspected SQL error, failed to get RowsAffected",
			errors.Join(err, types.ErrInspectedSQL))
	}
	if rows == 0 {
		return types.NewAppErr(" inspected SQL error, no change statuses deliveries",
			errors.Join(err, types.ErrInspectedSQL))
	}
	return nil
}
