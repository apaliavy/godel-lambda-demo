package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/apaliavy/godel-lambda-demo/app/config"
	"github.com/apaliavy/godel-lambda-demo/app/model"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	f, err := os.Open("assets/out/users-2021-09-07T16:52:42+03:00.json")
	if err != nil {
		logger.WithError(err).Fatal("failed to read users file")
		return
	}

	defer f.Close()

	dbConf, err := config.LoadDBSettings()
	if err != nil {
		logger.WithError(err).Fatal("failed to load DB settings")
		return
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbConf.BuildDNS())
	if err != nil {
		logger.WithError(err).Fatal("failed to establish DB connection")
		return
	}

	users := make([]model.User, 0)
	if err := json.NewDecoder(f).Decode(&users); err != nil {
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
