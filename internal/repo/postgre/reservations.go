package postgre

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiloMIA/kami-test-task/internal/model"
	"go.uber.org/zap"
)

type ReservationRepository struct {
	dbpool *pgxpool.Pool
	logger *zap.Logger
}

func NewReservationRepository(dbpool *pgxpool.Pool, logger *zap.Logger) *ReservationRepository {
	return &ReservationRepository{
		dbpool: dbpool,
		logger: logger,
	}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	query := `INSERT INTO reservations (room_id, start_time, end_time) VALUES ($1, $2, $3)`
	_, err := r.dbpool.Exec(ctx, query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		r.logger.Error("Failed to create reservation", zap.Error(err))
	}
	return err
}

func (r *ReservationRepository) GetByRoomID(ctx context.Context, roomID string) ([]model.Reservation, error) {
	query := `SELECT id, room_id, start_time, end_time FROM reservations WHERE room_id = $1 ORDER BY start_time`
	rows, err := r.dbpool.Query(ctx, query, roomID)
	if err != nil {
		r.logger.Error("Failed to get reservations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var reservation model.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (r *ReservationRepository) CheckConflict(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM reservations WHERE room_id = $1 AND ($2 < end_time AND $3 > start_time)`
	var count int
	err := r.dbpool.QueryRow(ctx, query, roomID, startTime, endTime).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to check conflicts", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
