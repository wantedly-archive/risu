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
		cacheDir = DefaultTarCacheDir
	}

	if _, err := os.Stat(cacheDir); err != nil {
		os.MkdirAll(cacheDir, 0755)
	}

	return &LocalFsCache{cacheDir}
}

func (c *LocalFsCache) Get(key string) (string, error) {
	archivedCacheFilePath := getArchivedCacheFilePath(c.cacheDir, key)
	inflateDirPath := getInflateDirPath(DefaultInflatedCacheDir, key)

	if _, err := os.Stat(archivedCacheFilePath); err != nil {
		return "", nil
	}

	if err := InflateTarGz(archivedCacheFilePath, inflateDirPath); err != nil {
		return "", err
	}

	return inflateDirPath, nil
}

func (c *LocalFsCache) Put(key, directory string) error {
	if err := DeflateTarGz(getArchivedCacheFilePath(c.cacheDir, key), directory); err != nil {
		return err
	}

	return nil
}
