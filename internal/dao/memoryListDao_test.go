package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"ngini.com/test-api/internal/dao"
	"ngini.com/test-api/internal/model"
	"testing"
)

func Test_MemoryListDAO_GetOrders(t *testing.T) {

	t.Run("get orders should get a list", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		// When
		orders, err := testDAO.GetOrders(ctx)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 3, len(orders))
	})

	t.Run("get order by existing ID should return an object", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		// When
		order, err := testDAO.GetOrder(ctx, "1")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, "1", order.ID)
	})

	t.Run("get order by non-existing ID should return nothing", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		// When
		order, err := testDAO.GetOrder(ctx, "4")

		// Then
		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("get order by existing Slug should return an object", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		// When
		order, err := testDAO.GetOrderBySlug(ctx, "roomy")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, "2", order.ID)
		assert.Equal(t, "roomy", order.Slug)
	})

	t.Run("get order by non-existing Slug should return nothing", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		// When
		order, err := testDAO.GetOrderBySlug(ctx, "test-slug")

		// Then
		assert.Error(t, err)
		assert.Nil(t, order)
	})
}

func Test_MemoryListDAO_AddOrder(t *testing.T) {

	t.Run("add order with non-existing ID should return added object", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID:   "586",
			Name: "test",
			Slug: "test-slug",
		}

		// When
		returnedOrder, err := testDAO.AddOrder(ctx, order)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, returnedOrder)
		assert.Equal(t, order.ID, returnedOrder.ID)
		assert.Equal(t, order.Name, returnedOrder.Name)
		assert.Equal(t, order.Slug, returnedOrder.Slug)
	})

	t.Run("add order with existing ID should return error", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID:   "1",
			Name: "test",
			Slug: "test-slug",
		}

		// When
		returnedOrder, err := testDAO.AddOrder(ctx, order)

		// Then
		assert.Error(t, err)
		assert.Nil(t, returnedOrder)
	})

}

func Test_MemoryListDAO_UpdateOrder(t *testing.T) {
	t.Run("update order with non-existing ID should return error", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID:   "587",
			Name: "test",
			Slug: "test-slug",
		}

		// When
		returnedOrder, err := testDAO.UpdateOrder(ctx, order.ID, order)

		// Then
		assert.Error(t, err)
		assert.Nil(t, returnedOrder)
	})

	t.Run("update order with existing ID should not return error", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID:   "1",
			Name: "test",
			Slug: "test-slug",
		}

		// When
		returnedOrder, err := testDAO.UpdateOrder(ctx, order.ID, order)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, returnedOrder)
		assert.Equal(t, order.ID, returnedOrder.ID)
		assert.Equal(t, order.Name, returnedOrder.Name)
		assert.Equal(t, order.Slug, returnedOrder.Slug)
	})
}

func Test_MemoryListDAO_DeleteOrder(t *testing.T) {

	t.Run("delete order with existing ID should not return error", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID: "1",
		}

		// When
		returnedOrder, err := testDAO.DeleteOrder(ctx, order.ID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, returnedOrder)
		assert.Equal(t, order.ID, returnedOrder.ID)
		assert.NotNil(t, returnedOrder.Name)
		assert.NotNil(t, returnedOrder.Slug)
	})

	t.Run("delete order with non-existing ID should return error", func(t *testing.T) {
		// Given
		ctx := context.TODO()
		testDAO := dao.NewMemoryListDAO()

		order := model.Order{
			ID: "78",
		}

		// When
		returnedOrder, err := testDAO.DeleteOrder(ctx, order.ID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, returnedOrder)
	})

}
