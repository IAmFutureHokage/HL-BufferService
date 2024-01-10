package toModelTelegramConverter

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	types "github.com/IAmFutureHokage/HL-BufferService/pkg/types"
	"github.com/google/uuid"
)

func ToModelTelegramConvert(mainTg *model.Telegram, draftTg *types.Telegram) error {

	if mainTg == nil || draftTg == nil {
		return fmt.Errorf("Telefram data is not valid")
	}

	mainTg.PostCode = string(draftTg.PostCode)

	if mainTg.DateTime.IsZero() {
		now := time.Now()
		mainTg.DateTime = time.Date(now.Year(), now.Month(), int(draftTg.Date), int(draftTg.Time), 0, 0, 0, now.Location())
	} else {
		mainTg.DateTime = time.Date(mainTg.DateTime.Year(), mainTg.DateTime.Month(), int(draftTg.Date), int(draftTg.Time), 0, 0, 0, mainTg.DateTime.Location())
	}
	mainTg.DateTime = mainTg.DateTime.Truncate(time.Hour)

	mainTg.EndBlockNum = draftTg.EndBlockNum

	mainTg.IsDangerous = bool(draftTg.IsDangerous)

	if draftTg.WaterLevelOnTime != nil {
		mainTg.WaterLevelOnTime = sql.NullInt32{Int32: int32(*draftTg.WaterLevelOnTime), Valid: true}
	} else {
		mainTg.WaterLevelOnTime = sql.NullInt32{Valid: false}
	}

	if draftTg.DeltaWaterLevel != nil {
		mainTg.DeltaWaterLevel = sql.NullInt32{Int32: int32(*draftTg.DeltaWaterLevel), Valid: true}
	} else {
		mainTg.DeltaWaterLevel = sql.NullInt32{Valid: false}
	}

	if draftTg.WaterLevelOn20h != nil {
		mainTg.WaterLevelOn20h = sql.NullInt32{Int32: int32(*draftTg.WaterLevelOn20h), Valid: true}
	} else {
		mainTg.WaterLevelOn20h = sql.NullInt32{Valid: false}
	}

	if draftTg.Temperature != nil {
		mainTg.WaterTemperature = sql.NullFloat64{Float64: *draftTg.Temperature.WaterTemperature, Valid: true}
		mainTg.AirTemperature = sql.NullInt32{Int32: *draftTg.Temperature.AirTemperature, Valid: true}
	} else {
		mainTg.WaterTemperature = sql.NullFloat64{Valid: false}
		mainTg.AirTemperature = sql.NullInt32{Valid: false}
	}

	if draftTg.IcePhenomeniaState != nil {
		mainTg.IcePhenomeniaState = sql.NullByte{Byte: draftTg.IcePhenomeniaState.ToByte(), Valid: true}
	} else {
		mainTg.IcePhenomeniaState = sql.NullByte{Valid: false}
	}

	if len(draftTg.IcePhenomenia) == 0 {
		mainTg.IcePhenomenia = nil

	} else {
		mainTg.IcePhenomenia = make([]*model.Phenomenia, len(draftTg.IcePhenomenia))

		for i := 0; i < len(mainTg.IcePhenomenia); i++ {

			mainTg.IcePhenomenia[i].Id = uuid.New()
			mainTg.IcePhenomenia[i].TelegramId = mainTg.Id

			err := ToModelIcePhenomeniaConvert(mainTg.IcePhenomenia[i], draftTg.IcePhenomenia[i])
			if err != nil {
				return err
			}
		}
	}

	if draftTg.IceInfo != nil {
		mainTg.Ice = sql.NullInt32{Int32: *draftTg.IceInfo.Ice, Valid: true}
		mainTg.Snow = sql.NullByte{Byte: byte(*draftTg.IceInfo.Snow), Valid: true}
	} else {
		mainTg.Ice = sql.NullInt32{Valid: false}
		mainTg.Snow = sql.NullByte{Valid: false}
	}

	if draftTg.Waterflow != nil {
		mainTg.Waterflow = sql.NullFloat64{Float64: float64(*draftTg.Waterflow), Valid: true}
	} else {
		mainTg.Waterflow = sql.NullFloat64{Valid: false}
	}

	if draftTg.Precipitation != nil {
		mainTg.PrecipitationValue = sql.NullFloat64{Float64: float64(*draftTg.Precipitation.Value), Valid: true}
		mainTg.PrecipitationDuration = sql.NullByte{Byte: byte(*draftTg.Precipitation.Duration), Valid: true}
	} else {
		mainTg.PrecipitationValue = sql.NullFloat64{Valid: false}
		mainTg.PrecipitationDuration = sql.NullByte{Valid: false}
	}

	if draftTg.IsReservoirDate != nil {
		mainTg.ReservoirDate = sql.NullTime{Time: time.Date(mainTg.DateTime.Year(), mainTg.DateTime.Month(), int(*draftTg.IsReservoirDate), 0, 0, 0, 0, mainTg.DateTime.Location()), Valid: true}
	} else {
		mainTg.ReservoirDate = sql.NullTime{Valid: false}
	}
	return nil
}

func ToModelIcePhenomeniaConvert(mainPh *model.Phenomenia, draftPh *types.Phenomenia) error {

	if mainPh == nil || draftPh == nil {
		return fmt.Errorf("Invalid ice phenomenia")
	}

	mainPh.Phenomen = draftPh.Phenomen
	mainPh.IsUntensity = draftPh.IsUntensity

	if mainPh.IsUntensity {
		mainPh.Intensity = sql.NullByte{Byte: *draftPh.Intensity, Valid: true}
	} else {
		mainPh.Intensity = sql.NullByte{Valid: false}
	}

	return nil
}
