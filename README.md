# Fetch
Simple go library that hopes to comply with the [Fetch Spec](https://fetch.spec.whatwg.org/).

### Example

```go
package main

import (
	"fmt"

	"github.com/linden/fetch"
)

type Status struct {
	Page struct {
		URL string `json:"url"`
	}
	Status struct {
		Description string
	}
}

func main() {
	// make the request
	response, err := fetch.Fetch[Status]("https://www.githubstatus.com/api/v2/status.json", fetch.Options[fetch.Empty]{
        Headers: fetch.Headers{
			"User-Agent": "Example",
		},
		Method: "GET",
	})

	if err != nil {
		panic(err)
	}

	// check the status
	if response.Status != 200 {
		panic(fmt.Errorf("HTTP Error Code %v\n", response.Status))
	}

	// print the response
	fmt.Printf("status: %v, body: %+v", response.Status, response.Body)
}
```
