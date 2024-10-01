package helpers

import (
	"os"
)

func GetStripeAPIKey() string {
	return os.Getenv("STRIPE_API_KEY") // Ensure the API key is set in environment variables
}
