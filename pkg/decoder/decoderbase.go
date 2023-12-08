package decoder

import (
	"fmt"
	"regexp"
	"strconv"

	types "github.com/IAmFutureHokage/HL-Coder/pkg/types"
)

func checkCodeBlock(s string) error {

	matched, err := regexp.MatchString(`^[0-9/]{5}$`, s)
	if err != nil {
		return fmt.Errorf("error while matching regex: %v", err)
	}

	if !matched {
		return fmt.Errorf("the string must be exactly 5 characters long and consist of a digit followed by either four digits or four slashes")
	}

	return nil
}

func PostCodeDecoder(s string) (*types.PostCode, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	response := types.PostCode(s)
	return &response, nil
}

func DateAndTimeDecoder(s string) (*types.DateAndTime, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	day, err := strconv.Atoi(s[:2])
	if err != nil || day > 31 {
		return nil, fmt.Errorf("invalid day value")
	}

	hour, err := strconv.Atoi(s[2:4])
	if err != nil || hour > 23 {
		return nil, fmt.Errorf("invalid hour value")
	}

	return &types.DateAndTime{
		Date: byte(day),
		Time: byte(hour),
	}, nil
}

func IsDangerousDecoder(s string) (*types.IsDangerous, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}
	if s != "97701" {
		return nil, fmt.Errorf("invalid 977 data")
	}

	response := types.IsDangerous(true)
	return &response, nil
}

func WaterLevelOnTimeDecoder(s string) (*types.WaterLevelOnTime, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '1' {
		return nil, fmt.Errorf("first character must be '1'")
	}

	if s[1:] == "////" {
		response := types.WaterLevelOnTime(32767)
		return &response, nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid waterlavel value")
	}

	if waterlevel > 5000 && waterlevel < 6000 {
		waterlevel = 0 - waterlevel + 5000
	}

	response := types.WaterLevelOnTime(waterlevel)
	return &response, nil
}

func DeltaWaterLevelDecoder(s string) (*types.DeltaWaterLevel, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '2' {
		return nil, fmt.Errorf("first character must be '2'")
	}

	if s[1:] == "////" {
		response := types.DeltaWaterLevel(32767)
		return &response, nil
	}

	if s[4] != '1' && s[4] != '2' {
		return nil, fmt.Errorf("last character must be '1' or '2'")
	}

	delta, err := strconv.Atoi(s[1:4])
	if err != nil {
		return nil, fmt.Errorf("invalid waterlavel value")
	}

	if s[4] == '1' {
		delta = 0 + delta
	} else {
		delta = 0 - delta
	}

	response := types.DeltaWaterLevel(delta)
	return &response, nil
}

func WaterLevelOn20hDecoder(s string) (*types.WaterLevelOn20h, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '3' {
		return nil, fmt.Errorf("first character must be '3'")
	}

	if s[1:] == "////" {
		response := types.WaterLevelOn20h(32767)
		return &response, nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid waterlavel value")
	}

	if waterlevel > 5000 && waterlevel < 6000 {
		waterlevel = 0 - waterlevel + 5000
	}

	response := types.WaterLevelOn20h(waterlevel)
	return &response, nil
}

func TemperatureDecoder(s string) (*types.Temperature, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '4' {
		return nil, fmt.Errorf("first character must be '4'")
	}

	var waterTempPtr *float32
	if s[1:3] != "//" {
		waterTemp, err := strconv.Atoi(s[1:3])
		if err != nil {
			return nil, fmt.Errorf("Invalid water temperature value")
		}
		waterTempFloat := float32(waterTemp) / 10.0
		waterTempPtr = &waterTempFloat
	}

	var airTempPtr *int8
	if s[3:] != "//" {
		airTemp, err := strconv.Atoi(s[3:])
		if err != nil {
			return nil, fmt.Errorf("Invalid air temperature value")
		}
		if airTemp > 50 {
			airTemp = 0 - airTemp + 50
		}
		airTempInt8 := int8(airTemp)
		airTempPtr = &airTempInt8
	}

	return &types.Temperature{
		WaterTemperature: waterTempPtr,
		AirTemperature:   airTempPtr,
	}, nil
}

