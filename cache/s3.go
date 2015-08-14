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
	s3Client    *s3.S3
}

const (
	TemporaryCacheDir = "/tmp/risu/cache"
)

func NewS3Cache() Cache {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_REGION") == "" {
		// TODO: raise error if AWS credentials are not set
		return nil
	}

	s3Client := s3.New(nil)

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

	return &S3Cache{cacheBucket, s3Client}
}

func (c *S3Cache) Get(key string) (string, error) {
	temporaryCache := cacheFilePath(DefaultCacheDir, key)
	inflateDir := inflateDirPath(DefaultCacheDir, key)

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

	if err = ioutil.WriteFile(temporaryCache, data, 0644); err != nil {
		return "", err
	}

	if err = InflateTarGz(temporaryCache, inflateDir); err != nil {
		return "", err
	}

	return inflateDir, nil
}

func (c *S3Cache) Put(key, directory string) error {
	temporaryCache := cacheFilePath(DefaultCacheDir, key)

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
