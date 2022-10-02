package elastics

import (
	"github.com/elastic/go-elasticsearch/v8"
	"fmt"
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
