# Monthly Transaction Report Lambda

This Lambda function reads transaction data from an S3 bucket, performs calculations on the data, generates a monthly transaction report, and sends it via email using Amazon SES.

## Prerequisites

- AWS CLI installed and configured with Administrator permissions
- SAM CLI installed
- Golang installed
- Verify the sender and recipient email addresses in Amazon SES

## Setup

1. Build and deploy your Lambda function using the SAM CLI:

```bash
sam build
sam deploy --guided
```
The output will give you your API URL and your bucket name.

2. After the deployment, an S3 bucket will be automatically created. Note the bucket name.
3. Upload the `txns.csv` file to the newly created S3 bucket. 
```bash
aws s3 cp file.csv s3://bucket-name/ txns.csv
```
Change the bucket name, for yours and file.csv for the direction of your txns.csv file.

## Usage

To test the Lambda function, send an API Gateway request with the following path parameters:

- `sender_email`: The verified email address that will send the monthly transaction report
- `recipient_email`: The verified email address that will receive the monthly transaction report

Example API Gateway request:

https://your-api-gateway-url/path?sender_email=sender@example.com&recipient_email=recipient@example.com

