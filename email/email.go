package email

import (
	"bytes"
	_ "embed"
	"html/template"

	"github.com/KevJV07/SendEmailAPI/transactions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Embed the email_template.html file into the binary
//
//go:embed email_template.html
var emailTemplate string

// EmailData represents the data needed to generate an email
type EmailData struct {
	TotalBalance float64
	MonthNames   map[string]string
	MonthlyStats map[string]*struct {
		Balance      float64
		Transactions int
		CreditTotal  float64
		DebitTotal   float64
		CreditCount  int
		DebitCount   int
	}
}

// GenerateEmailContent constructs and returns an EmailData object with transaction information
func GenerateEmailContent(totalBalance float64, monthlyStats map[string]*transactions.MonthlyStat) EmailData {
	// Initialize MonthNames
	monthNames := map[string]string{
		"1":  "January",
		"2":  "February",
		"3":  "March",
		"4":  "April",
		"5":  "May",
		"6":  "June",
		"7":  "July",
		"8":  "August",
		"9":  "September",
		"10": "October",
		"11": "November",
		"12": "December",
	}

	emailData := EmailData{
		TotalBalance: totalBalance,
		MonthNames:   monthNames,
		MonthlyStats: make(map[string]*struct {
			Balance      float64
			Transactions int
			CreditTotal  float64
			DebitTotal   float64
			CreditCount  int
			DebitCount   int
		}),
	}

	for month, stat := range monthlyStats {
		emailData.MonthlyStats[month] = &struct {
			Balance      float64
			Transactions int
			CreditTotal  float64
			DebitTotal   float64
			CreditCount  int
			DebitCount   int
		}{
			Balance:      stat.Balance,
			Transactions: stat.Transactions,
			CreditTotal:  stat.CreditTotal,
			DebitTotal:   stat.DebitTotal,
			CreditCount:  stat.CreditCount,
			DebitCount:   stat.DebitCount,
		}
	}

	return emailData
}

// SendSESEmail sends an email using Amazon SES
func SendSESEmail(senderEmail, recipientEmail, subject string, emailData EmailData) error {
	// Create a new AWS session with the default configuration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return err
	}

	// Create an SES service instance
	svc := ses.New(sess)

	// Process the embedded template using EmailData
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, emailData)
	if err != nil {
		return err
	}
	renderedEmailBody := buf.String()

	// Create a SendEmail request object
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{

			ToAddresses: []*string{
				aws.String(recipientEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(renderedEmailBody),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(senderEmail),
	}

	// Send the mail using SES
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}
