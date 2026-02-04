package use_cases

import (
	"context"
	"database/sql"
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/store"
	"errors"
	"fmt"

	"golang.org/x/exp/slog"
)

type DeliveryService struct {
	repo *store.Store
}

func NewDeliveryService(repo *store.Store) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (d *DeliveryService) Save(delivery model.DeliveryItemsDomain) (*model.DeliveryItemsDomain, error) {
	const op = "service.use_cases.delivery.Save"
	sLogger := slog.With("op", op)
	sLogger.Info("Begin save delivery.", slog.Any("delivery", delivery))

	tx, err := d.repo.BeginTxx(context.Background())
	if err != nil {
		sLogger.Error("failed get tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Open transaction")

	defer func() {
		err := tx.Rollback()
		switch {
		case err == nil:
			sLogger.Info("transaction rolled back")
		case errors.Is(err, sql.ErrTxDone):
		default:
			sLogger.Error("rollback failed", sl.Err(err))
		}
	}()

	idsToCheck := make([]int64, 0, len(delivery.Items))
	for _, item := range delivery.Items {
		idsToCheck = append(idsToCheck, item.ProductID)
	}

	foundIds, err := d.repo.ProductRepo.GetExistIds(tx, idsToCheck)
	if err != nil {
		sLogger.Error("failed check existing products", sl.Err(err))
		return nil, err
	}

	missingIds := findMissingIds(idsToCheck, foundIds)
	if len(missingIds) > 0 {
		sLogger.Warn("No found product with ids ",
			slog.Any("missing ids", missingIds))

		return nil, types.NewAppErr(
			fmt.Sprintf("Products with ids %v not found",
				missingIds), types.ErrNotFound)
	}

	deliveryWithItemsDB := delivery.MapToDBWithItems()
	saved, err := d.repo.Delivery.SaveWithItems(tx, deliveryWithItemsDB)
	if err != nil {
		sLogger.Error("failed save delivery", sl.Err(err))
		return nil, err
	}
	slog.Info("Saved delivery", slog.Int64("deliveryId", saved.ID))

	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Commit transaction")
	delivery.ID = saved.ID

	return &delivery, nil
}

func (d *DeliveryService) GetById(deliveryId int64) (*model.DeliveryItemsDomain, error) {
	const op = "service.use_cases.delivery.GetWithItemsById"
	sLogger := slog.With("op", op, "deliveryId", deliveryId)

	tx, err := d.repo.BeginTxx(context.Background())
	if err != nil {
		sLogger.Error("failed get tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Open transaction")

	defer func() {
		err := tx.Rollback()
		switch {
		case err == nil:
			sLogger.Info("transaction rolled back")
		case errors.Is(err, sql.ErrTxDone):
		default:
			sLogger.Error("rollback failed", sl.Err(err))
		}
	}()

	deliveryItemsDB, err := d.repo.Delivery.GetWithItemsById(tx, deliveryId)
	if err != nil {
		sLogger.Error("failed to find delivery", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found delivery successful")

	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Commit transaction")

	domain := deliveryItemsDB.MapToDomain()
	return &domain, nil
}

func (d *DeliveryService) GetAll() (*model.DeliveryItemListDomain, error) {
	const op = "control.delivery.GetAll"
	sLogger := slog.With("op", op)

	tx, err := d.repo.BeginTxx(context.Background())
	if err != nil {
		sLogger.Error("failed get tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Open transaction")

	defer func() {
		err := tx.Rollback()
		switch {
		case err == nil:
			sLogger.Info("transaction rolled back")
		case errors.Is(err, sql.ErrTxDone):
		default:
			sLogger.Error("rollback failed", sl.Err(err))
		}
	}()

	deliveryListDB, err := d.repo.Delivery.GetAll(tx)
	if err != nil {
		sLogger.Error("failed to find deliveries", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found deliveries successful")

	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Commit transaction")

	groupedDB := groupItemsByDeliveryId(*deliveryListDB)
	var deliveryItemList model.DeliveryItemListDomain
	for _, d := range groupedDB {
		deliveryItemList.DeliveryItemsDomain = append(deliveryItemList.DeliveryItemsDomain, d.MapToDomain())
	}
	return &deliveryItemList, nil
}

func (d *DeliveryService) GetByStatus(status model.DeliveryStatus) (*model.DeliveryItemListDomain, error) {

	const op = "service.use_cases.delivery.LockAndGetDeliveries"
	sLogger := slog.With("op", op, "status", status)

	tx, err := d.repo.BeginTxx(context.Background())
	if err != nil {
		sLogger.Error("failed get tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Open transaction")

	defer func() {
		err := tx.Rollback()
		switch {
		case err == nil:
			sLogger.Info("transaction rolled back")
		case errors.Is(err, sql.ErrTxDone):
		default:
			sLogger.Error("rollback failed", sl.Err(err))
		}
	}()

	deliveryListDB, err := d.repo.Delivery.GetAllWithItemsByStatus(tx, status)
	if err != nil {
		sLogger.Error("failed to find delivery", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found delivery successful")

	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Commit transaction")

	groupedDB := groupItemsByDeliveryId(*deliveryListDB)
	var deliveryItemList model.DeliveryItemListDomain
	for _, d := range groupedDB {
		deliveryItemList.DeliveryItemsDomain = append(deliveryItemList.DeliveryItemsDomain, d.MapToDomain())
	}
	return &deliveryItemList, nil
}

func (d *DeliveryService) UpdateById(deliveryId int64, status model.UpdateStatus) error {
	return nil
}

func (d *DeliveryService) DeleteById(deliveryId int64) error {
	return nil
}

func findMissingIds(checkIds []int64, foundIds []int64) []int64 {
	foundMap := make(map[int64]bool, len(foundIds))
	for _, id := range foundIds {
		foundMap[id] = true
	}
	missingIds := make([]int64, 0, len(checkIds))
	for _, id := range checkIds {
		if !foundMap[id] {
			missingIds = append(missingIds, id)
		}
	}
	return missingIds
}

func groupItemsByDeliveryId(deliveryListDB model.DeliverListDB) []model.DeliveryWithItemsDB {
	mapItem := make(map[int64][]model.ItemsDB)
	for _, v := range deliveryListDB.ItemsDB {
		mapItem[v.DeliveryID] = append(mapItem[v.DeliveryID], v)
	}

	var deliveryWithItemsDB []model.DeliveryWithItemsDB
	for _, v := range deliveryListDB.DeliveriesDB {
		deliveryWithItemsDB = append(deliveryWithItemsDB, model.DeliveryWithItemsDB{
			DeliveryDB:      v,
			DeliveryItemsDB: mapItem[v.ID],
		})
	}
	return deliveryWithItemsDB
}
