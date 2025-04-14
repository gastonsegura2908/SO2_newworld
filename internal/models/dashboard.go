package models

/**
 * @struct Response_d
 * @brief General response structure for API responses.
 *
 * This structure represents a generic response with a code and a message.
 */
type Response_d struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
}

/**
 * @struct AdminDashboardResponse
 * @brief Response structure for the admin dashboard.
 *
 * This structure represents the data returned by the admin dashboard endpoint,
 * including various metrics related to orders and revenue.
 */
type AdminDashboardResponse struct {
	TotalOrders      int64 `json:"total_orders"`
	TotalRevenue     int   `json:"total_revenue"`
	PendingOrders    int64 `json:"pending_orders"`
	DeliveredOrders  int64 `json:"delivered_orders"`
	PreparingOrders  int64 `json:"preparing_orders"`
	ProcessingOrders int64 `json:"processing_orders"`
	ShippedOrders    int64 `json:"shipped_orders"`
}

/**
 * @struct OrderStatusUpdateRequest
 * @brief Request structure for updating the status of an order.
 *
 * This structure represents the data required to update the status of an order.
 */
type OrderStatusUpdateRequest struct {
	Status string `json:"status"`
}

/**
 * @struct OrderStatusUpdateResponse
 * @brief Response structure for updating the status of an order.
 *
 * This structure represents the data returned after updating the status of an order.
 */
type OrderStatusUpdateResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

/**
 * @struct Dashboard
 * @brief Structure representing dashboard metrics.
 *
 * This structure contains various metrics related to users and orders for the dashboard.
 */
type Dashboard struct {
	UsersCount       int `json:"users_count"`
	OrdersCount      int `json:"orders_count"`
	TotalRevenue     int `json:"total_revenue"`
	PendingOrders    int `json:"pending_orders"`
	DeliveredOrders  int `json:"delivered_orders"`
	PreparingOrders  int `json:"preparing_orders"`
	ProcessingOrders int `json:"processing_orders"`
	ShippedOrders    int `json:"shipped_orders"`
}
