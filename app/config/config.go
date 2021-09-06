package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Database struct {
	User     string `env:"user"`
	Password string `env:"password"`
	Host     string `env:"host"`
	Name     string `env:"name"`
	Port     int    `env:"port"`
	Options  string `env:"options"`
}

func (db *Database) BuildDNS() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.Options)
}

func LoadDBSettings() (*Database, error) {
	dbConf := &Database{}
	if err := env.Parse(dbConf); err != nil {
		return nil, err
	}

	return dbConf, nil
}
