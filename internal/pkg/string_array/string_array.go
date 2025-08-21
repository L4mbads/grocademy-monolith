package string_array

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lib/pq"
)

type StringArray pq.StringArray

func (s StringArray) Value() (driver.Value, error) {
	return pq.StringArray(s).Value()
}

func (s *StringArray) Scan(value interface{}) error {
	return (*pq.StringArray)(s).Scan(value)
}

func (s StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

func (s *StringArray) UnmarshalJSON(data []byte) error {
	var tmp []string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*s = tmp
	return nil
}
