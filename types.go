package gostrava

import "time"

type DateTime string

func (dt *DateTime) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(*dt))
}
