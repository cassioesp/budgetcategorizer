package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jbleduigou/budgetcategorizer/exporter"
	"github.com/jbleduigou/budgetcategorizer/parser"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	// Create all collaborators for command
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	parser := parser.NewParser()
	exporter := exporter.NewExporter()
	for _, record := range s3Event.Records {
		// Retrieve data from S3 event
		s3event := record.S3
		objectKey := strings.ReplaceAll(s3event.Object.Key, "input/", "")
		// Instantiate a command
		c := &command{s3event.Bucket.Name, objectKey, downloader, uploader, parser, exporter}
		// Execute the command
		c.execute()
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
