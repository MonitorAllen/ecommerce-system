package orders

import (
	"context"
	db "ecommerce-system/internal/db/sqlc"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product no stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderRequest) (db.Order, error)
}

type OrderService struct {
	repo *db.Queries
	db   *pgx.Conn
}

func NewOrderService(repo *db.Queries, db *pgx.Conn) *OrderService {
	return &OrderService{
		repo: repo,
		db:   db,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, tempOrder createOrderRequest) (db.Order, error) {
	// validate payload
	if tempOrder.CustomerID == 0 {
		return db.Order{}, errors.New("invalid customer id")
	}

	if len(tempOrder.Items) == 0 {
		return db.Order{}, errors.New("invalid item number")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return db.Order{}, err
	}
	defer tx.Rollback(ctx)

	queriesWithTx := s.repo.WithTx(tx)

	order, err := queriesWithTx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return db.Order{}, err
	}

	for _, item := range tempOrder.Items {
		product, err := queriesWithTx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return db.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return db.Order{}, ErrProductNoStock
		}

		_, err = queriesWithTx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCenters,
		})
		if err != nil {
			return db.Order{}, err
		}

		quantity := product.Quantity - item.Quantity

		arg := db.UpdateProductStockParams{
			Quantity: quantity,
			ID:       product.ID,
		}

		err = queriesWithTx.UpdateProductStock(ctx, arg)
		if err != nil {
			return db.Order{}, fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return db.Order{}, err
	}

	return order, nil
}
