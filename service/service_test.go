package service

import (
	"testing"

	"example.com/takehome/model"
)

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

func TestZeroTotalReceipt(t *testing.T) {
	pointService := NewPointService()

	receipt := model.Receipt{
		Retailer:     "traderjoes",
		PurchaseDate: "2022-02-04",
		PurchaseTime: "06:00",
		Items: []model.Item{
			{ShortDescription: "Knorr Creamy Chicken", Price: "0.00"},
		},
		Total: "0.00",
	}

	id, err := pointService.ProcessReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to process receipt: %v", err)
	}

	points, err := pointService.GetPoint(id)
	if err != nil {
		t.Fatalf("Failed to get points: %v", err)
	}

	expectedPoints := 0
	if points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, points)
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

	id, err := pointService.ProcessReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to process receipt: %v", err)
	}

	_, err = pointService.GetPoint(id)
	if err != nil {
		t.Errorf("Failed to get points: %v", err)
	}
}
