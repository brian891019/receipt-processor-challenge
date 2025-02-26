package main

import (
	"log"
	"net/http"

	"example.com/takehome/handler"
	"example.com/takehome/service"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// init service
	pointService := service.NewPointService()

	// init handler
	h := handler.NewHandler(pointService)

	// New router
	router := httprouter.New()

	// Setup endpoints
	router.GET("/receipts/:id/points", h.GetPoints)
	router.POST("/receipts/process", h.ProcessReceipt)
	log.Printf("Starting server on %s\n", ":8080")
	http.ListenAndServe(":8080", router)
}
