package service

import (
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/chendicao/receipt-processor/models"
)

func CalculatePoints(receipt *models.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points += 1
		}
	}

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += 5 * (len(receipt.Items) / 2)

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
	for _, item := range receipt.Items {
		// Directly trimming and checking length for better efficiency
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			additionalPoints := int(math.Ceil(item.Price * 0.2))
			points += additionalPoints
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDay, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDay.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.Hour() > 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}
