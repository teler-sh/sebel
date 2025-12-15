package cache

import (
	"os"
	"path/filepath"
)

// BuildCacheFilepath constructs a cache filepath in the user's cache directory.
//
// Returns empty string if the path cannot be created.
func BuildCacheFilepath(dir, file string) string {
	if dir == "" || file == "" {
		return ""
	}

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}

	fullDir := filepath.Join(userCacheDir, dir)
	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		return ""
	}

	return filepath.Join(fullDir, file)
}
