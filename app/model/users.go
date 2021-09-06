package model

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Birthday  string `json:"birthday"`
	Active    bool   `json:"active"`
}

func UsersFromXLS(r io.Reader) ([]User, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	// Get all the rows in the vegan section.
	rows, err := f.GetRows("Sheet1")
	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		id, err := strconv.ParseFloat(row[0], 32)
		if err != nil {
			fmt.Printf("got an error parsing ID: %s\n", err)
			continue
		}

		active, err := strconv.ParseBool(row[4])
		if err != nil {
			fmt.Printf("got an error parsing Active: %s\n", err)
			continue
		}

		users = append(users, User{
			ID:        int(id),
			Firstname: row[1],
			Lastname:  row[2],
			Birthday:  row[3],
			Active:    active,
		})
	}

	return users, nil
}

func UsersFromCSV(r io.Reader) ([]User, error) {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for idx, row := range records {
		if idx == 0 {
			continue
		}

		id, err := strconv.ParseFloat(row[0], 32)
		if err != nil {
			fmt.Printf("got an error parsing ID: %s\n", err)
			continue
		}

		active, err := strconv.ParseBool(row[4])
		if err != nil {
			fmt.Printf("got an error parsing Active: %s\n", err)
			continue
		}

		users = append(users, User{
			ID:        int(id),
			Firstname: row[1],
			Lastname:  row[2],
			Birthday:  row[3],
			Active:    active,
		})
	}

	return users, nil
}
