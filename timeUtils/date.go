package timeUtils

import (
	"time"
)

const timeFormat = "2006-01-02"

type Date time.Time

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
}

func NewDateFromTime(t time.Time) Date {
	return Date(CleanTime(t))
}

func (d Date) T() time.Time {
	return time.Time(d)
}

func (d Date) String() string {
	return time.Time(d).Format(timeFormat)
}

func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(d).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(`"`+timeFormat+`"`, string(b))
	if err != nil {
		return err
	}
	*d = Date(date)
	return nil
}
