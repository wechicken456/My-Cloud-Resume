package service

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/model"
	"time"

	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
)

type NotificationService struct {
	sesClient *ses.Client
	snsClient *sns.Client
	config    *config.Config
}

func NewNotificationService(sesClient *ses.Client, snsClient *sns.Client, config *config.Config) *NotificationService {
	return &NotificationService{
		sesClient: sesClient,
		snsClient: snsClient,
		config:    config,
	}
}

func (ns *NotificationService) SendEmailNotification(ctx context.Context, payload *model.NotificationPayload) error {
	if ns.config.NotificationDstEmail == "" || ns.config.NotificationSrcEmail == "" {
		return nil // No emails configured, skip sending
	}

	// check payload type and construct email content
	estLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		payload.Timestamp = payload.Timestamp.In(estLocation) // Convert to EST
	}
	var subject, body string
	switch payload.Type {
	case "like":
		subject = "New Like 👍 on Your Resume!"
		body = fmt.Sprintf(`
Someone liked your resume!

Time: %s
Source: %s

Visit your resume: https://www.pwnph0fun.com

Very nice,
Your Cloud Resume Lambda
`, payload.Timestamp.Format("January 2, 2006 at 3:04 PM EST"), payload.Source)

	case "contact":
		subject = "New Contact Form Submission on Your Cloud Resume!"
		name, _ := payload.Data["name"].(string)
		email, _ := payload.Data["email"].(string)
		message, _ := payload.Data["message"].(string)

		body = fmt.Sprintf(`
New contact form submission received:

Name: %s
Email: %s
Message: %s

Received at: %s
Source: %s

Reply to this person directly at: %s

Very nice,
Your Cloud Resume Lambda
`, name, email, message, payload.Timestamp.Format("January 2, 2006 at 3:04 PM EST"), payload.Source, email)

	default:
		return fmt.Errorf("unknown notification type: %s", payload.Type)
	}

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				ns.config.NotificationDstEmail,
			},
		},

		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
				Body: &types.Body{
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(body),
					},
				},
			},
		},

		FromEmailAddress: &ns.config.NotificationSrcEmail,
	}

	_, err = ns.sesClient.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send SES email: %w", err)
	}

	return nil
}

func (ns *NotificationService) SendSMSNotification(ctx context.Context, payload *model.NotificationPayload) error {
	return nil
}
