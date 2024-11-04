package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kitchens-io/kitchens-api/pkg/ptr"
	"github.com/segmentio/ksuid"
)

type S3FileManager struct {
	s3     *s3.Client
	Bucket string
}

func NewS3FileManager(s3 *s3.Client, bucket string) *S3FileManager {
	return &S3FileManager{
		s3:     s3,
		Bucket: bucket,
	}
}

func (u *S3FileManager) UploadFromHeader(ctx context.Context, file *multipart.FileHeader, prefix string) (string, error) {
	// Validation
	if prefix[0:1] == "/" {
		return "", fmt.Errorf("prefix should not start with '/', received: '%s'", prefix)
	}

	parts := strings.Split(file.Filename, ".")
	if len(parts) == 1 {
		return "", fmt.Errorf("file missing extension, received: '%s'", file.Filename)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return "", err
	}

	target := prefix + ksuid.New().String()
	fileType := parts[len(parts)-1]

	key := fmt.Sprintf("%s.%s", target, fileType)
	_, err = u.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: ptr.String(u.Bucket),
		Key:    ptr.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (u *S3FileManager) UploadFromHeaders(ctx context.Context, files []*multipart.FileHeader, prefix string) ([]string, error) {
	keys := make([]string, len(files))

	// Validation
	if prefix[0:1] == "/" {
		return nil, fmt.Errorf("prefix should not start with '/', received: '%s'", prefix)
	}

	uuid := ksuid.New().String()

	for i, file := range files {
		parts := strings.Split(file.Filename, ".")
		if len(parts) == 1 {
			return nil, fmt.Errorf("file missing extension, received: '%s'", file.Filename)
		}

		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, src); err != nil {
			return nil, err
		}

		fileType := parts[len(parts)-1]

		target := prefix + uuid
		if len(files) > 1 {
			target = fmt.Sprintf("%s_%d", prefix+uuid, i)
		}

		key := fmt.Sprintf("%s.%s", target, fileType)
		_, err = u.s3.PutObject(ctx, &s3.PutObjectInput{
			Bucket: ptr.String(u.Bucket),
			Key:    ptr.String(key),
			Body:   bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (u *S3FileManager) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	result, err := u.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: ptr.String(u.Bucket),
		Key:    ptr.String(key),
	})
	if err != nil {
		return nil, "", err
	}

	return result.Body, *result.ContentType, nil
}

func (u *S3FileManager) Ping(ctx context.Context) error {
	result, err := u.s3.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	for _, bucket := range result.Buckets {
		if *bucket.Name == u.Bucket {
			return nil
		}
	}

	return fmt.Errorf("'%s' not found in list buckets response", u.Bucket)
}
