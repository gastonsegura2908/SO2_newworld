package models

import "gorm.io/gorm"

/**
 * @struct Offer
 * @brief Structure representing an offer.
 *
 * This structure represents an offer in the system, including details such as
 * name, quantity, price, and category.
 */
type Offer struct {
	gorm.Model
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

/**
 * @struct OffersRequest
 * @brief Request structure for fetching offers.
 *
 * This structure represents the data required to request offers, including a token for authentication.
 */
type OffersRequest struct {
	Token string `json:"token"`
}

/**
 * @struct OffersResponse
 * @brief Response structure for fetching offers.
 *
 * This structure represents the data returned when fetching offers, including a code and a list of offers.
 */
type OffersResponse struct {
	Code    string  `json:"code"`
	Message []Offer `json:"message"`
}
