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
	temporaryCacheDir := getArchivedCacheFilePath("/tmp/risu/", key)
	archivedCacheFilePath := getArchivedCacheFilePath(c.cacheDir, key)

	if err := DeflateTarGz(temporaryCacheDir, directory); err != nil {
		return err
	}

	if err := os.Rename(temporaryCacheDir, archivedCacheFilePath); err != nil {
		return err
	}

	return nil
}
