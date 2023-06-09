package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrlMetricdbMetric = "/api/metricdb/metric"
	apiUrlMetricdbQuery  = "/api/metricdb/query"

	aggrWriterSplit = "_aggr_"
	aggrWriterRegex = "^.+" + aggrWriterSplit + "[0-9]+$"
)

// metricdbMetricResponse is generated by GET to apstra apiUrlMetricdbMetric API endpoint
type metricdbMetricResponse struct {
	Items []MetricdbMetric `json:"items"`
}

type MetricdbMetric struct {
	Application string `json:"application"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

type metricdbQuery struct {
	Application string    `json:"application"` // required
	Namespace   string    `json:"namespace"`   // required
	Name        string    `json:"name"`        // required
	BeginTime   time.Time `json:"begin_time"`  // required
	EndTime     time.Time `json:"end_time"`    // required
	//Filter      string `json:"filter"`
	//Paging      struct {
	//	PageNumber int       `json:"page_number"`
	//	EndTime    time.Time `json:"end_time"`
	//	PageSize   int       `json:"page_size"`
	//	BeginTime  time.Time `json:"begin_time"`
	//} `json:"paging"`
	//Caching   struct {
	//	Clear     bool `json:"clear"`
	//	KeepAlive bool `json:"keep_alive"`
	//	Timeout   int  `json:"timeout"`
	//} `json:"caching"`
	//Aggregation struct {
	//	Metrics struct {
	//		AdditionalProp1 string `json:"additionalProp1"`
	//		AdditionalProp2 string `json:"additionalProp2"`
	//		AdditionalProp3 string `json:"additionalProp3"`
	//	} `json:"metrics"`
	//	Period int `json:"period"`
	//} `json:"aggregation"`
}

type MetricDbQueryRequest struct {
	metric MetricdbMetric
	begin  time.Time
	end    time.Time
}

type MetricDbQueryResponse struct {
	Status struct {
		TotalCount    int       `json:"total_count"`
		BeginTime     time.Time `json:"begin_time"`
		EndTime       time.Time `json:"end_time"`
		ResultCode    string    `json:"result_code"`
		LastTimestamp time.Time `json:"last_timestamp"`
	} `json:"status"`
	Items []json.RawMessage `json:"items"`
}

func (o *Client) getMetricdbMetrics(ctx context.Context) (*metricdbMetricResponse, error) {
	apstraUrl, err := url.Parse(apiUrlMetricdbMetric)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	response := &metricdbMetricResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

func (o *Client) queryMetricdb(ctx context.Context, begin, end time.Time, metric MetricdbMetric) (*MetricDbQueryResponse, error) {
	apstraUrl, err := url.Parse(apiUrlMetricdbQuery)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	response := &MetricDbQueryResponse{}
	q := metricdbQuery{
		Application: metric.Application,
		Namespace:   metric.Namespace,
		Name:        metric.Name,
		BeginTime:   begin,
		EndTime:     end,
	}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    q,
		apiResponse: response,
	})
}

func useAggregation(in string) (bool, string, int, error) {
	aggrWriter, err := regexp.MatchString(aggrWriterRegex, in)
	if err != nil {
		return false, "", 0, err
	}
	if !aggrWriter {
		return false, in, 0, nil
	}

	split := strings.Split(in, aggrWriterSplit)
	seconds, err := strconv.Atoi(split[1])
	if err != nil {
		return false, "", 0, err
	}

	return true, split[0], seconds, nil
}
