package timeUtils

import (
	"database/sql/driver"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime{Time: time.Time{}}
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime{Time: now}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(t.Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, err := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	if err != nil {
		return err
	}
	*t = LocalTime{Time: tTime}
	return nil
}

func (t LocalTime) String() string {
	return t.Format(TimeFormat)
}
