package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultInitDirName = ".taskforge"
	composeFileName    = "docker-compose.yml"
	envFileName        = ".env"
)

var (
	initForce      bool
	initForeground bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a local Taskforge stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseDir, err := resolveInitDir()
		if err != nil {
			return err
		}

		if err := os.MkdirAll(baseDir, 0o700); err != nil {
			return err
		}

		composePath := filepath.Join(baseDir, composeFileName)
		envPath := filepath.Join(baseDir, envFileName)

		values, extras, err := readEnvFile(envPath)
		if err != nil {
			if os.IsNotExist(err) {
				values = map[string]string{}
				extras = nil
			} else {
				return err
			}
		}

		if initForce {
			values = map[string]string{}
			extras = nil
		}

		required := requiredEnvValues(values)
		for key, val := range required {
			if values[key] == "" {
				values[key] = val
			}
		}

		if err := writeEnvFile(envPath, values, extras); err != nil {
			return err
		}

		if err := writeComposeFile(composePath, initForce); err != nil {
			return err
		}

		if err := maybeUpdateConfig(values); err != nil {
			return err
		}

		if err := runDockerCompose(baseDir, composePath, "pull"); err != nil {
			return err
		}

		if err := runDockerCompose(baseDir, composePath, "up", "-d", "postgres", "redis"); err != nil {
			return err
		}

		if err := waitForPostgres(baseDir, composePath); err != nil {
			return err
		}

		if err := runDockerCompose(
			baseDir,
			composePath,
			"run",
			"--rm",
			"--workdir",
			"/app/apps/server",
			"server",
			"node_modules/.bin/prisma",
			"migrate",
			"deploy",
		); err != nil {
			return err
		}

		upArgs := []string{"up"}
		if !initForeground {
			upArgs = append(upArgs, "-d")
		}
		if err := runDockerCompose(baseDir, composePath, upArgs...); err != nil {
			return err
		}

		fmt.Printf("Taskforge initialized in %s\n", baseDir)
		fmt.Printf("Compose file: %s\n", composePath)
		fmt.Printf("Env file: %s\n", envPath)
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing init files")
	initCmd.Flags().BoolVar(&initForeground, "foreground", false, "Run services in foreground")
}

func resolveInitDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultInitDirName), nil
}

func requiredEnvValues(existing map[string]string) map[string]string {
	adminToken := existing["TASKFORGE_ADMIN_TOKEN"]
	if adminToken == "" {
		adminToken = randomHex(32)
	}

	secretKey := existing["TASKFORGE_SECRET_KEY"]
	if secretKey == "" {
		secretKey = randomHex(32)
	}

	return map[string]string{
		"DATABASE_URL":          "postgresql://taskforge:taskforge@postgres:5432/taskforge",
		"REDIS_URL":             "redis://redis:6379",
		"TASKFORGE_ADMIN_TOKEN": adminToken,
		"TASKFORGE_SECRET_KEY":  secretKey,
		"PORT":                  "3000",
	}
}

func randomHex(bytes int) string {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func readEnvFile(path string) (map[string]string, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	values := make(map[string]string)
	var extras []string

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			extras = append(extras, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			extras = append(extras, line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			extras = append(extras, line)
			continue
		}
		values[key] = value
	}

	return values, extras, nil
}

func writeEnvFile(path string, values map[string]string, extras []string) error {
	order := []string{
		"DATABASE_URL",
		"REDIS_URL",
		"TASKFORGE_ADMIN_TOKEN",
		"TASKFORGE_SECRET_KEY",
		"PORT",
	}

	var lines []string
	lines = append(lines, "# Taskforge local environment")
	for _, key := range order {
		if val, ok := values[key]; ok && val != "" {
			lines = append(lines, fmt.Sprintf("%s=%s", key, val))
		}
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	known := make(map[string]struct{}, len(order))
	for _, key := range order {
		known[key] = struct{}{}
	}

	for _, key := range keys {
		if _, ok := known[key]; ok {
			continue
		}
		val := values[key]
		if val == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, val))
	}

	if len(extras) > 0 {
		lines = append(lines, extras...)
	}

	content := strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0o600)
}

func writeComposeFile(path string, force bool) error {
	if !force {
		data, err := os.ReadFile(path)
		if err == nil {
			if !strings.Contains(string(data), "container_name: taskforge-") {
				return nil
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	content := strings.TrimLeft(`name: taskforge

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: taskforge
      POSTGRES_PASSWORD: taskforge
      POSTGRES_DB: taskforge
    ports:
      - "5432:5432"
    volumes:
      - taskforge_pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U taskforge -d taskforge"]
      interval: 5s
      timeout: 3s
      retries: 10

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - taskforge_redisdata:/data
    command: ["redis-server", "--appendonly", "yes"]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 10

  server:
    image: gentij/taskforge-server:latest
    env_file:
      - .env
    ports:
      - "3000:3000"
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  worker:
    image: gentij/taskforge-worker:latest
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

volumes:
  taskforge_pgdata:
  taskforge_redisdata:
`, "\n")

	return os.WriteFile(path, []byte(content), 0o644)
}

func maybeUpdateConfig(values map[string]string) error {
	cfg, path, err := loadConfig()
	if err != nil {
		return err
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = defaultServerURL
	}

	if strings.TrimSpace(cfg.Token) == "" {
		if token := strings.TrimSpace(values["TASKFORGE_ADMIN_TOKEN"]); token != "" {
			cfg.Token = token
			return saveConfig(path, cfg)
		}
	}

	return nil
}

func runDockerCompose(baseDir string, composePath string, args ...string) error {
	if hasDockerComposePlugin() {
		fullArgs := append([]string{"compose", "-f", composePath}, args...)
		cmd := exec.Command("docker", fullArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = baseDir
		return cmd.Run()
	}

	fullArgs := append([]string{"-f", composePath}, args...)
	cmd := exec.Command("docker-compose", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = baseDir
	return cmd.Run()
}

func waitForPostgres(baseDir string, composePath string) error {
	for attempt := 1; attempt <= 20; attempt++ {
		err := runDockerCompose(
			baseDir,
			composePath,
			"exec",
			"-T",
			"postgres",
			"pg_isready",
			"-U",
			"taskforge",
			"-d",
			"taskforge",
		)
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("postgres did not become ready in time")
}

func hasDockerComposePlugin() bool {
	cmd := exec.Command("docker", "compose", "version")
	return cmd.Run() == nil
}
