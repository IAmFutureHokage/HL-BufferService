package model

import (
	"database/sql"
	"time"
)

type Telegram struct {
	Id                         string
	GroupId                    string
	TelegramCode               string
	PostCode                   string
	Date                       time.Time
	Time                       byte
	EndBlockNum                byte
	IsDangerous                bool
	WaterLevelOnTime           sql.NullInt32
	DeltaWaterLevel            sql.NullInt32
	WaterLevelOn20h            sql.NullInt32
	WaterTemperature           sql.NullFloat64
	AirTemperature             sql.NullInt16
	IcePhenomeniaState         byte
	IcePhenomenia              []*Phenomenia
	Ice                        sql.NullInt16
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
	Phenomen    byte
	IsUntensity bool
	Intensity   sql.NullByte
}
