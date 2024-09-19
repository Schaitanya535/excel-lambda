package sesmailer

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const region = "us-west-2"

// SendEmail sends an email with SES containing the S3 link
func SendEmail(sender, recipient, s3URL string) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	svc := ses.New(sess)

	subject := "Your Data File is Ready"
	htmlBody := fmt.Sprintf("<h1>Your Data File</h1><p>Download the file from <a href='%s'>this link</a>. The link expires in 4 hours.</p>", s3URL)
	textBody := fmt.Sprintf("Your file is ready. Download it from this link: %s. The link expires in 4 hours.", s3URL)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	_, err := svc.SendEmail(input)
	return err
}
