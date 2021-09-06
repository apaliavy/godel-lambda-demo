package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/apaliavy/godel-lambda-demo/app/model"
)

func main() {
	// make the handler available for RPC call by AWS Lambda
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// create AWS session
	sess := session.Must(session.NewSession())
	// create S3 client
	svc := s3.New(sess)

	logger := logrus.New()

	for _, record := range s3Event.Records {
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Key:    aws.String(record.S3.Object.Key),
			Bucket: aws.String(record.S3.Bucket.Name),
		})
		if err != nil {
			logger.WithError(err).Fatal("failed to get object from S3 bucket")
			return
		}

		users, err := model.UsersFromXLS(obj.Body)
		if err != nil {
			logger.WithError(err).Fatal("failed to read users csv file")
			return
		}

		buf := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(buf).Encode(users); err != nil {
			logger.WithError(err).Fatal("failed to convert users to JSON")
			return
		}

		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket:               aws.String("godel-demo-outbound"),
			Key:                  aws.String(fmt.Sprintf("users-%s.json", time.Now().Format(time.RFC3339))),
			Body:                 bytes.NewReader(buf.Bytes()),
			ContentLength:        aws.Int64(int64(buf.Len())),
			ContentType:          aws.String("application/json"),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
			ACL:                  aws.String("private"),
		})

		if err != nil {
			logger.WithError(err).Fatal("failed to save users file")
			return
		}

		logger.Info("successfully processed users info")
	}
}
