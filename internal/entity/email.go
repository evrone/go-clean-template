package entity

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var sesClient *sesv2.Client

func InitSES(ctx context.Context) error {
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if region == "" {
		return errors.New("AWS_REGION is required")
	}

	cfgOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if accessKeyID != "" && secretAccessKey != "" {
		cfgOpts = append(cfgOpts,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					accessKeyID,
					secretAccessKey,
					"",
				),
			),
		)
	}

	cfg, err := config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return err
	}

	sesClient = sesv2.NewFromConfig(cfg)
	return nil
}

func SendEmail(
	ctx context.Context,
	toEmails []string,
	subject string,
	body string,
) error {

	if sesClient == nil {
		return errors.New("SES client not initialized")
	}

	fromEmail := os.Getenv("NOTICE_EMAIL_ADDRESS")
	if fromEmail == "" {
		return errors.New("NOTICE_EMAIL_ADDRESS is not set")
	}

	_, err := sesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(fromEmail),
		Destination: &types.Destination{
			ToAddresses: toEmails,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Text: &types.Content{Data: aws.String(body)},
				},
			},
		},
	})

	return err
}
