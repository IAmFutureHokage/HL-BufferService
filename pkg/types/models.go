package types

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

type Telegram struct {
	PostCode PostCode
	DateAndTime
	IsDangerous            IsDangerous
	WaterLevelOnTime       *WaterLevelOnTime
	DeltaWaterLevel        *DeltaWaterLevel
	WaterLevelOn20h        *WaterLevelOn20h
	Temperature            *Temperature
	IcePhenomeniaState     *IcePhenomeniaState
	IcePhenomenia          []*Phenomenia
	IceInfo                *IceInfo
	Waterflow              *Waterflow
	Precipitation          *Precipitation
	IsReservoir            *IsReservoir
	Reservoir              *Reservoir
	IsReservoirWaterInflow *IsReservoirWaterInflow
	ReservoirWaterInflow   *ReservoirWaterInflow
}
