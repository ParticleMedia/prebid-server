package secretutil

import (
	"fmt"
	"os"
	"strings"
)

const (
	secretPath = "/etc/secrets"
)

func GetSecret(key string) string {
	bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", secretPath, key))
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(bytes), "\n")
}
