package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Headers map[string]string

type Options struct {
	Method  string
	Headers Headers
	Body    string
}

type Object map[string]interface{}

type Response struct {
	Body     io.ReadCloser
	BodyUsed bool

	Headers    map[string][]string
	Status     int
	StatusText string
	URL        string
}

var client *http.Client

func init() {
	client = &http.Client{}
}

func (response *Response) Text() (body string, err error) {
	if response.BodyUsed != false {
		return "", errors.New("response body already used")
	}

	plain, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	response.BodyUsed = true

	return string(plain), nil
}

func (response *Response) JSON() (Body Object, err error) {
	if response.BodyUsed != false {
		return Object{}, errors.New("response body already used")
	}

	plain, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return Object{}, err
	}

	response.BodyUsed = true

	var object Object

	err = json.Unmarshal(plain, &object)

	if err != nil {
		return Object{}, err
	}

	return object, nil
}

func Fetch(address string, options Options) (body Response, err error) {
	if options.Method == "" {
		options.Method = "GET"
	}

	request, err := http.NewRequest(options.Method, address, bytes.NewBuffer([]byte(options.Body)))

	if err != nil {
		return Response{}, err
	}

	if len(options.Headers) > 0 {
		for key, value := range options.Headers {
			request.Header.Add(key, value)
		}
	}

	response, err := client.Do(request)

	if err != nil {
		return Response{}, err
	}

	return Response{
		Bodyy:     response.Body,
		BodyUsed: false,

		Headers:    response.Header,
		Status:     response.StatusCode,
		StatusText: response.Status,
		URL:        address,
	}, nil
}

func SetClient(fresh *http.Client) {
	client = fresh
}
