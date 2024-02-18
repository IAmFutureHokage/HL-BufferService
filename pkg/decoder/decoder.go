package decoder

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"regexp"

	types "github.com/IAmFutureHokage/HL-BufferService/pkg/decoder/types"
)

type Telegram struct {
	PostCode types.PostCode
	types.DateAndTime
	IsDangerous                types.IsDangerous
	WaterLevelOnTime           *types.WaterLevelOnTime
	DeltaWaterLevel            *types.DeltaWaterLevel
	WaterLevelOn20h            *types.WaterLevelOn20h
	Temperature                *types.Temperature
	IcePhenomeniaState         *types.IcePhenomeniaState
	IcePhenomenia              []*types.Phenomenia
	IceInfo                    *types.IceInfo
	Waterflow                  *types.Waterflow
	Precipitation              *types.Precipitation
	IsReservoirDate            *types.IsReservoirDate
	Reservoir                  *types.Reservoir
	IsReservoirWaterInflowDate *types.IsReservoirWaterInflowDate
	ReservoirWaterInflow       *types.ReservoirWaterInflow
}

func NewTelegram(s string) (*Telegram, error) {

	codeBlocks, err := parseString(s)
	if err != nil {
		return nil, err
	}
	telegram := &Telegram{}

	var isReservoir, isResevoirInflow = false, false

	for i, block := range codeBlocks {

		if i == 0 {
			err := telegram.postCodeInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if i == 1 {
			err := telegram.dateAndTimeInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if i == 2 && block[:3] == "977" {
			err := telegram.isDangerousInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '1' && !isReservoir && !isResevoirInflow {
			err := telegram.waterLevelOnTimeInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '2' && !isReservoir && !isResevoirInflow {
			err := telegram.deltaWaterLevelInit(block)
			if err != nil {
				return nil, err
			}
			continue

		}
		if block[0] == '3' && !isReservoir && !isResevoirInflow {
			err := telegram.waterLevelOn20hInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}

		if block[0] == '4' && !isReservoir && !isResevoirInflow {
			err := telegram.temperatureInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}

		if block[0] == '5' && !isReservoir && !isResevoirInflow {
			err := telegram.phenomeniaAppend(block)
			if err != nil {
				return nil, err
			}
			continue
		}

		if block[0] == '6' && !isReservoir && !isResevoirInflow {
			err := telegram.icePhenomeniaStateInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}

		if block[0] == '7' && !isReservoir && !isResevoirInflow {
			err := telegram.iceInfoInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '8' && !isReservoir && !isResevoirInflow {
			err := telegram.waterflowInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '0' && !isReservoir && !isResevoirInflow {
			err := telegram.precipitationInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[:3] == "944" && !isReservoir && !isResevoirInflow {
			err := telegram.isReservoirInit(block)
			if err != nil {
				return nil, err
			}
			isReservoir = true
			continue
		}
		if block[0] == '1' && isReservoir && !isResevoirInflow {
			err := telegram.headwaterLevelInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '2' && isReservoir && !isResevoirInflow {
			err := telegram.averageReservoirLevelInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '4' && isReservoir && !isResevoirInflow {
			err := telegram.downstreamLevelInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '7' && isReservoir && !isResevoirInflow {
			err := telegram.reservoirVolumeInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[:3] == "955" && !isResevoirInflow {
			err := telegram.reservoirWaterInflowInit(block)
			if err != nil {
				return nil, err
			}
			isResevoirInflow = true
			continue
		}
		if block[0] == '4' && isResevoirInflow {
			err := telegram.inflowInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
		if block[0] == '7' && isResevoirInflow {
			err := telegram.resetInit(block)
			if err != nil {
				return nil, err
			}
			continue
		}
	}

	return telegram, nil
}

func NewtelegramsSlice(s string) ([]*Telegram, error) {

	var telegrams = splitSequence(s)
	var decodedTelegrams []*Telegram

	for _, telegramStr := range telegrams {
		decoded, err := NewTelegram(telegramStr)
		if err != nil {
			return nil, err
		}
		decodedTelegrams = append(decodedTelegrams, decoded)
	}
	if len(decodedTelegrams) == 0 {
		return nil, fmt.Errorf("Incorrect data")
	}
	return decodedTelegrams, nil
}

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

func (t *Telegram) postCodeInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	t.PostCode = types.PostCode(s)
	return nil
}

func (t *Telegram) dateAndTimeInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	day, err := strconv.Atoi(s[:2])
	if err != nil || day > 31 {
		return fmt.Errorf("invalid day value")
	}

	hour, err := strconv.Atoi(s[2:4])
	if err != nil || hour > 23 {
		return fmt.Errorf("invalid hour value")
	}

	endBlockNum, err := strconv.Atoi(s[4:])
	if err != nil || endBlockNum < 0 || endBlockNum > 7 {
		return fmt.Errorf("invalid endBlockNum value")
	}

	t.DateAndTime = types.DateAndTime{
		Date:        byte(day),
		Time:        byte(hour),
		EndBlockNum: byte(endBlockNum),
	}

	return nil
}

func (t *Telegram) isDangerousInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}
	if s != "97701" {
		return fmt.Errorf("invalid 977 data")
	}

	t.IsDangerous = types.IsDangerous(true)
	return nil
}

func (t *Telegram) waterLevelOnTimeInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '1' {
		return fmt.Errorf("first character must be '1'")
	}

	if s[1:] == "////" {
		waterlevelOnT := types.WaterLevelOnTime(types.CouldNotMeasure)
		t.WaterLevelOnTime = &waterlevelOnT
		return nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return fmt.Errorf("invalid waterlavel value")
	}

	if waterlevel > 5000 && waterlevel < 6000 {
		waterlevel = 0 - waterlevel + 5000
	}

	waterlevelOnT := types.WaterLevelOnTime(waterlevel)
	t.WaterLevelOnTime = &waterlevelOnT

	return nil
}

func (t *Telegram) deltaWaterLevelInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '2' {
		return fmt.Errorf("first character must be '2'")
	}

	if s[1:] == "////" {
		deltaWl := types.DeltaWaterLevel(types.CouldNotMeasure)
		t.DeltaWaterLevel = &deltaWl
		return nil
	}

	if s[4] != '1' && s[4] != '2' {
		return fmt.Errorf("last character must be '1' or '2'")
	}

	delta, err := strconv.Atoi(s[1:4])
	if err != nil {
		return fmt.Errorf("invalid waterlavel value")
	}

	if s[4] == '1' {
		delta = 0 + delta
	} else {
		delta = 0 - delta
	}

	deltaWL := types.DeltaWaterLevel(delta)
	t.DeltaWaterLevel = &deltaWL

	return nil
}

func (t *Telegram) waterLevelOn20hInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '3' {
		return fmt.Errorf("first character must be '3'")
	}

	if s[1:] == "////" {
		waterlevel20H := types.WaterLevelOn20h(types.CouldNotMeasure)
		t.WaterLevelOn20h = &waterlevel20H
		return nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return fmt.Errorf("invalid waterlavel value")
	}

	if waterlevel > 5000 && waterlevel < 6000 {
		waterlevel = 0 - waterlevel + 5000
	}

	waterlevel20H := types.WaterLevelOn20h(waterlevel)
	t.WaterLevelOn20h = &waterlevel20H

	return nil
}

func (t *Telegram) temperatureInit(s string) error {

	err := checkCodeBlock(s)

	if err != nil {
		return err
	}

	if s[0] != '4' {
		return fmt.Errorf("first character must be '4'")
	}

	var waterTempPtr *float64

	if s[1:3] == "//" {
		waterTemp := float64(types.CouldNotMeasure)
		waterTempPtr = &waterTemp
	} else {
		waterTemp, err := strconv.Atoi(s[1:3])
		if err != nil {
			return fmt.Errorf("Invalid water temperature value")
		}
		waterTempFloat := float64(waterTemp) / 10.0
		waterTempPtr = &waterTempFloat
	}

	var airTempPtr *int32

	if s[3:] == "//" {
		airTemp := int32(types.CouldNotMeasure)
		airTempPtr = &airTemp
	} else {
		airTemp, err := strconv.Atoi(s[3:])
		if err != nil {
			return fmt.Errorf("Invalid air temperature value")
		}
		if airTemp > 50 {
			airTemp = 0 - airTemp + 50
		}
		airTempInt := int32(airTemp)
		airTempPtr = &airTempInt
	}

	t.Temperature = &types.Temperature{
		WaterTemperature: waterTempPtr,
		AirTemperature:   airTempPtr,
	}

	return nil
}

func (t *Telegram) phenomeniaAppend(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '5' {
		return fmt.Errorf("first character must be '5'")
	}

	state := types.IcePhenomeniaState(0)

	if s[1:] == "////" {
		t.IcePhenomeniaState = &state
		return nil
	}

	firstPhenomenia, err := strconv.Atoi(s[1:3])
	if err != nil {
		return fmt.Errorf("Ivalid phenomenia value")
	}

	secondPhenomenia, err := strconv.Atoi(s[3:])
	if err != nil {
		return fmt.Errorf("Ivalid phenomenia value")
	}

	t.IcePhenomeniaState = &state

	if firstPhenomenia == secondPhenomenia {

		phenomen := types.Phenomenia{
			Phenomen:    byte(firstPhenomenia),
			IsUntensity: false,
			Intensity:   nil,
		}

		t.IcePhenomenia = append(t.IcePhenomenia, &phenomen)
		return nil
	}

	if secondPhenomenia < 11 {

		secondPhenomeniaByte := byte(secondPhenomenia)

		phenomen := types.Phenomenia{
			Phenomen:    byte(firstPhenomenia),
			IsUntensity: true,
			Intensity:   &secondPhenomeniaByte,
		}

		t.IcePhenomenia = append(t.IcePhenomenia, &phenomen)
		return nil
	}

	phenomens := []*types.Phenomenia{
		{
			Phenomen:    byte(firstPhenomenia),
			IsUntensity: false,
		},
		{
			Phenomen:    byte(secondPhenomenia),
			IsUntensity: false,
		},
	}

	t.IcePhenomenia = append(t.IcePhenomenia, phenomens...)

	return nil
}

func (t *Telegram) icePhenomeniaStateInit(s string) error {

	if s != "60000" {
		return fmt.Errorf("Ivalid 6 group")
	}

	icePhenomeniaState := types.IcePhenomeniaState(1)
	t.IcePhenomeniaState = &icePhenomeniaState
	return nil
}

func (t *Telegram) iceInfoInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '7' {
		return fmt.Errorf("first character must be '7'")
	}

	var iceHeightPtr *int32
	if s[1:4] != "///" {
		iceHeight, err := strconv.Atoi(s[1:4])
		if err != nil {
			return fmt.Errorf("Invalid ice height value")
		}
		iceHeightUint := int32(iceHeight)
		iceHeightPtr = &iceHeightUint
	} else {
		iceHeight := int32(types.CouldNotMeasure)
		iceHeightPtr = &iceHeight
	}

	var snowHeightPtr *types.SnowHeight
	if s[4] != '/' {
		snowHeight, err := strconv.Atoi(s[4:])
		if err != nil {
			return fmt.Errorf("Invalid snow height value")
		}
		snowHeightbyte := types.SnowHeight(byte(snowHeight))
		snowHeightPtr = &snowHeightbyte
	} else {
		snowHeight := types.SnowHeight(types.CouldNotMeasureByte)
		snowHeightPtr = &snowHeight
	}

	t.IceInfo = &types.IceInfo{
		Ice:  iceHeightPtr,
		Snow: snowHeightPtr,
	}

	return nil
}

func (t *Telegram) waterflowInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '8' {
		return fmt.Errorf("first characters must be '8'")
	}

	if s[1:] == "////" {
		waterflow := types.Waterflow(types.CouldNotMeasure)
		t.Waterflow = &waterflow
		return nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return fmt.Errorf("Ivalid factor waterflow value")
	}

	flow, err := strconv.Atoi(s[2:])
	if err != nil {
		return fmt.Errorf("Ivalid volume value")
	}

	floatFlow := float64(flow)
	for i := 0; i < factor; i++ {
		floatFlow *= 10
	}
	floatFlow = math.Round(floatFlow) / 1000.0

	waterflow := types.Waterflow(floatFlow)
	t.Waterflow = &waterflow

	return nil
}

func (t *Telegram) precipitationInit(s string) error {

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '0' {
		return fmt.Errorf("first character must be '0'")
	}

	var valuePtr *float64
	if s[1:4] != "///" {
		value, err := strconv.ParseFloat(s[1:4], 32)
		if err != nil {
			return fmt.Errorf("Invalid precipitation value")
		}

		if value >= 990 {
			value = math.Round((value - 990)) / 10
		}

		valueFloat := float64(value)
		valuePtr = &valueFloat
	} else {
		valueFloat := float64(types.CouldNotMeasure)
		valuePtr = &valueFloat
	}

	var durationPtr *types.PrecipitationDuration

	if s[4:] != "/" {
		duration, err := strconv.Atoi(s[4:])
		if err != nil || duration < 0 || duration > 4 {
			return fmt.Errorf("Invalid duration value")
		}

		durationPrecip := types.PrecipitationDuration(duration)
		durationPtr = &durationPrecip
	} else {
		durationPrecip := types.PrecipitationDuration(types.CouldNotMeasureByte)
		durationPtr = &durationPrecip
	}

	t.Precipitation = &types.Precipitation{
		Value:    valuePtr,
		Duration: durationPtr,
	}

	return nil
}

func (t *Telegram) isReservoirInit(s string) error {
	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[:3] != "944" {
		return fmt.Errorf("Ivalid reservoir data")
	}

	date, err := strconv.Atoi(s[3:])
	if err != nil || date > 31 {
		return fmt.Errorf("invalid day value")
	}

	reservoirDate := types.IsReservoirDate(date)
	t.IsReservoirDate = &reservoirDate
	t.Reservoir = &types.Reservoir{}

	return nil
}

func (t *Telegram) headwaterLevelInit(s string) error {

	if t.Reservoir == nil {
		return fmt.Errorf("Reservoir data not init")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '1' {
		return fmt.Errorf("first character must be '1'")
	}

	if s[1:] == "////" {
		headwaterLevel := types.HeadwaterLevel(types.CouldNotMeasure)
		t.Reservoir.HeadwaterLevel = &headwaterLevel

		return nil
	}

	headwaterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return fmt.Errorf("Ivalid headwater level value")
	}

	headwaterLevel := types.HeadwaterLevel(headwaterlevel)
	t.Reservoir.HeadwaterLevel = &headwaterLevel

	return nil
}

func (t *Telegram) averageReservoirLevelInit(s string) error {

	if t.Reservoir == nil {
		return fmt.Errorf("Reservoir data not init")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '2' {
		return fmt.Errorf("first character must be '2'")
	}

	if s[1:] == "////" {
		averageReservoirLevel := types.AverageReservoirLevel(types.CouldNotMeasure)
		t.Reservoir.AverageReservoirLevel = &averageReservoirLevel

		return nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return fmt.Errorf("Ivalid avarage waterlevel value")
	}

	averageReservoirLevel := types.AverageReservoirLevel(waterlevel)
	t.Reservoir.AverageReservoirLevel = &averageReservoirLevel

	return nil
}

func (t *Telegram) downstreamLevelInit(s string) error {

	if t.Reservoir == nil {
		return fmt.Errorf("Reservoir data not init")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '4' {
		return fmt.Errorf("first characters must be '4'")
	}

	if s[1:] == "////" {
		downstreamLevel := types.DownstreamLevel(types.CouldNotMeasure)
		t.Reservoir.DownstreamLevel = &downstreamLevel

		return nil
	}

	waterlevel, err := strconv.Atoi(s[1:])
	if err != nil {
		return fmt.Errorf("Ivalid downstream level value")
	}

	downstreamLevel := types.DownstreamLevel(waterlevel)
	t.Reservoir.DownstreamLevel = &downstreamLevel

	return nil
}

func (t *Telegram) reservoirVolumeInit(s string) error {

	if t.Reservoir == nil {
		return fmt.Errorf("Reservoir data not init")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '7' {
		return fmt.Errorf("first characters must be '7'")
	}

	if s[1:] == "////" {
		reservoirVolume := types.ReservoirVolume(types.CouldNotMeasure)
		t.Reservoir.ReservoirVolume = &reservoirVolume

		return nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return fmt.Errorf("Ivalid factor volume value")
	}

	volume, err := strconv.Atoi(s[2:])
	if err != nil {
		return fmt.Errorf("Ivalid volume value")
	}

	floatVolume := float64(volume)
	for i := 0; i < factor; i++ {
		floatVolume *= 10
	}
	floatVolume = math.Round(floatVolume) / 1000

	reservoirVolume := types.ReservoirVolume(floatVolume)
	t.Reservoir.ReservoirVolume = &reservoirVolume

	return nil
}

func (t *Telegram) reservoirWaterInflowInit(s string) error {
	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[:3] != "955" {
		return fmt.Errorf("Ivalid ReservoirWaterInflow data")
	}

	date, err := strconv.Atoi(s[3:])
	if err != nil || date > 31 {
		return fmt.Errorf("invalid day value")
	}

	reservoirWaterInflow := types.IsReservoirWaterInflowDate(date)
	t.IsReservoirWaterInflowDate = &reservoirWaterInflow
	t.ReservoirWaterInflow = &types.ReservoirWaterInflow{}

	return nil
}

func (t *Telegram) inflowInit(s string) error {

	if t.ReservoirWaterInflow == nil {
		return fmt.Errorf("Reservoir Water Inflow data not init")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '4' {
		return fmt.Errorf("first characters must be '4'")
	}

	if s[1:] == "////" {
		inflowResp := types.Inflow(types.CouldNotMeasure)
		t.ReservoirWaterInflow.Inflow = &inflowResp
		return nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return fmt.Errorf("Ivalid factor inflow value")
	}

	inflow, err := strconv.Atoi(s[2:])
	if err != nil {
		return fmt.Errorf("Ivalid oflow value")
	}

	floatInflow := float64(inflow)
	for i := 0; i < factor; i++ {
		floatInflow *= 10
	}

	floatInflow = math.Round(floatInflow) / 1000

	inflowResp := types.Inflow(floatInflow)
	t.ReservoirWaterInflow.Inflow = &inflowResp

	return nil
}

func (t *Telegram) resetInit(s string) error {

	if t.ReservoirWaterInflow == nil {
		return fmt.Errorf("955-блок не инициализирован")
	}

	err := checkCodeBlock(s)
	if err != nil {
		return err
	}

	if s[0] != '7' {
		return fmt.Errorf("Ошибка в блоке %s, блок должен начинаться с 7", s)
	}

	if s[1:] == "////" {
		resetResp := types.Reset(types.CouldNotMeasure)
		t.ReservoirWaterInflow.Reset = &resetResp

		return nil
	}

	factor, err := strconv.Atoi(s[1:2])
	if err != nil || factor < 1 || factor > 5 {
		return fmt.Errorf("Ошибка в блоке %s, Фактор должен быть от 1 до 5", s)
	}

	reset, err := strconv.Atoi(s[2:])
	if err != nil {
		return fmt.Errorf("Ivalid reset value")
	}

	floatReset := float64(reset)
	for i := 0; i < factor; i++ {
		floatReset *= 10
	}
	floatReset = math.Round(floatReset) / 1000

	resetResp := types.Reset(floatReset)
	t.ReservoirWaterInflow.Reset = &resetResp

	return nil
}

func parseString(input string) ([]string, error) {

	input = strings.Split(input, "=")[0]

	substrings := strings.Fields(input)

	for _, s := range substrings {
		if len(s) != 5 {
			return nil, fmt.Errorf("Неверный формат: подстрока '%s' должна иметь 5 символов", s)
		}
	}

	return substrings, nil
}

func splitSequence(s string) []string {

	blocks := strings.Fields(s)

	if len(blocks) < 2 {
		return nil
	}

	var sequences []string
	firstBlock := blocks[0]
	endBluckNum := blocks[1][4:]

	currentSequence := []string{firstBlock}

	for _, block := range blocks[1:] {
		if strings.HasPrefix(block, "922") {
			if len(currentSequence) > 1 {
				sequences = append(sequences, strings.Join(currentSequence, " "))
			}
			modifiedSecondBlock := block[3:5] + "08" + endBluckNum
			currentSequence = []string{firstBlock, modifiedSecondBlock}
		} else {
			currentSequence = append(currentSequence, block)
		}
	}

	sequences = append(sequences, strings.Join(currentSequence, " "))

	return sequences
}
