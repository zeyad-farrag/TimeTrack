package boot

import (
	"strings"
	"testing"
)

func TestValidateRequiredEnvReportsEachMissingVariable(t *testing.T) {
	for _, missing := range requiredEnvVars {
		t.Run(missing, func(t *testing.T) {
			setAllRequiredEnv(t)
			t.Setenv(missing, "")

			err := ValidateRequiredEnv()
			if err == nil {
				t.Fatalf("ValidateRequiredEnv() error = nil, want missing %s", missing)
			}
			if !strings.Contains(err.Error(), missing) {
				t.Fatalf("ValidateRequiredEnv() error = %q, want variable name %s", err.Error(), missing)
			}
			if got := MissingEnvVar(err); got != missing {
				t.Fatalf("MissingEnvVar() = %q, want %q", got, missing)
			}
		})
	}
}

func TestValidateRequiredEnvRejectsInvalidOrgCreationEnabled(t *testing.T) {
	setAllRequiredEnv(t)
	t.Setenv("ORG_CREATION_ENABLED", "yes")

	err := ValidateRequiredEnv()
	if err == nil {
		t.Fatal("ValidateRequiredEnv() error = nil, want invalid ORG_CREATION_ENABLED")
	}
	if got := MissingEnvVar(err); got != "ORG_CREATION_ENABLED" {
		t.Fatalf("MissingEnvVar() = %q, want ORG_CREATION_ENABLED", got)
	}
}

func setAllRequiredEnv(t *testing.T) {
	t.Helper()
	values := map[string]string{
		"DEFAULT_ORG_SLUG":         "acme",
		"DEFAULT_ORG_NAME":         "Acme Inc",
		"ORG_CREATION_ENABLED":     "false",
		"TEAM_APP_SHARED_SECRET":   "0123456789abcdef0123456789abcdef",
		"TEAM_APP_URL":             "https://team.multica.uittai.com",
		"MULTICA_BASE_URL":         "https://multica.uittai.com",
		"MULTICA_WS_URL":           "wss://multica.uittai.com/ws",
		"TEAM_APP_SYSTEM_USER_PAT": "mpat_example",
		"TEAM_APP_SYSTEM_USER_ID":  "00000000-0000-0000-0000-000000000000",
	}
	for name, value := range values {
		t.Setenv(name, value)
	}
}
