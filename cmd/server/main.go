package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/cmd/controllers"
	_ "github.com/ICOMP-UNC/newworld-gastonsegura2908.git/docs"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/models"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/repository"
	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"

	"github.com/ICOMP-UNC/newworld-gastonsegura2908.git/database"

	"github.com/joho/godotenv"
)

/**
 * @brief Structure representing the response from the supplies endpoint.
 */
type SuppliesResponse struct {
	Food struct {
		Fruits     int `json:"fruits"`
		Meat       int `json:"meat"`
		Vegetables int `json:"vegetables"`
		Water      int `json:"water"`
	} `json:"food"`
	Medicine struct {
		Analgesics  int `json:"analgesics"`
		Antibiotics int `json:"antibiotics"`
		Bandages    int `json:"bandages"`
	} `json:"medicine"`
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample Swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	envPath := "/root/.env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", envPath)
	} else {
		log.Printf("Loading .env file from %s", envPath)
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}

	database.ResetDatabase()

	app := fiber.New()

	db := database.InitDB()

	suppliesURL := os.Getenv("SUPPLIES_URL")
	offers, err := fetchSupplies(suppliesURL)
	if err != nil {
		log.Fatalf("Error fetching supplies: %v", err)
	}

	err = storeSuppliesInDB(db, offers)
	if err != nil {
		log.Fatalf("Error storing supplies in database: %v", err)
	}

	// Start a goroutine that sends POST requests every 30 seconds
	go startPeriodicUpdates(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	offerRepo := repository.NewOfferRepository(db)
	offerService := service.NewOfferService(offerRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Local",
	}))

	app.Use(cors.New())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	controllers.RegisterRoutes(app, userService, offerService, orderService)
	log.Fatal(app.Listen(":" + port))

}

/**
 * @brief Fetches the supplies data from the given URL.
 *
 * @param url The URL to fetch the supplies data from.
 * @return A slice of models.Offer containing the fetched supplies, and an error if there was an issue.
 */
func fetchSupplies(url string) ([]models.Offer, error) {
	var suppliesResponse SuppliesResponse
	var offers []models.Offer

	retryCount := 5
	retryInterval := 2 * time.Second

	for i := 0; i < retryCount; i++ {
		resp, err := http.Get(url)
		if err != nil {
			if i == retryCount-1 {
				return offers, fmt.Errorf("failed to fetch supplies after %d attempts: %w", retryCount, err)
			}
			log.Printf("Error fetching supplies (attempt %d/%d): %v. Retrying in %s...", i+1, retryCount, err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return offers, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return offers, fmt.Errorf("failed to read response body: %w", err)
		}

		if err := json.Unmarshal(body, &suppliesResponse); err != nil {
			return offers, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		offers = append(offers, models.Offer{Name: "fruits", Quantity: suppliesResponse.Food.Fruits, Price: rand.Intn(100) + 1, Category: "food"})
		offers = append(offers, models.Offer{Name: "meat", Quantity: suppliesResponse.Food.Meat, Price: rand.Intn(100) + 1, Category: "food"})
		offers = append(offers, models.Offer{Name: "vegetables", Quantity: suppliesResponse.Food.Vegetables, Price: rand.Intn(100) + 1, Category: "food"})
		offers = append(offers, models.Offer{Name: "water", Quantity: suppliesResponse.Food.Water, Price: rand.Intn(100) + 1, Category: "food"})
		offers = append(offers, models.Offer{Name: "analgesics", Quantity: suppliesResponse.Medicine.Analgesics, Price: rand.Intn(100) + 1, Category: "medicine"})
		offers = append(offers, models.Offer{Name: "antibiotics", Quantity: suppliesResponse.Medicine.Antibiotics, Price: rand.Intn(100) + 1, Category: "medicine"})
		offers = append(offers, models.Offer{Name: "bandages", Quantity: suppliesResponse.Medicine.Bandages, Price: rand.Intn(100) + 1, Category: "medicine"})

		break
	}

	return offers, nil
}

/**
 * @brief Stores the provided supplies offers in the database.
 *
 * @param db The database connection.
 * @param offers The offers to be stored in the database.
 * @return An error if there was an issue during the database operation.
 */
func storeSuppliesInDB(db *gorm.DB, offers []models.Offer) error {
	for _, offer := range offers {
		if err := db.Create(&offer).Error; err != nil {
			return err
		}
	}
	return nil
}

/**
 * @brief Starts a periodic update that sends the current state of supplies to the server every 30 seconds.
 *
 * @param db The database connection.
 */
func startPeriodicUpdates(db *gorm.DB) {
	ticker := time.NewTicker(30 * time.Second) // 30 seconds.
	defer ticker.Stop()

	for range ticker.C {
		offers, err := getSuppliesFromDB(db)
		if err != nil {
			log.Printf("Failed to get supplies from DB: %v", err)
			continue
		}
		sendUpdatesToServer(offers)
	}
}

/**
 * @brief Fetches the current supplies from the database.
 *
 * @param db The database connection.
 * @return A slice of models.Offer containing the current supplies, and an error if there was an issue during the database operation.
 */
func getSuppliesFromDB(db *gorm.DB) ([]models.Offer, error) {
	var offers []models.Offer
	if err := db.Find(&offers).Error; err != nil {
		return offers, err
	}
	return offers, nil
}

/**
 * @brief Sends the current state of offers to the C++ server.
 *
 * @param offers The offers to be sent to the server.
 */
func sendUpdatesToServer(offers []models.Offer) {
	urltosend := os.Getenv("SEND_SUPPLIES_URL") // "http://127.0.0.1:8080"
	client := &http.Client{}

	for _, offer := range offers {
		data := map[string]interface{}{
			"command":  "apirest",
			"supply":   offer.Name,
			"quantity": offer.Quantity,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal data: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", urltosend, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to create POST request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send POST request: %v", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			continue
		}

		log.Printf("Server response: %s", body)
	}
}
