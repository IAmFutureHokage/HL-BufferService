package repository

import (
	"context"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HydrologyBufferRepository struct {
	dbPool *pgxpool.Pool
}

func NewHydrologyBufferRepository(pool *pgxpool.Pool) *HydrologyBufferRepository {
	return &HydrologyBufferRepository{dbPool: pool}
}

func (r *HydrologyBufferRepository) GetTelegramByID(ctx context.Context, id uuid.UUID) (model.Telegram, error) {

	var telegram model.Telegram

	return telegram, nil
}

func (r *HydrologyBufferRepository) GetAll(ctx context.Context) ([]model.Telegram, error) {

	telegrams := make([]model.Telegram, 0)

	for i := 0; i < 3; i++ {

		telegram := model.Telegram{}
		telegrams = append(telegrams, telegram)
	}

	return telegrams, nil
}
