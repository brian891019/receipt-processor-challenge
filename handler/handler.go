package handler

import (
	"encoding/json"
	"net/http"

	"example.com/takehome/model"
	"example.com/takehome/service"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	pointService service.PointService
}

func NewHandler(pointService service.PointService) handler {
	return handler{
		pointService: pointService,
	}
}

func (h handler) ProcessReceipt(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// unmarshal request to model.Receipt
	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "JSON file error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id, err := h.pointService.ProcessReceipt(receipt)

	if err != nil {
		if err == model.ErrInvalidReceipt || err == model.ErrCalculatePoint {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(model.IDResponse{ID: id})
}

func (h handler) GetPoints(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// unmarshal request
	id := param.ByName("id")
	points, err := h.pointService.GetPoint(id)
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(model.PointsResponse{Points: points})
}
