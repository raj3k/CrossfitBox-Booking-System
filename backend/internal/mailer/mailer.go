package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	models "crossfitbox.booking.system/internal/data"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return nil
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return nil
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

const (
	// The subject line for the email.
	Subject = "Welcome to CrossBoxFit!"

	// The HTML body for the email.
	HtmlBodyTemplate = "<p>Hi,</p>" +
		"<p>Thanks for signing up for a CrossBoxFit account. We're excited to have you on board!</p> <p>For future reference, your user ID number is %s.</p>" +
		"<p>Thanks,</p>" +
		"<p>The CrossBoxFit Team</p>"

	//The email body for recipients with non-HTML email clients.
	TextBodyTemplate = "Hi,\n" +

		"Thanks for signing up for a CrossBoxFit account. We're exited to have you on board!\n" +

		"For future reference, your user ID number is %s.\n" +

		"Thanks,\n" +

		"The CrossBoxFit Team"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func (m Mailer) SendMailAWS(recipient string, data interface{}) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "cfbox-ses",
	})
	if err != nil {
		return err
	}

	HtmlBody := fmt.Sprintf(HtmlBodyTemplate, data.(*models.User).ID)
	TextBody := fmt.Sprintf(TextBodyTemplate, data.(*models.User).ID)

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{
		Region: aws.String("eu-central-1"),
	})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(m.sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return aerr
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return aerr
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return aerr
			default:
				return aerr
			}
		}
		return err
	}

	fmt.Println("Email Sent to address: " + recipient)
	fmt.Println(result)
	return nil
}
