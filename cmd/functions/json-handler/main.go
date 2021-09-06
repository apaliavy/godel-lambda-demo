package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"github.com/apaliavy/godel-lambda-demo/app/config"
	"github.com/apaliavy/godel-lambda-demo/app/model"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) {
	// create AWS session
	sess := session.Must(session.NewSession())
	// create S3 client
	svc := s3.New(sess)

	logger := logrus.New()

	dbConf, err := config.LoadDBSettings()
	if err != nil {
		logger.WithError(err).Fatal("failed to load DB settings")
		return
	}

	conn, err := pgx.Connect(ctx, dbConf.BuildDNS())
	if err != nil {
		logger.WithError(err).Fatal("failed to establish DB connection")
		return
	}

	for _, record := range s3Event.Records {
		key, err := url.PathUnescape(record.S3.Object.Key)
		if err != nil {
			logger.WithField("key", record.S3.Object.Key).WithError(err).Error("failed to decode object key")
			continue
		}

		obj, err := svc.GetObject(&s3.GetObjectInput{
			Key:    aws.String(key),
			Bucket: aws.String(record.S3.Bucket.Name),
		})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"key":    record.S3.Object.Key,
				"bucket": record.S3.Bucket.Name,
			}).WithError(err).Fatal("failed to get object from S3 bucket")
			return
		}

		users := make([]model.User, 0)
		if err := json.NewDecoder(obj.Body).Decode(&users); err != nil {
			logger.WithError(err).Fatal("failed to decode users info")
			return
		}

		for _, u := range users {
			_, err := conn.Exec(
				ctx,
				`INSERT INTO users (firstname, lastname, birthday, active) VALUES ($1, $2, $3, $4)`,
				u.Firstname, u.Lastname, u.Birthday, u.Active,
			)
			if err != nil {
				logger.WithField("user", u).WithError(err).Error("failed to save user to DB")
				continue
			}
		}

		logger.Info("successfully saved users")
	}
}
