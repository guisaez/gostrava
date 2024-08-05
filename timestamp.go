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

// MarshalJSON converts a TimeStamp into a JSON string in RFC3339 format.
func (t TimeStamp) MarshalJSON() ([]byte, error) {
	timeStr := time.Time(t)
	return json.Marshal(timeStr)
}
