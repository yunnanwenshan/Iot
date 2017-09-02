package model

import (
	"time"
	"errors"
)

type OftenTime time.Time

func NewOftenTime() OftenTime {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2000-01-01 00:00:00", time.Local)
	return OftenTime(t)
}

func (self OftenTime) String() string {
	t := time.Time(self)
	if t.IsZero() {
		return "0000-00-00 00:00:00"
	}
	return t.Format("2006-01-02 15:04:05")
}

func (self OftenTime) MarshalBinary() ([]byte, error) {
	return time.Time(self).MarshalBinary()
}

func (self OftenTime) MarshalJSON() ([]byte, error) {
	t := time.Time(self)
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

func (self OftenTime) MarshalText() ([]byte, error) {
	return time.Time(self).MarshalText()
}

func (this *OftenTime) UnmarshalBinary(data []byte) error {
	t := time.Time(*this)
	return t.UnmarshalBinary(data)
}

func (this *OftenTime) UnmarshalJSON(data []byte) (err error) {
	t := time.Time(*this)
	return t.UnmarshalJSON(data)
}

func (this *OftenTime) UnmarshalText(data []byte) (err error) {
	t := time.Time(*this)
	return t.UnmarshalText(data)
}