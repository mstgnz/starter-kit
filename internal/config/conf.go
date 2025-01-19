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

	"github.com/cemilsahin/arabamtaksit/internal/conn"
	"github.com/cemilsahin/arabamtaksit/internal/response"
	"github.com/cemilsahin/arabamtaksit/pkg/mstgnz/cache"
	"github.com/cemilsahin/arabamtaksit/pkg/mstgnz/gobuilder"
	"github.com/cemilsahin/arabamtaksit/pkg/mstgnz/mail"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

var (
	once     sync.Once
	mu       sync.Mutex
	instance *config
)

// context key type
type CKey string

type config struct {
	DB        *conn.DB
	Mail      *mail.Mail
	Cron      *cron.Cron
	Cache     *cache.Cache
	Redis     *conn.Redis
	Builder   *gobuilder.GoBuilder
	Validator *validator.Validate
	SecretKey string
	QUERY     map[string]string
	Running   int
	Shutting  bool
	LangCodes []string
}

func App() *config {
	once.Do(func() {
		instance = &config{
			DB:        &conn.DB{},
			Redis:     &conn.Redis{},
			Cron:      cron.New(),
			Cache:     cache.NewCache(),
			Builder:   gobuilder.NewGoBuilder(gobuilder.Postgres),
			Validator: validator.New(),
			// the secret key will change every time the application is restarted.
			SecretKey: os.Getenv("JWT_SECRET"), //RandomString(8),
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
		instance.Redis.ConnectRedis()
	})
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

func CalcPaginate(page, total, limit int) response.Paginate {
	size := (total + limit - 1) / limit
	current := Clamp(page, 1, size)
	return response.Paginate{
		Total:    total,
		Size:     size,
		Row:      limit,
		Current:  current,
		Previous: Clamp(current-1, 1, size),
		Next:     Clamp(current+1, 1, size),
	}
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
