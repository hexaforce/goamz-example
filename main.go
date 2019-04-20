package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/goamz/goamz/dynamodb"

	"github.com/BurntSushi/toml"
	"github.com/hexaforce/goamz-example/example"
	"github.com/hexaforce/goamz-example/model"
)

type Config struct {
	Region   string
	DynamoDB example.DynamoDBConfig
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

	client := example.NewDynamoDBClient(config.DynamoDB)
	client.Init(config.Region)

	table1 := client.GetTable(config.DynamoDB.TableName1)
	example1Key := strconv.Itoa(example1.ExampleID)
	example1Attribute := client.AttributeMapping(example1)
	dynamodbExample(table1, example1Key, example1Attribute)

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
