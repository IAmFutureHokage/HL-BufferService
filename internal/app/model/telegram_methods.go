package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
	decoder_types "github.com/IAmFutureHokage/HL-BufferService/pkg/decoder/types"

	"github.com/google/uuid"
)

func (r *Telegram) Update(draftTg *decoder.Telegram) error {

	if r == nil {
		return errors.New("nil pointer to Telegram")
	}

	if draftTg == nil {
		return fmt.Errorf("Telefram data is not valid")
	}

	r.PostCode = string(draftTg.PostCode)

	if r.DateTime.IsZero() {
		now := time.Now()
		r.DateTime = time.Date(now.Year(), now.Month(), int(draftTg.Date), int(draftTg.Time), 0, 0, 0, now.Location())
	} else {
		r.DateTime = time.Date(r.DateTime.Year(), r.DateTime.Month(), int(draftTg.Date), int(draftTg.Time), 0, 0, 0, r.DateTime.Location())
	}
	r.DateTime = r.DateTime.Truncate(time.Hour)

	r.EndBlockNum = draftTg.EndBlockNum

	r.IsDangerous = bool(draftTg.IsDangerous)

	if draftTg.WaterLevelOnTime != nil {
		r.WaterLevelOnTime = sql.NullInt32{Int32: int32(*draftTg.WaterLevelOnTime), Valid: true}
	} else {
		r.WaterLevelOnTime = sql.NullInt32{Valid: false}
	}

	if draftTg.DeltaWaterLevel != nil {
		r.DeltaWaterLevel = sql.NullInt32{Int32: int32(*draftTg.DeltaWaterLevel), Valid: true}
	} else {
		r.DeltaWaterLevel = sql.NullInt32{Valid: false}
	}

	if draftTg.WaterLevelOn20h != nil {
		r.WaterLevelOn20h = sql.NullInt32{Int32: int32(*draftTg.WaterLevelOn20h), Valid: true}
	} else {
		r.WaterLevelOn20h = sql.NullInt32{Valid: false}
	}

	if draftTg.Temperature != nil {
		r.WaterTemperature = sql.NullFloat64{Float64: *draftTg.Temperature.WaterTemperature, Valid: true}
		r.AirTemperature = sql.NullInt32{Int32: *draftTg.Temperature.AirTemperature, Valid: true}
	} else {
		r.WaterTemperature = sql.NullFloat64{Valid: false}
		r.AirTemperature = sql.NullInt32{Valid: false}
	}

	if draftTg.IcePhenomeniaState != nil {
		r.IcePhenomeniaState = sql.NullByte{Byte: draftTg.IcePhenomeniaState.ToByte(), Valid: true}
	} else {
		r.IcePhenomeniaState = sql.NullByte{Valid: false}
	}

	if len(draftTg.IcePhenomenia) == 0 {
		r.IcePhenomenia = nil

	} else {
		r.IcePhenomenia = make([]*Phenomenia, len(draftTg.IcePhenomenia))

		for i := 0; i < len(r.IcePhenomenia); i++ {

			r.IcePhenomenia[i] = &Phenomenia{}
			r.IcePhenomenia[i].Id = uuid.New()
			r.IcePhenomenia[i].TelegramId = r.Id

			err := r.IcePhenomenia[i].ToModelIcePhenomeniaConvert(draftTg.IcePhenomenia[i])
			if err != nil {
				return err
			}
		}
	}

	if draftTg.IceInfo != nil {
		r.Ice = sql.NullInt32{Int32: *draftTg.IceInfo.Ice, Valid: true}
		r.Snow = sql.NullByte{Byte: byte(*draftTg.IceInfo.Snow), Valid: true}
	} else {
		r.Ice = sql.NullInt32{Valid: false}
		r.Snow = sql.NullByte{Valid: false}
	}

	if draftTg.Waterflow != nil {
		r.Waterflow = sql.NullFloat64{Float64: float64(*draftTg.Waterflow), Valid: true}
	} else {
		r.Waterflow = sql.NullFloat64{Valid: false}
	}

	if draftTg.Precipitation != nil {
		r.PrecipitationValue = sql.NullFloat64{Float64: float64(*draftTg.Precipitation.Value), Valid: true}
		r.PrecipitationDuration = sql.NullByte{Byte: byte(*draftTg.Precipitation.Duration), Valid: true}
	} else {
		r.PrecipitationValue = sql.NullFloat64{Valid: false}
		r.PrecipitationDuration = sql.NullByte{Valid: false}
	}

	if draftTg.IsReservoirDate != nil {
		r.ReservoirDate = sql.NullTime{Time: time.Date(r.DateTime.Year(), r.DateTime.Month(), int(*draftTg.IsReservoirDate), 0, 0, 0, 0, r.DateTime.Location()), Valid: true}
	} else {
		r.ReservoirDate = sql.NullTime{Valid: false}
	}

	if draftTg.Reservoir != nil && draftTg.Reservoir.HeadwaterLevel != nil {
		r.HeadwaterLevel = sql.NullInt32{Int32: int32(*draftTg.Reservoir.HeadwaterLevel), Valid: true}
	} else {
		r.HeadwaterLevel = sql.NullInt32{Valid: false}
	}

	if draftTg.Reservoir != nil && draftTg.Reservoir.AverageReservoirLevel != nil {
		r.AverageReservoirLevel = sql.NullInt32{Int32: int32(*draftTg.Reservoir.AverageReservoirLevel), Valid: true}
	} else {
		r.AverageReservoirLevel = sql.NullInt32{Valid: false}
	}

	if draftTg.Reservoir != nil && draftTg.Reservoir.DownstreamLevel != nil {
		r.DownstreamLevel = sql.NullInt32{Int32: int32(*draftTg.Reservoir.DownstreamLevel), Valid: true}
	} else {
		r.DownstreamLevel = sql.NullInt32{Valid: false}
	}

	if draftTg.Reservoir != nil && draftTg.Reservoir.ReservoirVolume != nil {
		r.ReservoirVolume = sql.NullFloat64{Float64: float64(*draftTg.Reservoir.ReservoirVolume), Valid: true}
	} else {
		r.ReservoirVolume = sql.NullFloat64{Valid: false}
	}

	if draftTg.IsReservoirWaterInflowDate != nil {
		r.IsReservoirWaterInflowDate = sql.NullTime{Time: time.Date(r.DateTime.Year(), r.DateTime.Month(), int(*draftTg.IsReservoirWaterInflowDate), 0, 0, 0, 0, r.DateTime.Location()), Valid: true}
	} else {
		r.IsReservoirWaterInflowDate = sql.NullTime{Valid: false}
	}

	if draftTg.ReservoirWaterInflow != nil && draftTg.ReservoirWaterInflow.Inflow != nil {
		r.Inflow = sql.NullFloat64{Float64: float64(*draftTg.ReservoirWaterInflow.Inflow), Valid: true}
	} else {
		r.Inflow = sql.NullFloat64{Valid: false}
	}

	if draftTg.ReservoirWaterInflow != nil && draftTg.ReservoirWaterInflow.Reset != nil {
		r.Reset = sql.NullFloat64{Float64: float64(*draftTg.ReservoirWaterInflow.Reset), Valid: true}
	} else {
		r.Reset = sql.NullFloat64{Valid: false}
	}

	return nil
}

func (r *Phenomenia) ToModelIcePhenomeniaConvert(draftPh *decoder_types.Phenomenia) error {

	if draftPh == nil {
		return fmt.Errorf("Invalid ice phenomenia")
	}

	r.Phenomen = draftPh.Phenomen
	r.IsUntensity = draftPh.IsUntensity

	if r.IsUntensity {
		r.Intensity = sql.NullByte{Byte: *draftPh.Intensity, Valid: true}
	} else {
		r.Intensity = sql.NullByte{Valid: false}
	}

	return nil
}
