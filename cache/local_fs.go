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
	cache := cachePath(c.cacheDir, key)
	inflateDir := inflateDirPath(c.cacheDir, key)

	if _, err := os.Stat(cache); err != nil {
		return "", nil
	}

	if err := InflateTarGz(cache, inflateDir); err != nil {
		return "", err
	}

	return inflateDir, nil
}

func (c *LocalFsCache) Put(key, directory string) error {
	if err := DeflateTarGz(cachePath(c.cacheDir, key), directory); err != nil {
		return err
	}

	return nil
}

func cachePath(cacheDir, key string) string {
	return cacheDir + string(filepath.Separator) + key + ".tar.gz"
}

func inflateDirPath(cacheDir, key string) string {
	return cacheDir + string(filepath.Separator) + key
}
