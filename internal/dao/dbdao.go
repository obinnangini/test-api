package dao

import (
	"context"
	"fmt"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"ngini.com/test-api/internal/model"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type DBDAO struct {
	conn *pgx.Conn
}

func NewDBDAO() *DBDAO {
	return &DBDAO{}
}

func (d *DBDAO) InitConnection(ctx context.Context, dbUrl string) error {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}

	d.conn = conn

	var orderMemoryDB = []*model.Order{
		{ID: "1", Name: "Car", Slug: "goes-fast"},
		{ID: "2", Name: "House", Slug: "roomy"},
		{ID: "3", Name: "Watch", Slug: "timey-wimey"},
	}

	ordersFromDB, err := d.GetOrders(ctx)

	if err != nil {
		fmt.Printf("GetOrders Error: %v", err)
		return err
	}

	if len(ordersFromDB) == 0 {
		fmt.Println("Seeding DB")
		for _, order := range orderMemoryDB {
			_, err = d.AddOrder(ctx, *order)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DBDAO) GetOrders(ctx context.Context) ([]*model.Order, error) {

	var err error
	queryString, args, err := psql.Select(
		sm.Columns(d.getColumns()...),
		sm.From("orders"),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Select LIST Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	//orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])

	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []*model.Order
	var order *model.Order
	for rows.Next() {
		order, err = d.scanRow(rows)
		if order != nil {
			orders = append(orders, order)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (d *DBDAO) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	var err error
	queryString, args, err := psql.Select(
		sm.Columns(d.getColumns()...),
		sm.From("orders"),
		sm.Where(psql.Quote("user_id").EQ(psql.Arg(orderID))),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Select One Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	//orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])

	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []model.Order
	var order *model.Order
	for rows.Next() {
		order, err = d.scanRow(rows)
		if order != nil {
			orders = append(orders, *order)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(orders) > 1 {
		return nil, fmt.Errorf("expected 1 order, got %d", len(orders))
	}

	if len(orders) > 0 {
		return &orders[0], nil
	}

	return nil, nil
}

func (d *DBDAO) AddOrder(ctx context.Context, order model.Order) (*model.Order, error) {

	var err error
	var queryString string
	var args []interface{}
	var orderFromDB *model.Order

	orderFromDB, err = d.GetOrder(ctx, order.ID)

	if err != nil {
		fmt.Printf("Insert Query build error: %v", err)
		return nil, err
	}

	if orderFromDB != nil {
		return nil, fmt.Errorf("order with ID: %s already exists", order.ID)
	}

	queryString, args, err = psql.Insert(
		im.Into("orders"),
		im.Values(psql.Arg(order.ID, order.Name, order.Slug)),
		im.Returning(d.getColumns()...),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Insert Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []model.Order
	var orderfromDB *model.Order
	if rows.Next() {
		orderfromDB, err = d.scanRow(rows)
		if orderfromDB != nil {
			orders = append(orders, *orderfromDB)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if orderfromDB != nil {
		return orderfromDB, nil
	}

	return nil, nil

}

func (d *DBDAO) UpdateOrder(ctx context.Context, orderID string, order model.Order) (*model.Order, error) {
	var err error
	var queryString string
	var args []interface{}

	queryString, args, err = psql.Update(
		um.Table("orders"),
		um.Where(psql.Quote("user_id").EQ(psql.Arg(orderID))),
		um.SetCol("name").ToArg(order.Name),
		um.SetCol("slug").ToArg(order.Slug),
		um.Returning(d.getColumns()...),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Update Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []model.Order
	var orderfromDB *model.Order
	if rows.Next() {
		orderfromDB, err = d.scanRow(rows)
		if orderfromDB != nil {
			orders = append(orders, *orderfromDB)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if orderfromDB != nil {
		return orderfromDB, nil
	}

	return nil, nil

}

func (d *DBDAO) DeleteOrder(ctx context.Context, orderID string) (*model.Order, error) {

	var err error
	var queryString string
	var args []interface{}

	queryString, args, err = psql.Delete(
		dm.From("orders"),
		dm.Where(psql.Quote("user_id").EQ(psql.Arg(orderID))),
		dm.Returning(d.getColumns()...),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Delete Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []model.Order
	var orderfromDB *model.Order
	for rows.Next() {
		orderfromDB, err = d.scanRow(rows)
		if orderfromDB != nil {
			orders = append(orders, *orderfromDB)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if orderfromDB != nil {
		return orderfromDB, nil
	}

	return nil, nil

}

func (d *DBDAO) GetOrderBySlug(ctx context.Context, slug string) (*model.Order, error) {
	var err error
	queryString, args, err := psql.Select(
		sm.Columns(d.getColumns()...),
		sm.From("orders"),
		sm.Where(psql.Quote("slug").EQ(psql.Arg(slug))),
	).Build(ctx)

	if err != nil {
		fmt.Printf("Select One Query build error: %v", err)
		return nil, err
	}

	fmt.Printf("QUERY: %s with args: %v \n", queryString, args)

	rows, err := d.conn.Query(ctx, queryString, args...)
	defer rows.Close()
	//orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])

	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}

	var orders []model.Order
	var order *model.Order
	for rows.Next() {
		order, err = d.scanRow(rows)
		if order != nil {
			orders = append(orders, *order)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(orders) > 1 {
		return nil, fmt.Errorf("expected 1 order, got %d", len(orders))
	}

	if len(orders) > 0 {
		return &orders[0], nil
	}

	return nil, nil
}

func (d *DBDAO) getColumns() []any {
	return []any{
		"user_id",
		"name",
		"slug",
	}
}

func (d *DBDAO) scanRow(rows pgx.Rows) (*model.Order, error) {
	var order model.Order
	err := rows.Scan(&order.ID, &order.Name, &order.Slug)
	if err != nil {
		return nil, err
	} else {
		return &order, nil
	}
}

func (d *DBDAO) shutDownConnection(ctx context.Context) error {
	if d.conn == nil {
		return nil
	}
	err := d.conn.Close(ctx)

	if err != nil {
		return err
	}
	return nil
}
