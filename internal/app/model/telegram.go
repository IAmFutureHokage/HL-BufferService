package model

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

type Telegram struct {
	Id                         uuid.UUID
	GroupId                    uuid.UUID
	TelegramCode               string
	PostCode                   string
	DateTime                   time.Time
	EndBlockNum                byte
	IsDangerous                bool
	WaterLevelOnTime           sql.NullInt32
	DeltaWaterLevel            sql.NullInt32
	WaterLevelOn20h            sql.NullInt32
	WaterTemperature           sql.NullFloat64
	AirTemperature             sql.NullInt32
	IcePhenomeniaState         sql.NullByte
	IcePhenomenia              []*Phenomenia
	Ice                        sql.NullInt32
	Snow                       sql.NullByte
	Waterflow                  sql.NullFloat64
	PrecipitationValue         sql.NullFloat64
	PrecipitationDuration      sql.NullByte
	ReservoirDate              sql.NullTime
	HeadwaterLevel             sql.NullInt32
	AverageReservoirLevel      sql.NullInt32
	DownstreamLevel            sql.NullInt32
	ReservoirVolume            sql.NullFloat64
	IsReservoirWaterInflowDate sql.NullTime
	Inflow                     sql.NullFloat64
	Reset                      sql.NullFloat64
}

type Phenomenia struct {
	Id          uuid.UUID
	TelegramId  uuid.UUID
	Phenomen    byte
	IsUntensity bool
	Intensity   sql.NullByte
}
