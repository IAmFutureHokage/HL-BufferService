package kafka_dto

import (
	"encoding/json"
	"time"
)

type WaterLevel struct {
	PostCode   string    `json:"post_code"`
	Date       time.Time `json:"date"`
	WaterLevel int32     `json:"water_level"`
}

type WaterLevelRecords struct {
	Waterlevels []WaterLevel `json:"waterlevels"`
}

func NewWaterLevelRecords(capacity int) *WaterLevelRecords {
	return &WaterLevelRecords{
		Waterlevels: make([]WaterLevel, 0, capacity),
	}
}

func (w WaterLevelRecords) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(w.Waterlevels)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
