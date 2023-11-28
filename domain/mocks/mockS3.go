package mocks

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type MockS3Client struct {
	s3iface.S3API
	PutObjectCallCount int
	PutObjectInput     *s3.PutObjectInput
	PutObjectError     error
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	m.PutObjectCallCount++
	m.PutObjectInput = input
	return &s3.PutObjectOutput{}, m.PutObjectError
}
