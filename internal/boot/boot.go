package boot

import (
	"fmt"
	"os"
	"strings"
)

var requiredEnvVars = []string{
	"DEFAULT_ORG_SLUG",
	"DEFAULT_ORG_NAME",
	"ORG_CREATION_ENABLED",
	"TEAM_APP_SHARED_SECRET",
	"TEAM_APP_URL",
	"MULTICA_BASE_URL",
	"MULTICA_WS_URL",
	"TEAM_APP_SYSTEM_USER_PAT",
	"TEAM_APP_SYSTEM_USER_ID",
}

type EnvError struct {
	Name string
}

func (e EnvError) Error() string {
	return fmt.Sprintf("missing or invalid required environment variable %s", e.Name)
}

func MissingEnvVar(err error) string {
	if envErr, ok := err.(EnvError); ok {
		return envErr.Name
	}
	return ""
}

func ValidateRequiredEnv() error {
	for _, name := range requiredEnvVars {
		value, ok := os.LookupEnv(name)
		if !ok || strings.TrimSpace(value) == "" {
			return EnvError{Name: name}
		}
		if name == "ORG_CREATION_ENABLED" && !isBool(value) {
			return EnvError{Name: name}
		}
	}
	return nil
}

func isBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "false":
		return true
	default:
		return false
	}
}
