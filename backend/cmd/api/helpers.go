package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"crossfitbox.booking.system/internal/cookies"
	"crossfitbox.booking.system/internal/data"
	"crossfitbox.booking.system/internal/validator"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// Retrieve the "id" URL parameter from the current request context,
// then convert it to an integer and return it. If the operation isn't successful, return 0 and an error
func (app *application) readIDParam(r *http.Request) (*uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, errors.New("invalid id parameter")
	}

	return &id, nil
}

type envelope map[string]interface{}

// writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// readJSON() helper for decoding JSON into target destination and error handling
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.Trim(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field: %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only a single JSON value")
	}
	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// The readCSV() helper reads a string value from the query string and the splits it
// into a slice on the comma character. If no matching could be found, uit returns
// provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query and converts it to an integer
// before returning. If no matching key could be found, it returns the provided default value.
// If the value couldn't be converted to an integer, then error message is provided to Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) storeInRedis(prefix string, hash string, userID uuid.UUID, expiration time.Duration) error {
	ctx := context.Background()
	err := app.redisClient.Set(
		ctx,
		fmt.Sprintf("%s%s", prefix, userID),
		hash,
		expiration,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) getFromRedis(key string) (*string, error) {
	ctx := context.Background()

	hash, err := app.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &hash, nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

func (app *application) extractParamsFromSession(r *http.Request) (*data.UserID, *int, error) {
	gobEncodedValue, err := cookies.ReadEncrypted(r, "sessionid", app.config.secret.secretKey)
	if err != nil {
		var errorData error
		var status int
		switch {
		case errors.Is(err, http.ErrNoCookie):
			status = http.StatusUnauthorized
			errorData = errors.New("you are not authorized to access this resource")
		case errors.Is(err, cookies.ErrInvalidValue):
			app.logger.PrintError(err, nil)
			status = http.StatusBadRequest
			errorData = errors.New("invalid cookie")
		default:
			status = http.StatusInternalServerError
			errorData = errors.New("something happened getting your cookie data")
		}
		return nil, &status, errorData
	}

	var userID data.UserID

	reader := strings.NewReader(gobEncodedValue)
	if err := gob.NewDecoder(reader).Decode(&userID); err != nil {
		status := http.StatusInternalServerError
		return nil, &status, errors.New("something happened getting your cookie data")
	}

	return &userID, nil, nil
}
