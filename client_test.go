package redisclient

import (
	"testing"
	"os"
	"strconv"
	"time"
	"sync"
	"log"
)

const performanceIterations  = 1000
const stringKey = "sixthMonth"
const dictKey = "planets"
const listKey = "cars"
const nonExistentKey = "nonExistent"

const itemNotFoundMsg = "Cache item not found"

var client *Client

func TestMain(m *testing.M) {

	client = New("http://localhost:8080/api/v1/")
	code := m.Run()
	os.Exit(code)
}

func TestPutAndGetString(t *testing.T) {

	value := "June"

	storedValue, err := client.PutWithExpire(stringKey, value, 10)
	assertEquals(t, nil, err)
	assertEquals(t, value, storedValue)

	storedValue, err = client.Get(stringKey)
	assertEquals(t, nil, err)
	assertEquals(t, value, storedValue)
}

func TestPutAndGetDict(t *testing.T) {

	value := map[string]string{"planet1": "Mercury", "planet2": "Venus", "planet3": "Earth"}

	storedValue, err := client.Put(dictKey, value)
	assertEquals(t, nil, err)

	storedValue, err = client.GetDictValue(dictKey, "planet1")
	assertEquals(t, nil, err)
	assertEquals(t, "Mercury", storedValue)

}

func TestPutAndGetList(t *testing.T) {

	value := [3]string{"Toyota", "Opel", "Ford"}

	storedValue, err := client.Put(listKey, value)
	assertEquals(t, nil, err)

	storedValue, err = client.GetListElement(listKey, 1)
	assertEquals(t, nil, err)
	assertEquals(t, "Opel", storedValue)
}

func TestGetNonExistentItem(t *testing.T) {

	storedValue, err := client.Get(nonExistentKey)
	assertEquals(t, itemNotFoundMsg, err.Error())
	assertEquals(t, nil, storedValue)
}

func TestKeys(t *testing.T) {

	keys, err := client.Keys()
	assertEquals(t, nil, err)
	assertTrue(t, len(keys) > 0)
}

func TestDelete(t *testing.T) {

	_, err := client.Delete(listKey)
	assertEquals(t, nil, err)

	storedValue, err := client.Get(listKey)
	assertEquals(t, itemNotFoundMsg, err.Error())
	assertEquals(t, nil, storedValue)
}

func TestExpireAndCheckTtl(t *testing.T) {

	_, err := client.Expire(stringKey, 10)
	assertEquals(t, nil, err)

	ttl, err := client.GetTtl(stringKey)
	assertEquals(t, nil, err)
	assertTrue(t, ttl > 0)
}

func TestNonExistentExpire(t *testing.T) {

	_, err := client.Expire(nonExistentKey, 10)
	assertEquals(t, itemNotFoundMsg, err.Error())
}

func TestPersistItemTtl(t *testing.T) {

	_, err := client.GetTtl(dictKey)
	assertEquals(t, itemNotFoundMsg, err.Error())
}

func TestNonExistentTtl(t *testing.T) {

	_, err := client.GetTtl(nonExistentKey)
	assertEquals(t, itemNotFoundMsg, err.Error())
}

func TestPutAndGetPerformance(t *testing.T) {

	data := make([]string, performanceIterations)
	for i := range data {
		data[i] = stringKey + strconv.Itoa(i)
	}

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(performanceIterations)
	for i := range data {
		value := data[i]
		go func() {
			client.Put(value, value)
			client.Get(stringKey + strconv.Itoa(i))
			wg.Done()
		}()
	}
	wg.Wait()

	keys, _ := client.Keys()
	log.Printf("Keys size: %v", len(keys))

	for i := range keys {
		client.Delete(keys[i])
	}

	elapsed := time.Since(start)
	log.Printf("%v times put and get took %v", performanceIterations, elapsed)
}

func assertEquals(t *testing.T, expected interface{}, actual interface{}) {

	if expected == actual {
		return
	}

	t.Fatalf("Expected: '%v'. Actual: '%v'", expected, actual)
}

func assertTrue(t *testing.T, condition bool) {

	if condition {
		return
	}

	t.Fatal("Expected 'true' but was 'false'")
}