func PhenomeniaDecoder(s string) ([]*types.Phenomenia, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '5' {
		return nil, fmt.Errorf("first character must be '5'")
	}

	if s[1:] == "////" {
		return []*types.Phenomenia{nil}, nil
	}

	firstPhenomenia, err := strconv.Atoi(s[1:3])
	if err != nil {
		return nil, fmt.Errorf("Ivalid phenomenia value")
	}

	secondPhenomenia, err := strconv.Atoi(s[3:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid phenomenia value")
	}

	if firstPhenomenia == secondPhenomenia {
		return []*types.Phenomenia{
			{
				Phenomen:    byte(firstPhenomenia),
				IsUntensity: false,
				Intensity:   nil,
			},
		}, nil
	}

	if secondPhenomenia < 11 {
		secondPhenomeniaByte := byte(secondPhenomenia)
		return []*types.Phenomenia{
			{
				Phenomen:    byte(firstPhenomenia),
				IsUntensity: true,
				Intensity:   &secondPhenomeniaByte,
			},
		}, nil
	}

	return []*types.Phenomenia{
		{
			Phenomen:    byte(firstPhenomenia),
			IsUntensity: false,
		},
		{
			Phenomen:    byte(secondPhenomenia),
			IsUntensity: false,
		},
	}, nil
}

func IcePhenomeniaStateDecoder(s string) (*types.IcePhenomeniaState, error) {

	if s != "60000" {
		return nil, fmt.Errorf("Ivalid 6 group")
	}

	response := types.IcePhenomeniaState(2)
	return &response, nil
}

func IceInfoDecoder(s string) (*types.IceInfo, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '7' {
		return nil, fmt.Errorf("first character must be '7'")
	}

	var iceHeightPtr *uint16
	if s[1:4] != "///" {
		iceHeight, err := strconv.Atoi(s[1:4])
		if err != nil {
			return nil, fmt.Errorf("Invalid ice height value")
		}
		iceHeightUint := uint16(iceHeight)
		iceHeightPtr = &iceHeightUint
	}

	var snowHeightPtr *types.SnowHeight
	if s[4] != '/' {
		snowHeight, err := strconv.Atoi(s[4:])
		if err != nil {
			return nil, fmt.Errorf("Invalid snow height value")
		}
		snowHeightbyte := types.SnowHeight(byte(snowHeight))
		snowHeightPtr = &snowHeightbyte
	}

	return &types.IceInfo{
		Ice:  iceHeightPtr,
		Snow: snowHeightPtr,
	}, nil
}

func WaterflowDecoder(s string) (*types.Waterflow, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '8' {
		return nil, fmt.Errorf("first characters must be '8'")
	}

	if s[1:] == "////" {
		response := types.Waterflow(4294967295)
		return &response, nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return nil, fmt.Errorf("Ivalid factor waterflow value")
	}

	flow, err := strconv.Atoi(s[2:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid volume value")
	}

	floatFlow := float64(flow) / 1000.0
	for i := 0; i < factor; i++ {
		floatFlow *= 10
	}

	response := types.Waterflow(floatFlow)
	return &response, nil
}

func PrecipitationDecoder(s string) (*types.Precipitation, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '0' {
		return nil, fmt.Errorf("first character must be '0'")
	}

	var valuePtr *float32
	if s[1:4] != "///" {
		value, err := strconv.ParseFloat(s[1:4], 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid precipitation value")
		}

		if value >= 990 {
			value = (value - 990) / 10
		}

		valueFloat32 := float32(value)
		valuePtr = &valueFloat32
	}

	var durationPtr *types.PrecipitationDuration
	if s[4:] != "/" {
		duration, err := strconv.Atoi(s[4:])
		if err != nil || duration < 0 || duration > 4 {
			return nil, fmt.Errorf("Invalid duration value")
		}

		durationPrecip := types.PrecipitationDuration(duration)
		durationPtr = &durationPrecip
	}

	return &types.Precipitation{
		Value:    valuePtr,
		Duration: durationPtr,
	}, nil
}

func IsReservoirDecoder(s string) (*types.IsReservoir, error) {
	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[:3] != "944" {
		return nil, fmt.Errorf("Ivalid reservoir data")
	}

	date, err := strconv.Atoi(s[3:])
	if err != nil || date > 31 {
		return nil, fmt.Errorf("invalid day value")
	}

	return &types.IsReservoir{
		State: true,
		Date:  byte(date),
	}, nil
}

func HeadwaterLevelDecoder(s string) (*types.HeadwaterLevel, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '1' {
		return nil, fmt.Errorf("first character must be '1'")
	}

	if s[1:] == "////" {
		response := types.HeadwaterLevel(4294967295)
		return &response, nil
	}

	headwaterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid headwater level value")
	}

	response := types.HeadwaterLevel(headwaterlevel)
	return &response, nil
}

func AverageReservoirLevelDecoder(s string) (*types.AverageReservoirLevel, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '2' {
		return nil, fmt.Errorf("first character must be '2'")
	}

	if s[1:] == "////" {
		response := types.AverageReservoirLevel(4294967295)
		return &response, nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid avarage waterlevel value")
	}

	response := types.AverageReservoirLevel(waterlevel)
	return &response, nil
}

func DownstreamLevelDecoder(s string) (*types.DownstreamLevel, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '4' {
		return nil, fmt.Errorf("first characters must be '4'")
	}

	if s[1:] == "////" {
		response := types.DownstreamLevel(4294967295)
		return &response, nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid downstream level value")
	}

	response := types.DownstreamLevel(waterlevel)
	return &response, nil
}

func ReservoirVolumeDecoder(s string) (*types.ReservoirVolume, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '7' {
		return nil, fmt.Errorf("first characters must be '7'")
	}

	if s[1:] == "////" {
		response := types.ReservoirVolume(200000)
		return &response, nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return nil, fmt.Errorf("Ivalid factor volume value")
	}

	volume, err := strconv.Atoi(s[2:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid volume value")
	}

	floatVolume := float64(volume) / 1000.0
	for i := 0; i < factor; i++ {
		floatVolume *= 10
	}

	response := types.ReservoirVolume(floatVolume)
	return &response, nil
}

func IsReservoirWaterInflowDecoder(s string) (*types.IsReservoirWaterInflow, error) {
	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[:3] != "955" {
		return nil, fmt.Errorf("Ivalid ReservoirWaterInflow data")
	}

	date, err := strconv.Atoi(s[3:])
	if err != nil || date > 31 {
		return nil, fmt.Errorf("invalid day value")
	}

	return &types.IsReservoirWaterInflow{
		IsReservoirWaterInflow: true,
		Date:                   byte(date),
	}, nil
}

func InflowDecoder(s string) (*types.Inflow, error) {

	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '4' {
		return nil, fmt.Errorf("first characters must be '4'")
	}

	if s[1:] == "////" {
		response := types.Inflow(4294967295)
		return &response, nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return nil, fmt.Errorf("Ivalid factor inflow value")
	}

	inflow, err := strconv.Atoi(s[2:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid oflow value")
	}

	floatInflow := float64(inflow) / 1000.0
	for i := 0; i < factor; i++ {
		floatInflow *= 10
	}

	response := types.Inflow(floatInflow)
	return &response, nil
}

func ResetDecoder(s string) (*types.Reset, error) {
	err := checkCodeBlock(s)
	if err != nil {
		return nil, err
	}

	if s[0] != '7' {
		return nil, fmt.Errorf("first characters must be '7'")
	}

	if s[1:] == "////" {
		response := types.Reset(4294967295)
		return &response, nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return nil, fmt.Errorf("Ivalid factor inflow value")
	}

	reset, err := strconv.Atoi(s[2:])
	if err != nil {
		return nil, fmt.Errorf("Ivalid reset value")
	}

	floatReset := float64(reset) / 1000.0
	for i := 0; i < factor; i++ {
		floatReset *= 10
	}

	response := types.Reset(floatReset)
	return &response, nil
}
