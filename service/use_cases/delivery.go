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

func (d *DeliveryService) Save(deliveryToSave model.DeliveryToSave) (*model.DeliveryDomain, error) {
	const op = "service.use_cases.delivery.Save"
	sLogger := slog.With("op", op)
	sLogger.Info("Begin save delivery.", slog.Any("delivery", deliveryToSave))

	tx, err := d.repo.BeginTx(context.Background())
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

	idsToCheck := make([]int64, 0, len(deliveryToSave.Items))
	for _, item := range deliveryToSave.Items {
		idsToCheck = append(idsToCheck, item.ProductId)
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

	deliveryDomain := deliveryToSave.MapToDomain()
	deliveryDB := deliveryDomain.MapToDB()

	saved, err := d.repo.Delivery.SaveDelivery(tx, deliveryDB)
	if err != nil {
		sLogger.Error("failed save delivery", sl.Err(err))
		return nil, err
	}
	slog.Info("Saved delivery", slog.Int64("deliveryId", saved.Id))

	deliveryProductsDB := deliveryDomain.MapToDeliveryProductsDB(saved.Id)
	if err := d.repo.Delivery.SaveDeliveryProducts(tx, deliveryProductsDB); err != nil {
		slog.Error("failed save deliveryProductsDB", sl.Err(err))
		return nil, err
	}
	slog.Info("Saved deliveryProductsDB")
	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.Err(err))
		return nil, err
	}
	sLogger.Info("Commit transaction")
	deliveryDomain.Id = saved.Id

	return deliveryDomain, nil
}

func (d *DeliveryService) GetById(deliveryId int64) (*model.DeliveryDomain, error) {
	const op = "service.use_cases.delivery.GetById"
	sLogger := slog.With("op", op, "deliveryId", deliveryId)

	foundDB, err := d.repo.Delivery.GetById(nil, deliveryId)
	if err != nil {
		sLogger.Error("failed to find delivery", slog.Int64("deliveryId", deliveryId), sl.Err(err))
		return nil, err
	}
	foundDomain := foundDB.MapToDomain()
	sLogger.Info("found delivery")
	return &foundDomain, nil
}
func (d *DeliveryService) GetAll() ([]model.DeliveryDomain, error) {
	return nil, nil

}
func (d *DeliveryService) GetByStatus(status string) (*model.DeliveryDomain, error) {
	return nil, nil
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
