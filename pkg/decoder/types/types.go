package decoder_types

import "math"

const (
	CouldNotMeasure          = math.MinInt32
	CouldNotMeasureByte byte = 100
)

type Reservoir struct {
	HeadwaterLevel        *HeadwaterLevel
	AverageReservoirLevel *AverageReservoirLevel
	DownstreamLevel       *DownstreamLevel
	ReservoirVolume       *ReservoirVolume
}

type ReservoirWaterInflow struct {
	Inflow *Inflow
	Reset  *Reset
}

type PostCode string

type IsDangerous bool

type DateAndTime struct {
	Date        byte
	Time        byte
	EndBlockNum byte
}

type WaterLevelOnTime int32

type DeltaWaterLevel int32

type WaterLevelOn20h int32

type Temperature struct {
	WaterTemperature *float64
	AirTemperature   *int32
}

type Phenomenia struct {
	Phenomen    byte
	IsUntensity bool
	Intensity   *byte
}

type IceInfo struct {
	Ice  *int32
	Snow *SnowHeight
}

type Waterflow float64
type Precipitation struct {
	Value    *float64
	Duration *PrecipitationDuration
}

type IsReservoirDate byte

type HeadwaterLevel int32

type AverageReservoirLevel int32

type DownstreamLevel int32

type ReservoirVolume float64

type IsReservoirWaterInflowDate byte

type Inflow float64

type Reset float64
