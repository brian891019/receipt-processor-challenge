package service

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"example.com/takehome/model"
	"github.com/google/uuid"
)

type PointService interface {
	ProcessReceipt(receipt model.Receipt) (string, error)
	GetPoint(id string) (int, error)
}

type pointService struct {
	IdToPointMap map[string]int
}

func NewPointService() PointService {
	return &pointService{
		IdToPointMap: make(map[string]int),
	}
}

// processes receipt and stores the points in memory.
func (s *pointService) ProcessReceipt(receipt model.Receipt) (string, error) {
	id := uuid.New().String()
	receipt_valid, err := s.ValidateReceipt(receipt)
	if err != nil {
		return id, errors.New("error validating receipt")
	}
	if !receipt_valid {
		return id, errors.New("receipt contains error field")
	}

	points, err := s.calculatePoints(receipt)
	if err != nil {
		return id, errors.New("error calculating points")
	}
	if receipt.Retailer == "" {
		return "", errors.New("invalid receipt data")
	}
	s.IdToPointMap[id] = points
	return id, nil
}

// check if any field on the receipt is empty, if it is throw an error
func (s *pointService) ValidateReceipt(receipt model.Receipt) (bool, error) {
	total := receipt.Total
	purchase_date := receipt.PurchaseDate
	purchase_time := receipt.PurchaseTime
	retailer := receipt.Retailer

	if total == "" || purchase_time == "" || purchase_date == "" || retailer == "" {
		return false, errors.New("any receipt field shouldn't be empty")
	}

	for _, item := range receipt.Items {
		if item.Price == "" || item.ShortDescription == "" {
			return false, errors.New("any item field shouldn't be empty")
		}
	}

	return true, nil
}

// retrieves the points for a given receipt ID.
func (s *pointService) GetPoint(id string) (int, error) {
	points, ok := s.IdToPointMap[id]
	if !ok {
		return 0, model.ErrNotFound
	}
	return points, nil
}

func (s *pointService) calculatePoints(receipt model.Receipt) (int, error) {
	points := 0

	// Parse total to float
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return points, errors.New("invalid total")
	}
	// 25 points if the total is a multiple of 0.25.
	if int(total*100)%25 == 0 {
		points += 25
	}
	// 50 points if the total is a round dollar amount with no cents.
	if total == float64(int(total)) {
		points += 50
	}

	// One point for every alphanumeric character in the retailer name.
	for _, char := range receipt.Retailer {
		if alphanumeric(char) {
			points++
		}
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3,
	//multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLength%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return points, errors.New("invalid item Price")

			}
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return points, errors.New("invalid purchaseTime")
	}
	two_pm, _ := time.Parse("15:04", "14:00")
	four_pm, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(two_pm) && purchaseTime.Before(four_pm) {
		points += 10
	}

	// 6 points if the day in the purchase date is odd.
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return points, errors.New("invalid PurchaseDate")
	}
	if date.Day()%2 != 0 {
		points += 6
	}

	return points, nil
}

// alphanumeric checks if a character is alphanumeric.
func alphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z')
}
