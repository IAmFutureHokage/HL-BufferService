package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HydrologyBufferStorage struct {
	dbPool *pgxpool.Pool
}

func NewHydrologyBufferStorage(pool *pgxpool.Pool) *HydrologyBufferStorage {
	return &HydrologyBufferStorage{dbPool: pool}
}

func (r *HydrologyBufferStorage) AddTelegram(ctx context.Context, data []model.Telegram) error {

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

func (r *HydrologyBufferStorage) GetTelegramByID(ctx context.Context, id uuid.UUID) (model.Telegram, error) {

	selectTelegramBuilder := goqu.From("telegram").
		Where(goqu.Ex{"id": id}).
		Limit(1)

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

	selectPhenomeniaBuilder := goqu.From("phenomenia").
		Where(goqu.Ex{"telegramid": id})

	sql, args, err = selectPhenomeniaBuilder.ToSQL()
	if err != nil {
		return model.Telegram{}, err
	}

	rows, err := r.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return model.Telegram{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var phenomenia model.Phenomenia
		err := rows.Scan(
			&phenomenia.Id,
			&phenomenia.TelegramId,
			&phenomenia.Phenomen,
			&phenomenia.IsUntensity,
			&phenomenia.Intensity,
		)
		if err != nil {
			return model.Telegram{}, err
		}
		telegram.IcePhenomenia = append(telegram.IcePhenomenia, &phenomenia)
	}

	return telegram, nil
}

func (r *HydrologyBufferStorage) RemoveTelegrams(ctx context.Context, ids []uuid.UUID) error {

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

	result, err := tx.Exec(ctx, "DELETE FROM phenomenia WHERE telegramId = ANY($1)", ids)
	if err != nil {
		return err
	}

	result, err = tx.Exec(ctx, "DELETE FROM telegram WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("no matching rows in telegram")
	}

	return nil
}

func (r *HydrologyBufferStorage) GetAll(ctx context.Context) ([]model.Telegram, error) {

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

		phenomeniaBuilder := goqu.From("phenomenia").
			Where(goqu.Ex{"telegramid": telegram.Id})

		phenomeniaSQL, phenomeniaArgs, err := phenomeniaBuilder.ToSQL()
		if err != nil {
			return nil, err
		}

		phenomeniaRows, err := r.dbPool.Query(ctx, phenomeniaSQL, phenomeniaArgs...)
		if err != nil {
			return nil, err
		}
		defer phenomeniaRows.Close()

		for phenomeniaRows.Next() {
			var phenomenia model.Phenomenia
			err := phenomeniaRows.Scan(
				&phenomenia.Id,
				&phenomenia.TelegramId,
				&phenomenia.Phenomen,
				&phenomenia.IsUntensity,
				&phenomenia.Intensity,
			)
			if err != nil {
				return nil, err
			}

			telegram.IcePhenomenia = append(telegram.IcePhenomenia, &phenomenia)
		}

		if err := phenomeniaRows.Err(); err != nil {
			return nil, err
		}

		telegrams = append(telegrams, telegram)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return telegrams, nil
}

func (r *HydrologyBufferStorage) UpdateTelegram(ctx context.Context, updatedTelegram model.Telegram) error {

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
			"groupid":                    updatedTelegram.GroupId,
			"telegramcode":               updatedTelegram.TelegramCode,
			"postcode":                   updatedTelegram.PostCode,
			"datetime":                   updatedTelegram.DateTime,
			"endblocknum":                updatedTelegram.EndBlockNum,
			"isdangerous":                updatedTelegram.IsDangerous,
			"waterlevelontime":           updatedTelegram.WaterLevelOnTime,
			"deltawaterlevel":            updatedTelegram.DeltaWaterLevel,
			"waterlevelon20h":            updatedTelegram.WaterLevelOn20h,
			"watertemperature":           updatedTelegram.WaterTemperature,
			"airtemperature":             updatedTelegram.AirTemperature,
			"icephenomeniastate":         updatedTelegram.IcePhenomeniaState,
			"ice":                        updatedTelegram.Ice,
			"snow":                       updatedTelegram.Snow,
			"waterflow":                  updatedTelegram.Waterflow,
			"precipitationvalue":         updatedTelegram.PrecipitationValue,
			"precipitationduration":      updatedTelegram.PrecipitationDuration,
			"reservoirdate":              updatedTelegram.ReservoirDate,
			"headwaterlevel":             updatedTelegram.HeadwaterLevel,
			"averagereservoirlevel":      updatedTelegram.AverageReservoirLevel,
			"downstreamlevel":            updatedTelegram.DownstreamLevel,
			"reservoirvolume":            updatedTelegram.ReservoirVolume,
			"isreservoirwaterinflowdate": updatedTelegram.IsReservoirWaterInflowDate,
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
				"telegramid":  updatedTelegram.Id,
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

	return nil
}

func (r *HydrologyBufferStorage) GetTelegramsById(ctx context.Context, ids []uuid.UUID) ([]model.Telegram, error) {

	return nil, nil
}
