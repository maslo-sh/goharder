package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

const (
	region = "eu-central-1"
)

type CloudWatchConfiguration struct {
	Client        *cloudwatchlogs.Client
	LogGroupName  string
	LogStreamName string
}

func NewCloudWatchConfiguration(groupName, streamName string) *CloudWatchConfiguration {
	// Initialize AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("failed to load configuration, %v", err)
	}

	log.Printf("CloudWatch logging configuration completed successfully")

	return &CloudWatchConfiguration{
		cloudwatchlogs.NewFromConfig(cfg),
		groupName,
		streamName}
}

func (c *CloudWatchConfiguration) InitLogStore() {
	err := c.CreateLogGroup()
	if err != nil {
		fmt.Printf("error when creating log group: %v\n", err)
	}

	err = c.CreateLogStream()
	if err != nil {
		fmt.Printf("error when creating log stream: %v\n", err)
	}
}

func (c *CloudWatchConfiguration) CreateLogGroup() error {
	// Create or update the log group
	_, err := c.Client.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(c.LogGroupName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudWatchConfiguration) CreateLogStream() error {
	// Create or update the log stream
	_, err := c.Client.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(c.LogGroupName),
		LogStreamName: aws.String(c.LogStreamName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudWatchConfiguration) SendLog(message string) error {
	_, err := c.Client.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(c.LogGroupName),
		LogStreamName: aws.String(c.LogStreamName),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(nowMillis()),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func nowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
