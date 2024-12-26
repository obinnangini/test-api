package dao

import (
	"context"
	"errors"
	"maps"
	"ngini.com/test-api/internal/model"
	"slices"
)

type MemoryMapDAO struct {
	dbMap map[string]*model.Order
}

var orderMemoryMap = map[string]*model.Order{
	"1": {ID: "1", Name: "Car", Slug: "goes-fast"},
	"2": {ID: "2", Name: "House", Slug: "roomy"},
	"3": {ID: "3", Name: "Watch", Slug: "timey-wimey"},
}

func NewMemoryMapDAO() *MemoryMapDAO {

	map2 := map[string]*model.Order{}

	for key, value := range orderMemoryMap {
		map2[key] = value
	}
	memoryMapDao := MemoryMapDAO{
		dbMap: map2,
	}
	return &memoryMapDao
}

func (m *MemoryMapDAO) GetOrders(ctx context.Context) ([]*model.Order, error) {
	list := slices.Collect(maps.Values(m.dbMap))
	return list, nil
}

func (m *MemoryMapDAO) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	order := m.dbMap[orderID]
	if order != nil {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (m *MemoryMapDAO) AddOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	if len(order.Name) == 0 {
		return nil, errors.New("name is required")
	}

	if _, ok := m.dbMap[order.ID]; ok {
		return nil, errors.New("order already exists")
	}

	m.dbMap[order.ID] = &order

	return &order, nil
}

func (m *MemoryMapDAO) UpdateOrder(ctx context.Context, orderID string, order model.Order) (*model.Order, error) {
	if len(order.Name) == 0 {
		return nil, errors.New("name is required")
	}

	_, found := m.dbMap[orderID]

	if found {
		m.dbMap[orderID] = &order
		return &order, nil
	}

	return nil, errors.New("order not found")
}

func (m *MemoryMapDAO) DeleteOrder(ctx context.Context, orderID string) (*model.Order, error) {
	if _, ok := m.dbMap[orderID]; ok {
		order := m.dbMap[orderID]
		delete(m.dbMap, orderID)
		return order, nil
	}

	return nil, errors.New("order not found")
}

func (m *MemoryMapDAO) GetOrderBySlug(ctx context.Context, slug string) (*model.Order, error) {
	for _, order := range slices.Collect(maps.Values(m.dbMap)) {
		if order.Slug == slug {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}
