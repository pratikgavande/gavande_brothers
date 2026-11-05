package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnv(t *testing.T) {
	unsetEnv(t, "PORT")
	unsetEnv(t, "MYSQL_DSN")

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "billing-app-config-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("remove temporary directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	env := "PORT=9090\nMYSQL_DSN=test:test@tcp(localhost:3306)/testdb\n"
	if err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte(env), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.MySQLDSN != "test:test@tcp(localhost:3306)/testdb" {
		t.Errorf("MySQLDSN = %q, want value from .env", cfg.MySQLDSN)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	value, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		var err error
		if existed {
			err = os.Setenv(key, value)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Errorf("restore %s: %v", key, err)
		}
	})
}
