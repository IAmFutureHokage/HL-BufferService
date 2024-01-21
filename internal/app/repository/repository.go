package repository

import (
	"context"
	"fmt"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HydrologyBufferRepository struct {
	dbPool *pgxpool.Pool
}

func NewHydrologyBufferRepository(pool *pgxpool.Pool) *HydrologyBufferRepository {
	return &HydrologyBufferRepository{dbPool: pool}
}

func (r *HydrologyBufferRepository) AddTelegram(ctx context.Context, data []model.Telegram) error {

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			fmt.Errorf("rollback erorr")
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Errorf("commit error")
		}
	}()

	for _, telegram := range data {
		telegramInsert := goqu.Insert("telegram").Rows(
			goqu.Record{
				"id":                         telegram.Id,
				"groupId":                    telegram.GroupId,
				"telegramCode":               telegram.TelegramCode,
				"postCode":                   telegram.PostCode,
				"dateTime":                   telegram.DateTime,
				"endBlockNum":                telegram.EndBlockNum,
				"isDangerous":                telegram.IsDangerous,
				"waterLevelOnTime":           telegram.WaterLevelOnTime,
				"deltaWaterLevel":            telegram.DeltaWaterLevel,
				"waterLevelOn20h":            telegram.WaterLevelOn20h,
				"waterTemperature":           telegram.WaterTemperature,
				"airTemperature":             telegram.AirTemperature,
				"icePhenomeniaState":         telegram.IcePhenomeniaState,
				"ice":                        telegram.Ice,
				"snow":                       telegram.Snow,
				"waterflow":                  telegram.Waterflow,
				"precipitationValue":         telegram.PrecipitationValue,
				"precipitationDuration":      telegram.PrecipitationDuration,
				"reservoirDate":              telegram.ReservoirDate,
				"headwaterLevel":             telegram.HeadwaterLevel,
				"averageReservoirLevel":      telegram.AverageReservoirLevel,
				"downstreamLevel":            telegram.DownstreamLevel,
				"reservoirVolume":            telegram.ReservoirVolume,
				"isReservoirWaterInflowDate": telegram.IsReservoirWaterInflowDate,
				"inflow":                     telegram.Inflow,
				"reset":                      telegram.Reset,
			},
		)

		sql, args, err := telegramInsert.ToSQL()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}

		for _, phenomen := range telegram.IcePhenomenia {
			phenomeniaInsert := goqu.Insert("phenomenia").Rows(
				goqu.Record{
					"id":          phenomen.Id,
					"telegramId":  phenomen.TelegramId,
					"phenomen":    phenomen.Phenomen,
					"isUntensity": phenomen.IsUntensity,
					"intensity":   phenomen.Intensity,
				},
			)

			sql, args, err := phenomeniaInsert.ToSQL()
			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, sql, args...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *HydrologyBufferRepository) GetTelegramByID(ctx context.Context, id uuid.UUID) (model.Telegram, error) {

	selectTelegramBuilder := goqu.From("telegram").Where(goqu.Ex{"id": id}).Limit(1)
	sql, args, err := selectTelegramBuilder.ToSQL()
	if err != nil {
		return model.Telegram{}, err
	}

	row := r.dbPool.QueryRow(ctx, sql, args...)

	var telegram model.Telegram
	err = row.Scan(
		&telegram.Id,
		&telegram.GroupId,
		&telegram.TelegramCode,
		&telegram.PostCode,
		&telegram.DateTime,
		&telegram.EndBlockNum,
		&telegram.IsDangerous,
		&telegram.WaterLevelOnTime,
		&telegram.DeltaWaterLevel,
		&telegram.WaterLevelOn20h,
		&telegram.WaterTemperature,
		&telegram.AirTemperature,
		&telegram.IcePhenomeniaState,
		&telegram.Ice,
		&telegram.Snow,
		&telegram.Waterflow,
		&telegram.PrecipitationValue,
		&telegram.PrecipitationDuration,
		&telegram.ReservoirDate,
		&telegram.HeadwaterLevel,
		&telegram.AverageReservoirLevel,
		&telegram.DownstreamLevel,
		&telegram.ReservoirVolume,
		&telegram.IsReservoirWaterInflowDate,
		&telegram.Inflow,
		&telegram.Reset,
	)
	if err != nil {
		return model.Telegram{}, err
	}

	selectPhenomeniaBuilder := goqu.From("phenomenia").Where(goqu.Ex{"telegramId": id})
	sql, args, err = selectPhenomeniaBuilder.ToSQL()
	if err != nil {
		return model.Telegram{}, err
	}

	rows, err := r.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return model.Telegram{}, err
	}
	defer rows.Close()

	var phenomeniaSlice []model.Phenomenia
	var phenomeniaCount int

	for rows.Next() {
		var phenomen model.Phenomenia
		err := rows.Scan(
			&phenomen.Id, &phenomen.TelegramId, &phenomen.Phenomen, &phenomen.IsUntensity, &phenomen.Intensity,
		)
		if err != nil {
			return model.Telegram{}, err
		}
		phenomeniaSlice = append(phenomeniaSlice, phenomen)
		phenomeniaCount++
	}

	phenomeniaArray := make([]*model.Phenomenia, phenomeniaCount)
	for i, p := range phenomeniaSlice {
		phenomeniaArray[i] = &p
	}

	telegram.IcePhenomenia = phenomeniaArray

	return telegram, nil
}

func (r *HydrologyBufferRepository) RemoveTelegram(ctx context.Context, id uuid.UUID) error {

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			fmt.Errorf("rollback error")
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Errorf("commit error")
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM phenomenia WHERE telegramId = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM telegram WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *HydrologyBufferRepository) GetAll(ctx context.Context) ([]model.Telegram, error) {

	selectBuilder := goqu.From("telegram")

	sql, args, err := selectBuilder.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telegrams []model.Telegram

	for rows.Next() {
		var telegram model.Telegram
		err := rows.Scan(
			&telegram.Id,
			&telegram.GroupId,
			&telegram.TelegramCode,
			&telegram.PostCode,
			&telegram.DateTime,
			&telegram.EndBlockNum,
			&telegram.IsDangerous,
			&telegram.WaterLevelOnTime,
			&telegram.DeltaWaterLevel,
			&telegram.WaterLevelOn20h,
			&telegram.WaterTemperature,
			&telegram.AirTemperature,
			&telegram.IcePhenomeniaState,
			&telegram.Ice,
			&telegram.Snow,
			&telegram.Waterflow,
			&telegram.PrecipitationValue,
			&telegram.PrecipitationDuration,
			&telegram.ReservoirDate,
			&telegram.HeadwaterLevel,
			&telegram.AverageReservoirLevel,
			&telegram.DownstreamLevel,
			&telegram.ReservoirVolume,
			&telegram.IsReservoirWaterInflowDate,
			&telegram.Inflow,
			&telegram.Reset,
		)
		if err != nil {
			return nil, err
		}

		telegrams = append(telegrams, telegram)
	}

	return telegrams, nil
}
