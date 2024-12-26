package dao

import (
	"context"
	"ngini.com/test-api/internal/model"
)

type DAO interface {
	GetOrders(ctx context.Context) ([]*model.Order, error)
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	AddOrder(ctx context.Context, order model.Order) (*model.Order, error)
	UpdateOrder(ctx context.Context, orderID string, order model.Order) (*model.Order, error)
	DeleteOrder(ctx context.Context, orderID string) (*model.Order, error)
	GetOrderBySlug(ctx context.Context, slug string) (*model.Order, error)
}
