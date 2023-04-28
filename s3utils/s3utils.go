package s3utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ReadDataFromS3 reads data from an S3 bucket and returns it as a string.
func ReadDataFromS3() (string, error) {

	// Get the environment variables from the template.yaml
	bucket := os.Getenv("BUCKET_NAME")
	uploadedKey := os.Getenv("TXNS_FILE")

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create a new S3 service client
	s3Svc := s3.New(sess)

	// Get the object from the S3 bucket
	resp, err := s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(uploadedKey),
	})

	if err != nil {
		return "", fmt.Errorf("failed to read data from S3: %v", err)
	}

	// Ensure the response body is closed after the function returns
	defer resp.Body.Close()

	// Read the data from the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read data from S3 object body: %v", err)
	}

	// Return the data as a string
	return string(body), nil
}
