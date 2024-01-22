package model

import (
	"encoding/json"
	"time"
)

type WaterLevel struct {
	Date       time.Time
	WaterLevel int32
}

func (w WaterLevel) Serialize() ([]byte, error) {
	return json.Marshal(w)
}
