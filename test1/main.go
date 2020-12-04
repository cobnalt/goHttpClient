package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	client := NewClient("")

	ch1 := make(chan interface{})
	ch2 := make(chan error)

	urls := []string{
		"https://petstore.swagger.io/v2/pet",
		"https://swapi.dev/api/people",
	}
	for _, url := range urls {
		wg.Add(1)
		go client.GetById(context.Background(), url, 1, ch1, ch2)

	}

	for item := range ch1 {
		fmt.Println("\n\n", item)
	}

	fmt.Println("\n\n", "range ends")

	wg.Wait()
	fmt.Println("Все запросы выполнены")

}

const (
	BaseUrlV2 = "https://petstore.swagger.io/v2"
	BaseUrlV1 = "https://swapi.dev/api"
)

type SW struct {
	Name       string   `json:"name"`
	Height     string   `json:"height"`
	Mass       string   `json:"mass"`
	Hair_color string   `json:"hair_color"`
	Skin_color string   `json:"skin_color"`
	Eye_color  string   `json:"eye_color"`
	Birth_year string   `json:"birth_year"`
	Gender     string   `json:"gender"`
	Homeworld  string   `json:"homeworld"`
	Films      []string `json:"films"`
}
type Pet struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type Client struct {
	BaseUrl    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseUrl: BaseUrlV1,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetById(ctx context.Context, url string, id int, ch1 chan interface{}, ch2 chan error) {
	defer wg.Done()
	url1 := fmt.Sprintf("%s/%d/", url, id)
	//url := fmt.Sprintf("%s/pet/%d/", c.BaseUrl, id)
	fmt.Println(url1)
	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		ch1 <- nil
		ch2 <- err
	}

	req = req.WithContext(ctx)

	var res interface{}

	if err := c.sendRequest(req, &res); err != nil {
		ch1 <- nil
		ch2 <- err
	}
	//fmt.Println("res1\n\n", res)
	ch1 <- res
	ch2 <- nil
	close(ch1)
	close(ch2)
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	err = json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
