package dispatch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/brianaffirm/towr/internal/config"
)

// AcquireLaunchLock acquires a file lock for serializing Claude launches.
// Returns a cleanup function to release the lock.
func AcquireLaunchLock() (func(), error) {
	lockPath := filepath.Join(config.TowrHome(), ".claude-launch.lock")

	// Try to acquire for up to 60 seconds.
	for i := 0; i < 60; i++ {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			fmt.Fprintf(f, "%d", os.Getpid())
			f.Close()
			return func() { os.Remove(lockPath) }, nil
		}
		// Check if lock is stale (>90 seconds old).
		if info, statErr := os.Stat(lockPath); statErr == nil {
			if time.Since(info.ModTime()) > 90*time.Second {
				os.Remove(lockPath)
				continue
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("timeout acquiring launch lock")
}
