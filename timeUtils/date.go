package timeUtils

import (
	"database/sql/driver"
	"time"
)

const DateFormat = "2006-01-02"

type Date time.Time

func (t *Date) UnmarshalJSON(data []byte) error {
	if len(data) == 2 {
		*t = Date(time.Time{})
		return nil
	}

	now, err := time.Parse(`"`+DateFormat+`"`, string(data))
	if err != nil {
		return err
	}

	*t = Date(now)
	return nil
}

func (t Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}

func (t Date) Value() (driver.Value, error) {
	if t.String() == "0001-01-01" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(DateFormat)), nil
}

func (t *Date) Scan(v interface{}) error {
	tTime, err := time.Parse(DateFormat, v.(time.Time).String())
	if err != nil {
		return err
	}

	*t = Date(tTime)
	return nil
}

func (t Date) String() string {
	return time.Time(t).Format(DateFormat)
}

func (t Date) Time() time.Time {
	return time.Time(t)
}
