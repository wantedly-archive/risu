package cache

import (
	"os"
)

type LocalFsCache struct {
	cacheDir string
}

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
	cachePath := cacheFilePath(c.cacheDir, key)
	inflateDir := inflateDirPath(c.cacheDir, key)

	if _, err := os.Stat(cachePath); err != nil {
		return "", nil
	}

	if err := InflateTarGz(cachePath, inflateDir); err != nil {
		return "", err
	}

	return inflateDir, nil
}

func (c *LocalFsCache) Put(key, directory string) error {
	if err := DeflateTarGz(cacheFilePath(c.cacheDir, key), directory); err != nil {
		return err
	}

	return nil
}
