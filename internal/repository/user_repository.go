package repository

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"gorm.io/gorm"
)

/**
 * @brief UserRepository interface defines methods for user-related database operations.
 */
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserToken(user *models.User) error
	GetUserByToken(token string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUserByEmail(email string) error
	GetDB() *gorm.DB
}

/**
 * @brief userRepository struct provides the implementation of UserRepository.
 */
type userRepository struct {
	db *gorm.DB
}

/**
 * @brief NewUserRepository creates a new instance of userRepository.
 *
 * @param db The database connection.
 * @return A new UserRepository instance.
 */
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

/**
 * @brief Creates a new user in the database.
 *
 * @param user The user model to be created.
 * @return An error if the creation fails.
 */
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

/**
 * @brief Retrieves the database connection.
 *
 * @return The database connection.
 */
func (r *userRepository) GetDB() *gorm.DB {
	return r.db
}

/**
 * @brief Retrieves a user by their email.
 *
 * @param email The email of the user.
 * @return The user model and an error if the retrieval fails.
 */
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

/**
 * @brief Updates a user's token in the database.
 *
 * @param user The user model with the updated token.
 * @return An error if the update fails.
 */
func (r *userRepository) UpdateUserToken(user *models.User) error {
	return r.db.Save(user).Error
}

/**
 * @brief Retrieves a user by their token.
 *
 * @param token The token of the user.
 * @return The user model and an error if the retrieval fails.
 */
func (r *userRepository) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

/**
 * @brief Retrieves all users from the database.
 *
 * @return A slice of User models and an error if the retrieval fails.
 */
func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

/**
 * @brief Deletes a user by their email.
 *
 * @param email The email of the user to be deleted.
 * @return An error if the deletion fails.
 */
func (r *userRepository) DeleteUserByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.User{}).Error
}
