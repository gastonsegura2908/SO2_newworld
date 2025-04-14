package models

import "gorm.io/gorm"

/**
 * @struct Order
 * @brief Structure representing an order.
 *
 * This structure represents an order in the system, including its status, total amount,
 * and associated order items.
 */
type Order struct {
	gorm.Model
	Status     string
	Total      int
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

/**
 * @struct OrderItem
 * @brief Structure representing an item in an order.
 *
 * This structure represents an item within an order, including its product ID, quantity,
 * and price.
 */
type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
	Price     int
}

/**
 * @struct CheckoutRequest
 * @brief Request structure for checking out an order.
 *
 * This structure represents the data required to complete a checkout process, including
 * a list of order items.
 */
type CheckoutRequest struct {
	OrderItems []OrderItemRequest `json:"orderItems"`
}

/**
 * @struct OrderItemRequest
 * @brief Structure representing an order item in a checkout request.
 *
 * This structure represents an item in a checkout request, including its product ID and quantity.
 */
type OrderItemRequest struct {
	ProductID uint `json:"productID"`
	Quantity  int  `json:"quantity"`
}

/**
 * @struct CheckoutResponse
 * @brief Response structure for the checkout process.
 *
 * This structure represents the data returned after completing a checkout process,
 * including a response code, message, and the order ID.
 */
type CheckoutResponse struct {
	Code    string          `json:"code"`
	Message CheckoutMessage `json:"message"`
	OrderID uint            `json:"order_id"`
}

/**
 * @struct CheckoutMessage
 * @brief Structure representing the message content in a checkout response.
 *
 * This structure represents the details of a checkout response, including total amount
 * and status.
 */
type CheckoutMessage struct {
	Total  int    `json:"total"`
	Status string `json:"status"`
}

/**
 * @struct OrderStatusResponse
 * @brief Response structure for querying the status of an order.
 *
 * This structure represents the data returned when querying the status of an order,
 * including a response code and the status.
 */
type OrderStatusResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}
