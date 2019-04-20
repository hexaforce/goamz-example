package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goamz/goamz/dynamodb"
	"github.com/goamz/goamz/s3"
	"github.com/goamz/goamz/sqs"

	"github.com/BurntSushi/toml"
	"github.com/hexaforce/goamz-example/example"
	"github.com/hexaforce/goamz-example/model"
)

type Config struct {
	Region   string
	DynamoDB example.DynamoDBConfig
	SQS      example.SQSConfig
	S3       example.S3Config
}

func main() {

	var config Config
	if _, err := toml.DecodeFile("./env/local.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	example1 := model.Example{
		ExampleID:   1,
		CustomerID:  1111111,
		ProductID:   100,
		ProductItem: []string{"AAA", "BBB", "CCC"},
		OrderDate:   time.Now(),
	}

	// ################## S3 ########################

	// deprecated!! goamz -> signature v2
	// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingAWSSDK.html#UsingAWSSDK-sig2-deprecation

	clientS := example.News3Client(config.S3)
	clientS.Init(config.Region)
	bucket1 := clientS.GetBucket(config.S3.BucketName1)

	data, _ := json.Marshal(example1)
	// Put example
	if err := bucket1.Put("example.json", data, "text/plain", s3.BucketOwnerFull, s3.Options{}); err == nil {
		// Get example
		if content, err := bucket1.Get("example.json"); err == nil {
			log.Println(content)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	// File upload.
	if file, err := os.Open("./example.txt"); err == nil {
		defer file.Close()
		if fileInfo, err := file.Stat(); err == nil {
			err = bucket1.PutReader("example.txt", file, fileInfo.Size(), "text/plain", s3.PublicRead, s3.Options{})
			if err != nil {
				panic(err.Error())
			}
		}
	}

	// ################## DynamoDB ########################
	clientD := example.NewDynamoDBClient(config.DynamoDB)
	clientD.Init(config.Region)

	table1 := clientD.GetTable(config.DynamoDB.TableName1)
	example1Key := strconv.Itoa(example1.ExampleID)
	example1Attribute := clientD.AttributeMapping(example1)
	dynamodbPutGetExample(table1, example1Key, example1Attribute)

	// ################## SQS ########################

	// goamz ->  FIFO queue Notsupported!!

	clientQ := example.NewSQSClient(config.SQS)
	clientQ.Init(config.Region)
	queue1 := clientQ.GetQueue(config.SQS.QueueName1)

	receiveChan := make(chan sqs.Message)
	defer close(receiveChan)

	// 1 Long polling.
	go sqsExampleLongpolling(queue1, receiveChan)

	// 2 Message handling.
	go sqsExampleMessagehandling(queue1, receiveChan)

	// 3 Message enqueue.
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			example1.OrderDate = t
			x, _ := json.Marshal(example1)
			queue1.SendMessage(string(x))
		}
	}

}
func sqsExampleLongpolling(queue *sqs.Queue, receiveChan chan sqs.Message) {
	for {
		if resp, err := queue.ReceiveMessage(10); err == nil { // Max is 10
			for _, message := range resp.Messages {
				receiveChan <- message
			}
		} else {
			log.Println(err)
			break
		}
	}
}

func sqsExampleMessagehandling(queue *sqs.Queue, receiveChan chan sqs.Message) {
	for {
		select {
		case message, ok := <-receiveChan:
			if !ok {
				log.Println("receiveChan closed.")
				break
			}
			go func() {
				log.Println(message.Body)
				if _, err := queue.DeleteMessage(&message); err == nil {
					log.Println("delete message.")
				} else {
					log.Println(err)
				}
			}()
		}
	}
}

func dynamodbPutGetExample(t *dynamodb.Table, k string, attributes []dynamodb.Attribute) {

	// PutItem example
	if ok, err := t.PutItem(k, "", attributes); ok {

		key := &dynamodb.Key{
			HashKey:  k,
			RangeKey: "",
		}

		// GetItem example
		result, _ := t.GetItem(key)
		x, _ := json.Marshal(result)
		log.Println(string(x))

	} else {
		log.Println(err)
	}
}
