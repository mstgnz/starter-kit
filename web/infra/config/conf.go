package config

import (
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
)

type CKey string

type Config struct {
	IsAuth    bool
	Validator *validator.Validate
	Token     string
	Lang      string
	Langs     []string
	Routes    map[string]map[string]string
}

var (
	mu       sync.Mutex
	instance *Config
)

func App() *Config {
	if instance == nil {
		instance = &Config{
			Validator: validator.New(),
			Lang:      "tr",
			Langs:     []string{"tr", "en"},
			Routes:    make(map[string]map[string]string),
		}
	}
	return instance
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
