package encoder

import (
	"errors"
	"fmt"
	"strings"

	types "github.com/IAmFutureHokage/HL-BufferService/pkg/decoder/types"
)

func PostCodeEncoder(p *types.PostCode) (string, error) {

	if p == nil {
		return "", errors.New("PostCode is nil")
	}

	return string(*p), nil
}

func DateAndTimeEncoder(d *types.DateAndTime) (string, error) {

	if d == nil {
		return "", errors.New("DateAndTime is nil")
	}

	if d.Date > 31 {
		return "", fmt.Errorf("invalid day value: %d", d.Date)
	}

	if d.Time > 23 {
		return "", fmt.Errorf("invalid hour value: %d", d.Time)
	}

	if d.EndBlockNum > 7 {
		return "", fmt.Errorf("invalid endblock value: %d", d.EndBlockNum)
	}

	return fmt.Sprintf("%02d%02d%01d", d.Date, d.Time, d.EndBlockNum), nil
}

func IsDangerousEncoder(d *types.IsDangerous) (string, error) {

	if d == nil {
		return "", errors.New("IsDangerous is nil")
	}

	if *d {
		return "97701", nil
	}

	return "", nil
}

func WaterLevelOnTimeEncoder(w *types.WaterLevelOnTime) (string, error) {

	if w == nil {
		return "", nil
	}

	waterlevel := int(*w)

	if waterlevel == types.CouldNotMeasure {
		return "1////", nil
	}

	if waterlevel < 0 {
		waterlevel = 5000 - waterlevel
	}

	return fmt.Sprintf("1%04d", waterlevel), nil
}

func DeltaWaterLevelEncoder(d *types.DeltaWaterLevel) (string, error) {

	if d == nil {
		return "", nil
	}

	delta := int(*d)

	if delta == types.CouldNotMeasure {
		return "2////", nil
	}

	sign := '1'

	if delta < 0 {
		sign = '2'
		delta = -delta
	}

	return fmt.Sprintf("2%03d%c", delta, sign), nil
}

func WaterLevelOn20hEncoder(w *types.WaterLevelOn20h) (string, error) {

	if w == nil {
		return "", nil
	}

	waterlevel := int(*w)

	if waterlevel == types.CouldNotMeasure {
		return "3////", nil
	}

	if waterlevel < 0 {
		waterlevel = 5000 - waterlevel
	}

	return fmt.Sprintf("3%04d", waterlevel), nil
}

func TemperatureEncoder(t *types.Temperature) (string, error) {

	if t == nil {
		return "", nil
	}

	var waterTempStr, airTempStr = "//", "//"

	if *t.WaterTemperature != float64(types.CouldNotMeasure) {
		waterTemp := int(*t.WaterTemperature * 10)
		waterTempStr = fmt.Sprintf("%02d", waterTemp)
	}

	if *t.AirTemperature != types.CouldNotMeasure {
		airTemp := int(*t.AirTemperature)
		if airTemp < 0 {
			airTemp = 50 - airTemp
		}
		airTempStr = fmt.Sprintf("%02d", airTemp)
	}

	return fmt.Sprintf("4%s%s", waterTempStr, airTempStr), nil
}

func IcePhenomeniaEncoder(state *types.IcePhenomeniaState, phenomenias []*types.Phenomenia) (string, error) {

	if state == nil {
		return "", nil
	}

	if *state == types.IcePhenomeniaState(0) && len(phenomenias) == 0 {
		return "5////", nil
	}

	if *state == 1 {
		return "", nil
	}

	var encodedStrings []string

	for i := 0; i < len(phenomenias); {
		current := phenomenias[i]

		encodedString := fmt.Sprintf("5%02d", current.Phenomen)
		if current.IsUntensity {
			intensityStr := "  "
			if current.Intensity != nil {
				intensityStr = fmt.Sprintf("%02d", *current.Intensity)
			}
			encodedString += intensityStr
			i++
		} else {
			nextPhenomenon := current.Phenomen
			if i+1 < len(phenomenias) && !phenomenias[i+1].IsUntensity {
				nextPhenomenon = phenomenias[i+1].Phenomen
				i++
			}
			encodedString += fmt.Sprintf("%02d", nextPhenomenon)
			i++
		}

		encodedStrings = append(encodedStrings, encodedString)
	}

	return strings.Join(encodedStrings, " "), nil
}

func IcePhenomeniaStateEncoder(iceState *types.IcePhenomeniaState) (string, error) {

	if *iceState == 1 {
		return "60000", nil
	}

	return "", nil
}

func IceInfoEncoder(iceInfo *types.IceInfo) (string, error) {

	if iceInfo == nil {
		return "", nil
	}

	var iceHeightStr, snowHeightStr = "///", "/"

	if *iceInfo.Ice != types.CouldNotMeasure {
		iceHeightStr = fmt.Sprintf("%03d", *iceInfo.Ice)
	}

	if *iceInfo.Snow != types.SnowHeight(types.CouldNotMeasureByte) {
		snowHeightStr = fmt.Sprintf("%d", *iceInfo.Snow)
	}

	return fmt.Sprintf("7%s%s", iceHeightStr, snowHeightStr), nil
}

