package dto

import (
	"encoding/json"
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

func (d Date) MarshalJSON() ([]byte, error) {
	result := d.Time.Format(time.DateOnly)
	return json.Marshal(result)
}

func (d *Date) toString(bytes []byte) string {
	return strings.Replace(string(bytes), "\"", "", -1)
}

func ToTime(date *Date) *time.Time {
	if date == nil {
		return nil
	}
	return &date.Time
}
