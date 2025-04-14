package repository

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"gorm.io/gorm"
)

/**
 * @interface OfferRepository
 * @brief Interface for managing offers in the repository.
 *
 * This interface defines methods for performing CRUD operations on offers.
 */
type OfferRepository interface {
	GetOffers() ([]models.Offer, error)
	CreateOffer(offer *models.Offer) error
	GetOfferByID(id uint) (*models.Offer, error)
	ExistsOffer(id uint) (bool, error)
	GetOfferQuantity(id uint) (int, error)
	UpdateOfferQuantity(id uint, newQuantity int) error
}

/**
 * @struct offerRepository
 * @brief Struct implementing the OfferRepository interface.
 *
 * This struct contains methods for managing offers in the repository using
 * GORM for database operations.
 */
type offerRepository struct {
	db *gorm.DB
}

/**
 * @function NewOfferRepository
 * @brief Creates a new instance of OfferRepository.
 *
 * This function initializes a new `offerRepository` with the provided database connection and returns it as an `OfferRepository` interface.
 *
 * @param db Pointer to a `gorm.DB` instance used for database operations.
 * @return An instance of `OfferRepository` that uses the provided `gorm.DB`.
 */
func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db: db}
}

/**
 * @brief Creates a new offer in the repository.
 *
 * This method inserts a new `models.Offer` record into the database.
 *
 * @param offer Pointer to the `models.Offer` to be created.
 * @return An `error` indicating if any error occurred during the creation.
 */
func (r *offerRepository) CreateOffer(offer *models.Offer) error {
	return r.db.Create(offer).Error
}

/**
 * @brief Retrieves an offer by its ID from the repository.
 *
 * This method queries the database for an `models.Offer` with the specified
 * ID and returns it.
 *
 * @param id The ID of the offer to be retrieved.
 * @return A pointer to the `models.Offer` if found.
 * @return An `error` indicating if any error occurred during the retrieval.
 */
func (r *offerRepository) GetOfferByID(id uint) (*models.Offer, error) {
	var offer models.Offer
	err := r.db.Where("id = ?", id).First(&offer).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

/**
 * @brief Retrieves all offers from the repository.
 *
 * This method queries the database for all `models.Offer` records and returns them.
 *
 * @return A slice of `models.Offer` containing all the offers.
 * @return An `error` indicating if any error occurred during the retrieval.
 */
func (r *offerRepository) GetOffers() ([]models.Offer, error) {
	var offers []models.Offer
	result := r.db.Find(&offers)
	return offers, result.Error
}

/**
 * @brief Checks if an offer exists by its ID.
 *
 * This method queries the database to check if an offer with the specified ID exists.
 *
 * @param id The ID of the offer to be checked.
 * @return A boolean indicating if the offer exists.
 * @return An `error` indicating if any error occurred during the check.
 */
func (r *offerRepository) ExistsOffer(id uint) (bool, error) {
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
 * This method queries the database to get the quantity of an offer with the specified ID.
 *
 * @param id The ID of the offer.
 * @return The quantity of the offer.
 * @return An `error` indicating if any error occurred during the retrieval.
 */
func (r *offerRepository) GetOfferQuantity(id uint) (int, error) {
	var offer models.Offer
	err := r.db.Model(&models.Offer{}).Where("id = ?", id).Select("quantity").First(&offer).Error
	if err != nil {
		return 0, err
	}
	return offer.Quantity, nil
}

/**
 * @brief Updates the quantity of an offer by its ID.
 *
 * This method updates the quantity of an `models.Offer` with the specified ID.
 *
 * @param id The ID of the offer.
 * @param newQuantity The new quantity of the offer.
 * @return An `error` indicating if any error occurred during the update.
 */
func (r *offerRepository) UpdateOfferQuantity(id uint, newQuantity int) error {
	return r.db.Model(&models.Offer{}).Where("id = ?", id).Update("quantity", newQuantity).Error
}
