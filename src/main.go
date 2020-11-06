package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)

		svc := cloudfront.New(session.New())

		params := &cloudfront.CreateInvalidationInput{
			DistributionId: aws.String(os.Getenv("DistributionId")),
			InvalidationBatch: &cloudfront.InvalidationBatch{
				CallerReference: aws.String(record.EventTime.String()),
				Paths: &cloudfront.Paths{
					Quantity: aws.Int64(1),
					Items: []*string{
						aws.String("/" + s3.Object.Key),
					},
				},
			},
		}
		resp, err := svc.CreateInvalidation(params)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(resp)

	}

}

func main() {
	lambda.Start(handler)
}
