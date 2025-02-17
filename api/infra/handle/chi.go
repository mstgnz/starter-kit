package handle

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/starter-kit/api/infra/config"
	"github.com/mstgnz/starter-kit/api/infra/response"
	"github.com/mstgnz/starter-kit/api/infra/validate"
)

type Request any
type Response any

func Handle[Req Request, Res Response](handler func(ctx context.Context, req *Req) Res) config.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req Req

		// body parser
		if err := response.ReadJSON(w, r, &req); err != nil {
			return response.WriteJSON(w, http.StatusBadRequest, response.Response{Code: http.StatusBadRequest, Success: false, Message: err.Error()})
		}

		// params parser
		if rctx := chi.RouteContext(r.Context()); rctx != nil {
			if err := parseParams(rctx, &req); err != nil {
				return response.WriteJSON(w, http.StatusBadRequest, response.Response{Code: http.StatusBadRequest, Success: false, Message: err.Error()})
			}
		}

		// query parser
		if err := parseQuery(r.URL.Query(), &req); err != nil {
			return response.WriteJSON(w, http.StatusBadRequest, response.Response{Code: http.StatusBadRequest, Success: false, Message: err.Error()})
		}

		// header parser
		if err := parseHeader(r.Header, &req); err != nil {
			return response.WriteJSON(w, http.StatusBadRequest, response.Response{Code: http.StatusBadRequest, Success: false, Message: err.Error()})
		}

		// validation
		err := validate.Validate(req)
		if err != nil {
			return response.WriteJSON(w, http.StatusUnprocessableEntity, response.Response{Code: http.StatusUnprocessableEntity, Success: false, Message: err.Error()})
		}

		// timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		res := handler(ctx, &req)

		result, ok := any(res).(response.Response)
		if !ok {
			return response.WriteJSON(w, http.StatusInternalServerError, response.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "Invalid response type",
			})
		}

		return response.WriteJSON(w, result.Code, result)
	}
}

func parseParams(rctx *chi.Context, req interface{}) error {
	v := reflect.ValueOf(req).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("param"); tag != "" {
			if value := rctx.URLParam(tag); value != "" {
				if err := setFieldValue(v.Field(i), value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parseQuery(query url.Values, req interface{}) error {
	v := reflect.ValueOf(req).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("query"); tag != "" {
			if value := query.Get(tag); value != "" {
				if err := setFieldValue(v.Field(i), value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func parseHeader(header http.Header, req interface{}) error {
	v := reflect.ValueOf(req).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("header"); tag != "" {
			if value := header.Get(tag); value != "" {
				if err := setFieldValue(v.Field(i), value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}
