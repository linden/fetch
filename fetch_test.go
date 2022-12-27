package fetch

import (
	"fmt"
	"net/http"
	"testing"
)

type Example struct {
	Hello string
}

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

	response, err := Fetch[string]("http://localhost:8080/plain", Options{})

	test.Run("get response", func(test *testing.T) {
		if err != nil {
			test.Fatalf("HTTP response failed: %v\n", err)
		} else {
			test.Logf("HTTP response received\n")
		}
	})

	test.Run("check response body", func(test *testing.T) {
		if response.Body != "Hello World" {
			test.Fatalf("HTTP response body expected \"Hello World\" but got \"%s\"", response.Body)
		} else {
			test.Logf("HTTP response was correct\n")
		}
	})
}

func TestAdvanced(test *testing.T) {
	go host()

	response, err := Fetch[Example]("http://localhost:8080/json", Options{
		Method: "POST",
		Headers: Headers{
			"User-Agent": "example",
			"Example":    []string{"Hello", "World"},
		},
		Body: "Hello World",
	})

	test.Run("get response", func(test *testing.T) {
		if err != nil {
			test.Fatalf("HTTP response failed: %v\n", err)
		} else {
			test.Logf("HTTP response received\n")
		}
	})

	test.Run("check response body", func(test *testing.T) {
		if response.Body.Hello != "World" {
			test.Fatalf("HTTP response body expected { \"World\": \"World\" } but got %+v\n", response.Body)
		} else {
			test.Logf("HTTP response was correct\n")
		}
	})
}

func BenchmarkGet(benchmark *testing.B) {
	go host()

	for index := 0; index < benchmark.N; index++ {
		Fetch[string]("http://localhost:8080/plain", Options{})
	}
}
