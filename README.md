# Client for the Redis-like in-memory cache

##Installation and upgrade

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
var client *redisclient.Client
	
client = redisclient.New("http://localhost:8080/api/v1/")

err := client.Connect(host, port)

if err != nil {
  log.Fatalf("Connect failed: %s\n", err.Error())
  return
}

log.Println("Connected to the redis server.")


```