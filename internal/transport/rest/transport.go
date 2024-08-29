package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiloMIA/kami-test-task/internal/model"
	"github.com/kiloMIA/kami-test-task/internal/service"
)

type ReservationTransport struct {
	service *service.ReservationService
}

func NewReservationTransport(service *service.ReservationService) *ReservationTransport {
	return &ReservationTransport{service: service}
}

func (t *ReservationTransport) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation model.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reservation.StartTime.After(reservation.EndTime) || reservation.StartTime.Equal(reservation.EndTime) {
		http.Error(w, "Invalid time range", http.StatusBadRequest)
		return
	}

	if err := t.service.CreateReservation(r.Context(), &reservation); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *ReservationTransport) GetReservations(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	reservations, err := t.service.GetReservationsByRoomID(r.Context(), roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}
