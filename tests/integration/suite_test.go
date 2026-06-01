//go:build integration

// Package integration_test is an end-to-end test suite that exercises the full
// HTTP → link-api → link-rpc (gRPC) → MySQL/Redis chain.
//
// Both services are started as subprocesses by TestMain; all tests talk to them
// over real TCP sockets.
//
// Prerequisites:
//
//	make infra-up && make migrate-up
//
// Run:
//
//	go test -tags=integration -v ./tests/integration/...
//
// Skip (e.g. in CI without infra):
//
//	INTEGRATION_SKIP=true go test -tags=integration ./...
package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ─── package-level state set by TestMain ─────────────────────────────────────

var (
	apiBase           string // e.g. "http://127.0.0.1:18080"
	testAdminUsername = "admin"
	testAdminPassword = "zerolink"
)

// ─── TestMain ────────────────────────────────────────────────────────────────

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION_SKIP") == "true" {
		fmt.Println("integration: skipped via INTEGRATION_SKIP=true")
		os.Exit(0)
	}

	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration: cannot find repo root:", err)
		os.Exit(1)
	}

	rpcPort, err := getFreePort()
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration: getFreePort rpc:", err)
		os.Exit(1)
	}
	apiPort, err := getFreePort()
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration: getFreePort api:", err)
		os.Exit(1)
	}
	apiBase = fmt.Sprintf("http://127.0.0.1:%d", apiPort)

	tmpDir, err := os.MkdirTemp("", "zero-link-integration-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration: MkdirTemp:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	rpcConf, apiConf := writeTempConfigs(tmpDir, rpcPort, apiPort)

	// Start link-rpc first (link-api connects to it on startup).
	rpcCmd := startService(repoRoot, "./services/link-rpc", rpcConf)
	defer func() { _ = rpcCmd.Process.Kill() }()

	if err := waitForTCP(fmt.Sprintf("127.0.0.1:%d", rpcPort), 30*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "integration: link-rpc did not start:", err)
		os.Exit(1)
	}

	// Start link-api.
	apiCmd := startService(repoRoot, "./services/link-api", apiConf)
	defer func() { _ = apiCmd.Process.Kill() }()

	// /readyz calls link-rpc health gRPC — proves the full chain is wired.
	if err := waitForHTTP(apiBase+"/readyz", 30*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "integration: link-api /readyz failed:", err)
		os.Exit(1)
	}

	// Smoke-check DB connectivity — migration must have run before tests.
	db, err := sql.Open("mysql", testDSN())
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration: sql.Open:", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pingErr := db.PingContext(ctx)
	cancel()
	db.Close()
	if pingErr != nil {
		fmt.Fprintln(os.Stderr, "integration: DB not reachable:", pingErr)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func testDSN() string {
	if v := os.Getenv("MYSQL_DSN"); v != "" {
		return v
	}
	return "zerolink:zerolink@tcp(127.0.0.1:3306)/zero_link?charset=utf8mb4&parseTime=true&loc=Local"
}

func testRedisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "127.0.0.1:6379"
}

// getFreePort asks the OS for an unused TCP port on loopback.
func getFreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// findRepoRoot walks up from cwd until it finds the go.mod declaring
// "module github.com/aliaxy/zero-link".
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(dir + "/go.mod")
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(line) == "module github.com/aliaxy/zero-link" {
					return dir, nil
				}
			}
		}
		parent := dir[:strings.LastIndex(dir, "/")]
		if parent == dir || parent == "" {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod with module github.com/aliaxy/zero-link not found")
}

// writeTempConfigs writes YAML config files for both services into dir.
// Returns (rpcConfigPath, apiConfigPath).
func writeTempConfigs(dir string, rpcPort, apiPort int) (string, string) {
	redis := testRedisAddr()
	dsn := testDSN()

	rpcYAML := fmt.Sprintf(`Name: link-rpc-test
ListenOn: 127.0.0.1:%d
Mode: dev
Log:
  Mode: console
  Level: error
Dependencies:
  MySQL:
    Endpoint: 127.0.0.1:3306
    DataSource: %s
  Redis:
    Endpoint: %s
CacheRedis:
  - Host: %s
    Type: node
Analytics:
  IPSalt: integration-test-salt
`, rpcPort, dsn, redis, redis)

	apiYAML := fmt.Sprintf(`Name: link-api-test
Host: 127.0.0.1
Port: %d
Mode: dev
Log:
  Mode: console
  Level: error
MaxConns: 100
MaxBytes: 1048576
Auth:
  Secret: test-jwt-secret-integration
  TokenTTLSeconds: 3600
Redis:
  Host: %s
  Type: node
LinkRPC:
  Target: 127.0.0.1:%d
RateLimit:
  RedirectPerIPPerSecond: 10000
  LoginPerIPPerMinute: 10000
Cors:
  AllowOrigins:
    - "*"
`, apiPort, redis, rpcPort)

	rpcConf := dir + "/link-rpc-test.yaml"
	apiConf := dir + "/link-api-test.yaml"
	if err := os.WriteFile(rpcConf, []byte(rpcYAML), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "integration: write rpc config:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(apiConf, []byte(apiYAML), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "integration: write api config:", err)
		os.Exit(1)
	}
	return rpcConf, apiConf
}

// startService runs "go run <pkg> -f <conf>" as a child process with its
// stdout/stderr attached to os.Stderr so service logs appear in go test output.
func startService(repoRoot, pkg, conf string) *exec.Cmd {
	cmd := exec.Command("go", "run", pkg, "-f", conf)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "integration: startService", pkg, err)
		os.Exit(1)
	}
	return cmd
}

// waitForTCP polls addr until a TCP connection succeeds or timeout expires.
// Used to wait for link-rpc's gRPC port before starting link-api.
func waitForTCP(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("TCP %s unreachable after %s", addr, timeout)
}

// waitForHTTP polls url with GET until HTTP 200 is returned or timeout expires.
func waitForHTTP(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := noRedirectHTTPClient.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("waitForHTTP %s timed out: %w", url, lastErr)
}

// uniqueName returns a test-scoped string safe for use in URLs and usernames.
func uniqueName(t *testing.T) string {
	t.Helper()
	safe := strings.NewReplacer("/", "", " ", "", "-", "_").Replace(t.Name())
	if len(safe) > 12 {
		safe = safe[len(safe)-12:]
	}
	return fmt.Sprintf("%s%d", safe, time.Now().UnixNano()%1_000_000)
}
