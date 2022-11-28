package storage

import (
	"os"

	"github.com/zikster3262/shared-lib/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateNewClient() *s3.Client {

	minioURL := os.Getenv("S3_URL")
	region := os.Getenv("AWS_REGION")
	user := os.Getenv("AWS_ACCESS_KEY")
	pass := os.Getenv("AWS_SECRET_KEY")

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			SigningRegion:     region,
			URL:               minioURL,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:                      region,
		Credentials:                 credentials.NewStaticCredentialsProvider(user, pass, ""),
		EndpointResolverWithOptions: resolver,
	}

	utils.LogWithInfo("connected to aws bucket", "aws-s3")

	return s3.NewFromConfig(cfg)

}
