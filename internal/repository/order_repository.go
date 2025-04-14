package repository

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"gorm.io/gorm"
)

/**
 * @brief OrderRepository interface defines methods for order-related database operations.
 */
type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uint) (*models.Order, error)
	UpdateOrderStatus(id uint, status string) error
	CountOrders() (int64, error)
	CalculateTotalRevenue() (int, error)
	CountOrdersByStatus(status string) (int64, error)
	ExistsOffer(id uint) (bool, error)
	GetOfferQuantity(id uint) (int, error)
	UpdateOfferQuantity(id uint, quantity int) error
	GetOfferPrice(id uint) (int, error)
	GetAllOffers() ([]models.Offer, error)
	GetAllOrders() ([]models.Order, error)
}

/**
 * @brief orderRepository struct provides the implementation of OrderRepository.
 */
type orderRepository struct {
	db *gorm.DB
}

/**
 * @brief NewOrderRepository creates a new instance of orderRepository.
 *
 * @param db The database connection.
 * @return A new OrderRepository instance.
 */
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

/**
 * @brief Creates a new order in the database.
 *
 * @param order The order model to be created.
 * @return An error if the creation fails.
 */
func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

/**
 * @brief Retrieves an order by its ID.
 *
 * @param id The ID of the order.
 * @return The order model and an error if the retrieval fails.
 */
func (r *orderRepository) GetOrderById(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

/**
 * @brief Updates the status of an order by its ID.
 *
 * @param id The ID of the order.
 * @param status The new status to be updated.
 * @return An error if the update fails.
 */
func (r *orderRepository) UpdateOrderStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

/**
 * @brief Counts the total number of orders.
 *
 * @return The total number of orders and an error if the count fails.
 */
func (r *orderRepository) CountOrders() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * @brief Calculates the total revenue from all orders.
 *
 * @return The total revenue and an error if the calculation fails.
 */
func (r *orderRepository) CalculateTotalRevenue() (int, error) {
	var total int
	if err := r.db.Model(&models.Order{}).Select("SUM(total)").Row().Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

/**
 * @brief Counts the total number of orders by their status.
 *
 * @param status The status of the orders to be counted.
 * @return The total number of orders by the given status and an error if the count fails.
 */
func (r *orderRepository) CountOrdersByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Order{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

/**
 * @brief Checks if an offer exists by its ID.
 *
 * @param id The ID of the offer.
 * @return A boolean indicating if the offer exists and an error if the check fails.
 */
func (r *orderRepository) ExistsOffer(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Offer{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/**
 * @brief Retrieves the quantity of an offer by its ID.
 *
 * @param id The ID of the offer.
 * @return The quantity of the offer and an error if the retrieval fails.
 */
func (r *orderRepository) GetOfferQuantity(id uint) (int, error) {
	var offer models.Offer
	err := r.db.Select("quantity").Where("id = ?", id).First(&offer).Error
	if err != nil {
		return 0, err
	}
	return offer.Quantity, nil
}

/**
 * @brief Updates the quantity of an offer by its ID.
 *
 * @param id The ID of the offer.
 * @param quantity The new quantity to be updated.
 * @return An error if the update fails.
 */
func (r *orderRepository) UpdateOfferQuantity(id uint, quantity int) error {
	return r.db.Model(&models.Offer{}).Where("id = ?", id).Update("quantity", quantity).Error
}

/**
 * @brief Retrieves the price of an offer by its ID.
 *
 * @param id The ID of the offer.
 * @return The price of the offer and an error if the retrieval fails.
 */
func (r *orderRepository) GetOfferPrice(id uint) (int, error) {
	var offer models.Offer
	err := r.db.Select("price").Where("id = ?", id).First(&offer).Error
	if err != nil {
		return 0, err
	}
	return offer.Price, nil
}

/**
 * @brief Retrieves all offers from the database.
 *
 * This method queries the database for all records of type `models.Offer`
 * and returns them. If an error occurs during the query, it returns the error.
 *
 * @return A slice of `models.Offer` containing all the offers.
 * @return An `error` indicating if any error occurred during the database query.
 */
func (r *orderRepository) GetAllOffers() ([]models.Offer, error) {
	var offers []models.Offer
	if err := r.db.Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

/**
 * @brief Retrieves all orders from the database, including their items.
 *
 * This method queries the database for all records of type `models.Order`,
 * including associated `OrderItems`, and returns them. If an error occurs during
 * the query, except for `gorm.ErrRecordNotFound`, it returns the error.
 *
 * @return A slice of `models.Order` containing all the orders.
 * @return An `error` indicating if any error occurred during the database query.
 */
func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("OrderItems").Find(&orders).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orders, nil
}
