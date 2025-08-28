package request

import (
	"encoding/json"
	"io"
)

func Serialization(body io.Reader, v any) error {

	content, erro := io.ReadAll(body)

	if erro != nil {
		return nil
	}

	if erro := json.Unmarshal(content, v); erro != nil {
		return erro
	}

	return nil
}
