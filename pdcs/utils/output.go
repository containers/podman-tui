package utils

import (
	"bytes"
	"encoding/json"
)

// GetJSONOutput converts interface to json output.
func GetJSONOutput(v interface{}) (string, error) {
	result := ""

	output, err := json.Marshal(v)
	if err != nil {
		return result, err
	}

	var out bytes.Buffer

	err = json.Indent(&out, output, "", "  ")
	if err != nil {
		return result, err
	}

	result = out.String()

	return result, nil
}
