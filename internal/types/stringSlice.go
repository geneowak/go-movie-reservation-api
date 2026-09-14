package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		// other types might return it as a string
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("cannot scan %T into StringSlice", value)
		}
		bytes = []byte(str)
	}

	return json.Unmarshal(bytes, s)
}
