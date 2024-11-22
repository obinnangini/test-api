package memorydao

import (
	"errors"
	"ngini.com/test-api/internal/model"
)

//var orderMemoryDB = make([]*model.Order, 3)

var orderMemoryDB = []*model.Order{
	{ID: "1", Name: "Car", Slug: "goes-fast"},
	{ID: "2", Name: "House", Slug: "roomy"},
	{ID: "3", Name: "Watch", Slug: "timey-wimey"},
}

func GetOrders() []*model.Order {
	return orderMemoryDB
}

func GetOrder(orderID string) (*model.Order, error) {
	for _, order := range orderMemoryDB {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}

func AddOrder(order *model.Order) error {
	if len(order.Name) == 0 {
		return errors.New("name is required")
	}

	if orderCheck, _ := GetOrder(order.ID); orderCheck == nil {
		orderMemoryDB = append(orderMemoryDB, order)
	} else {
		return errors.New("order already exists")
	}
	return nil
}

func UpdateOrder(orderID string, order *model.Order) error {
	if len(order.Name) == 0 {
		return errors.New("name is required")
	}

	if _, err := DeleteOrder(orderID); err != nil {
		return err
	} else {
		orderMemoryDB = append(orderMemoryDB, order)
	}

	//if order, _ := GetOrder(orderID); order != nil {
	//	_, err := DeleteOrder(order.ID)
	//	if err != nil {
	//		return err
	//	}
	//	orderMemoryDB = append(orderMemoryDB, order)
	//} else {
	//	return errors.New("order does not exist")
	//}
	return nil
}

func DeleteOrder(orderID string) (*model.Order, error) {
	index := len(orderMemoryDB)
	for i, order := range orderMemoryDB {
		if order.ID == orderID {
			index = i
		}
	}
	if index < len(orderMemoryDB) {
		order := orderMemoryDB[index]
		orderMemoryDB = append(orderMemoryDB[:index], orderMemoryDB[index+1:]...)
		return order, nil
	} else {
		return nil, errors.New("order not found")
	}
}

func GetOrderBySlug(slug string) (*model.Order, error) {
	for _, order := range orderMemoryDB {
		if order.Slug == slug {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}
