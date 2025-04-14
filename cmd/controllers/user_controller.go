package controllers

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/middleware"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/service"
	"github.com/gofiber/fiber/v2"
)

var (
	userService  service.UserService
	offerService service.OfferService
	orderService service.OrderService
)

/**
 * @brief Registers the routes for the application.
 *
 * @param app The Fiber application instance.
 * @param us The user service to handle user-related operations.
 * @param os The offer service to handle offer-related operations.
 * @param ords The order service to handle order-related operations.
 */
func RegisterRoutes(app *fiber.App, us service.UserService, os service.OfferService, ords service.OrderService) {
	userService = us
	offerService = os
	orderService = ords

	app.Post("/auth/register", Register)
	app.Post("/auth/login", Login)
	app.Get("/auth/offers", middleware.Protected(), GetOffers)
	app.Post("/auth/checkout", middleware.Protected(), Checkout)
	app.Get("/auth/orders/:id", middleware.Protected(), GetOrderStatus)
	app.Get("/admin/dashboard", middleware.Protected(), AdminDashboard)
	app.Patch("/admin/orders/:id", middleware.Protected(), UpdateOrderStatus)
	app.Get("/admin/users", middleware.Protected(), GetAllBuyers)
	app.Delete("/admin/users", middleware.Protected(), RemoveCustomer)

}

// @Summary Register a new user
// @Description Register a new user with the given details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterUserRequest true "User"
// @Success 201 {object} models.Response "User added"
// @Failure 400 {object} models.Response "Bad request"
// @Failure 500 {object} models.Response "Bad server"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	request := new(models.RegisterUserRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	user := &models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	err := userService.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{Code: "201", Message: "User added"})
}

// @Summary Login a user
// @Description Login a user with the given credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse "token"
// @Failure 400 {object} models.Response "Bad request"
// @Failure 500 {object} models.Response "Bad server"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	login := new(models.LoginRequest)
	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	token, err := userService.LoginUser(login)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{Code: "200", Auth: token})
}

// @Summary Get available offers
// @Description Get all available offers
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT <token>"
// @Security ApiKeyAuth
// @Success 200 {object} models.OffersResponse "offers"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /auth/offers [get]
func GetOffers(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAuthenticated(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	offers, err := offerService.GetOffers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.OffersResponse{Code: "200", Message: offers})
}

// @Summary Checkout
// @Description Buy a list of orders. If you want to add more products, here is an example of the structure to follow:{"orderItems": [ { "productID": 1, "quantity": 2 }, { "productID": 2, "quantity": 1 }, { "productID": 3, "quantity": 5 } ] }
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param checkout body models.CheckoutRequest true "Checkout Request"
// @Param Authorization header string true "JWT <token>"
// @Success 200 {object} models.CheckoutResponse
// @Failure 400 {object} models.Response "Bad request"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /auth/checkout [post]
func Checkout(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAuthenticated(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	checkout := new(models.CheckoutRequest)
	if err := c.BodyParser(checkout); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	orderID, total, err := orderService.Checkout(checkout)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	return c.Status(fiber.StatusOK).JSON(models.CheckoutResponse{
		Code: "200",
		Message: models.CheckoutMessage{
			Total:  total,
			Status: "pending",
		},
		OrderID: orderID,
	})
}

// @Summary Get status of a specific order
// @Description Get the status of a specific order by id
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param Authorization header string true "JWT <token>"
// @Success 200 {object} models.OrderStatusResponse "status"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /auth/orders/{id} [get]
func GetOrderStatus(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAuthenticated(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	orderId := c.Params("id")
	status, err := orderService.GetOrderStatus(orderId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.OrderStatusResponse{Code: "200", Status: status})
}

// @Summary Admin dashboard
// @Description Get the admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT <token>"
// @Success 200 {object} models.Response_d
// @Failure 401 {object} models.Response_d "Unauthorized"
// @Failure 500 {object} models.Response_d "Bad server"
// @Router /admin/dashboard [get]
func AdminDashboard(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAdmin(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response_d{Code: "401", Message: "Unauthorized"})
	}

	dashboard, offers, orders, err := orderService.GetAdminDashboard()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response_d{Code: "500", Message: "Server error"})
	}

	response := models.Response_d{
		Code: "200",
		Message: map[string]interface{}{
			"dashboard": dashboard,
			"offers":    offers,
			"orders":    orders,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary Update the status of a specific order
// @Description Update the status of a specific order by id. Valid statuses are "preparing", "processing", "shipped", "delivered".
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT <token>"
// @Param id path string true "Order ID"
// @Param updateRequest body models.OrderStatusUpdateRequest true "Order Status Update Request"
// @Success 200 {object} models.OrderStatusUpdateResponse "status"
// @Failure 400 {object} models.Response "Bad request"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /admin/orders/{id} [patch]
func UpdateOrderStatus(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAdmin(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	orderId := c.Params("id")
	updateRequest := new(models.OrderStatusUpdateRequest)
	if err := c.BodyParser(updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	validStatuses := map[string]bool{
		"preparing":  true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
	}

	if !validStatuses[updateRequest.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	status, err := orderService.UpdateOrderStatus(orderId, updateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.OrderStatusUpdateResponse{Code: "200", Status: status})
}

// @Summary Get all buyers
// @Description Get all buyers, only for admins
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT <token>"
// @Success 200 {object} models.UsersResponse "users"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /admin/users [get]
func GetAllBuyers(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAdmin(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	users, err := userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.UsersResponse{Code: "200", Users: users})
}

// @Summary Remove a customer
// @Description Remove a customer by email, only for admins
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT <token>"
// @Param email body models.DeleteUserRequest true "Delete User Request"
// @Success 200 {object} models.Response "success"
// @Failure 400 {object} models.Response "Bad request"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Bad server"
// @Router /admin/users [delete]
func RemoveCustomer(c *fiber.Ctx) error {
	userRepo := repository.NewUserRepository(userService.GetDB())

	if !middleware.IsAdmin(c, userRepo) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Code: "401", Message: "Unauthorized"})
	}

	request := new(models.DeleteUserRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{Code: "400", Message: "Bad request"})
	}

	err := userService.DeleteUserByEmail(request.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{Code: "200", Message: "success"})
}
