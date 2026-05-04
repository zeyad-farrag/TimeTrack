package boot

import (
	"errors"
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

func TestValidateRequiredEnvTreatsWhitespaceOnlyAsMissing(t *testing.T) {
	whitespaceCandidates := []string{
		"DEFAULT_ORG_SLUG",
		"DEFAULT_ORG_NAME",
		"TEAM_APP_SHARED_SECRET",
		"TEAM_APP_URL",
		"MULTICA_BASE_URL",
		"MULTICA_WS_URL",
		"TEAM_APP_SYSTEM_USER_PAT",
		"TEAM_APP_SYSTEM_USER_ID",
		"DATABASE_URL",
	}
	for _, name := range whitespaceCandidates {
		t.Run(name, func(t *testing.T) {
			setAllRequiredEnv(t)
			t.Setenv(name, "   \t  ")

			err := ValidateRequiredEnv()
			if err == nil {
				t.Fatalf("ValidateRequiredEnv() error = nil, want %s rejected as whitespace-only", name)
			}
			if got := MissingEnvVar(err); got != name {
				t.Fatalf("MissingEnvVar() = %q, want %q", got, name)
			}
		})
	}
}

func TestValidateRequiredEnvRejectsInvalidOrgCreationEnabled(t *testing.T) {
	invalid := []string{"yes", "no", "1", "0", "TRUEISH", "off", "on", " "}
	for _, value := range invalid {
		t.Run(value, func(t *testing.T) {
			setAllRequiredEnv(t)
			t.Setenv("ORG_CREATION_ENABLED", value)

			err := ValidateRequiredEnv()
			if err == nil {
				t.Fatalf("ValidateRequiredEnv() error = nil, want invalid ORG_CREATION_ENABLED %q", value)
			}
			if got := MissingEnvVar(err); got != "ORG_CREATION_ENABLED" {
				t.Fatalf("MissingEnvVar() = %q, want ORG_CREATION_ENABLED", got)
			}
		})
	}
}

func TestValidateRequiredEnvAcceptsOrgCreationEnabledCaseInsensitively(t *testing.T) {
	valid := []string{"true", "false", "TRUE", "FALSE", "True", "False", "tRuE", "fAlSe"}
	for _, value := range valid {
		t.Run(value, func(t *testing.T) {
			setAllRequiredEnv(t)
			t.Setenv("ORG_CREATION_ENABLED", value)

			if err := ValidateRequiredEnv(); err != nil {
				t.Fatalf("ValidateRequiredEnv() error = %v, want nil for ORG_CREATION_ENABLED=%q", err, value)
			}
		})
	}
}

func TestValidateRequiredEnvSucceedsWhenAllVarsSet(t *testing.T) {
	setAllRequiredEnv(t)
	if err := ValidateRequiredEnv(); err != nil {
		t.Fatalf("ValidateRequiredEnv() error = %v, want nil", err)
	}
}

func TestMissingEnvVarReturnsEmptyForNonEnvError(t *testing.T) {
	if got := MissingEnvVar(errors.New("unrelated")); got != "" {
		t.Fatalf("MissingEnvVar(non-EnvError) = %q, want empty string", got)
	}
	if got := MissingEnvVar(nil); got != "" {
		t.Fatalf("MissingEnvVar(nil) = %q, want empty string", got)
	}
}

func TestEnvErrorIncludesNameInMessage(t *testing.T) {
	e := EnvError{Name: "TEAM_APP_URL"}
	if !strings.Contains(e.Error(), "TEAM_APP_URL") {
		t.Fatalf("EnvError.Error() = %q, want it to mention the variable name", e.Error())
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
		"DATABASE_URL":             "postgresql://postgres:postgres@db:5432/team_app?sslmode=disable",
	}
	for name, value := range values {
		t.Setenv(name, value)
	}
}
