package example

import (
	"log"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/sqs"
)

type SQSConfig struct {
	QueueName1 string
	QueueName2 string
}

type SQSClient struct {
	Config SQSConfig
	auth   aws.Auth
	region aws.Region
	server *sqs.SQS
}

func NewSQSClient(c SQSConfig) *SQSClient {
	return &SQSClient{
		Config: c,
	}
}

//Init
func (c *SQSClient) Init(region string) {

	var err error
	// The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
	if c.auth, err = aws.EnvAuth(); err != nil {
		// $HOME/.aws/credentials. The AWS_PROFILE environment variables is used to select the profile.
		if c.auth, err = aws.SharedAuth(); err != nil {
			if region != "localhost" {
				panic(err)
			}
		}
	}

	var ok bool
	if c.region, ok = aws.Regions[region]; !ok {
		c.region = aws.APNortheast
		c.region.Name = "localhost"
		c.region.SQSEndpoint = "http://localhost:9324"
	}

	c.server = sqs.New(c.auth, c.region)
}

func (c *SQSClient) GetQueue(queueName string) *sqs.Queue {
	if q, err := c.server.GetQueue(queueName); err == nil {
		return q
	} else {
		log.Println(err)
	}
	if q, err := c.server.CreateQueueWithAttributes(
		c.Config.QueueName1,
		map[string]string{"ReceiveMessageWaitTimeSeconds": "20"}, // for long-polling
	); err == nil {
		return q
	} else {
		log.Println(err)
	}
	return nil
}
