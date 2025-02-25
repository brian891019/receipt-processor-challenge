package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/takehome/model"
	"example.com/takehome/service"

	"github.com/julienschmidt/httprouter"
)

func TestProcessReceiptHandler(t *testing.T) {
	router := httprouter.New()
	pointService := service.NewPointService()
	h := NewHandler(pointService)

	router.POST("/receipts/process", h.ProcessReceipt)

	receipt := model.Receipt{
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

	receipt_marshaled, _ := json.Marshal(receipt)
	request := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receipt_marshaled))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)

	response := w.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	var result_id model.IDResponse
	err := json.NewDecoder(response.Body).Decode(&result_id)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result_id.ID == "" {
		t.Errorf("Expected non-empty id")
	}
}

func TestGetPointsHandler(t *testing.T) {
	pointService := service.NewPointService()
	h := NewHandler(pointService)

	router := httprouter.New()
	router.GET("/receipts/:id/points", h.GetPoints)

	receipt := model.Receipt{
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

	id, err := pointService.ProcessReceipt(receipt)
	if err != nil {
		t.Fatalf("Failed to process receipt: %v", err)
	}

	// Use the ID to get points
	req := httptest.NewRequest("GET", "/receipts/"+id+"/points", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var pointsResponse model.PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointsResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedPoints := 28
	if pointsResponse.Points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, pointsResponse.Points)
	}
}