func WaterflowEncoder(waterflow *types.Waterflow) (string, error) {

	if waterflow == nil {
		return "", nil
	}

	flow := float64(*waterflow)

	if flow == float64(types.CouldNotMeasure) {
		return "8////", nil
	}

	var factor int

	for flow > 1 {
		flow /= 10
		factor++
	}

	scaledFlow := int(flow * 1000)

	if factor < 1 || factor > 5 {
		return "", fmt.Errorf("invalid waterflow value for encoding: %v", *waterflow)
	}

	return fmt.Sprintf("8%d%03d", factor, scaledFlow), nil
}

func PrecipitationEncoder(precip *types.Precipitation) (string, error) {

	if precip == nil {
		return "", nil
	}

	var valueStr, durationStr = "///", "/"

	if *precip.Value != float64(types.CouldNotMeasure) {
		value := float64(*precip.Value)
		if value < 1 {
			value = (value * 10) + 990
		}
		valueStr = fmt.Sprintf("%03d", int(value))
	}

	if *precip.Duration != types.PrecipitationDuration(types.CouldNotMeasureByte) {
		durationStr = fmt.Sprintf("%d", *precip.Duration)
	}

	return fmt.Sprintf("0%s%s", valueStr, durationStr), nil
}

func IsReservoirEncoder(reservoirDate *types.IsReservoirDate) (string, error) {

	if reservoirDate == nil {
		return "", nil
	}

	if *reservoirDate > 31 {
		return "", fmt.Errorf("invalid day value: %d", reservoirDate)
	}

	return fmt.Sprintf("944%02d", *reservoirDate), nil
}

func HeadwaterLevelEncoder(headwater *types.HeadwaterLevel) (string, error) {

	if headwater == nil {
		return "", nil
	}

	headwaterLevel := int(*headwater)

	if headwaterLevel == types.CouldNotMeasure {
		return "1////", nil
	}

	return fmt.Sprintf("1%04d", headwaterLevel), nil
}

func AverageReservoirLevelEncoder(averageLevel *types.AverageReservoirLevel) (string, error) {

	if averageLevel == nil {
		return "", nil
	}

	averageWaterLevel := int(*averageLevel)

	if averageWaterLevel == types.CouldNotMeasure {
		return "2////", nil
	}

	return fmt.Sprintf("2%04d", averageWaterLevel), nil
}

func DownstreamLevelEncoder(downstreamLevel *types.DownstreamLevel) (string, error) {

	if downstreamLevel == nil {
		return "", nil
	}

	waterLevel := int(*downstreamLevel)

	if waterLevel == types.CouldNotMeasure {
		return "4////", nil
	}

	return fmt.Sprintf("4%04d", waterLevel), nil
}

func ReservoirVolumeEncoder(reservoirVolume *types.ReservoirVolume) (string, error) {

	if reservoirVolume == nil {
		return "", nil
	}

	volume := float64(*reservoirVolume)
	var factor int

	if volume == float64(types.CouldNotMeasure) {
		return "7////", nil
	}

	for volume > 1 {
		volume /= 10
		factor++
	}

	scaledVolume := int(volume * 1000)

	if factor < 1 || factor > 5 {
		return "", fmt.Errorf("invalid reservoir volume value for encoding: %v", *reservoirVolume)
	}

	return fmt.Sprintf("7%d%03d", factor, uint32(scaledVolume)), nil
}

func IsReservoirWaterInflowEncoder(inflowDate *types.IsReservoirWaterInflowDate) (string, error) {

	if inflowDate == nil {
		return "", nil
	}

	if *inflowDate > 31 {
		return "", fmt.Errorf("invalid day value: %d", *inflowDate)
	}

	return fmt.Sprintf("955%02d", *inflowDate), nil
}

func InflowEncoder(inflow *types.Inflow) (string, error) {

	if inflow == nil {
		return "", nil
	}

	flow := float64(*inflow)
	var factor int

	if flow == float64(types.CouldNotMeasure) {
		return "4////", nil
	}

	for flow > 1 {
		flow /= 10
		factor++
	}

	scaledFlow := int(flow * 1000)

	if factor < 1 || factor > 5 {
		return "", fmt.Errorf("invalid inflow value for encoding: %v", *inflow)
	}

	return fmt.Sprintf("4%d%03d", factor, uint32(scaledFlow)), nil
}

func ResetEncoder(reset *types.Reset) (string, error) {

	if reset == nil {
		return "", nil
	}

	value := float64(*reset)
	var factor int

	if value == float64(types.CouldNotMeasure) {
		return "7////", nil
	}

	for value > 1 {
		value /= 10
		factor++
	}

	if factor < 1 || factor > 5 {
		return "", fmt.Errorf("invalid reset value for encoding: %v", *reset)
	}

	scaledValue := int(value * 1000)

	return fmt.Sprintf("7%d%03d", factor, uint32(scaledValue)), nil
}
