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
			fmt.Println("rollback error:", err)
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Println("commit error:", err)
		}
	}()

	for _, telegram := range data {
		telegramInsert := goqu.Insert("telegram").Rows(
			goqu.Record{
				"id":                         telegram.Id,
				"groupid":                    telegram.GroupId,
				"telegramcode":               telegram.TelegramCode,
				"postcode":                   telegram.PostCode,
				"datetime":                   telegram.DateTime,
				"endblocknum":                telegram.EndBlockNum,
				"isdangerous":                telegram.IsDangerous,
				"waterlevelontime":           telegram.WaterLevelOnTime,
				"deltawaterlevel":            telegram.DeltaWaterLevel,
				"waterlevelon20h":            telegram.WaterLevelOn20h,
				"watertemperature":           telegram.WaterTemperature,
				"airtemperature":             telegram.AirTemperature,
				"icephenomeniastate":         telegram.IcePhenomeniaState,
				"ice":                        telegram.Ice,
				"snow":                       telegram.Snow,
				"waterflow":                  telegram.Waterflow,
				"precipitationvalue":         telegram.PrecipitationValue,
				"precipitationduration":      telegram.PrecipitationDuration,
				"reservoirdate":              telegram.ReservoirDate,
				"headwaterlevel":             telegram.HeadwaterLevel,
				"averagereservoirlevel":      telegram.AverageReservoirLevel,
				"downstreamlevel":            telegram.DownstreamLevel,
				"reservoirvolume":            telegram.ReservoirVolume,
				"isreservoirwaterinflowdate": telegram.IsReservoirWaterInflowDate,
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
					"telegramid":  phenomen.TelegramId,
					"phenomen":    phenomen.Phenomen,
					"isuntensity": phenomen.IsUntensity,
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

	return telegram, nil
}

func (r *HydrologyBufferRepository) RemoveTelegrams(ctx context.Context, ids []uuid.UUID) error {

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			fmt.Println("rollback error:", err)
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Println("commit error:", err)
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM phenomenia WHERE telegramId = ANY($1)", ids)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM telegram WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}

	return nil
}

func (r *HydrologyBufferRepository) GetAll(ctx context.Context) ([]model.Telegram, error) {

	var rowCount int
	selectBuilder := goqu.From("telegram")

	sql, args, err := selectBuilder.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rowCount++
	}

	defer rows.Close()

	rows, err = r.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	telegrams := make([]model.Telegram, rowCount)

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

func (r *HydrologyBufferRepository) UpdateTelegram(ctx context.Context, updatedTelegram model.Telegram) error {

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			fmt.Println("rollback error:", err)
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			fmt.Println("commit error:", err)
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM phenomenia WHERE telegramId = $1", updatedTelegram.Id)
	if err != nil {
		return err
	}

	telegramUpdate := goqu.Update("telegram").
		Set(goqu.Record{
			"groupId":                    updatedTelegram.GroupId,
			"telegramCode":               updatedTelegram.TelegramCode,
			"postCode":                   updatedTelegram.PostCode,
			"dateTime":                   updatedTelegram.DateTime,
			"endBlockNum":                updatedTelegram.EndBlockNum,
			"isDangerous":                updatedTelegram.IsDangerous,
			"waterLevelOnTime":           updatedTelegram.WaterLevelOnTime,
			"deltaWaterLevel":            updatedTelegram.DeltaWaterLevel,
			"waterLevelOn20h":            updatedTelegram.WaterLevelOn20h,
			"waterTemperature":           updatedTelegram.WaterTemperature,
			"airTemperature":             updatedTelegram.AirTemperature,
			"icePhenomeniaState":         updatedTelegram.IcePhenomeniaState,
			"ice":                        updatedTelegram.Ice,
			"snow":                       updatedTelegram.Snow,
			"waterflow":                  updatedTelegram.Waterflow,
			"precipitationValue":         updatedTelegram.PrecipitationValue,
			"precipitationDuration":      updatedTelegram.PrecipitationDuration,
			"reservoirDate":              updatedTelegram.ReservoirDate,
			"headwaterLevel":             updatedTelegram.HeadwaterLevel,
			"averageReservoirLevel":      updatedTelegram.AverageReservoirLevel,
			"downstreamLevel":            updatedTelegram.DownstreamLevel,
			"reservoirVolume":            updatedTelegram.ReservoirVolume,
			"isReservoirWaterInflowDate": updatedTelegram.IsReservoirWaterInflowDate,
			"inflow":                     updatedTelegram.Inflow,
			"reset":                      updatedTelegram.Reset,
		}).
		Where(goqu.Ex{"id": updatedTelegram.Id})

	sql, args, err := telegramUpdate.ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	for _, phenomen := range updatedTelegram.IcePhenomenia {
		phenomeniaInsert := goqu.Insert("phenomenia").Rows(
			goqu.Record{
				"id":          phenomen.Id,
				"telegramId":  updatedTelegram.Id,
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
	return nil
}
