package example

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

type S3Config struct {
	BucketName1 string
	BucketName2 string
}

type S3Client struct {
	Config S3Config
	auth   aws.Auth
	region aws.Region
	server *s3.S3
}

func News3Client(c S3Config) *S3Client {
	return &S3Client{
		Config: c,
	}
}

//Init
func (c *S3Client) Init(region string) {

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
		c.region.S3Endpoint = "http://localhost:9000"
		c.region.S3BucketEndpoint = "http://localhost:9000"

		c.auth.AccessKey = "AKIAIOSFODNN7EXAMPLE"
		c.auth.SecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	}

	c.server = s3.New(c.auth, c.region)

}

func (c *S3Client) GetBucket(bucketName string) *s3.Bucket {
	return c.server.Bucket(bucketName)
}
