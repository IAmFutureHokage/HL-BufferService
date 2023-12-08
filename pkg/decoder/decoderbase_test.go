package decoder

import (
	"reflect"
	"testing"

	types "github.com/IAmFutureHokage/HL-Coder/pkg/types"
)

func TestCheckCodeBlock(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid code",
			input:   "12345",
			wantErr: false,
		},
		{
			name:    "Valid code with slashes",
			input:   "1////",
			wantErr: false,
		},
		{
			name:    "Invalid code - too short",
			input:   "1234",
			wantErr: true,
		},
		{
			name:    "Invalid code - too long",
			input:   "123456",
			wantErr: true,
		},
		{
			name:    "Invalid code - wrong characters",
			input:   "12a34",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCodeBlock(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkCodeBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostCodeDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.PostCode
		wantErr bool
	}{
		{
			name:    "Valid PostCode",
			input:   "12345",
			want:    func() *types.PostCode { pc := types.PostCode("12345"); return &pc }(),
			wantErr: false,
		},
		{
			name:    "Invalid PostCode - wrong format",
			input:   "abcde",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostCodeDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostCodeDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostCodeDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateAndTimeDecoder(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *types.DateAndTime
		wantErr bool
	}{
		{
			name:  "Valid Date and Time",
			input: "23101",
			want: &types.DateAndTime{
				Date: 23,
				Time: 10,
			},
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "2a101",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid day",
			input:   "32101",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid hour",
			input:   "23241",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DateAndTimeDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateAndTimeDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && (*got != *tt.want) {
				t.Errorf("DateAndTimeDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestIsDangerousDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.IsDangerous
		wantErr bool
	}{
		{
			name:    "Valid IsDangerous",
			input:   "97701",
			want:    func() *types.IsDangerous { b := types.IsDangerous(true); return &b }(),
			wantErr: false,
		},
		{
			name:    "Invalid IsDangerous - wrong format",
			input:   "97700",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid IsDangerous - invalid code",
			input:   "12345",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDangerousDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDangerousDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("IsDangerousDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestWaterLevelOnTimeDecoder(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *types.WaterLevelOnTime
		wantErr bool
	}{
		{
			name:    "Valid Water Level",
			input:   "15500",
			want:    func() *types.WaterLevelOnTime { wl := types.WaterLevelOnTime(-500); return &wl }(),
			wantErr: false,
		},
		{
			name:    "Valid Water Level 2",
			input:   "10150",
			want:    func() *types.WaterLevelOnTime { wl := types.WaterLevelOnTime(150); return &wl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Water Level - wrong format",
			input:   "2////",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid Water Level - NaN value",
			input:   "1////",
			want:    func() *types.WaterLevelOnTime { wl := types.WaterLevelOnTime(32767); return &wl }(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterLevelOnTimeDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterLevelOnTimeDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("WaterLevelOnTimeDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestDeltaWaterLevelDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.DeltaWaterLevel
		wantErr bool
	}{
		{
			name:    "Valid Delta Water Level - Positive",
			input:   "20101",
			want:    func() *types.DeltaWaterLevel { dl := types.DeltaWaterLevel(10); return &dl }(),
			wantErr: false,
		},
		{
			name:    "Valid Delta Water Level - Negative",
			input:   "20102",
			want:    func() *types.DeltaWaterLevel { dl := types.DeltaWaterLevel(-10); return &dl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Delta Water Level - wrong format",
			input:   "3////",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid Delta Water Level - NaN value",
			input:   "2////",
			want:    func() *types.DeltaWaterLevel { dl := types.DeltaWaterLevel(32767); return &dl }(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeltaWaterLevelDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeltaWaterLevelDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("DeltaWaterLevelDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestWaterLevelOn20hDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.WaterLevelOn20h
		wantErr bool
	}{
		{
			name:    "Valid Water Level",
			input:   "35000",
			want:    func() *types.WaterLevelOn20h { wl := types.WaterLevelOn20h(5000); return &wl }(),
			wantErr: false,
		},
		{
			name:    "Valid Water Level - negative",
			input:   "35500",
			want:    func() *types.WaterLevelOn20h { wl := types.WaterLevelOn20h(-500); return &wl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Water Level - wrong format",
			input:   "4////",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid Water Level - NaN value",
			input:   "3////",
			want:    func() *types.WaterLevelOn20h { wl := types.WaterLevelOn20h(32767); return &wl }(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterLevelOn20hDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterLevelOn20hDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("WaterLevelOn20hDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestTemperatureDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.Temperature
		wantErr bool
	}{
		{
			name:  "Valid Temperature - both temperatures present",
			input: "41520",
			want: &types.Temperature{
				WaterTemperature: func() *float32 { f := float32(1.5); return &f }(),
				AirTemperature:   func() *int8 { i := int8(20); return &i }(),
			},
			wantErr: false,
		},
		{
			name:  "Valid Temperature - only water temperature",
			input: "415//",
			want: &types.Temperature{
				WaterTemperature: func() *float32 { f := float32(1.5); return &f }(),
				AirTemperature:   nil,
			},
			wantErr: false,
		},
		{
			name:  "Valid Temperature - only air temperature",
			input: "4//02",
			want: &types.Temperature{
				WaterTemperature: nil,
				AirTemperature:   func() *int8 { i := int8(2); return &i }(),
			},
			wantErr: false,
		},
		{
			name:    "Invalid Temperature - wrong format",
			input:   "5abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Temperature - non-numeric water temperature",
			input:   "4ab//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Temperature - non-numeric air temperature",
			input:   "4//ab",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemperatureDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemperatureDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.WaterTemperature != nil && tt.want.WaterTemperature != nil && *got.WaterTemperature != *tt.want.WaterTemperature {
					t.Errorf("TemperatureDecoder() WaterTemperature = %v, want %v", *got.WaterTemperature, *tt.want.WaterTemperature)
				}
				if got.AirTemperature != nil && tt.want.AirTemperature != nil && *got.AirTemperature != *tt.want.AirTemperature {
					t.Errorf("TemperatureDecoder() AirTemperature = %v, want %v", *got.AirTemperature, *tt.want.AirTemperature)
				}
			}
		})
	}
}

func TestPhenomeniaDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*types.Phenomenia
		wantErr bool
	}{
		{
			name:  "Valid Phenomenia - single value",
			input: "51234",
			want: []*types.Phenomenia{
				{
					Phenomen:    12,
					IsUntensity: false,
				},
				{
					Phenomen:    34,
					IsUntensity: false,
				},
			},
			wantErr: false,
		},
		{
			name:  "Valid Phenomenia - same values",
			input: "51212",
			want: []*types.Phenomenia{
				{
					Phenomen:    12,
					IsUntensity: false,
				},
			},
			wantErr: false,
		},
		{
			name:  "Valid Phenomenia - intensity",
			input: "51210",
			want: []*types.Phenomenia{
				{
					Phenomen:    12,
					IsUntensity: true,
					Intensity:   func() *byte { b := byte(10); return &b }(),
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid Phenomenia - wrong format",
			input:   "6abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Phenomenia - non-numeric value",
			input:   "5ab//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Valid Phenomenia - no data",
			input:   "5////",
			want:    []*types.Phenomenia{nil},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PhenomeniaDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PhenomeniaDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("PhenomeniaDecoder() got %v items, want %v items", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != nil && tt.want[i] != nil {
					if got[i].Phenomen != tt.want[i].Phenomen || got[i].IsUntensity != tt.want[i].IsUntensity {
						t.Errorf("PhenomeniaDecoder() item %d = %v, want %v", i, got[i], tt.want[i])
					}
					if got[i].Intensity != nil && tt.want[i].Intensity != nil && *got[i].Intensity != *tt.want[i].Intensity {
						t.Errorf("PhenomeniaDecoder() item %d Intensity = %v, want %v", i, *got[i].Intensity, *tt.want[i].Intensity)
					}
				}
			}
		})
	}
}

func TestIcePhenomeniaStateDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.IcePhenomeniaState
		wantErr bool
	}{
		{
			name:    "Valid Ice Phenomenia State",
			input:   "60000",
			want:    func() *types.IcePhenomeniaState { s := types.IcePhenomeniaState(2); return &s }(),
			wantErr: false,
		},
		{
			name:    "Invalid Ice Phenomenia State - wrong value",
			input:   "60001",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IcePhenomeniaStateDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IcePhenomeniaStateDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("IcePhenomeniaStateDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestIceInfoDecoder(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *types.IceInfo
		wantErr bool
	}{
		{
			name:  "Valid Ice Info - both ice and snow height present",
			input: "70102",
			want: &types.IceInfo{
				Ice:  func() *uint16 { h := uint16(10); return &h }(),
				Snow: func() *types.SnowHeight { sh := types.SnowHeight(2); return &sh }(),
			},
			wantErr: false,
		},
		{
			name:  "Valid Ice Info - only ice height",
			input: "7010/",
			want: &types.IceInfo{
				Ice:  func() *uint16 { h := uint16(10); return &h }(),
				Snow: nil,
			},
			wantErr: false,
		},
		{
			name:  "Valid Ice Info - only snow height",
			input: "7///2",
			want: &types.IceInfo{
				Ice:  nil,
				Snow: func() *types.SnowHeight { sh := types.SnowHeight(2); return &sh }(),
			},
			wantErr: false,
		},
		{
			name:    "Invalid Ice Info - wrong format",
			input:   "8abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Ice Info - non-numeric ice height",
			input:   "7ab//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Ice Info - non-numeric snow height",
			input:   "7///a",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IceInfoDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IceInfoDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.Ice != nil && tt.want.Ice != nil && *got.Ice != *tt.want.Ice {
					t.Errorf("IceInfoDecoder() Ice = %v, want %v", *got.Ice, *tt.want.Ice)
				}
				if got.Snow != nil && tt.want.Snow != nil && *got.Snow != *tt.want.Snow {
					t.Errorf("IceInfoDecoder() Snow = %v, want %v", *got.Snow, *tt.want.Snow)
				}
			}
		})
	}
}

func TestWaterflowDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.Waterflow
		wantErr bool
	}{
		{
			name:    "Valid Waterflow",
			input:   "83234",
			want:    func() *types.Waterflow { wf := types.Waterflow(234); return &wf }(),
			wantErr: false,
		},
		{
			name:    "Valid Waterflow - NaN value",
			input:   "8////",
			want:    func() *types.Waterflow { wf := types.Waterflow(4294967295); return &wf }(),
			wantErr: false,
		},
		{
			name:    "Invalid Waterflow - wrong format",
			input:   "9abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Waterflow - non-numeric factor",
			input:   "8a234",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Waterflow - non-numeric flow",
			input:   "812ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Waterflow - invalid factor",
			input:   "80678",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterflowDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterflowDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("WaterflowDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestPrecipitationDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.Precipitation
		wantErr bool
	}{
		{
			name:  "Valid Precipitation - value and duration",
			input: "09994",
			want: &types.Precipitation{
				Value:    func() *float32 { v := float32(0.9); return &v }(),
				Duration: func() *types.PrecipitationDuration { d := types.PrecipitationDuration(4); return &d }(),
			},
			wantErr: false,
		},
		{
			name:  "Valid Precipitation - only value",
			input: "0999/",
			want: &types.Precipitation{
				Value:    func() *float32 { v := float32(0.9); return &v }(),
				Duration: nil,
			},
			wantErr: false,
		},
		{
			name:  "Valid Precipitation - only duration",
			input: "0///4",
			want: &types.Precipitation{
				Value:    nil,
				Duration: func() *types.PrecipitationDuration { d := types.PrecipitationDuration(4); return &d }(),
			},
			wantErr: false,
		},
		{
			name:  "Valid Precipitation - nil data",
			input: "0////",
			want: &types.Precipitation{
				Value:    nil,
				Duration: nil,
			},
			wantErr: false,
		},
		{
			name:    "Invalid Precipitation - wrong format",
			input:   "1abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Precipitation - non-numeric value",
			input:   "0ab//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Precipitation - non-numeric duration",
			input:   "0////a",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Precipitation - invalid duration",
			input:   "0////5",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrecipitationDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrecipitationDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.Value != nil && tt.want.Value != nil && *got.Value != *tt.want.Value {
					t.Errorf("PrecipitationDecoder() Value = %v, want %v", *got.Value, *tt.want.Value)
				}
				if got.Duration != nil && tt.want.Duration != nil && *got.Duration != *tt.want.Duration {
					t.Errorf("PrecipitationDecoder() Duration = %v, want %v", *got.Duration, *tt.want.Duration)
				}
			}
		})
	}
}

func TestIsReservoirDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.IsReservoir
		wantErr bool
	}{
		{
			name:  "Valid IsReservoir",
			input: "94415",
			want: &types.IsReservoir{
				State: true,
				Date:  15,
			},
			wantErr: false,
		},
		{
			name:    "Invalid IsReservoir - wrong prefix",
			input:   "94515",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid IsReservoir - non-numeric date",
			input:   "944ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid IsReservoir - invalid date",
			input:   "94432",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsReservoirDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReservoirDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && (*got != *tt.want) {
				t.Errorf("IsReservoirDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestHeadwaterLevelDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.HeadwaterLevel
		wantErr bool
	}{
		{
			name:    "Valid Headwater Level",
			input:   "11234",
			want:    func() *types.HeadwaterLevel { hl := types.HeadwaterLevel(1234); return &hl }(),
			wantErr: false,
		},
		{
			name:    "Valid Headwater Level - NaN value",
			input:   "1////",
			want:    func() *types.HeadwaterLevel { hl := types.HeadwaterLevel(4294967295); return &hl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Headwater Level - wrong format",
			input:   "2abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Headwater Level - non-numeric value",
			input:   "1abcd",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HeadwaterLevelDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HeadwaterLevelDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("HeadwaterLevelDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestAverageReservoirLevelDecoder(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *types.AverageReservoirLevel
		wantErr bool
	}{
		{
			name:    "Valid Average Reservoir Level",
			input:   "21234",
			want:    func() *types.AverageReservoirLevel { lvl := types.AverageReservoirLevel(1234); return &lvl }(),
			wantErr: false,
		},
		{
			name:    "Valid Average Reservoir Level - NaN value",
			input:   "2////",
			want:    func() *types.AverageReservoirLevel { lvl := types.AverageReservoirLevel(4294967295); return &lvl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Average Reservoir Level - wrong format",
			input:   "3abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Average Reservoir Level - non-numeric value",
			input:   "2abcd",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AverageReservoirLevelDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AverageReservoirLevelDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("AverageReservoirLevelDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestDownstreamLevelDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.DownstreamLevel
		wantErr bool
	}{
		{
			name:    "Valid Downstream Level",
			input:   "41234",
			want:    func() *types.DownstreamLevel { lvl := types.DownstreamLevel(1234); return &lvl }(),
			wantErr: false,
		},
		{
			name:    "Valid Downstream Level - NaN value",
			input:   "4////",
			want:    func() *types.DownstreamLevel { lvl := types.DownstreamLevel(4294967295); return &lvl }(),
			wantErr: false,
		},
		{
			name:    "Invalid Downstream Level - wrong format",
			input:   "5abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Downstream Level - non-numeric value",
			input:   "4abcd",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownstreamLevelDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownstreamLevelDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("DownstreamLevelDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestReservoirVolumeDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.ReservoirVolume
		wantErr bool
	}{
		{
			name:    "Valid Reservoir Volume",
			input:   "75234",
			want:    func() *types.ReservoirVolume { vol := types.ReservoirVolume(23400); return &vol }(),
			wantErr: false,
		},
		{
			name:    "Valid Reservoir Volume - NaN value",
			input:   "7////",
			want:    func() *types.ReservoirVolume { vol := types.ReservoirVolume(200000); return &vol }(),
			wantErr: false,
		},
		{
			name:    "Invalid Reservoir Volume - wrong format",
			input:   "8abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reservoir Volume - non-numeric factor",
			input:   "7a234",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reservoir Volume - non-numeric volume",
			input:   "712ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reservoir Volume - invalid factor",
			input:   "70678",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReservoirVolumeDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReservoirVolumeDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("ReservoirVolumeDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestIsReservoirWaterInflowDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.IsReservoirWaterInflow
		wantErr bool
	}{
		{
			name:  "Valid IsReservoirWaterInflow",
			input: "95515",
			want: &types.IsReservoirWaterInflow{
				IsReservoirWaterInflow: true,
				Date:                   15,
			},
			wantErr: false,
		},
		{
			name:    "Invalid IsReservoirWaterInflow - wrong prefix",
			input:   "95615",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid IsReservoirWaterInflow - non-numeric date",
			input:   "955ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid IsReservoirWaterInflow - invalid date",
			input:   "95532",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsReservoirWaterInflowDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReservoirWaterInflowDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && (*got != *tt.want) {
				t.Errorf("IsReservoirWaterInflowDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestInflowDecoder(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		want    *types.Inflow
		wantErr bool
	}{
		{
			name:    "Valid Inflow",
			input:   "42234",
			want:    func() *types.Inflow { inf := types.Inflow(23.4); return &inf }(),
			wantErr: false,
		},
		{
			name:    "Valid Inflow - NaN value",
			input:   "4////",
			want:    func() *types.Inflow { inf := types.Inflow(4294967295); return &inf }(),
			wantErr: false,
		},
		{
			name:    "Invalid Inflow - wrong format",
			input:   "5abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Inflow - non-numeric factor",
			input:   "4a234",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Inflow - non-numeric inflow",
			input:   "412ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Inflow - invalid factor",
			input:   "40678",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InflowDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("InflowDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("InflowDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestResetDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *types.Reset
		wantErr bool
	}{
		{
			name:    "Valid Reset",
			input:   "74234",
			want:    func() *types.Reset { r := types.Reset(2340); return &r }(),
			wantErr: false,
		},
		{
			name:    "Valid Reset - NaN value",
			input:   "7////",
			want:    func() *types.Reset { r := types.Reset(4294967295); return &r }(),
			wantErr: false,
		},
		{
			name:    "Invalid Reset - wrong format",
			input:   "8abcd",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reset - non-numeric factor",
			input:   "7a234",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reset - non-numeric reset value",
			input:   "712ab",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid Reset - invalid factor",
			input:   "70678",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResetDecoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("ResetDecoder() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
