package redisclient

import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"strings"
	"strconv"
	"errors"
	"fmt"
)

const dictKeyParam = "dictKey"
const valueParam = "value"
const expireParam = "expire"
const listIndexParam = "listIndex"
const errorParam = "error"
const messageParam = "message"

const valuesUrlWithParamTemplate = "%vvalues/%v?%v=%v"

type Client struct {
	url string
	httpClient http.Client
}

func New(url string) *Client {
	if strings.HasSuffix(url, "/") == false {
		url += "/"
	}
	return &Client{url, http.Client{Timeout: time.Second * 2}}
}

func (c *Client) Connect() error {
	_, err := c.httpCall(http.MethodGet, c.url + "ping",  nil)
	return err
}

//put key-value pair into the cache
func (c *Client) Put(key string, value interface{}) (interface{}, error) {
	return c.httpCall(http.MethodPut, c.url + "values/" + key,  value)
}

//put key-value pair with expire into the cache
func (c *Client) PutWithExpire(key string, value interface{}, expire int) (interface{}, error) {
	return c.httpCall(http.MethodPut,
		fmt.Sprintf(valuesUrlWithParamTemplate, c.url, key, expireParam, strconv.Itoa(expire)),  value)
}

//get list of cache keys
func (c *Client) Keys() ([]string, error) {
	var result []string
	respValue, err := c.httpCall(http.MethodGet, c.url + "keys", nil)
	if err == nil {
		keys := respValue.([]interface{})
		result = make([]string, len(keys))
		for i := range keys {
			result[i] = keys[i].(string)
		}
	}
	return result, err
}

//get cache value by key
func (c *Client) Get(key string) (interface{}, error) {
	return c.httpCall(http.MethodGet, c.url + "values/" + key,  nil)
}

//get i item from cache list value
func (c *Client) GetListElement(key string, listIndex int) (interface{}, error) {
	return c.httpCall(http.MethodGet,
		fmt.Sprintf(valuesUrlWithParamTemplate, c.url, key, listIndexParam, strconv.Itoa(listIndex)),  nil)
}

//get item by key from dict cache value
func (c *Client) GetDictValue(key string, dictKey string) (interface{}, error) {
	return c.httpCall(http.MethodGet,
		fmt.Sprintf(valuesUrlWithParamTemplate, c.url, key, dictKeyParam, dictKey),  nil)
}

//set a timeout on key in seconds
func (c *Client) Expire(key string, expire int) (interface{}, error) {
	return c.httpCall(http.MethodPut, c.url + "expire/" + key,  expire)
}

//returns the remaining time to live of a key that has a timeout
func (c *Client) GetTtl(key string) (int, error) {
	var ttl int
	resp, err := c.httpCall(http.MethodGet, c.url + "ttl/" + key,  nil)
	if (err == nil) && (resp != nil) {
		ttl = int(resp.(float64))
	}
	return ttl, err
}

//remove key-value pair from the cache
func (c *Client) Delete(key string) (interface{}, error) {
	return c.httpCall(http.MethodDelete, c.url + "values/" + key, nil)
}

//persist cache data
func (c *Client) Persist() (interface{}, error) {
	return c.httpCall(http.MethodPost, c.url + "persist", nil)
}

//reload persisted data
func (c *Client) Reload() (interface{}, error) {
	return c.httpCall(http.MethodPost, c.url + "reload", nil)
}

func (c *Client) httpCall(method, url string, value interface{}) (interface{}, error) {
	var result interface{}
	b := new(bytes.Buffer)
	if value != nil {
		json.NewEncoder(b).Encode(value)
	}
	req, err := http.NewRequest(method, url,  b)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			err = doErr
		} else {
			var respValue map[string]interface{}
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&respValue)
			if err == nil {
				if respValue[errorParam] != nil {
					err = errors.New(respValue[errorParam].(string))
				} else if respValue[valueParam] != nil {
					result = respValue[valueParam]
				} else if respValue[messageParam] != nil {
					result = respValue[messageParam]
				}
			}
		}
	}
	return result, err
}