package service

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/middleware"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/**
 * @interface UserService
 * @brief Interface for user-related services.
 *
 * This interface defines methods for managing users, including user creation, login, retrieval, deletion, and database access.
 */
type UserService interface {
	CreateUser(user *models.User) error
	LoginUser(login *models.LoginRequest) (string, error)
	GetAllUsers() ([]models.User, error)
	DeleteUserByEmail(email string) error
	GetDB() *gorm.DB
}

/**
 * @struct userService
 * @brief Implementation of the UserService interface.
 *
 * This structure provides the implementation of the methods defined in the `UserService` interface.
 */
type userService struct {
	userRepository repository.UserRepository
}

/**
 * @brief Creates a new UserService instance.
 *
 * @param userRepo The user repository to use for database operations.
 * @return A new UserService instance.
 */
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

/**
 * @brief Returns the database instance used by the service.
 *
 * @return A pointer to the gorm.DB instance.
 */
func (s *userService) GetDB() *gorm.DB {
	return s.userRepository.GetDB()
}

/**
 * @brief Creates a new user with the given data.
 *
 * @param user The user data to create the new user.
 * @return An error if the user creation fails.
 */
func (s *userService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepository.CreateUser(user)
}

/**
 * @brief Logs in a user with the given login data.
 *
 * @param login The login request containing the user's email and password.
 * @return A JWT token string if the login is successful, or an error if it fails.
 */
func (s *userService) LoginUser(login *models.LoginRequest) (string, error) {
	user, err := s.userRepository.GetUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", err
	}

	role := "normal"
	if user.Username == "UBUNTU" && login.Password == "UBUNTU" {
		role = "Admin"
	}

	token, err := middleware.GenerateJWT(user.Email, role)
	if err != nil {
		return "", err
	}

	user.Token = token
	if err := s.userRepository.UpdateUserToken(user); err != nil {
		return "", err
	}

	return token, nil
}

/**
 * @brief Retrieves all users from the database.
 *
 * @return A slice of User models and an error if the retrieval fails.
 */
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.GetAllUsers()
}

/**
 * @brief Deletes a user by their email.
 *
 * @param email The email of the user to delete.
 * @return An error if the deletion fails.
 */
func (s *userService) DeleteUserByEmail(email string) error {
	return s.userRepository.DeleteUserByEmail(email)
}
