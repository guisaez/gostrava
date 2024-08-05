package gostrava

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimeStamp time.Time

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	var timeStr string
	err := json.Unmarshal(data, &timeStr)
	if err != nil {
		return fmt.Errorf("timestamp unmarshal error: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return fmt.Errorf("timestamp unmarshal error: %s", err)
	}

	*t = TimeStamp(parsedTime)

	return nil
}