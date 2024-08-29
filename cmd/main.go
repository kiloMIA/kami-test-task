// cmd/main.go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiloMIA/kami-test-task/internal/logger"
	"github.com/kiloMIA/kami-test-task/internal/repo"
	"github.com/kiloMIA/kami-test-task/internal/repo/postgre"
	"github.com/kiloMIA/kami-test-task/internal/service"
	"github.com/kiloMIA/kami-test-task/internal/transport/rest"
)

func main() {
	log := logger.CreateLogger()
	defer log.Sync()

	dbpool := postgre.ConnectDB(log)
	defer dbpool.Close()

	repository := repo.NewRepository(dbpool, log)
	reservationService := service.NewReservationService(repository)

	transport := rest.NewReservationTransport(reservationService)

	r := chi.NewRouter()
	r.Post("/reservations", transport.CreateReservation)
	r.Get("/reservations/{room_id}", transport.GetReservations)

	log.Info("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
