package repository

import (
	"database/sql"
	"order-crm/internal/model"
	"order-crm/internal/model/dto"
)

type OrderRepository interface {
	GetAllOrders() ([]model.Order, error)
	GetOrderById(id int) (*dto.FullOrder, error)
	CreateOrderWithItems(order *model.Order, items []model.OrderItem) error
	UpdateOrderStatus(id int, statusId int) error
	AddPaymentToOrder(payment *model.Payment) error
	AddOrderItem(item *model.OrderItem) error
	DeleteOrderItem(orderId, orderItemId int) error
	DeleteOrder(id int) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAllOrders() ([]model.Order, error) {
	query := `SELECT * FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		err = rows.Scan(&o.ID, &o.Label, &o.IDStatus, &o.IDClient, &o.Amount)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepository) GetOrderById(id int) (*dto.FullOrder, error) {
	// * Получаем сам заказ
	orderQuery := `SELECT * FROM orders WHERE id = $1`

	var o dto.FullOrder
	err := r.db.QueryRow(orderQuery, id).Scan(&o.ID, &o.Label, &o.IDStatus, &o.IDClient, &o.Amount)
	if err != nil {
		return nil, err
	}

	// * Получаем позиции заказа
	itemsQuery := `SELECT * FROM order_items WHERE id_order = $1 ORDER BY id`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		err = rows.Scan(&item.ID, &item.Label, &item.IDOrder, &item.Amount)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	// * Получаем информацию о платеже
	paymentsQuery := `SELECT * FROM payments WHERE id_order = $1 ORDER BY id`

	rows, err = r.db.Query(paymentsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment model.Payment
		err = rows.Scan(&payment.ID, &payment.IDOrder, &payment.IDPaymentType, &payment.Amount)
		if err != nil {
			return nil, err
		}
		o.Payments = append(o.Payments, payment)
	}

	return &o, nil
}

func (r *orderRepository) CreateOrderWithItems(order *model.Order, items []model.OrderItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (label, id_status, id_client, amount) VALUES ($1, $2, $3, $4) RETURNING id`

	err = tx.QueryRow(query, order.Label, order.IDStatus, order.IDClient, order.Amount).Scan(&order.ID)
	if err != nil {
		return err
	}

	// Добавляем позиции
	for _, item := range items {
		_, err = tx.Exec(`INSERT INTO order_items (label, id_order, amount) VALUES ($1, $2, $3)`,
			item.Label, order.ID, item.Amount)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *orderRepository) UpdateOrderStatus(id int, statusId int) error {
	query := `UPDATE orders SET id_status = $1 WHERE id = $2`

	_, err := r.db.Exec(query, statusId, id)
	return err
}

func (r *orderRepository) AddPaymentToOrder(payment *model.Payment) error {
	query := `INSERT INTO payments (id_order, id_payment_type, amount) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, payment.IDOrder, payment.IDPaymentType, payment.Amount)
	return err
}

func (r *orderRepository) AddOrderItem(item *model.OrderItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryItem := `INSERT INTO order_items (label, id_order, amount) VALUES ($1, $2, $3)`
	queryOrder := `UPDATE orders SET amount = $1 WHERE id = $2`

	_, err = tx.Exec(queryItem, item.Label, item.IDOrder, item.Amount)
	if err != nil {
		return err
	}

	var amount float64
	err = tx.QueryRow(`SELECT SUM(amount) FROM order_items WHERE id_order = $1`, item.IDOrder).Scan(&amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryOrder, amount, item.IDOrder)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *orderRepository) DeleteOrderItem(orderId, orderItemId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var amount float64

	query := `DELETE FROM order_items WHERE id_order = $1 AND id = $2`
	queryOrder := `UPDATE orders SET amount = $1 WHERE id = $2`

	_, err = tx.Exec(query, orderId, orderItemId)
	if err != nil {
		return err
	}

	err = tx.QueryRow(`SELECT coalesce(SUM(amount), 0) FROM order_items WHERE id_order = $1`, orderId).Scan(&amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryOrder, amount, orderItemId)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *orderRepository) DeleteOrder(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
