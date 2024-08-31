package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiloMIA/kami-test-task/internal/logger"
	"github.com/kiloMIA/kami-test-task/internal/model"
	"github.com/kiloMIA/kami-test-task/internal/repo"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	connString := "postgres://user:password@localhost:5432/db?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	_, err = dbpool.Exec(context.Background(), "TRUNCATE TABLE reservations")
	if err != nil {
		t.Fatalf("Failed to clean reservations table: %v", err)
	}

	return dbpool
}

func TestCreateReservation_Success(t *testing.T) {
	dbpool := setupTestDB(t)
	defer dbpool.Close()

	log := logger.CreateLogger()
	repository := repo.NewRepository(dbpool, log)
	service := NewReservationService(repository)

	reservation := &model.Reservation{
		RoomID:    "Room1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := service.CreateReservation(context.Background(), reservation)
	assert.NoError(t, err, "Expected no error when creating a reservation without conflicts")
}

func TestCreateReservation_Conflict(t *testing.T) {
	dbpool := setupTestDB(t)
	defer dbpool.Close()

	log := logger.CreateLogger()
	repository := repo.NewRepository(dbpool, log)
	service := NewReservationService(repository)

	reservation1 := &model.Reservation{
		RoomID:    "Room1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	err := service.CreateReservation(context.Background(), reservation1)
	assert.NoError(t, err)

	reservation2 := &model.Reservation{
		RoomID:    "Room1",
		StartTime: time.Now().Add(90 * time.Minute),
		EndTime:   time.Now().Add(3 * time.Hour),
	}
	err = service.CreateReservation(context.Background(), reservation2)
	assert.Error(t, err, "Expected an error when creating a conflicting reservation")
}

func TestConcurrentReservations(t *testing.T) {
	dbpool := setupTestDB(t)
	defer dbpool.Close()

	log := logger.CreateLogger()
	repository := repo.NewRepository(dbpool, log)
	service := NewReservationService(repository)

	createReservation := func(roomID string, startTime, endTime time.Time) error {
		reservation := &model.Reservation{
			RoomID:    roomID,
			StartTime: startTime,
			EndTime:   endTime,
		}
		return service.CreateReservation(context.Background(), reservation)
	}

	var err1, err2 error
	done := make(chan bool)

	go func() {
		err1 = createReservation("Room1", time.Now().Add(1*time.Hour), time.Now().Add(2*time.Hour))
		done <- true
	}()

	go func() {
		err2 = createReservation("Room1", time.Now().Add(90*time.Minute), time.Now().Add(3*time.Hour))
		done <- true
	}()

	<-done
	<-done

	if err1 == nil {
		assert.Error(t, err2, "Expected one of the concurrent reservations to fail")
	} else {
		assert.Error(t, err1, "Expected one of the concurrent reservations to fail")
	}
}
