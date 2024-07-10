package gostrava

import (
	"encoding/json"
	"time"
)

type TimeStamp struct {
	Time time.Time
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	var timeStr string
	err := json.Unmarshal(data, &timeStr)
	if err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}

	t.Time = parsedTime

	return nil
}
