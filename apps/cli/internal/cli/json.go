package cli

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

func readJSONFile(path string) (any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func readOptionalJSONFile(path string) (any, error) {
	if strings.TrimSpace(path) == "" {
		return map[string]any{}, nil
	}

	return readJSONFile(path)
}
