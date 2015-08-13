package cache

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Cache struct {
	cacheBucket string
	s3Client    *s3.S3
}

const (
	TemporaryCacheDir = "/tmp/risu/cache"
)

func NewS3Cache() Cache {
	// TODO: raise error if AWS credentials are not set

	s3Client := s3.New(nil)

	var cacheBucket string

	if os.Getenv("RISU_CACHE_BUCKET") != "" {
		cacheBucket = os.Getenv("RISU_CACHE_BUCKET")
	}

	// TODO: raise error if bucket not found

	return &S3Cache{cacheBucket, s3Client}
}

func (c *S3Cache) Get(key string) (string, error) {
	// TODO
	return "", nil
}

func (c *S3Cache) Put(key, directory string) error {
	temporaryCache := temporaryCachePath(key)

	if err := DeflateTarGz(temporaryCache, directory); err != nil {
		return err
	}

	file, err := os.Open(temporaryCache)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = c.s3Client.PutObject(
		&s3.PutObjectInput{
			ACL:    aws.String("private"),
			Bucket: aws.String(c.cacheBucket),
			Body:   file,
			Key:    aws.String(key),
		})

	if err != nil {
		return err
	}

	return nil
}

func temporaryCachePath(key string) string {
	return TemporaryCacheDir + string(filepath.Separator) + key + ".tar.gz"
}
