package cache

import (
	"os"
	"path/filepath"
)

type LocalFsCache struct {
	cacheDir string
}

const (
	DefaultCacheDir = "/tmp/risu/cache"
)

func NewLocalFsCache() Cache {
	var cacheDir string

	if os.Getenv("RISU_CACHE_DIR") != "" {
		cacheDir = os.Getenv("RISU_CACHE_DIR")
	}

	if cacheDir == "" {
		cacheDir = DefaultCacheDir
	}

	if _, err := os.Stat(cacheDir); err != nil {
		os.MkdirAll(cacheDir, 0755)
	}

	return &LocalFsCache{cacheDir}
}

func (c *LocalFsCache) Get(key string) (string, error) {
	return "", nil
}

func (c *LocalFsCache) Put(key, directory string) error {
	if err := DeflateTarGz(cachePath(key), directory); err != nil {
		return err
	}

	return nil
}

func cachePath(key string) string {
	return DefaultCacheDir + string(filepath.Separator) + key + ".tar.gz"
}
