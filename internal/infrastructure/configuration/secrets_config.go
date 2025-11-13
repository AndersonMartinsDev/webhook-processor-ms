package configuration

import (
	"fmt"
	"os"
	"strings"
)

func GetSecret(env string) string {
	filePath := os.Getenv(env)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Error ao tentar recuperar secret! %s", filePath))
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		panic("Error ao recuperar valor do secret")
	}

	return strings.TrimSpace(string(content))
}
