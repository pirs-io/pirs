package structs

import (
	"encoding/json"
)

// ToMap returns map[string]interface{} representation of given struct
func ToMap(in interface{}) map[string]interface{} {
	var result map[string]interface{}

	inrec, err := json.Marshal(in)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(inrec, &result)
	if err != nil {
		return nil
	}

	return result
}
