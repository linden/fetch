package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Headers map[string]any

type Options[T any] struct {
	Method  string
	Headers Headers
	Body    T
}

type Response[T any] struct {
	Body T

	Headers    Headers
	Status     int
	StatusText string
	URL        string
}

type Empty struct{}

var Client = &http.Client{}

func Fetch[T any, B any](address string, options Options[B]) (Response[T], error) {
	if options.Method == "" {
		options.Method = "GET"
	}

	var payload []byte

	switch value := any(options.Body).(type) {
	case []byte:
		payload = value

	case string:
		payload = []byte(value)

	case Empty:
		// do nothing

	default:
		// TODO: check if `T` is non-json
		var err error

		payload, err = json.Marshal(value)

		if err != nil {
			return Response[T]{}, err
		}
	}

	request, err := http.NewRequest(options.Method, address, bytes.NewBuffer(payload))

	if err != nil {
		return Response[T]{}, err
	}

	if len(options.Headers) > 0 {
		for key, unknown := range options.Headers {
			switch value := unknown.(type) {
			case string:
				request.Header.Add(key, value)

			case []string:
				for _, cursor := range value {
					request.Header.Add(key, cursor)
				}

			default:
				return Response[T]{}, fmt.Errorf("%T is not a supported header type", value)
			}
		}
	}

	response, err := Client.Do(request)

	if err != nil {
		return Response[T]{}, err
	}

	plain, err := io.ReadAll(response.Body)

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
		// TODO: check if `T` is non-json
		err = json.Unmarshal(plain, &body)

		if err != nil {
			return Response[T]{}, err
		}
	}

	headers := make(map[string]any)

	for key, value := range response.Header {
		if len(value) == 1 {
			headers[key] = any(value[0])
		} else {
			headers[key] = any(value)
		}
	}

	return Response[T]{
		Body: body,

		Headers:    headers,
		Status:     response.StatusCode,
		StatusText: response.Status,
		URL:        address,
	}, nil
}
