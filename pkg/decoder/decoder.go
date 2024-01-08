package decoder

import (
	"strings"

	types "github.com/IAmFutureHokage/HL-BufferService/pkg/types"
)

func FullDecoder(s string) ([]*types.Telegram, error) {

	var telegrams = splitSequence(s)
	var decodedTelegrams []*types.Telegram

	for _, telegramStr := range telegrams {
		decoded, err := Decoder(telegramStr)
		if err != nil {
			return nil, err
		}
		decodedTelegrams = append(decodedTelegrams, decoded)
	}

	return decodedTelegrams, nil
}

func Decoder(s string) (*types.Telegram, error) {

	codeBlocks := parseString(s)
	telegram := &types.Telegram{}

	var isReservoir, isResevoirInflow = false, false

	for i, block := range codeBlocks {

		if i == 0 {
			postCode, err := PostCodeDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.PostCode = *postCode
			continue
		}
		if i == 1 {
			dateAndTime, err := DateAndTimeDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.DateAndTime = *dateAndTime
			continue
		}
		if i == 2 && block[:3] == "977" {
			isDangerous, err := IsDangerousDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.IsDangerous = *isDangerous
			continue
		}
		if block[0] == '1' && !isReservoir && !isResevoirInflow {
			waterLevel, err := WaterLevelOnTimeDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.WaterLevelOnTime = waterLevel
			continue
		}
		if block[0] == '2' && !isReservoir && !isResevoirInflow {
			delta, err := DeltaWaterLevelDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.DeltaWaterLevel = delta
			continue

		}
		if block[0] == '3' && !isReservoir && !isResevoirInflow {
			waterLevel, err := WaterLevelOn20hDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.WaterLevelOn20h = waterLevel
			continue
		}

		if block[0] == '4' && !isReservoir && !isResevoirInflow {
			temperature, err := TemperatureDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Temperature = temperature
			continue
		}

		if block[0] == '5' && !isReservoir && !isResevoirInflow {
			state, phenomenia, err := PhenomeniaDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.IcePhenomeniaState = state
			telegram.IcePhenomenia = append(telegram.IcePhenomenia, phenomenia...)
			continue

		}

		if block[0] == '6' && !isReservoir && !isResevoirInflow {
			state, err := IcePhenomeniaStateDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.IcePhenomeniaState = state
			continue
		}

		if block[0] == '7' && !isReservoir && !isResevoirInflow {
			info, err := IceInfoDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.IceInfo = info
			continue
		}
		if block[0] == '8' && !isReservoir && !isResevoirInflow {
			waterflow, err := WaterflowDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Waterflow = waterflow
			continue
		}
		if block[0] == '0' && !isReservoir && !isResevoirInflow {
			prec, err := PrecipitationDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Precipitation = prec
			continue
		}
		if block[:3] == "944" && !isReservoir && !isResevoirInflow {
			date, err := IsReservoirDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Reservoir = &types.Reservoir{}
			isReservoir = true
			telegram.IsReservoirDate = date
			continue
		}
		if block[0] == '1' && isReservoir && !isResevoirInflow {
			data, err := HeadwaterLevelDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Reservoir.HeadwaterLevel = data
			continue
		}
		if block[0] == '2' && isReservoir && !isResevoirInflow {
			data, err := AverageReservoirLevelDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Reservoir.AverageReservoirLevel = data
			continue
		}
		if block[0] == '4' && isReservoir && !isResevoirInflow {
			data, err := DownstreamLevelDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Reservoir.DownstreamLevel = data
			continue
		}
		if block[0] == '7' && isReservoir && !isResevoirInflow {
			data, err := ReservoirVolumeDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.Reservoir.ReservoirVolume = data
			continue
		}
		if block[:3] == "955" && !isResevoirInflow {
			date, err := IsReservoirWaterInflowDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.ReservoirWaterInflow = &types.ReservoirWaterInflow{}
			isResevoirInflow = true
			telegram.IsReservoirWaterInflowDate = date
			continue
		}
		if block[0] == '4' && isResevoirInflow {
			data, err := InflowDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.ReservoirWaterInflow.Inflow = data
			continue
		}
		if block[0] == '7' && isResevoirInflow {
			data, err := ResetDecoder(block)
			if err != nil {
				return nil, err
			}
			telegram.ReservoirWaterInflow.Reset = data
			continue
		}
	}

	return telegram, nil
}

func parseString(input string) []string {

	substrings := strings.Fields(input)

	var result []string
	for _, s := range substrings {
		s = strings.TrimSuffix(s, "=")
		result = append(result, s)
	}

	return result
}

func splitSequence(s string) []string {

	blocks := strings.Fields(s)

	if len(blocks) < 2 {
		return nil
	}

	var sequences []string
	firstBlock := blocks[0]

	currentSequence := []string{firstBlock}

	for _, block := range blocks[1:] {
		if strings.HasPrefix(block, "922") {
			if len(currentSequence) > 1 {
				sequences = append(sequences, strings.Join(currentSequence, " "))
			}
			modifiedSecondBlock := block[3:5] + "081"
			currentSequence = []string{firstBlock, modifiedSecondBlock}
		} else {
			currentSequence = append(currentSequence, block)
		}
	}

	sequences = append(sequences, strings.Join(currentSequence, " "))

	return sequences
}
