# Client for the Redis-like in-memory cache

## Installation and upgrade

```
    # Use 'go get' to install or upgrade (-u) the redis package
    go get -u github.com/dkazakevich/redisclient
    
    # Use `go test` to run tests
    go test github.com/dkazakevich/redisclient
```

## Usage

Use `import` to use `redisclient` in your program:

```
import (
  "github.com/dkazakevich/redisclient"
)
```

The `redisclient.New()` function returns a `*redisclient.Client` pointer that you can use
to interact with a redis server.

This example shows how to use redisclient.

```go
import (
	"github.com/dkazakevich/redisclient"
	"log"
)

func main() {

	var client *redisclient.Client

	client = redisclient.New("http://localhost:8080/api/v1/")

	err := client.Connect()

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
	}

	log.Println("Connected to the redis server.")

	client.PutWithExpire("hello", "Hello world!", 10)
	printPair(client.Get("hello"))

	client.Put("list", [6]int{2, 3, 5, 7, 11, 13})
	printPair(client.Get("list"))
	printPair(client.GetListElement("list", 0))

	client.Put("planets", map[string]string{"planet1": "Mercury", "planet2": "Venus", "planet3": "Earth"})
	printPair(client.Get("planets"))
	printPair(client.GetDictValue("planets", "planet3"))
	printPair(client.GetDictValue("planets", "planet4"))

	printPair(client.Keys())

	printPair(client.Delete("planets"))
	printPair(client.Keys())

	printPair(client.GetTtl("hello"))
	printPair(client.GetTtl("list"))

	printPair(client.Expire("list", 10))

	client.Persist()
}

func printPair(value interface{}, err error) {
	if err == nil {
		log.Printf("Value: '%v'\n", value)
	} else {
		log.Printf("Error: %v\n", err)
	}
}
```