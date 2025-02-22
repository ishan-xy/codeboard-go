package models

import (
	"encoding/json"
)

type DifficultyLevel string

func (dl *DifficultyLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*dl = DifficultyLevel(s)
	return nil
}

func (dl DifficultyLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(dl))
}