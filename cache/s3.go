package cache

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Cache struct {
	cacheBucket string
	cacheDir    string
	s3Client    *s3.S3
}

func NewS3Cache() Cache {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_REGION") == "" {
		// TODO: raise error if AWS credentials are not set
		return nil
	}

	s3Client := s3.New(nil)

	var cacheDir string

	if os.Getenv("RISU_CACHE_DIR") != "" {
		cacheDir = os.Getenv("RISU_CACHE_DIR")
	}

	if cacheDir == "" {
		cacheDir = DefaultTarCacheDir
	}

	if os.Getenv("RISU_CACHE_BUCKET") == "" {
		// TODO: raise error if RISU_CACHE_BUCKET is not set
		return nil
	}

	cacheBucket := os.Getenv("RISU_CACHE_BUCKET")

	_, err := s3Client.HeadBucket(
		&s3.HeadBucketInput{
			Bucket: aws.String(cacheBucket),
		})

	if err != nil {
		// TODO: raise error if bucket not found
		return nil
	}

	return &S3Cache{cacheBucket, cacheDir, s3Client}
}

func (c *S3Cache) Get(key string) (string, error) {
	archivedCacheFilePath := getArchivedCacheFilePath(c.cacheDir, key)
	inflateDirPath := getInflateDirPath(DefaultInflatedCacheDir, key)

	_, err := c.s3Client.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(c.cacheBucket),
			Key:    aws.String(key),
		})

	if err != nil {
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == 404 {
			return "", nil
		}

		return "", err
	}

	resp, err := c.s3Client.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(c.cacheBucket),
			Key:    aws.String(key),
		})

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if err = ioutil.WriteFile(archivedCacheFilePath, data, 0644); err != nil {
		return "", err
	}

	if err = InflateTarGz(archivedCacheFilePath, inflateDirPath); err != nil {
		return "", err
	}

	return inflateDirPath, nil
}

func (c *S3Cache) Put(key, directory string) error {
	temporaryCacheDir := getArchivedCacheFilePath("/tmp/risu/", key)

	if err := DeflateTarGz(temporaryCacheDir, directory); err != nil {
		return err
	}

	file, err := os.Open(temporaryCacheDir)

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
