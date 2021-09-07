package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/apaliavy/godel-lambda-demo/app/model"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	f, err := os.Open("assets/in/users.csv")
	if err != nil {
		logger.WithError(err).Fatal("failed to read users csv file")
		return
	}

	defer f.Close()

	users, err := model.UsersFromCSV(f)
	if err != nil {
		logger.WithError(err).Fatal("failed to convert users")
		return
	}

	outname := fmt.Sprintf("assets/out/users-%s.json", time.Now().Format(time.RFC3339))
	f2, err := os.OpenFile(outname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 06444)
	if err != nil {
		logger.WithError(err).Fatal("failed to write users json file")
		return
	}

	defer f2.Close()

	if err := json.NewEncoder(f2).Encode(users); err != nil {
		logger.WithError(err).Fatal("failed to encode users")
		return
	}
}
