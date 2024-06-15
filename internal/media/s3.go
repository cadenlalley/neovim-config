package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/ptr"
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
	if prefix[0:1] == "/" {
		return "", fmt.Errorf("prefix should not start with '/', received: '%s'", prefix)
	}

	parts := strings.Split(file.Filename, ".")
	fileType := parts[len(parts)-1]

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

	key := fmt.Sprintf("%s.%s", target, fileType)
	_, err = u.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: ptr.String("kitchens-app-local-us-east-1"),
		Key:    ptr.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (u *S3FileManager) Ping(ctx context.Context) error {
	_, err := u.s3.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	return nil
}
