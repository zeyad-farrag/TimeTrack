package boot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repo root not found walking up from %s", file)
		}
		dir = parent
	}
}

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

// T-001 / AC1 — every Architecture-pinned backend library is a direct require in go.mod.
// AC1 lists two sub-package paths (`.../sqlc/cmd/sqlc`, `.../client_golang/prometheus`); the
// module-level paths are what go.mod records.
func TestGoModDirectRequiresIncludeAllArchitecturePinnedModules(t *testing.T) {
	required := []string{
		"github.com/go-chi/chi/v5",
		"github.com/jackc/pgx/v5",
		"github.com/sqlc-dev/sqlc",
		"github.com/teambition/rrule-go",
		"github.com/gorilla/websocket",
		"github.com/google/uuid",
		"github.com/prometheus/client_golang",
	}
	body := readRepoFile(t, "go.mod")
	lines := strings.Split(body, "\n")
	for _, mod := range required {
		// AC1 requires DIRECT requires. A line that contains the module path AND
		// a `v…` version token AND does NOT carry an `// indirect` marker counts
		// as direct. Handles both `require <mod> v…` and grouped require blocks.
		found := false
		needle := mod + " v"
		for _, line := range lines {
			if !strings.Contains(line, needle) {
				continue
			}
			if strings.Contains(line, "// indirect") {
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf("go.mod missing DIRECT require for %s (indirect-only references do not satisfy AC1)", mod)
		}
	}
}

// T-002 / AC1 — tools.go blank-imports the literal AC1 paths so they survive `go mod tidy`.
func TestToolsGoBlankImportsCoverEveryAC1Path(t *testing.T) {
	ac1Imports := []string{
		"github.com/go-chi/chi/v5",
		"github.com/jackc/pgx/v5",
		"github.com/sqlc-dev/sqlc/cmd/sqlc",
		"github.com/teambition/rrule-go",
		"github.com/gorilla/websocket",
		"github.com/google/uuid",
		"github.com/prometheus/client_golang/prometheus",
	}
	body := readRepoFile(t, "tools.go")
	if !strings.Contains(body, "//go:build tools") {
		t.Error("tools.go missing //go:build tools constraint")
	}
	for _, path := range ac1Imports {
		quoted := `"` + path + `"`
		if !strings.Contains(body, quoted) {
			t.Errorf("tools.go missing blank import %s", quoted)
		}
	}
}

// T-003 / AC2 — every AC2 frontend dependency is pinned in frontend/package.json.
func TestFrontendPackageJSONPinsRequiredDependencies(t *testing.T) {
	required := []string{
		"@tanstack/react-query",
		"@base-ui/react",
		"zustand",
		"date-fns",
	}
	body := readRepoFile(t, filepath.Join("frontend", "package.json"))
	var pkg struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(body), &pkg); err != nil {
		t.Fatalf("parse frontend/package.json: %v", err)
	}
	for _, dep := range required {
		if _, ok := pkg.Dependencies[dep]; !ok {
			t.Errorf("frontend/package.json missing dependency %q", dep)
		}
	}
}

// T-004 / AC3 — required directories and files from the Architecture project tree are on disk.
func TestRepoTreeContainsArchitectureRequiredPaths(t *testing.T) {
	type pathSpec struct {
		path string
		dir  bool
	}
	required := []pathSpec{
		{"cmd/server/main.go", false},
		{"cmd/server/router.go", false},
		{"internal/auth", true},
		{"internal/events", true},
		{"internal/handler", true},
		{"internal/middleware", true},
		{"internal/multica", true},
		{"internal/service/anomaly", true},
		{"internal/service/approval", true},
		{"internal/service/autofill", true},
		{"internal/service/capacity", true},
		{"internal/service/org", true},
		{"internal/service/rrule", true},
		{"pkg/db/queries", true},
		{"pkg/db/generated", true},
		{"migrations", true},
		{"scripts/team-app-cli", true},
		{"frontend/app", true},
		{"frontend/packages/core/team", true},
		{"frontend/packages/views/team", true},
		{"nginx/team.conf", false},
		{"docker-compose.yml", false},
		{"Dockerfile", false},
		{"sqlc.yaml", false},
	}
	root := repoRoot(t)
	for _, p := range required {
		info, err := os.Stat(filepath.Join(root, p.path))
		if err != nil {
			t.Errorf("missing %s: %v", p.path, err)
			continue
		}
		if p.dir && !info.IsDir() {
			t.Errorf("%s exists but is not a directory", p.path)
		}
		if !p.dir && info.IsDir() {
			t.Errorf("%s exists but is not a regular file", p.path)
		}
	}
}

