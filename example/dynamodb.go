package example

import (
	"strconv"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"github.com/hexaforce/goamz-example/model"
)

type DynamoDBConfig struct {
	TableName1 string
	TableName2 string
}

type DynamoDBClient struct {
	Config DynamoDBConfig
	auth   aws.Auth
	region aws.Region
	server *dynamodb.Server
}

func NewDynamoDBClient(c DynamoDBConfig) *DynamoDBClient {
	return &DynamoDBClient{
		Config: c,
	}
}

//Init
func (c *DynamoDBClient) Init(region string) {

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
		c.region.DynamoDBEndpoint = "http://localhost:8000"
	}

	c.server = &dynamodb.Server{
		Auth:   c.auth,
		Region: c.region,
	}

}

func (c *DynamoDBClient) GetTable(tableName string) *dynamodb.Table {
	if tables, err := c.server.ListTables(); err == nil {
		for _, table := range tables {
			if tableName == table {
				if tableDescription, err := c.server.DescribeTable(table); err == nil {
					if pk, err := tableDescription.BuildPrimaryKey(); err == nil {
						return c.server.NewTable(tableDescription.TableName, pk)
					}
				}
			}
		}
	}
	return nil
}

func (c *DynamoDBClient) AttributeMapping(example model.Example) (attributes []dynamodb.Attribute) {
	attributes = make([]dynamodb.Attribute, 5)
	attributes[0] = dynamodb.Attribute{
		Type:  dynamodb.TYPE_NUMBER,
		Name:  "example_id",
		Value: strconv.Itoa(example.ExampleID),
	}
	attributes[1] = dynamodb.Attribute{
		Type:  dynamodb.TYPE_NUMBER,
		Name:  "customer_id",
		Value: strconv.Itoa(example.CustomerID),
	}
	attributes[2] = dynamodb.Attribute{
		Type:  dynamodb.TYPE_NUMBER,
		Name:  "product_id",
		Value: strconv.Itoa(example.ProductID),
	}
	attributes[3] = dynamodb.Attribute{
		Type:      dynamodb.TYPE_STRING_SET,
		Name:      "product_item",
		SetValues: example.ProductItem,
	}
	attributes[4] = dynamodb.Attribute{
		Type:  dynamodb.TYPE_STRING,
		Name:  "order_date",
		Value: example.OrderDate.Format("2006-01-02T15:04:05.999999"),
	}
	return
}
