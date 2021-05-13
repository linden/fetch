# Fetch
Simple go library that hopes to comply with the [Fetch Spec](https://fetch.spec.whatwg.org/).

### Example

```go
package main

import (
	"github.com/linden/fetch"

	"fmt"
)

func main() {
	//make the inital request
	response, err := fetch.Fetch("https://www.githubstatus.com/api/v2/status.json", fetch.Options{
        	Headers: fetch.Headers{
			"User-Agent": "Example",
		},
		Method: "GET",
	})

	if err != nil {
		panic(err)
	}

	//check the HTTP status
	if response.Status != 200 {
		panic(fmt.Errorf("HTTP Error Code %v\n", response.Status))
	}

	//cast to a map[string]interface{}
	body, err := response.JSON()

	if err != nil {
		panic(err)
	}

	fmt.Printf("status: %v, body: %+v", response.Status, body)
}
```
