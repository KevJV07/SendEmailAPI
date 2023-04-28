package main

import (
	"github.com/KevJV07/SendEmailAPI/email"
	"github.com/KevJV07/SendEmailAPI/s3utils"
	"github.com/KevJV07/SendEmailAPI/transactions"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is the main Lambda function entry point.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Read email-related values from the path parameters
	senderEmail := request.PathParameters["sender_email"]
	recipientEmail := request.PathParameters["recipient_email"]

	// Read data from the S3 bucket
	data, err := s3utils.ReadDataFromS3()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error reading data from S3"}, nil
	}

	// Perform transaction calculations
	monthlyStats, totalBalance, err := transactions.PerformTransactionCalculations(data)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error performing transaction calculations"}, nil
	}

	// Generate email content
	emailData := email.GenerateEmailContent(totalBalance, monthlyStats)

	// Send email
	err = email.SendSESEmail(senderEmail, recipientEmail, "Monthly Transaction Report", emailData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error sending email"}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Email sent successfully",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
