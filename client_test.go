package redisclient

import (
	"testing"
	"sync"
	"time"
	"strconv"
	"log"
	"os"
)

const performIterations  = 1000

var client *Client

func TestMain(m *testing.M) {
	client = NewClient("http://localhost:8080/api/v1/")
	code := m.Run()
	os.Exit(code)
}

func TestPutPerformance(t *testing.T) {

	data := make([]string, performIterations)
	for i := range data {
		data[i] = "value" + strconv.Itoa(i)
	}

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(performIterations)
	for i := range data {
		value := data[i]
		go func() {
			client.Put(value, value)
			client.Get("value"+strconv.Itoa(i))
			wg.Done()
		}()
	}
	wg.Wait()

	keys, _ := client.Keys()
	log.Printf("Keys size: %v; keys: %v\n", len(keys), keys)

	elapsed := time.Since(start)
	log.Printf("%v times put took %v", performIterations, elapsed)
}