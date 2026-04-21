package dto

import "order-crm/internal/model"

type CreateOrderRequest struct {
	Label    string            `json:"label"`
	IDClient int               `json:"id_client"`
	Items    []model.OrderItem `json:"items"`
}
type UpdateOrderStatusRequest struct {
	IDStatus int `json:"id_status"`
}

type AddPaymentRequest struct {
	IDPaymentType int     `json:"id_payment_type"`
	Amount        float64 `json:"amount"`
}

type OrderItemRequest struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

type DeleteOrderItemRequest struct {
	OrderId     int `json:"order_id"`
	OrderItemId int `json:"order_item_id"`
}

type FullOrder struct {
	ID       int     `json:"id"`
	Label    string  `json:"label"`
	IDStatus int     `json:"id_status"`
	IDClient int     `json:"id_client"`
	Amount   float64 `json:"amount"`

	Items    []model.OrderItem `json:"items"`
	Payments []model.Payment   `json:"payments"`
}
