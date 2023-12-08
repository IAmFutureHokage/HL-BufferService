package types

type PostCode string

type IsDangerous bool

type DateAndTime struct {
	Date byte
	Time byte
}

type WaterLevelOnTime int16

type DeltaWaterLevel int16

type WaterLevelOn20h int16

type Temperature struct {
	WaterTemperature *float32
	AirTemperature   *int8
}

type Phenomenia struct {
	Phenomen    byte
	IsUntensity bool
	Intensity   *byte
}

type IceInfo struct {
	Ice  *uint16
	Snow *SnowHeight
}

type Waterflow float32
type Precipitation struct {
	Value    *float32
	Duration *PrecipitationDuration
}

type IsReservoir struct {
	State bool
	Date  byte
}

type HeadwaterLevel uint32

type AverageReservoirLevel uint32

type DownstreamLevel uint32

type ReservoirVolume float32

type IsReservoirWaterInflow struct {
	IsReservoirWaterInflow bool
	Date                   byte
}

type Inflow float32

type Reset float32

type PrevDay struct {
	IsNextDay bool
	Date      byte
}
