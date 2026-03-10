package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

// TowrHome returns the root towr directory, defaulting to ~/.towr.
// Respects the TOWR_HOME environment variable for overriding.
func TowrHome() string {
	if v := os.Getenv("TOWR_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback — should not happen in practice.
		return filepath.Join(os.TempDir(), ".towr")
	}
	return filepath.Join(home, ".towr")
}

// RepoStatePath returns the per-repo state directory: ~/.towr/repos/<hash>/
func RepoStatePath(repoRoot string) string {
	return filepath.Join(TowrHome(), "repos", RepoHash(repoRoot))
}

// WorktreeRoot returns the directory where worktrees are stored.
// Defaults to ~/.towr/worktrees but can be overridden in config.
func WorktreeRoot() string {
	return filepath.Join(TowrHome(), "worktrees")
}

// RepoHash produces a short deterministic hash of a repo root path,
// used to namespace per-repo state without filesystem-unfriendly characters.
func RepoHash(repoRoot string) string {
	h := sha256.Sum256([]byte(repoRoot))
	return fmt.Sprintf("%x", h[:8]) // 16 hex chars — short but collision-resistant
}

// EnsureTowrDirs creates the core towr directory structure if it doesn't exist.
func EnsureTowrDirs() error {
	dirs := []string{
		TowrHome(),
		filepath.Join(TowrHome(), "repos"),
		WorktreeRoot(),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}
	return nil
}

// GlobalConfigPath returns the path to the global config file.
func GlobalConfigPath() string {
	return filepath.Join(TowrHome(), "global-config.toml")
}

// RepoConfigPath returns the path to a repo-specific config file.
func RepoConfigPath(repoRoot string) string {
	return filepath.Join(RepoStatePath(repoRoot), "config.toml")
}
