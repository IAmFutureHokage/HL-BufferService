package postgres

import (
	"context"
	"database/sql"
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

func (r *HydrologyBufferStorage) GetTelegramByID(ctx context.Context, id uuid.UUID) (*model.Telegram, error) {

	selectBuilder := goqu.
		From("telegram").
		Select(
			goqu.I("telegram.id"),
			goqu.I("telegram.groupid"),
			goqu.I("telegram.telegramcode"),
			goqu.I("telegram.postcode"),
			goqu.I("telegram.datetime"),
			goqu.I("telegram.endblocknum"),
			goqu.I("telegram.isdangerous"),
			goqu.I("telegram.waterlevelontime"),
			goqu.I("telegram.deltawaterlevel"),
			goqu.I("telegram.waterlevelon20h"),
			goqu.I("telegram.watertemperature"),
			goqu.I("telegram.airtemperature"),
			goqu.I("telegram.icephenomeniastate"),
			goqu.I("telegram.ice"),
			goqu.I("telegram.snow"),
			goqu.I("telegram.waterflow"),
			goqu.I("telegram.precipitationvalue"),
			goqu.I("telegram.precipitationduration"),
			goqu.I("telegram.reservoirdate"),
			goqu.I("telegram.headwaterlevel"),
			goqu.I("telegram.averagereservoirlevel"),
			goqu.I("telegram.downstreamlevel"),
			goqu.I("telegram.reservoirvolume"),
			goqu.I("telegram.isreservoirwaterinflowdate"),
			goqu.I("telegram.inflow"),
			goqu.I("telegram.reset"),
			goqu.I("phenomenia.id").As("phenomenia_id"),
			goqu.I("phenomenia.telegramid"),
			goqu.I("phenomenia.phenomen"),
			goqu.I("phenomenia.isuntensity"),
			goqu.I("phenomenia.intensity"),
		).
		LeftJoin(
			goqu.I("phenomenia"),
			goqu.On(goqu.Ex{"telegram.id": goqu.I("phenomenia.telegramid")}),
		).
		Where(goqu.Ex{"telegram.id": id})

	sqlScript, args, err := selectBuilder.ToSQL()
	if err != nil {
		return &model.Telegram{}, err
	}

	rows, err := r.dbPool.Query(ctx, sqlScript, args...)
	if err != nil {
		return &model.Telegram{}, err
	}
	defer rows.Close()

	var tg model.Telegram

	for rows.Next() {
		var telegram model.Telegram
		var phenomeniaId *uuid.UUID
		var phenomeniaTelegramId uuid.UUID
		var phenomeniaPhenomen sql.NullByte
		var phenomeniaIsUntensity sql.NullBool
		var phenomeniaIntensity sql.NullByte

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
			&phenomeniaId,
			&phenomeniaTelegramId,
			&phenomeniaPhenomen,
			&phenomeniaIsUntensity,
			&phenomeniaIntensity,
		)
		if err != nil {
			return &model.Telegram{}, err
		}

		if tg.Id != id {
			tg = telegram
		}

		if phenomeniaId != nil {
			tg.IcePhenomenia = append(tg.IcePhenomenia, &model.Phenomenia{
				Id:          *phenomeniaId,
				TelegramId:  phenomeniaTelegramId,
				Phenomen:    phenomeniaPhenomen.Byte,
				IsUntensity: phenomeniaIsUntensity.Bool,
				Intensity:   phenomeniaIntensity,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return &model.Telegram{}, err
	}

	return &tg, nil
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

	result, err := tx.Exec(ctx, "DELETE FROM telegram WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("no matching rows in telegram")
	}

	return nil
}

func (r *HydrologyBufferStorage) GetAll(ctx context.Context) (*[]model.Telegram, error) {

	selectBuilder := goqu.
		From("telegram").
		Select(
			goqu.I("telegram.id"),
			goqu.I("telegram.groupid"),
			goqu.I("telegram.telegramcode"),
			goqu.I("telegram.postcode"),
			goqu.I("telegram.datetime"),
			goqu.I("telegram.endblocknum"),
			goqu.I("telegram.isdangerous"),
			goqu.I("telegram.waterlevelontime"),
			goqu.I("telegram.deltawaterlevel"),
			goqu.I("telegram.waterlevelon20h"),
			goqu.I("telegram.watertemperature"),
			goqu.I("telegram.airtemperature"),
			goqu.I("telegram.icephenomeniastate"),
			goqu.I("telegram.ice"),
			goqu.I("telegram.snow"),
			goqu.I("telegram.waterflow"),
			goqu.I("telegram.precipitationvalue"),
			goqu.I("telegram.precipitationduration"),
			goqu.I("telegram.reservoirdate"),
			goqu.I("telegram.headwaterlevel"),
			goqu.I("telegram.averagereservoirlevel"),
			goqu.I("telegram.downstreamlevel"),
			goqu.I("telegram.reservoirvolume"),
			goqu.I("telegram.isreservoirwaterinflowdate"),
			goqu.I("telegram.inflow"),
			goqu.I("telegram.reset"),
			goqu.I("phenomenia.id").As("phenomenia_id"),
			goqu.I("phenomenia.telegramid"),
			goqu.I("phenomenia.phenomen"),
			goqu.I("phenomenia.isuntensity"),
			goqu.I("phenomenia.intensity"),
		).
		LeftJoin(
			goqu.I("phenomenia"),
			goqu.On(goqu.Ex{"telegram.id": goqu.I("phenomenia.telegramid")}),
		)

	sqlScript, args, err := selectBuilder.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(ctx, sqlScript, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telegrams []model.Telegram

	for rows.Next() {
		var telegram model.Telegram
		var phenomeniaId *uuid.UUID
		var phenomeniaTelegramId uuid.UUID
		var phenomeniaPhenomen sql.NullByte
		var phenomeniaIsUntensity sql.NullBool
		var phenomeniaIntensity sql.NullByte

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
			&phenomeniaId,
			&phenomeniaTelegramId,
			&phenomeniaPhenomen,
			&phenomeniaIsUntensity,
			&phenomeniaIntensity,
		)
		if err != nil {
			return nil, err
		}

		if len(telegrams) == 0 || telegram.Id != telegrams[len(telegrams)-1].Id {
			telegrams = append(telegrams, telegram)
		}

		if phenomeniaId != nil {
			telegrams[len(telegrams)-1].IcePhenomenia = append(telegrams[len(telegrams)-1].IcePhenomenia, &model.Phenomenia{
				Id:          *phenomeniaId,
				TelegramId:  phenomeniaTelegramId,
				Phenomen:    phenomeniaPhenomen.Byte,
				IsUntensity: phenomeniaIsUntensity.Bool,
				Intensity:   phenomeniaIntensity,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &telegrams, nil
}

func (r *HydrologyBufferStorage) UpdateTelegram(ctx context.Context, updatedTelegram *model.Telegram) error {

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

func (r *HydrologyBufferStorage) GetTelegramsById(ctx context.Context, ids []uuid.UUID) (*[]model.Telegram, error) {

	selectBuilder := goqu.
		Select(
			goqu.I("t.id"),
			goqu.I("t.groupid"),
			goqu.I("t.telegramcode"),
			goqu.I("t.postcode"),
			goqu.I("t.datetime"),
			goqu.I("t.endblocknum"),
			goqu.I("t.isdangerous"),
			goqu.I("t.waterlevelontime"),
			goqu.I("t.deltawaterlevel"),
			goqu.I("t.waterlevelon20h"),
			goqu.I("t.watertemperature"),
			goqu.I("t.airtemperature"),
			goqu.I("t.icephenomeniastate"),
			goqu.I("t.ice"),
			goqu.I("t.snow"),
			goqu.I("t.waterflow"),
			goqu.I("t.precipitationvalue"),
			goqu.I("t.precipitationduration"),
			goqu.I("t.reservoirdate"),
			goqu.I("t.headwaterlevel"),
			goqu.I("t.averagereservoirlevel"),
			goqu.I("t.downstreamlevel"),
			goqu.I("t.reservoirvolume"),
			goqu.I("t.isreservoirwaterinflowdate"),
			goqu.I("t.inflow"),
			goqu.I("t.reset"),
			goqu.I("p.id"),
			goqu.I("p.telegramid"),
			goqu.I("p.phenomen"),
			goqu.I("p.isuntensity"),
			goqu.I("p.intensity"),
		).
		From(goqu.I("telegram").As("t")).
		LeftJoin(
			goqu.From(goqu.I("phenomenia")).As("p"),
			goqu.On(goqu.Ex{"t.id": goqu.I("p.telegramid")}),
		).
		Where(goqu.I("t.id").In(ids))

	sqlScript, args, err := selectBuilder.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(ctx, sqlScript, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	telegrams := make([]model.Telegram, 0, len(ids))

	for rows.Next() {
		var telegram model.Telegram
		var phenomeniaId *uuid.UUID
		var phenomeniaTelegramId uuid.UUID
		var phenomeniaPhenomen sql.NullByte
		var phenomeniaIsUntensity sql.NullBool
		var phenomeniaIntensity sql.NullByte

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
			&phenomeniaId,
			&phenomeniaTelegramId,
			&phenomeniaPhenomen,
			&phenomeniaIsUntensity,
			&phenomeniaIntensity,
		)
		if err != nil {
			return nil, err
		}

		if len(telegrams) == 0 || telegram.Id != telegrams[len(telegrams)-1].Id {
			telegrams = append(telegrams, telegram)
		}

		if phenomeniaId != nil {
			telegrams[len(telegrams)-1].IcePhenomenia = append(telegrams[len(telegrams)-1].IcePhenomenia, &model.Phenomenia{
				Id:          *phenomeniaId,
				TelegramId:  phenomeniaTelegramId,
				Phenomen:    phenomeniaPhenomen.Byte,
				IsUntensity: phenomeniaIsUntensity.Bool,
				Intensity:   phenomeniaIntensity,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &telegrams, nil
}
