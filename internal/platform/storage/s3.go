package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
)

// S3 implementation of storage
type S3 struct {
	Bucket *string
	Client *s3.S3
}

// NewS3 returns a pointer to new s3 instance
func NewS3(bucket string) *S3 {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(os.Getenv("AWS_DEFAULT_REGION"))

	return &S3{Bucket: aws.String(bucket), Client: s3.New(sess)}
}

// Read file from storage
func (s *S3) Read(file string) ([]byte, error) {
	if obj, err := s.Client.GetObject(&s3.GetObjectInput{Bucket: s.Bucket, Key: aws.String(file)}); err != nil {
		return []byte(""), err
	} else {
		return ioutil.ReadAll(obj.Body)
	}
}

// Write to storage
func (s *S3) Write(file string, content []byte, contentType string) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket:               s.Bucket,
		Key:                  aws.String(file),
		Body:                 bytes.NewReader(content),
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		ACL:                  aws.String("private"),
	})

	return err
}
