package s3uploader

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const region = "us-west-2" // Adjust as needed

// UploadToS3 uploads a file to S3 and returns the pre-signed URL
func UploadToS3(buf bytes.Buffer, fileName, bucketName string) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Upload the file to S3
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", err
	}

	// Generate a pre-signed URL for the file (4-hour expiry)
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	s3URL, err := req.Presign(4 * time.Hour)
	if err != nil {
		return "", err
	}

	return s3URL, nil
}
