package fetch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Headers map[string]string

type Options struct {
	Method  string
	Headers Headers
	Body    string
}

type Response[T any] struct {
	Body T

	Headers    map[string][]string
	Status     int
	StatusText string
	URL        string
}

var client *http.Client

func init() {
	client = &http.Client{}
}

func Fetch[T any](address string, options Options) (Response[T], error) {
	if options.Method == "" {
		options.Method = "GET"
	}

	request, err := http.NewRequest(options.Method, address, bytes.NewBuffer([]byte(options.Body)))

	if err != nil {
		return Response[T]{}, err
	}

	if len(options.Headers) > 0 {
		for key, value := range options.Headers {
			request.Header.Add(key, value)
		}
	}

	response, err := client.Do(request)

	if err != nil {
		return Response[T]{}, err
	}

	plain, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return Response[T]{}, err
	}

	var body T

	switch any(body).(type) {
	case string:
		body = any(string(plain)).(T)

	case []byte:
		body = any(plain).(T)

	default:
		//TODO: check if `T` is non-json

		err = json.Unmarshal(plain, &body)

		if err != nil {
			return Response[T]{}, err
		}
	}

	return Response[T]{
		Body: body,

		Headers:    response.Header,
		Status:     response.StatusCode,
		StatusText: response.Status,
		URL:        address,
	}, nil
}
