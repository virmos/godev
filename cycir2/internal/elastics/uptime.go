package elastics

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (es *ElasticRepository) CreateIndex(index string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	checkExistsReq := esapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := checkExistsReq.Do(ctx, es.Client)
	if err != nil {
		log.Fatalf("Error checking index existence: %s", err)
		return err
	}
	if res.StatusCode == 200 {
		log.Println("Index already exists.")
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"host": {
					"type": "keyword"
				},
				"status_code": {
					"type": "keyword"
				},
				"@timestamp": {
					"type": "date" 
				}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(string(mapping)),
	}
	_, err = req.Do(ctx, es.Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	return nil
}

// testing
func (es *ElasticRepository) GetAllReports(index string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{

			},
		},
	}

	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error querying elasticsearch: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	return r, nil
}

func (es *ElasticRepository) InsertHostStatusReport(index, hostName, statusCode, date string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	report := Report{
		Host:       hostName,
		StatusCode: statusCode,
		Time:       date,
	}
	data, err := json.Marshal(report)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	// Perform the request with the client.
	_, err = req.Do(ctx, es.Client)

	if err != nil {
		log.Fatalf("Error insert into elasticsearch: %s", err)
		return err
	}

	return nil
}

// startDate and Endate are in UTC format
func (es *ElasticRepository) GetYesterdayUptimeReport(index string) (map[string]Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": map[string]interface{}{
						"status_code": "200",
					},
				},
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"gte": "now-1d/d",
							
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"query": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "host",
				},
				"aggs": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"date_histogram": map[string]interface{}{
							"field":             "@timestamp",
							"calendar_interval": "hour",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error querying elasticsearch: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	timeReports := make(map[string]Report)
	for _, bucket := range r["aggregations"].(map[string]interface{})["query"].(map[string]interface{})["buckets"].([]interface{}) {
		host := bucket.(map[string]interface{})["key"]
		reports := bucket.(map[string]interface{})["timestamp"].(map[string]interface{})["buckets"]
		hoursHistogram := make([]int, 24, 24)

		for _, report := range reports.([]interface{}) {
			timestamp := report.(map[string]interface{})["key_as_string"]
			count := report.(map[string]interface{})["doc_count"]
			time, _ := DateFromUTCString(timestamp.(string))
			hoursHistogram[time.Hour()] = int(count.(float64))
		}
		timeReports[host.(string)] = Report{HoursHistogram: hoursHistogram}
	}

	return timeReports, nil
}

// startDate and Endate are in UTC format
func (es *ElasticRepository) GetYesterdayReport(index string) (map[string]Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"gte": "now-1d/d",
							
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"query": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "host",
				},
				"aggs": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"date_histogram": map[string]interface{}{
							"field":             "@timestamp",
							"calendar_interval": "hour",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error querying elasticsearch: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	timeReports := make(map[string]Report)
	for _, bucket := range r["aggregations"].(map[string]interface{})["query"].(map[string]interface{})["buckets"].([]interface{}) {
		host := bucket.(map[string]interface{})["key"]
		reports := bucket.(map[string]interface{})["timestamp"].(map[string]interface{})["buckets"]
		hoursHistogram := make([]int, 24, 24)

		for _, report := range reports.([]interface{}) {
			timestamp := report.(map[string]interface{})["key_as_string"]
			count := report.(map[string]interface{})["doc_count"]
			time, _ := DateFromUTCString(timestamp.(string))
			hoursHistogram[time.Hour()] = int(count.(float64))
		}
		timeReports[host.(string)] = Report{HoursHistogram: hoursHistogram}
	}

	return timeReports, nil
}

// startDate and Endate are in UTC format
func (es *ElasticRepository) GetRangeUptimeReport(index, hostName, startDate, endDate string) (map[string]Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match": map[string]interface{}{
						"status_code": "200",
					},
				},
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"gte": startDate,
							"lt":  endDate,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"query": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "host",
				},
				"aggs": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"date_histogram": map[string]interface{}{
							"field":             "@timestamp",
							"calendar_interval": "day",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error querying elasticsearch: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	timeReports := make(map[string]Report)
	for _, bucket := range r["aggregations"].(map[string]interface{})["query"].(map[string]interface{})["buckets"].([]interface{}) {
		host := bucket.(map[string]interface{})["key"]
		reports := bucket.(map[string]interface{})["timestamp"].(map[string]interface{})["buckets"]
		daysHistogram := make([]int, 31, 31)
		if host.(string) != hostName {
			continue
		}

		for _, report := range reports.([]interface{}) {
			timestamp := report.(map[string]interface{})["key_as_string"]
			count := report.(map[string]interface{})["doc_count"]
			time, _ := DateFromUTCString(timestamp.(string))
			daysHistogram[time.Hour()] = int(count.(float64))
		}
		timeReports[host.(string)] = Report{DaysHistogram: daysHistogram}
	}

	return timeReports, nil
}

// startDate and Endate are in UTC format
func (es *ElasticRepository) GetRangeReport(index, hostName, startDate, endDate string) (map[string]Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"gte": startDate,
							"lt":  endDate,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"query": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "host",
				},
				"aggs": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"date_histogram": map[string]interface{}{
							"field":             "@timestamp",
							"calendar_interval": "day",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error querying elasticsearch: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	timeReports := make(map[string]Report)
	for _, bucket := range r["aggregations"].(map[string]interface{})["query"].(map[string]interface{})["buckets"].([]interface{}) {
		host := bucket.(map[string]interface{})["key"]
		reports := bucket.(map[string]interface{})["timestamp"].(map[string]interface{})["buckets"]
		daysHistogram := make([]int, 31, 31)
		
		if host.(string) != hostName {
			continue
		}

		for _, report := range reports.([]interface{}) {
			timestamp := report.(map[string]interface{})["key_as_string"]
			count := report.(map[string]interface{})["doc_count"]
			time, _ := DateFromUTCString(timestamp.(string))
			daysHistogram[time.Hour()] = int(count.(float64))
		}
		timeReports[host.(string)] = Report{DaysHistogram: daysHistogram}
	}

	return timeReports, nil
}
