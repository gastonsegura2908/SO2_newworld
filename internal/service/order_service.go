package service

import (
	"fmt"
	"strconv"

	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
)

/**
 * @interface OrderService
 * @brief Interface for order-related services.
 *
 * This interface defines methods for managing orders, including checking out, retrieving order status, updating order status, and getting admin dashboard data.
 */
type OrderService interface {
	Checkout(order *models.CheckoutRequest) (uint, int, error)
	GetOrderStatus(id string) (string, error)
	UpdateOrderStatus(id string, status *models.OrderStatusUpdateRequest) (string, error)
	GetAdminDashboard() (models.AdminDashboardResponse, []models.Offer, []models.Order, error)
}

/**
 * @struct orderService
 * @brief Implementation of the OrderService interface.
 *
 * This struct implements the `OrderService` interface, providing methods for managing orders using the specified repository.
 */
type orderService struct {
	orderRepository repository.OrderRepository
}

/**
 * @brief Creates a new OrderService instance.
 *
 * @param orderRepo The order repository to use for database operations.
 * @return A new OrderService instance.
 */
func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepository: orderRepo}
}

/**
 * @brief Processes the checkout of an order.
 *
 * @param checkout The checkout request containing the order items.
 * @return The order ID, total amount, and an error if the checkout fails.
 */
func (s *orderService) Checkout(checkout *models.CheckoutRequest) (uint, int, error) {
	var total int
	neworder := models.Order{Status: "pending", Total: 0, OrderItems: []models.OrderItem{}}

	for _, item := range checkout.OrderItems {
		if item.Quantity <= 0 {
			return 0, 0, fmt.Errorf("invalid quantity for product %d: must be greater than zero", item.ProductID)
		}
		exists, err := s.orderRepository.ExistsOffer(item.ProductID)
		if err != nil || !exists {
			return 0, 0, fmt.Errorf("product %d does not exist", item.ProductID)
		}

		quantity, err := s.orderRepository.GetOfferQuantity(item.ProductID)
		if err != nil || quantity < item.Quantity {
			return 0, 0, fmt.Errorf("product %d not available in the requested quantity", item.ProductID)
		}

		price, err := s.orderRepository.GetOfferPrice(item.ProductID)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get price for product %d", item.ProductID)
		}

		newQuantity := quantity - item.Quantity
		if err := s.orderRepository.UpdateOfferQuantity(item.ProductID, newQuantity); err != nil {
			return 0, 0, fmt.Errorf("failed to update quantity for product %d", item.ProductID)
		}

		total += item.Quantity * price

		newItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		}
		neworder.OrderItems = append(neworder.OrderItems, newItem)
	}

	neworder.Total = total
	err := s.orderRepository.CreateOrder(&neworder)
	if err != nil {
		return 0, 0, err
	}
	return neworder.ID, total, nil
}

/**
 * @brief Retrieves the status of an order by its ID.
 *
 * @param id The order ID.
 * @return The status of the order and an error if the retrieval fails.
 */
func (s *orderService) GetOrderStatus(id string) (string, error) {
	orderID, _ := strconv.ParseUint(id, 10, 64)
	order, err := s.orderRepository.GetOrderById(uint(orderID))
	if err != nil {
		return "", err
	}

	return order.Status, nil
}

/**
 * @brief Updates the status of an order by its ID.
 *
 * @param id The order ID.
 * @param status The new status to update.
 * @return The updated status and an error if the update fails.
 */
func (s *orderService) UpdateOrderStatus(id string, status *models.OrderStatusUpdateRequest) (string, error) {
	orderID, _ := strconv.ParseUint(id, 10, 64)
	err := s.orderRepository.UpdateOrderStatus(uint(orderID), status.Status)
	if err != nil {
		return "", err
	}

	return status.Status, nil
}

/**
 * @brief Retrieves the admin dashboard data.
 *
 * @return The admin dashboard response, offers, orders, and an error if the retrieval fails.
 */
func (s *orderService) GetAdminDashboard() (models.AdminDashboardResponse, []models.Offer, []models.Order, error) {
	totalOrders, err := s.orderRepository.CountOrders()
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	totalRevenue, err := s.orderRepository.CalculateTotalRevenue()
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	pendingOrders, err := s.orderRepository.CountOrdersByStatus("pending")
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	deliveredOrders, err := s.orderRepository.CountOrdersByStatus("delivered")
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	preparingOrders, err := s.orderRepository.CountOrdersByStatus("preparing")
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	processingOrders, err := s.orderRepository.CountOrdersByStatus("processing")
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	shippedOrders, err := s.orderRepository.CountOrdersByStatus("shipped")
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	offers, err := s.orderRepository.GetAllOffers()
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	orders, err := s.orderRepository.GetAllOrders()
	if err != nil {
		return models.AdminDashboardResponse{}, nil, nil, err
	}

	dashboard := models.AdminDashboardResponse{
		TotalOrders:      totalOrders,
		TotalRevenue:     totalRevenue,
		PendingOrders:    pendingOrders,
		DeliveredOrders:  deliveredOrders,
		PreparingOrders:  preparingOrders,
		ProcessingOrders: processingOrders,
		ShippedOrders:    shippedOrders,
	}

	return dashboard, offers, orders, nil
}
