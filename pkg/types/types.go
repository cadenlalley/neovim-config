package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IntSlice []int

// Scan implements the sql.Scanner interface for the IntSlice type
func (is *IntSlice) Scan(src interface{}) error {
	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case nil:
		*is = IntSlice{}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	// Handle empty array case
	if len(bytes) == 0 || string(bytes) == "[]" {
		*is = IntSlice{}
		return nil
	}

	// Unmarshal the JSON array into the slice
	return json.Unmarshal(bytes, &is)
}

// Value implements the driver.Valuer interface for the IntSlice type
func (is IntSlice) Value() (driver.Value, error) {
	if is == nil || len(is) == 0 {
		return "[]", nil
	}
	return json.Marshal(is)
}
