package fetch

import (
	"fmt"
	"net/http"
	"testing"
)

var hosting bool

func host() {
	if hosting == true {
		return
	}

	http.HandleFunc("/plain", func(writer http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(writer, "Hello World")
	})

	http.HandleFunc("/json", func(writer http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(writer, "{ \"Hello\": \"World\" }\n")
	})

	hosting = true

	http.ListenAndServe(":8080", nil)
}

func TestSimple(test *testing.T) {
	go host()

	response, err := Fetch("http://localhost:8080/plain", Options{})

	test.Run("get response", func(test *testing.T) {
		if err != nil {
			test.Errorf("HTTP response failed: %v\n", err)
			return
		} else {
			test.Logf("HTTP response received\n")
		}
	})

	_, err = response.Text()

	test.Run("get response body", func(test *testing.T) {
		if err != nil {
			test.Errorf("HTTP response body failed: %v\n", err)
			return
		} else {
			test.Logf("HTTP response body read\n")
		}
	})

	_, err = response.Text()

	test.Run("close response body", func(test *testing.T) {
		if err == nil {
			test.Errorf("HTTP response body failed to close")
			return
		} else {
			test.Logf("HTTP response body closed\n")
		}
	})
}

func TestAdvanced(test *testing.T) {
	go host()

	response, err := Fetch("http://localhost:8080/json", Options{
		Method: "POST",
		Headers: Headers{
			"User-Agent": "example",
		},
		Body: "Hello World",
	})

	test.Run("get response", func(test *testing.T) {
		if err != nil {
			test.Errorf("HTTP response failed: %v\n", err)
			return
		} else {
			test.Logf("HTTP response received\n")
		}
	})

	_, err = response.JSON()

	test.Run("get response body", func(test *testing.T) {
		if err != nil {
			test.Errorf("HTTP response body failed: %v\n", err)
			return
		} else {
			test.Logf("HTTP response body read\n")
		}
	})

	_, err = response.JSON()

	test.Run("get response body", func(test *testing.T) {
		if err == nil {
			test.Errorf("HTTP response body failed: %v\n", err)
			return
		} else {
			test.Logf("HTTP response body read\n")
		}
	})
}

func BenchmarkGet(benchmark *testing.B) {
	go host()

	for index := 0; index < benchmark.N; index++ {
		Fetch("http://localhost:8080/plain", Options{})
	}
}
