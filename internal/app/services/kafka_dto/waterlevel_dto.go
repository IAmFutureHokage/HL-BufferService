package kafka_dto

import (
	"time"

	"github.com/mailru/easyjson"
)

// easyjson:json
type WaterLevel struct {
	PostCode   string    `json:"post_code"`
	Date       time.Time `json:"date"`
	WaterLevel int32     `json:"water_level"`
}

// easyjson:json
type WaterLevelRecords struct {
	Waterlevels []WaterLevel `json:"waterlevels"`
}

func NewWaterLevelRecords(capacity int) *WaterLevelRecords {
	return &WaterLevelRecords{
		Waterlevels: make([]WaterLevel, 0, capacity),
	}
}

func (w WaterLevelRecords) Serialize() ([]byte, error) {
	return easyjson.Marshal(w)
}
