# Client for the Redis-like in-memory cache

## Installation and upgrade

```
    # Use 'go get' to install or upgrade (-u) the redis package
    go get -u github.com/dkazakevich/redisclient
    
    # Use `go test` to run tests (the server should be run for it)
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
	"fmt"
)

func main() {

	client := redisclient.New("http://localhost:8081/api/v1/")

	err := client.Connect()

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
	}

	fmt.Println("Connected to the redis server.")

	fmt.Println("================ Put and get string ================")
	client.PutWithExpire("hello", "Hello world!", 10)
	printResult(client.Get("hello"))

	fmt.Println("================ Put and get list ==================")
	client.Put("list", [6]int{2, 3, 5, 7, 11, 13})
	printResult(client.Get("list"))
	printResult(client.GetListElement("list", 0))

	fmt.Println("================ Put and get dictionary ============")
	client.Put("planets", map[string]string{"planet1": "Mercury", "planet2": "Venus", "planet3": "Earth"})
	printResult(client.Get("planets"))
	printResult(client.GetDictValue("planets", "planet3"))
	printResult(client.GetDictValue("planets", "planet4"))

	fmt.Println("================ Cache keys ========================")
	printResult(client.Keys())

	fmt.Println("================ Delete 'planets' ==================")
	printResult(client.Delete("planets"))

	fmt.Println("================ Cache keys ========================")
	printResult(client.Keys())

	fmt.Println("================ Get 'hello' ttl  ==================")
	printResult(client.GetTtl("hello"))

	fmt.Println("================ Get 'nonExistent' ttl  ============")
	printResult(client.GetTtl("nonExistent"))

	fmt.Println("================ Set 'hello' expire  ===============")
	printResult(client.Expire("hello", 10))

	fmt.Println("================ Persist cache data  ================")
	printResult(client.Persist())
}

func printResult(value interface{}, err error) {
	if err == nil {
		fmt.Printf("Value: '%v'\n", value)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
```