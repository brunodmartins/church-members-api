package dto

import (
	"strings"
	"time"
)

// Date shall be used to unmarshall a date only field using time.DateOnly
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	result, err := time.Parse(time.DateOnly, d.toString(bytes))
	if err != nil {
		return err
	}
	d.Time = result
	return nil
}

func (d *Date) toString(bytes []byte) string {
	return strings.Replace(string(bytes), "\"", "", -1)
}