// T-005 / AC4 — docker-compose.yml declares db, api, frontend with required image and ports.
func TestDockerComposeDeclaresRequiredServices(t *testing.T) {
	body := readRepoFile(t, "docker-compose.yml")
	for _, header := range []string{"\n  db:", "\n  api:", "\n  frontend:"} {
		if !strings.Contains(body, header) {
			t.Errorf("docker-compose.yml missing service header %q", strings.TrimSpace(header))
		}
	}
	if !strings.Contains(body, "pgvector/pgvector:pg17") {
		t.Error("docker-compose.yml does not pin db image to pgvector/pgvector:pg17 (AR9)")
	}
	if !strings.Contains(body, "8080:8080") {
		t.Error("docker-compose.yml does not expose api on :8080")
	}
	if !strings.Contains(body, "3000:3000") {
		t.Error("docker-compose.yml does not expose frontend on :3000")
	}
	// AR4 — db database name must be parameterised away from Multica's DB; .env.example sets it
	// to team_app. docker-compose.yml passes POSTGRES_DB through, which keeps the boundary.
	if !strings.Contains(body, "POSTGRES_DB") {
		t.Error("docker-compose.yml does not declare POSTGRES_DB on the db service (AR4)")
	}
}

// T-006 / AC4 — nginx/team.conf routes /api/, /gates/, / correctly and has SSE buffering off.
func TestNginxConfRoutesRequiredPaths(t *testing.T) {
	body := readRepoFile(t, filepath.Join("nginx", "team.conf"))
	if !strings.Contains(body, "server_name team.multica.uittai.com") {
		t.Error("nginx/team.conf missing server_name team.multica.uittai.com")
	}
	for _, route := range []struct {
		location string
		upstream string
	}{
		{"location /api/", "http://api:8080"},
		{"location /gates/", "http://api:8080"},
		// Catch-all "location / {" — disambiguated from /api/, /gates/, /healthz by the trailing
		// space + brace.
		{"location / {", "http://frontend:3000"},
	} {
		idx := strings.Index(body, route.location)
		if idx == -1 {
			t.Errorf("nginx/team.conf missing %s", route.location)
			continue
		}
		// Inspect the block body (up to the matching closing brace) for the proxy_pass target.
		// Locations are top-level in this conf, so the next "}" closes the block.
		nextBrace := strings.Index(body[idx:], "}")
		if nextBrace == -1 {
			t.Errorf("nginx/team.conf: unterminated block after %s", route.location)
			continue
		}
		block := body[idx : idx+nextBrace]
		if !strings.Contains(block, "proxy_pass "+route.upstream) {
			t.Errorf("nginx/team.conf: %s does not proxy_pass to %s", route.location, route.upstream)
		}
	}
	// SSE prefix: buffering must be off and read timeout must be 1h.
	if !strings.Contains(body, "proxy_buffering off") {
		t.Error("nginx/team.conf missing proxy_buffering off (required for SSE)")
	}
	if !strings.Contains(body, "proxy_read_timeout 1h") {
		t.Error("nginx/team.conf missing proxy_read_timeout 1h (required for SSE)")
	}
}

