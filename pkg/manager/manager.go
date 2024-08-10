package manager

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/starter-kit/pkg/cache"
	"github.com/mstgnz/starter-kit/pkg/conn"
	"github.com/mstgnz/starter-kit/pkg/response"
)

type Manager struct {
	DB        *conn.DB
	Mail      *response.Mail
	Cache     *cache.Cache
	Kafka     *conn.Kafka
	Redis     *conn.Redis
	Validator *validator.Validate
	SecretKey string
	QUERY     map[string]string
}

var instance *Manager

func Init() *Manager {
	if instance == nil {
		instance = &Manager{
			DB:        &conn.DB{},
			Cache:     &cache.Cache{},
			Kafka:     &conn.Kafka{},
			Redis:     &conn.Redis{},
			Validator: &validator.Validate{},
			// the secret key will change every time the application is restarted.
			SecretKey: "asdf1234", //RandomString(8),
			Mail: &response.Mail{
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
	}
	return instance
}
