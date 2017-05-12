package mailroom

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func testS3(m *mailroom) error {
	params := &s3.HeadBucketInput{
		Bucket: aws.String(m.config.S3_Media_Bucket),
	}
	_, err := m.s3Client.HeadBucket(params)
	if err != nil {
		return err
	}

	return nil
}

func putS3File(m *mailroom, filename string, contentType string, contents []byte) (string, error) {
	path := filepath.Join(m.config.S3_Media_Prefix, filename[:4], filename)
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	params := &s3.PutObjectInput{
		Bucket:      aws.String(m.config.S3_Media_Bucket),
		Body:        bytes.NewReader(contents),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}
	_, err := m.s3Client.PutObject(params)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com%s", m.config.S3_Media_Bucket, path)
	return url, nil
}
