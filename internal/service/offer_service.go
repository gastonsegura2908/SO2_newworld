package service

import (
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
)

/**
 * @interface OfferService
 * @brief Interface for offer-related services.
 *
 * This interface defines methods for managing offers, including retrieving all offers.
 */
type OfferService interface {
	GetOffers() ([]models.Offer, error)
}

/**
 * @struct offerService
 * @brief Implementation of the OfferService interface.
 *
 * This struct implements the `OfferService` interface, providing methods for managing offers using the specified repository.
 */
type offerService struct {
	offerRepository repository.OfferRepository
}

/**
 * @brief Creates a new OfferService instance.
 *
 * @param offerRepo The offer repository to use for database operations.
 * @return A new OfferService instance.
 */
func NewOfferService(offerRepo repository.OfferRepository) OfferService {
	return &offerService{offerRepository: offerRepo}
}

/**
 * @brief Retrieves all offers from the database.
 *
 * @return A slice of Offer models and an error if the retrieval fails.
 */
func (s *offerService) GetOffers() ([]models.Offer, error) {
	return s.offerRepository.GetOffers()
}
