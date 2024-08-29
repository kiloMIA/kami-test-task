package service

import (
	"context"
	"errors"

	"github.com/kiloMIA/kami-test-task/internal/model"
	"github.com/kiloMIA/kami-test-task/internal/repo"
)

type ReservationService struct {
	repo repo.ReservationRepo
}

func NewReservationService(repo repo.ReservationRepo) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, reservation *model.Reservation) error {
	conflict, err := s.repo.CheckConflict(ctx, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	if conflict {
		return errors.New("reservation conflict")
	}

	return s.repo.Create(ctx, reservation)
}

func (s *ReservationService) GetReservationsByRoomID(ctx context.Context, roomID string) ([]model.Reservation, error) {
	return s.repo.GetByRoomID(ctx, roomID)
}
