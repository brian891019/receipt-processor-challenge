package service

import (
	"testing"

	"example.com/takehome/model"
)

func TestProcessReceipt(t *testing.T) {

	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []model.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}

	receipt2 := model.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []model.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "Klarbrunn 12-PK 12 FL OZ", Price: "12.00"},
		},
		Total: "35.35",
	}

	checkPoints(t, pointService, receipt, 109)
	checkPoints(t, pointService, receipt2, 28)
}

func checkPoints(t *testing.T, pointService PointService, receipt model.Receipt, expectedPoints int) {
	id, err := pointService.ProcessReceipt(receipt)

	if err != nil {
		t.Fatalf("Failed to process receipt: %v", err)
	}
	if id == "" {
		t.Errorf("Expected non-empty id, got empty")
	}

	points, err := pointService.GetPoint(id)
	if err != nil {
		t.Fatalf("Failed to get points: %v", err)
	}
	if points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, points)
	}
}

func TestInvalidReceipt(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer: "",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid receipt")
	}
}

func TestEmptyReceipt(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "",
		PurchaseDate: "",
		PurchaseTime: "",
		Items:        []model.Item{},
		Total:        "",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty receipt")
	}
}

func TestInvalidTotalReceipt(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-02-04",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "10.00"},
		},
		Total: "invalid",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid total receipt")
	}
}

func TestEmptyTotalReceipt(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-02-04",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "10.00"},
		},
		Total: "",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty total receipt")
	}
}

func TestInvalidDate(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "invalid",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid date")
	}
}

func TestEmptyDate(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty date")
	}
}

func TestInvalidTime(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "invalid",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid time")
	}
}
func TestEmptyTime(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty time")
	}
}
func TestEmptyDescription(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid short description")
	}
}

func TestInvalidPrice(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "invalid"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for invalid price")
	}
}
func TestEmptyPrice(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: ""},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty price")
	}
}
func TestEmptyRetailer(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-19",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		},
		Total: "10.00",
	}

	_, err := pointService.ProcessReceipt(receipt)
	if err == nil {
		t.Errorf("Expected error for empty retailer")
	}
}

func TestGetPoint(t *testing.T) {
	pointService := NewPointService()

	_, err := pointService.GetPoint("non_existing_id")
	if err == nil {
		t.Errorf("Expected error for not existing id")
	}
}
