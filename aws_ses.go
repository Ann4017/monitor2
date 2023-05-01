package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/ini.v1"
)

type C_ses struct {
	s_region     string
	s_access_Key string
	s_secret_key string
	pc_client    *ses.Client
	s_sender     string
	s_recipient  []string
	s_subject    string
	s_body       string
}

func (c *C_ses) Load_aws_config(_s_ini_file_path, _s_section string) error {
	file, err := ini.Load(_s_ini_file_path)
	if err != nil {
		return err
	}

	section, err := file.GetSection(_s_section)
	if err == nil {
		c.Ses_aws_config(section)
		return nil
	}

	section, err = file.NewSection(_s_section)
	if err != nil {
		return err
	}

	section.NewKey("region", "value")
	section.NewKey("access_Key", "value")
	section.NewKey("secret_key", "value")

	fmt.Printf("New section '%s' created in the config file with default values.\nPlease update the values.\n", _s_section)

	return nil
}

func (c *C_ses) Ses_aws_config(_section *ini.Section) {
	c.s_region = _section.Key("region").String()
	c.s_access_Key = _section.Key("access_Key").String()
	c.s_secret_key = _section.Key("secret_key").String()
}

func (c *C_ses) Write_email(_s_sender string, _s_recipient []string, _s_subject, _s_body string) {
	c.s_sender = _s_sender
	c.s_recipient = _s_recipient
	c.s_subject = _s_subject
	c.s_body = _s_body
}

func (c *C_ses) Set_config() error {
	cred := credentials.NewStaticCredentialsProvider(c.s_access_Key, c.s_secret_key, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion(c.s_region), config.WithCredentialsProvider(cred))
	if err != nil {
		return err
	}

	c.pc_client = ses.NewFromConfig(cfg)

	return nil
}

func (c *C_ses) Send_email(_pc_client *ses.Client, _s_sender string, _s_recipient []string, _s_subject string, _s_body string) error {
	input := ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: _s_recipient,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(_s_subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(_s_body),
				},
			},
		},
		Source: aws.String(_s_sender),
	}

	_, err := c.pc_client.SendEmail(context.Background(), &input)
	if err != nil {
		return err
	}

	return nil
}
