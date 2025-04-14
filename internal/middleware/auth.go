package middleware

import (
	"os"
	"time"

	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = os.Getenv("SECRET_KEY")

/**
 * @brief Middleware to protect routes with JWT authentication.
 *
 * This middleware checks for the presence of a valid JWT in the request's
 * Authorization header. If the JWT is missing, malformed, invalid, or expired,
 * the request is rejected with a 401 Unauthorized status.
 *
 * @return A fiber.Handler that checks JWT authentication.
 */
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT"})
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals("user", claims)
		return c.Next()
	}
}

/**
 * @brief Generates a JWT for the given email and role.
 *
 * @param email The email address of the user.
 * @param role The role of the user (e.g., "Admin", "User").
 * @return A signed JWT string or an error if the token generation fails.
 */
func GenerateJWT(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

/**
 * @brief Checks if the user is authenticated based on the JWT.
 *
 * @param c The Fiber context.
 * @param userRepo The user repository to query the user data.
 * @return True if the user is authenticated, false otherwise.
 */
func IsAuthenticated(c *fiber.Ctx, userRepo repository.UserRepository) bool {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	user, err := userRepo.GetUserByToken(tokenString)
	if err != nil || user == nil {
		return false
	}

	return true
}

/**
 * @brief Checks if the user is an admin based on the JWT.
 *
 * @param c The Fiber context.
 * @param userRepo The user repository to query the user data.
 * @return True if the user is an admin, false otherwise.
 */
func IsAdmin(c *fiber.Ctx, userRepo repository.UserRepository) bool {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["role"] != "Admin" {
		return false
	}

	user, err := userRepo.GetUserByToken(tokenString)
	if err != nil || user == nil {
		return false
	}

	return true
}
