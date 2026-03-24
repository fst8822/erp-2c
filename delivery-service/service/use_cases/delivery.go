package use_cases

import (
	"context"
	"delivery-service/lib/observability/app_metrics"
	"delivery-service/lib/sl"
	"delivery-service/lib/types"
	"delivery-service/model"
	clientProductrgrpc "delivery-service/server_grpc/proto/v1/product"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
)

type deliveryRepositoryInt interface {
	SaveWithItems(tx *sqlx.Tx, deliveryWithItems model.DeliveryWithItemsDB) (*model.DeliveryDB, error)
	GetWithItemsById(tx *sqlx.Tx, deliveryId int64) (model.DeliveryWithItemsDB, error)
	GetAll(tx *sqlx.Tx) (*model.DeliverListDB, error)
	GetAllWithItemsByStatus(tx *sqlx.Tx, status model.DeliveryStatus) (*model.DeliverListDB, error)
	UpdateById(tx *sqlx.Tx, deliveryId int64, status model.UpdateStatus) error
	DeleteById(tx *sqlx.Tx, deliveryId int64) error
}

type DeliveryService struct {
	deliveryRepo      deliveryRepositoryInt
	clientProductGRPC clientProductrgrpc.ProductServiceClient
}

func NewDeliveryService(
	deliveryRepo deliveryRepositoryInt,
	clientProductGRPC clientProductrgrpc.ProductServiceClient,
) *DeliveryService {
	return &DeliveryService{
		deliveryRepo:      deliveryRepo,
		clientProductGRPC: clientProductGRPC,
	}
}

func (d *DeliveryService) Save(delivery model.DeliveryItemsDomain) (*model.DeliveryItemsDomain, error) {
	const op = "service.use_cases.delivery.Save"
	sLogger := slog.With("op", op)
	sLogger.Info("Begin save delivery.", slog.Any("delivery", delivery))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	idsToCheck := make([]int64, 0, len(delivery.Items))
	for _, item := range delivery.Items {
		idsToCheck = append(idsToCheck, item.ProductID)
	}
	idsToCheckToSend := &clientProductrgrpc.ListProductIDs{ProductId: idsToCheck}
	foundIds, err := d.clientProductGRPC.GetExistIds(ctx, idsToCheckToSend)
	if err != nil {
		sLogger.Error("failed check existing products", sl.Err(err))
		return nil, err
	}

	missingIds := findMissingIds(idsToCheck, foundIds.GetProductId())
	if len(missingIds) > 0 {
		sLogger.Error("No found product with ids ", slog.Any("missing ids", missingIds))
		return nil, types.NewAppErr(
			fmt.Sprintf("Products with ids %v not found",
				missingIds), types.ErrNotFound)
	}

	deliveryWithItemsDB := delivery.MapToDBWithItems()
	saved, err := d.deliveryRepo.SaveWithItems(nil, deliveryWithItemsDB)
	if err != nil {
		sLogger.Error("failed save delivery", sl.Err(err))
		return nil, err
	}

	slog.Info("Saved delivery", slog.Int64("deliveryId", saved.ID))
	delivery.ID = saved.ID

	app_metrics.TotalCountDelivery.Inc()
	return &delivery, nil
}

func (d *DeliveryService) GetById(deliveryId int64) (*model.DeliveryItemsDomain, error) {
	const op = "service.use_cases.delivery.GetWithItemsById"
	sLogger := slog.With("op", op, "deliveryId", deliveryId)

	deliveryItemsDB, err := d.deliveryRepo.GetWithItemsById(nil, deliveryId)
	if err != nil {
		sLogger.Error("failed to find delivery", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found delivery successful")

	domain := deliveryItemsDB.MapToDomain()
	return &domain, nil
}

func (d *DeliveryService) GetAll() (*model.DeliveryItemListDomain, error) {
	const op = "control.delivery.GetAll"
	sLogger := slog.With("op", op)

	deliveryListDB, err := d.deliveryRepo.GetAll(nil)
	if err != nil {
		sLogger.Error("failed to find deliveries", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found deliveries successful")

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

	deliveryListDB, err := d.deliveryRepo.GetAllWithItemsByStatus(nil, status)
	if err != nil {
		sLogger.Error("failed to find delivery", sl.Err(err))
		return nil, err
	}
	sLogger.Info("found delivery successful")

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
