package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	excelwriter "github.com/Schaitanya535/excel-lambda/m/excel"
	sesmailer "github.com/Schaitanya535/excel-lambda/m/mail-service"
	s3uploader "github.com/Schaitanya535/excel-lambda/m/s3-upload"
)

type SQSMessageBody struct {
	Email   string            `json:"email"`
	Filters map[string]string `json:"filters"`
}

type ApiResponse struct {
	Name    string                   `json:"name"`
	Headers []string                 `json:"headers"`
	Data    []map[string]interface{} `json:"data"`
}

var (
	bucketName = os.Getenv("S3_BUCKET_NAME") // Bucket name stored in environment variable
	apiURL     = os.Getenv("API_URL")        // External API base URL stored in environment variable
)

// Lambda handler to process SQS events
func handler(sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var msgBody SQSMessageBody
		if err := json.Unmarshal([]byte(message.Body), &msgBody); err != nil {
			log.Printf("Failed to unmarshal SQS message body: %v", err)
			return err
		}

		// Step 1: Call the Nest.js API using query filters
		fullAPIURL := fmt.Sprintf("%s?%s", apiURL, buildQueryParams(msgBody.Filters))
		resp, err := http.Get(fullAPIURL)
		if err != nil {
			log.Printf("Failed to fetch data from API: %v", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read API response: %v", err)
			return err
		}

		var apiResponse ApiResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			log.Printf("Failed to unmarshal API response: %v", err)
			return err
		}

		// Step 2: Write data to an Excel file in memory
		var buf bytes.Buffer
		if err := excelwriter.WriteToExcelFile(apiResponse.Headers, apiResponse.Data, &buf); err != nil {
			log.Printf("Failed to write Excel file: %v", err)
			return err
		}

		// Step 3: Upload the Excel file to S3
		fileName := fmt.Sprintf("%s.xlsx", apiResponse.Name)
		s3URL, err := s3uploader.UploadToS3(buf, fileName, bucketName)
		if err != nil {
			log.Printf("Failed to upload file to S3: %v", err)
			return err
		}

		// Step 4: Send email using SES with pre-signed S3 URL
		if err := sesmailer.SendEmail("sender@example.com", msgBody.Email, s3URL); err != nil {
			log.Printf("Failed to send email: %v", err)
			return err
		}

		log.Printf("Successfully processed message for %s", msgBody.Email)
	}
	return nil
}

// buildQueryParams constructs the query string from filters
func buildQueryParams(filters map[string]string) string {
	query := ""
	for key, value := range filters {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", key, value)
	}
	return query
}

func main() {
	lambda.Start(handler)
}
