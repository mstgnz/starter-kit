package config

import (
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/starter-kit/api/internal/conn"
	"github.com/mstgnz/starter-kit/api/pkg/mstgnz/cache"
	"github.com/mstgnz/starter-kit/api/pkg/mstgnz/gobuilder"
	"github.com/mstgnz/starter-kit/api/pkg/mstgnz/mail"
)

type CKey string

type Config struct {
	DB        *conn.DB
	Mail      *mail.Mail
	Cache     *cache.Cache
	Builder   *gobuilder.GoBuilder
	Kafka     *conn.Kafka
	Redis     *conn.Redis
	Validator *validator.Validate
	SecretKey string
	Token     string
	QUERY     map[string]string
	Lang      string
	Langs     []string
	Routes    map[string]map[string]string
	Running   int
	Shutting  bool
}

var (
	mu       sync.Mutex
	instance *Config
)

func App() *Config {
	if instance == nil {
		instance = &Config{
			DB:        &conn.DB{},
			Cache:     cache.NewCache(),
			Builder:   gobuilder.NewGoBuilder(gobuilder.Postgres),
			Kafka:     &conn.Kafka{},
			Redis:     &conn.Redis{},
			Validator: validator.New(),
			// the secret key will change every time the application is restarted.
			SecretKey: "asdf1234", //RandomString(8),
			Lang:      "tr",
			Langs:     []string{"tr", "en"},
			Routes:    make(map[string]map[string]string),
			Mail: &mail.Mail{
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

func ShuttingWrapper(fn func()) {
	if !App().Shutting {
		fn()
	}
}

func IncrementRunning() {
	mu.Lock()
	App().Running++
	mu.Unlock()
}

func DecrementRunning() {
	mu.Lock()
	App().Running--
	mu.Unlock()
}

func StructToMap(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result
}

func GetIntQuery(r *http.Request, name string) int {
	pageStr := r.URL.Query().Get(name)
	if page, err := strconv.Atoi(pageStr); err == nil {
		return int(math.Abs(float64(page)))
	}
	return 1
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ActiveClass(a, b int) string {
	active := ""
	if a == b {
		active = "active"
	}
	return active
}

func WriteBody(r *http.Request) {
	if body, err := io.ReadAll(r.Body); err != nil {
		log.Println("WriteBody: ", err)
	} else {
		log.Println("WriteBody: ", string(body))
	}
}
