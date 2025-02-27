package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/takehome/model"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
)

func TestProcessReceiptHandler(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.POST("/receipts/process", h.ProcessReceipt)

	receipt := model.Receipt{}
	mockPointService.On("ProcessReceipt", mock.Anything).Return("mockReceiptID", nil)

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
	if result_id.ID != "mockReceiptID" {
		t.Fatalf("Receipt ID not match")
	}

	mockPointService.AssertExpectations(t)

}

func TestGetPointsHandler(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.GET("/receipts/:id/points", h.GetPoints)

	receiptID := "mockReceiptID"
	mockPointService.On("GetPoint", receiptID).Return(28, nil)

	req := httptest.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
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
	mockPointService.AssertExpectations(t)

}

func TestProcessErrorReceiptHandler(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.POST("/receipts/process", h.ProcessReceipt)

	receipt := model.Receipt{}
	mockPointService.On("ProcessReceipt", mock.Anything).Return("", model.ErrInvalidReceipt)

	receipt_marshaled, _ := json.Marshal(receipt)
	request := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receipt_marshaled))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	response := w.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status bad request error, got %d", response.StatusCode)
	}
	mockPointService.AssertExpectations(t)

}

func TestErrorGetPointsHandler(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.GET("/receipts/:id/points", h.GetPoints)

	receiptID := "invalidID"
	mockPointService.On("GetPoint", receiptID).Return(0, model.ErrNotFound)

	req := httptest.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status not found error, got %d", response.StatusCode)
	}
	mockPointService.AssertExpectations(t)

}

func TestInternalServerError(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.POST("/receipts/process", h.ProcessReceipt)

	receipt := model.Receipt{}
	mockPointService.On("ProcessReceipt", mock.Anything).Return("", errors.New("Internal Server Error"))

	receipt_marshaled, _ := json.Marshal(receipt)
	request := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receipt_marshaled))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	response := w.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status internal server error, got %d", response.StatusCode)
	}
	mockPointService.AssertExpectations(t)

}

func TestGetPointsError(t *testing.T) {
	mockPointService := new(MockPointService)
	router := httprouter.New()
	h := NewHandler(mockPointService)
	router.GET("/receipts/:id/points", h.GetPoints)

	receiptID := "invalidID"
	mockPointService.On("GetPoint", receiptID).Return(0, errors.New("Internal Server Error"))

	req := httptest.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status bad request error, got %d", response.StatusCode)
	}
	mockPointService.AssertExpectations(t)

}
