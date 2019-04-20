package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/goamz/goamz/dynamodb"
	"github.com/goamz/goamz/sqs"

	"github.com/BurntSushi/toml"
	"github.com/hexaforce/goamz-example/example"
	"github.com/hexaforce/goamz-example/model"
)

type Config struct {
	Region   string
	DynamoDB example.DynamoDBConfig
	SQS      example.SQSConfig
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

	// clientD := example.NewDynamoDBClient(config.DynamoDB)
	// clientD.Init(config.Region)

	// table1 := clientD.GetTable(config.DynamoDB.TableName1)
	// example1Key := strconv.Itoa(example1.ExampleID)
	// example1Attribute := clientD.AttributeMapping(example1)
	//dynamodbExample(table1, example1Key, example1Attribute)

	clientQ := example.NewSQSClient(config.SQS)
	clientQ.Init(config.Region)

	queue1 := clientQ.GetQueue(config.SQS.QueueName1)

	receiveChan := make(chan sqs.Message)
	defer close(receiveChan)

	// Long polling.
	go func() {
		for {
			if resp, err := queue1.ReceiveMessage(10); err == nil { // Max is 10
				for _, message := range resp.Messages {
					receiveChan <- message
				}
			} else {
				log.Println(err)
				break
			}
		}
	}()

	// Message handling.
	go func() {
		for {
			select {
			case message, ok := <-receiveChan:
				if !ok {
					log.Println("receiveChan closed.")
					break
				}
				go func() {
					x, _ := json.Marshal(message.Body)
					log.Println(string(x))
					if _, err := queue1.DeleteMessage(&message); err == nil {
						log.Println("delete message.")
					} else {
						log.Println(err)
					}
				}()
			}
		}
	}()

	// Enqueue message.
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

func dynamodbExample(t *dynamodb.Table, k string, attributes []dynamodb.Attribute) {

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
