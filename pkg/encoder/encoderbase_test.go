package encoder

import (
	"testing"

	"github.com/IAmFutureHokage/HL-Coder/pkg/types"
)

func TestPostCodeEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.PostCode
		want    string
		wantErr bool
	}{
		{
			name:    "Valid PostCode",
			input:   func() *types.PostCode { p := types.PostCode("12345"); return &p }(),
			want:    "12345",
			wantErr: false,
		},
		{
			name:    "Nil PostCode",
			input:   nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostCodeEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostCodeEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostCodeEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateAndTimeEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.DateAndTime
		want    string
		wantErr bool
	}{
		{
			name:    "Valid DateAndTime",
			input:   &types.DateAndTime{Date: 15, Time: 12},
			want:    "15121",
			wantErr: false,
		},
		{
			name:    "Nil DateAndTime",
			input:   nil,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid Date",
			input:   &types.DateAndTime{Date: 32, Time: 12},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid Time",
			input:   &types.DateAndTime{Date: 15, Time: 24},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DateAndTimeEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateAndTimeEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DateAndTimeEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDangerousEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.IsDangerous
		want    string
		wantErr bool
	}{
		{
			name:    "IsDangerous true",
			input:   func() *types.IsDangerous { b := types.IsDangerous(true); return &b }(),
			want:    "97701",
			wantErr: false,
		},
		{
			name:    "IsDangerous false",
			input:   func() *types.IsDangerous { b := types.IsDangerous(false); return &b }(),
			want:    "",
			wantErr: false,
		},
		{
			name:    "Nil IsDangerous",
			input:   nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDangerousEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDangerousEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDangerousEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWaterLevelOnTimeEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.WaterLevelOnTime
		want    string
		wantErr bool
	}{
		{
			name:    "Valid WaterLevelOnTime",
			input:   func() *types.WaterLevelOnTime { w := types.WaterLevelOnTime(1234); return &w }(),
			want:    "11234",
			wantErr: false,
		},
		{
			name:    "Special case WaterLevelOnTime",
			input:   func() *types.WaterLevelOnTime { w := types.WaterLevelOnTime(32767); return &w }(),
			want:    "1////",
			wantErr: false,
		},
		{
			name:    "Negative WaterLevelOnTime",
			input:   func() *types.WaterLevelOnTime { w := types.WaterLevelOnTime(-766); return &w }(),
			want:    "15766",
			wantErr: false,
		},
		{
			name:    "Nil WaterLevelOnTime",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterLevelOnTimeEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterLevelOnTimeEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WaterLevelOnTimeEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeltaWaterLevelEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.DeltaWaterLevel
		want    string
		wantErr bool
	}{
		{
			name:    "Positive DeltaWaterLevel",
			input:   func() *types.DeltaWaterLevel { d := types.DeltaWaterLevel(123); return &d }(),
			want:    "21231",
			wantErr: false,
		},
		{
			name:    "Negative DeltaWaterLevel",
			input:   func() *types.DeltaWaterLevel { d := types.DeltaWaterLevel(-23); return &d }(),
			want:    "20232",
			wantErr: false,
		},
		{
			name:    "Special case DeltaWaterLevel",
			input:   func() *types.DeltaWaterLevel { d := types.DeltaWaterLevel(32767); return &d }(),
			want:    "2////",
			wantErr: false,
		},
		{
			name:    "Nil DeltaWaterLevel",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeltaWaterLevelEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeltaWaterLevelEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeltaWaterLevelEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWaterLevelOn20hEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.WaterLevelOn20h
		want    string
		wantErr bool
	}{
		{
			name:    "Valid WaterLevelOn20h",
			input:   func() *types.WaterLevelOn20h { w := types.WaterLevelOn20h(1234); return &w }(),
			want:    "31234",
			wantErr: false,
		},
		{
			name:    "Special case WaterLevelOn20h",
			input:   func() *types.WaterLevelOn20h { w := types.WaterLevelOn20h(32767); return &w }(),
			want:    "3////",
			wantErr: false,
		},
		{
			name:    "Negative WaterLevelOn20h",
			input:   func() *types.WaterLevelOn20h { w := types.WaterLevelOn20h(-766); return &w }(),
			want:    "35766",
			wantErr: false,
		},
		{
			name:    "Nil WaterLevelOn20h",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterLevelOn20hEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterLevelOn20hEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WaterLevelOn20hEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemperatureEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.Temperature
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Temperature - both temperatures present",
			input:   &types.Temperature{WaterTemperature: func() *float32 { f := float32(1.5); return &f }(), AirTemperature: func() *int8 { i := int8(20); return &i }()},
			want:    "41520",
			wantErr: false,
		},
		{
			name:    "Valid Temperature - only water temperature",
			input:   &types.Temperature{WaterTemperature: func() *float32 { f := float32(1.5); return &f }(), AirTemperature: nil},
			want:    "415//",
			wantErr: false,
		},
		{
			name:    "Valid Temperature - only air temperature",
			input:   &types.Temperature{WaterTemperature: nil, AirTemperature: func() *int8 { i := int8(-5); return &i }()},
			want:    "4//55",
			wantErr: false,
		},
		{
			name:    "Valid Temperature - nill temperature",
			input:   &types.Temperature{WaterTemperature: nil, AirTemperature: nil},
			want:    "4////",
			wantErr: false,
		},
		{
			name:    "Nil Temperature",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemperatureEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemperatureEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TemperatureEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcePhenomeniaEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   []*types.Phenomenia
		want    string
		wantErr bool
	}{
		{
			name: "Valid IcePhenomenia - single phenomenia",
			input: []*types.Phenomenia{
				{Phenomen: 12, IsUntensity: false},
			},
			want:    "51212",
			wantErr: false,
		},
		{
			name: "Valid IcePhenomenia - multiple phenomenias",
			input: []*types.Phenomenia{
				{Phenomen: 12, IsUntensity: false},
				{Phenomen: 34, IsUntensity: false},
			},
			want:    "51234",
			wantErr: false,
		},
		{
			name: "Valid IcePhenomenia - with intensity",
			input: []*types.Phenomenia{
				{Phenomen: 12, IsUntensity: true, Intensity: func() *byte { b := byte(10); return &b }()},
			},
			want:    "51210",
			wantErr: false,
		},
		{
			name:    "First element nil",
			input:   []*types.Phenomenia{nil},
			want:    "5////",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IcePhenomeniaEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IcePhenomeniaEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IcePhenomeniaEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcePhenomeniaStateEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.IcePhenomeniaState
		want    string
		wantErr bool
	}{
		{
			name:    "Valid IcePhenomeniaState - 2",
			input:   func() *types.IcePhenomeniaState { s := types.IcePhenomeniaState(2); return &s }(),
			want:    "60000",
			wantErr: false,
		},
		{
			name:    "Valid IcePhenomeniaState - not 2",
			input:   func() *types.IcePhenomeniaState { s := types.IcePhenomeniaState(1); return &s }(),
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IcePhenomeniaStateEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IcePhenomeniaStateEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IcePhenomeniaStateEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIceInfoEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.IceInfo
		want    string
		wantErr bool
	}{
		{
			name:    "Valid IceInfo - both ice and snow height present",
			input:   &types.IceInfo{Ice: func() *uint16 { h := uint16(123); return &h }(), Snow: func() *types.SnowHeight { sh := types.SnowHeight(4); return &sh }()},
			want:    "71234",
			wantErr: false,
		},
		{
			name:    "Valid IceInfo - only ice height",
			input:   &types.IceInfo{Ice: func() *uint16 { h := uint16(123); return &h }(), Snow: nil},
			want:    "7123/",
			wantErr: false,
		},
		{
			name:    "Valid IceInfo - only snow height",
			input:   &types.IceInfo{Ice: nil, Snow: func() *types.SnowHeight { sh := types.SnowHeight(4); return &sh }()},
			want:    "7///4",
			wantErr: false,
		},
		{
			name:    "Valid IceInfo - only snow height",
			input:   &types.IceInfo{Ice: nil, Snow: nil},
			want:    "7////",
			wantErr: false,
		},
		{
			name:    "Nil IceInfo",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IceInfoEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IceInfoEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IceInfoEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWaterflowEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.Waterflow
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Waterflow",
			input:   func() *types.Waterflow { wf := types.Waterflow(1230); return &wf }(),
			want:    "84123",
			wantErr: false,
		},
		{
			name:    "Special case Waterflow - maximum value",
			input:   func() *types.Waterflow { wf := types.Waterflow(100001.0); return &wf }(),
			want:    "8////",
			wantErr: false,
		},
		{
			name:    "Invalid Waterflow - too small",
			input:   func() *types.Waterflow { wf := types.Waterflow(0.0001); return &wf }(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil Waterflow",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterflowEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaterflowEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WaterflowEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrecipitationEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.Precipitation
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Precipitation - both value and duration",
			input:   &types.Precipitation{Value: func() *float32 { v := float32(0.4); return &v }(), Duration: func() *types.PrecipitationDuration { d := types.PrecipitationDuration(4); return &d }()},
			want:    "09944",
			wantErr: false,
		},
		{
			name:    "Valid Precipitation - only value",
			input:   &types.Precipitation{Value: func() *float32 { v := float32(0.4); return &v }(), Duration: nil},
			want:    "0994/",
			wantErr: false,
		},
		{
			name:    "Valid Precipitation - only duration",
			input:   &types.Precipitation{Value: nil, Duration: func() *types.PrecipitationDuration { d := types.PrecipitationDuration(4); return &d }()},
			want:    "0///4",
			wantErr: false,
		},
		{
			name:    "Nil Precipitation",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrecipitationEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrecipitationEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrecipitationEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsReservoirEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.IsReservoir
		want    string
		wantErr bool
	}{
		{
			name:    "Valid IsReservoir",
			input:   &types.IsReservoir{Date: 15},
			want:    "94415",
			wantErr: false,
		},
		{
			name:    "Invalid IsReservoir - invalid day",
			input:   &types.IsReservoir{Date: 32},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil IsReservoir",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsReservoirEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReservoirEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsReservoirEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeadwaterLevelEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.HeadwaterLevel
		want    string
		wantErr bool
	}{
		{
			name:    "Valid HeadwaterLevel",
			input:   func() *types.HeadwaterLevel { h := types.HeadwaterLevel(1234); return &h }(),
			want:    "11234",
			wantErr: false,
		},
		{
			name:    "Special case HeadwaterLevel - maximum value",
			input:   func() *types.HeadwaterLevel { h := types.HeadwaterLevel(4294967295); return &h }(),
			want:    "1////",
			wantErr: false,
		},
		{
			name:    "Nil HeadwaterLevel",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HeadwaterLevelEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HeadwaterLevelEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HeadwaterLevelEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverageReservoirLevelEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.AverageReservoirLevel
		want    string
		wantErr bool
	}{
		{
			name:    "Valid AverageReservoirLevel",
			input:   func() *types.AverageReservoirLevel { a := types.AverageReservoirLevel(1234); return &a }(),
			want:    "21234",
			wantErr: false,
		},
		{
			name:    "Special case AverageReservoirLevel - maximum value",
			input:   func() *types.AverageReservoirLevel { a := types.AverageReservoirLevel(4294967295); return &a }(),
			want:    "2////",
			wantErr: false,
		},
		{
			name:    "Nil AverageReservoirLevel",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AverageReservoirLevelEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AverageReservoirLevelEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AverageReservoirLevelEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownstreamLevelEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.DownstreamLevel
		want    string
		wantErr bool
	}{
		{
			name:    "Valid DownstreamLevel",
			input:   func() *types.DownstreamLevel { d := types.DownstreamLevel(1234); return &d }(),
			want:    "41234",
			wantErr: false,
		},
		{
			name:    "Special case DownstreamLevel - maximum value",
			input:   func() *types.DownstreamLevel { d := types.DownstreamLevel(4294967295); return &d }(),
			want:    "4////",
			wantErr: false,
		},
		{
			name:    "Nil DownstreamLevel",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownstreamLevelEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownstreamLevelEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DownstreamLevelEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReservoirVolumeEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.ReservoirVolume
		want    string
		wantErr bool
	}{
		{
			name:    "Valid ReservoirVolume",
			input:   func() *types.ReservoirVolume { rv := types.ReservoirVolume(1230); return &rv }(),
			want:    "74123",
			wantErr: false,
		},
		{
			name:    "Special case ReservoirVolume - maximum value",
			input:   func() *types.ReservoirVolume { rv := types.ReservoirVolume(100001.0); return &rv }(),
			want:    "7////",
			wantErr: false,
		},
		{
			name:    "Invalid ReservoirVolume - too small",
			input:   func() *types.ReservoirVolume { rv := types.ReservoirVolume(0.0001); return &rv }(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil ReservoirVolume",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReservoirVolumeEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReservoirVolumeEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReservoirVolumeEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsReservoirWaterInflowEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.IsReservoirWaterInflow
		want    string
		wantErr bool
	}{
		{
			name:    "Valid IsReservoirWaterInflow",
			input:   &types.IsReservoirWaterInflow{Date: 15},
			want:    "95515",
			wantErr: false,
		},
		{
			name:    "Invalid IsReservoirWaterInflow - invalid day",
			input:   &types.IsReservoirWaterInflow{Date: 32},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil IsReservoirWaterInflow",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsReservoirWaterInflowEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReservoirWaterInflowEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsReservoirWaterInflowEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInflowEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.Inflow
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Inflow",
			input:   func() *types.Inflow { i := types.Inflow(23.5); return &i }(),
			want:    "42235",
			wantErr: false,
		},
		{
			name:    "Special case Inflow - maximum value",
			input:   func() *types.Inflow { i := types.Inflow(4294967295); return &i }(),
			want:    "4////",
			wantErr: false,
		},
		{
			name:    "Invalid Inflow - too small",
			input:   func() *types.Inflow { i := types.Inflow(0.0001); return &i }(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil Inflow",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InflowEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("InflowEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InflowEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResetEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   *types.Reset
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Reset",
			input:   func() *types.Reset { r := types.Reset(234.0); return &r }(),
			want:    "73234",
			wantErr: false,
		},
		{
			name:    "Special case Reset - maximum value",
			input:   func() *types.Reset { r := types.Reset(4294967295); return &r }(),
			want:    "7////",
			wantErr: false,
		},
		{
			name:    "Invalid Reset - too small",
			input:   func() *types.Reset { r := types.Reset(0.0001); return &r }(),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Nil Reset",
			input:   nil,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResetEncoder(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResetEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}
