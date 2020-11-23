package storage

import (
	"bytes"
	"io/ioutil"
	"net/http"
	pathutil "path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AmazonS3Backend is a storage backend for Amazon S3
type AmazonS3Backend struct {
	Client     *s3.S3
	Downloader *s3manager.Downloader
	Prefix     string
	Uploader   *s3manager.Uploader
	SSE        string
}

// NewAmazonS3BackendWithCredentials creates a new instance of AmazonS3Backend with credentials
func NewAmazonS3BackendWithCredentials(prefix string, region string, endpoint string, sse string, credentials *credentials.Credentials) *AmazonS3Backend {
	service := s3.New(session.New(), &aws.Config{
		Credentials:      credentials,
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(strings.HasPrefix(endpoint, "http://")),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
	})
	b := &AmazonS3Backend{
		Client:     service,
		Downloader: s3manager.NewDownloaderWithClient(service),
		Prefix:     cleanPrefix(prefix),
		Uploader:   s3manager.NewUploaderWithClient(service),
		SSE:        sse,
	}
	return b
}

// CreateBucket Create a S3 bucket, at prefix
func (b AmazonS3Backend) CreateBucket(bucket string) (Object, error) {
	var object Object
	s3Input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}
	s3Result, err := b.Client.CreateBucket(s3Input)
	if err != nil {
		return object, err
	}
	object.Path = s3Result.String()
	return object, nil
}

// DeleteBucket Create a S3 bucket, at prefix
func (b AmazonS3Backend) DeleteBucket(bucket string) error {
	s3Input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	_, err := b.Client.DeleteBucket(s3Input)
	return err
}

// ListObjects lists all objects in Amazon S3 bucket, at prefix
func (b AmazonS3Backend) ListObjects(bucket, prefix string) ([]Object, error) {
	var objects []Object
	prefix = pathutil.Join(b.Prefix, prefix)
	s3Input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	for {
		s3Result, err := b.Client.ListObjects(s3Input)
		if err != nil {
			return objects, err
		}
		for _, obj := range s3Result.Contents {
			path := removePrefixFromObjectPath(prefix, *obj.Key)
			if objectPathIsInvalid(path) {
				continue
			}
			object := Object{
				Path:         path,
				Content:      []byte{},
				LastModified: *obj.LastModified,
			}
			objects = append(objects, object)
		}
		if !*s3Result.IsTruncated {
			break
		}
		s3Input.Marker = s3Result.Contents[len(s3Result.Contents)-1].Key
	}
	return objects, nil
}

// GetObject retrieves an object from Amazon S3 bucket, at prefix
func (b AmazonS3Backend) GetObject(bucket, path string) (Object, error) {
	var object Object
	object.Path = path
	var content []byte
	s3Input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(pathutil.Join(b.Prefix, path)),
	}
	s3Result, err := b.Client.GetObject(s3Input)
	if err != nil {
		return object, err
	}
	content, err = ioutil.ReadAll(s3Result.Body)
	if err != nil {
		return object, err
	}
	object.Content = content
	object.LastModified = *s3Result.LastModified
	return object, nil
}

// PutObject uploads an object to Amazon S3 bucket, at prefix
func (b AmazonS3Backend) PutObject(bucket, path string, content []byte) error {
	s3Input := &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(pathutil.Join(b.Prefix, path)),
		Body:        bytes.NewBuffer(content),
		ContentType: aws.String(http.DetectContentType(content)),
	}

	if b.SSE != "" {
		s3Input.ServerSideEncryption = aws.String(b.SSE)
	}

	_, err := b.Uploader.Upload(s3Input)
	return err
}

// DeleteObject removes an object from Amazon S3 bucket, at prefix
func (b AmazonS3Backend) DeleteObject(bucket, path string) error {
	s3Input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(pathutil.Join(b.Prefix, path)),
	}
	_, err := b.Client.DeleteObject(s3Input)
	return err
}
