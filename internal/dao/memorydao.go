package dao

import (
	"context"
	"errors"
	"ngini.com/test-api/internal/model"
)

type MemoryDAO struct {
	dbList []*model.Order
}

var orderMemoryDB = []*model.Order{
	{ID: "1", Name: "Car", Slug: "goes-fast"},
	{ID: "2", Name: "House", Slug: "roomy"},
	{ID: "3", Name: "Watch", Slug: "timey-wimey"},
}

func NewMemoryDAO() *MemoryDAO {
	return &MemoryDAO{
		dbList: orderMemoryDB,
	}
}

//func NewDBDAO() *DBDAO {
//	return &DBDAO{}
//}

func (m *MemoryDAO) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return m.dbList, nil
}

func (m *MemoryDAO) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	for _, order := range m.dbList {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}

func (m *MemoryDAO) AddOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	if len(order.Name) == 0 {
		return nil, errors.New("name is required")
	}

	if orderCheck, _ := m.GetOrder(ctx, order.ID); orderCheck == nil {
		m.dbList = append(m.dbList, &order)
	} else {
		return nil, errors.New("order already exists")
	}
	return &order, nil
}

func (m *MemoryDAO) UpdateOrder(ctx context.Context, orderID string, order *model.Order) (*model.Order, error) {
	if len(order.Name) == 0 {
		return nil, errors.New("name is required")
	}

	if _, err := m.DeleteOrder(ctx, orderID); err != nil {
		return nil, err
	} else {
		m.dbList = append(m.dbList, order)
	}

	//if order, _ := m.GetOrder(orderID); order != nil {
	//	_, err := m.DeleteOrder(order.ID)
	//	if err != nil {
	//		return err
	//	}
	//	m.dbList = append(m.dbList, order)
	//} else {
	//	return errors.New("order does not exist")
	//}
	return order, nil
}

func (m *MemoryDAO) DeleteOrder(ctx context.Context, orderID string) (*model.Order, error) {
	index := len(m.dbList)
	for i, order := range m.dbList {
		if order.ID == orderID {
			index = i
		}
	}
	if index < len(m.dbList) {
		order := m.dbList[index]
		m.dbList = append(m.dbList[:index], m.dbList[index+1:]...)
		return order, nil
	} else {
		return nil, errors.New("order not found")
	}
}

func (m *MemoryDAO) GetOrderBySlug(ctx context.Context, slug string) (*model.Order, error) {
	for _, order := range m.dbList {
		if order.Slug == slug {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}
