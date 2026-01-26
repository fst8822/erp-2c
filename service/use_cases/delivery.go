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
	"log/slog"
)

type DeliveryService struct {
	repo *store.Store
}

func NewDeliveryService(repo *store.Store) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (d *DeliveryService) Save(deliveryToSave model.DeliveryToSave) (*model.DeliveryDomain, error) {
	const op = "control.delivery.Save"

	tx, err := d.repo.BeginTxx(context.Background())
	if err != nil {
		slog.Error("failed get tx", sl.ErrWithOP(err, op))
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("rollback failed", sl.ErrWithOP(err, op))
		}
	}()

	idsToCheck := make([]int64, 0, len(deliveryToSave.Items))
	for _, item := range deliveryToSave.Items {
		idsToCheck = append(idsToCheck, item.ProductId)
	}

	foundIds, err := d.repo.ProductRepo.GetExistIds(idsToCheck)
	if err != nil {
		slog.Error("failed check existing products", sl.ErrWithOP(err, op))
		return nil, err
	}
	missingIds := findMissingIds(idsToCheck, foundIds)
	if len(missingIds) > 0 {
		slog.Warn("No product with ids found",
			missingIds, slog.String("op", op),
			slog.Any("missing ids", missingIds))

		return nil, types.NewAppErr(
			fmt.Sprintf("Products with ids %v not found",
				missingIds), types.ErrNotFound)
	}

	deliveryDomain := deliveryToSave.MapToDomain()
	deliveryDB := deliveryDomain.MapToDB()

	saved, err := d.repo.Delivery.SaveDelivery(tx, deliveryDB)
	if err != nil {
		slog.Error("failed save delivery", sl.ErrWithOP(err, op))
		return nil, err
	}

	deliveryProductsDB := deliveryDomain.MapToDeliveryProducts(saved.Id)

	if err := d.repo.Delivery.SaveDeliveryProducts(tx, deliveryProductsDB); err != nil {
		slog.Error("failed save deliveryProductsDB", sl.ErrWithOP(err, op))
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		slog.Error("failed commit tx", sl.ErrWithOP(err, op))
		return nil, err
	}
	deliveryDomain.Id = saved.Id

	return deliveryDomain, nil
}

func (d *DeliveryService) GetById(deliveryId int64) (*model.DeliveryDomain, error) {
	return nil, nil
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
