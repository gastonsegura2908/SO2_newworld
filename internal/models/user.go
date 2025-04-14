package models

import "gorm.io/gorm"

/**
 * @struct User
 * @brief Structure representing a user in the system.
 *
 * This structure represents a user with attributes such as username, email,
 * password, and token.
 */
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

/**
 * @struct RegisterUserRequest
 * @brief Request structure for registering a new user.
 *
 * This structure represents the data required to register a new user,
 * including username, email, and password.
 */
type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
 * @struct LoginRequest
 * @brief Structure representing the request data for user login.
 *
 * This structure contains the email and password required for a user to log in.
 */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
 * @struct LoginResponse
 * @brief Structure representing the response data for user login.
 *
 * This structure contains the response code and authentication token returned after a successful login.
 */
type LoginResponse struct {
	Code string `json:"code"`
	Auth string `json:"auth"`
}

/**
 * @struct Response
 * @brief General response structure for API responses.
 *
 * This structure represents a generic response with a code and a message, used for various API responses.
 */
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

/**
 * @struct DeleteUserRequest
 * @brief Structure representing the request data for deleting a user.
 *
 * This structure contains the email of the user to be deleted.
 */
type DeleteUserRequest struct {
	Email string `json:"email"`
}

/**
 * @struct UsersResponse
 * @brief Structure representing the response data for fetching users.
 *
 * This structure contains the response code and a list of users returned by the API.
 */
type UsersResponse struct {
	Code  string `json:"code"`
	Users []User `json:"users"`
}
