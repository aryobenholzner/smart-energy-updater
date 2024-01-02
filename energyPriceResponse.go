package main

import "time"

type Response struct {
	Tariff   string         `json:"tariff"`
	Unit     string         `json:"unit"`
	Interval int            `json:"interval"`
	Data     []ResponseData `json:"data"`
}

type ResponseData struct {
	Date  DateTime `json:"date"`
	Value float64  `json:"value"`
}

type DateTime struct {
	time.Time
}

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05-07:00"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
