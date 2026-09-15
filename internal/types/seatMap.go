package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type SeatMap struct {
	Rows []struct {
		Row   string `json:"row"`
		Seats []struct {
			Number int    `json:"number"`
			Type   string `json:"type"`
		} `json:"seats"`
	} `json:"rows"`
	TotalSeats int    `json:"total_seats"`
	Screen     string `json:"screen"`
}

func (s SeatMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SeatMap) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("cannot scan %T into SeatMap", value)
		}
		bytes = []byte(str)
	}

	return json.Unmarshal(bytes, s)
}
