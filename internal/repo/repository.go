package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiloMIA/kami-test-task/internal/model"
	"github.com/kiloMIA/kami-test-task/internal/repo/postgre"
	"go.uber.org/zap"
)

type ReservationRepo interface {
	Create(ctx context.Context, reservation *model.Reservation) error
	GetByRoomID(ctx context.Context, roomID string) ([]model.Reservation, error)
	CheckConflict(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
}

type Repository struct {
	ReservationRepo
}

func NewRepository(dbpool *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		ReservationRepo: postgre.NewReservationRepository(dbpool, logger),
	}
}
