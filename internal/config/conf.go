package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/starter-kit/internal/conn"
	"github.com/mstgnz/starter-kit/pkg/mstgnz"
)

type CKey string

type Config struct {
	DB        *conn.DB
	Mail      *mstgnz.Mail
	Cache     *mstgnz.Cache
	Builder   *mstgnz.GoBuilder
	Kafka     *conn.Kafka
	Redis     *conn.Redis
	Validator *validator.Validate
	SecretKey string
	QUERY     map[string]string
	Lang      string
	Langs     []string
	Routes    map[string]map[string]string
}

var instance *Config

func App() *Config {
	if instance == nil {
		instance = &Config{
			DB:        &conn.DB{},
			Cache:     mstgnz.NewCache(),
			Builder:   mstgnz.NewGoBuilder(mstgnz.Postgres),
			Kafka:     &conn.Kafka{},
			Redis:     &conn.Redis{},
			Validator: validator.New(),
			// the secret key will change every time the application is restarted.
			SecretKey: "asdf1234", //RandomString(8),
			Lang:      "tr",
			Langs:     []string{"tr", "en"},
			Routes:    make(map[string]map[string]string),
			Mail: &mstgnz.Mail{
				From: os.Getenv("MAIL_FROM"),
				Name: os.Getenv("MAIL_FROM_NAME"),
				Host: os.Getenv("MAIL_HOST"),
				Port: os.Getenv("MAIL_PORT"),
				User: os.Getenv("MAIL_USER"),
				Pass: os.Getenv("MAIL_PASS"),
			},
		}
		// Connect to Postgres DB
		instance.DB.ConnectDatabase()
		//instance.Kafka.ConnectKafka()
		//instance.Redis.ConnectRedis()
	}
	return instance
}
