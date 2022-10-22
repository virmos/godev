package elastics

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"time"
)

type ElasticRepository struct {
	Client *elasticsearch.Client
}

func NewElasticRepository(client *elasticsearch.Client) *ElasticRepository {
	return &ElasticRepository{
		Client: client,
	}
}

type TestElasticRepository struct {
	
}

func NewTestElasticRepository() *TestElasticRepository {
	return &TestElasticRepository{
	}
}

// helper function
func DateFromUTCString(s string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, s)
	if err != nil {
		fmt.Println(err)
		yearOne := time.Date(0001, 11, 17, 20, 34, 58, 651387237, time.UTC)
		return yearOne, err
	}
	return t, err
}

type Report struct {
	Host           string   `json:"host"`
	StatusCode     string   `json:"status_code"`
	Time           string   `json:"@timestamp"`
	HoursHistogram []int    `json:"@hours_histogram"`
	DaysHistogram  []int    `json:"@days_histogram"`
	Histogram      []string `json:"@histogram"`
	Count          []string `json:"@count"`
}