// T-007 / AC5 — .env.example documents every variable that ValidateRequiredEnv enforces.
// Reuses the package-private requiredEnvVars to keep .env.example and boot validation in sync.
func TestEnvExampleListsEveryRequiredEnvVar(t *testing.T) {
	body := readRepoFile(t, ".env.example")
	for _, name := range requiredEnvVars {
		prefix := name + "="
		if !strings.Contains(body, prefix) {
			t.Errorf(".env.example missing assignment for required variable %s", name)
		}
	}
}

// T-008 / AC1+AC3 — sqlc.yaml mirrors server/sqlc.yaml shape (engine, paths, pgx/v5, json+empty).
func TestSqlcYamlPinsEngineAndPaths(t *testing.T) {
	body := readRepoFile(t, "sqlc.yaml")
	for _, fragment := range []string{
		"engine: postgresql",
		"queries: pkg/db/queries/",
		"schema: migrations/",
		"out: pkg/db/generated",
		"sql_package: pgx/v5",
		"emit_json_tags: true",
		"emit_empty_slices: true",
	} {
		if !strings.Contains(body, fragment) {
			t.Errorf("sqlc.yaml missing %q", fragment)
		}
	}
}

// T-009 / AC6 + AR22 — frontend ESLint config keeps the no-restricted-imports patterns for
// both core/team/** and views/team/** boundary blocks.
func TestESLintConfigEnforcesPackageBoundaries(t *testing.T) {
	body := readRepoFile(t, filepath.Join("frontend", "eslint.config.mjs"))
	if !strings.Contains(body, "packages/core/team/") {
		t.Fatal("eslint.config.mjs missing packages/core/team/ override block")
	}
	if !strings.Contains(body, "packages/views/team/") {
		t.Fatal("eslint.config.mjs missing packages/views/team/ override block")
	}
	for _, pattern := range []string{
		`"next"`,
		`"next/*"`,
		`"react-router-dom"`,
		`"react-dom"`,
		`"react-dom/*"`,
		`"@base-ui/react"`,
		`"tailwindcss"`,
		`"tailwindcss/*"`,
	} {
		if !strings.Contains(body, pattern) {
			t.Errorf("eslint.config.mjs missing restricted-imports pattern %s", pattern)
		}
	}
	if !strings.Contains(body, "no-restricted-globals") {
		t.Error("eslint.config.mjs missing no-restricted-globals rule (AR22 — no process/window/localStorage in core/team)")
	}
	for _, global := range []string{`"process"`, `"window"`, `"localStorage"`} {
		if !strings.Contains(body, global) {
			t.Errorf("eslint.config.mjs missing %s in no-restricted-globals", global)
		}
	}
}

// T-010 / AR9 — CI workflow uses the agreed runtime versions and runs the documented gates.
func TestCIWorkflowRunsRequiredGates(t *testing.T) {
	body := readRepoFile(t, filepath.Join(".github", "workflows", "ci.yml"))
	for _, fragment := range []string{
		"pgvector/pgvector:pg17",
		"go-version: 1.26.1",
		"node-version: 22",
		"go test ./...",
		"pnpm lint",
		"pnpm build",
	} {
		if !strings.Contains(body, fragment) {
			t.Errorf("ci.yml missing %q", fragment)
		}
	}
}

// T-011 / AR14 — sole-writer comments survive in the placeholder packages so future contributors
// see the invariant before they edit.
func TestSoleWriterCommentsPresentInPlaceholders(t *testing.T) {
	cases := []string{
		"internal/service/approval/doc.go",
		"internal/service/autofill/doc.go",
		"internal/multica/doc.go",
		"internal/handler/workspace_work_week.go",
	}
	for _, path := range cases {
		body := readRepoFile(t, path)
		if !strings.Contains(body, "SOLE WRITER") {
			t.Errorf("%s missing SOLE WRITER comment", path)
		}
	}
}
