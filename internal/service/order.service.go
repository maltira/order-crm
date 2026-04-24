package service

import (
	"order-crm/internal/model"
	"order-crm/internal/model/dto"
	"order-crm/internal/repository"
)

type OrderService interface {
	GetAllOrders() ([]model.Order, error)
	GetOrderById(id int) (*dto.FullOrder, error)
	CreateOrder(req *dto.CreateOrderRequest) (*model.Order, error)
	UpdateOrderStatus(id int, statusId int) error
	AddPaymentToOrder(id int, req *dto.AddPaymentRequest) error
	AddOrderItem(id int, req *dto.OrderItemRequest) error
	DeleteOrderItem(orderId, orderItemId int) error
	DeleteOrder(id int) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (sc *orderService) GetAllOrders() ([]model.Order, error) {
	return sc.repo.GetAllOrders()
}

func (sc *orderService) GetOrderById(id int) (*dto.FullOrder, error) {
	return sc.repo.GetOrderById(id)
}

func (sc *orderService) CreateOrder(req *dto.CreateOrderRequest) (*model.Order, error) {
	order := &model.Order{
		Label:    req.Label,
		IDClient: req.IDClient,
		IDStatus: 1, // Проект
	}

	var items []model.OrderItem
	var totalAmount float64 = 0

	for _, item := range req.Items {
		items = append(items, model.OrderItem{
			Label:  item.Label,
			Amount: item.Amount,
		})
		totalAmount += item.Amount
	}
	order.Amount = totalAmount

	err := sc.repo.CreateOrderWithItems(order, items)
	return order, err
}

func (sc *orderService) UpdateOrderStatus(id int, statusId int) error {
	return sc.repo.UpdateOrderStatus(id, statusId)
}

func (sc *orderService) AddPaymentToOrder(id int, req *dto.AddPaymentRequest) error {
	payment := &model.Payment{
		IDOrder:       id,
		IDPaymentType: req.IDPaymentType,
		Amount:        req.Amount,
	}
	return sc.repo.AddPaymentToOrder(payment)
}

func (sc *orderService) AddOrderItem(id int, req *dto.OrderItemRequest) error {
	item := &model.OrderItem{
		IDOrder: id,
		Label:   req.Label,
		Amount:  req.Amount,
	}
	return sc.repo.AddOrderItem(item)
}

func (sc *orderService) DeleteOrderItem(orderId, orderItemId int) error {
	return sc.repo.DeleteOrderItem(orderId, orderItemId)
}

func (sc *orderService) DeleteOrder(id int) error {
	return sc.repo.DeleteOrder(id)
}
