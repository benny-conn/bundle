package internal

import (
	"encoding/json"
	"os"
)

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
