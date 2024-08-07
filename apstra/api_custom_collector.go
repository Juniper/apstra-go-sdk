package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiUrlCollectors              = "/api/telemetry/collectors"
	apiUrlCollectorsByServiceName = apiUrlCollectors + apiUrlPathDelim + "%s"
)

type CollectorPlatform struct {
	OsType    CollectorOSType
	OsVersion CollectorOSVersion
	OsFamily  []CollectorOSFamily
	Model     string
}

func (o *CollectorPlatform) UnmarshalJSON(data []byte) error {
	var raw struct {
		OsType    string `json:"os_type"`
		OsVersion string `json:"os_version"`
		OsFamily  string `json:"family"`
		Model     string `json:"model"`
	}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	err = o.OsType.FromString(raw.OsType)
	if err != nil {
		return err
	}

	err = o.OsVersion.FromString(raw.OsVersion)
	if err != nil {
		return err
	}

	o.Model = raw.Model

	for _, v := range strings.Split(raw.OsFamily, ",") {
		var variant CollectorOSFamily
		err = variant.FromString(v)
		if err != nil {
			return err
		}
		o.OsFamily = append(o.OsFamily, variant)
	}

	return nil
}

func (o *CollectorPlatform) MarshalJSON() ([]byte, error) {
	var raw struct {
		OsType    string `json:"os_type"`
		OsVersion string `json:"os_version"`
		OsFamily  string `json:"family"`
		Model     string `json:"model"`
	}
	raw.OsType = o.OsType.String()
	raw.OsVersion = o.OsVersion.String()
	raw.Model = o.Model
	raw.OsFamily = o.OsFamily[0].String()
	for _, v := range o.OsFamily[1:] {
		raw.OsFamily = raw.OsFamily + "," + v.String()
	}
	return json.Marshal(raw)
}

type CollectorQuery struct {
	Accessors map[string]string `json:"accessors"`
	Keys      map[string]string `json:"keys"`
	Value     string            `json:"value"`
	Filter    string            `json:"filter"`
}
type Collector struct {
	ServiceName             string
	Platform                CollectorPlatform   `json:"platform"`
	SourceType              CollectorSourceType `json:"source_type"`
	Cli                     string              `json:"cli"`
	Query                   CollectorQuery      `json:"query"`
	RelaxedSchemaValidation bool                `json:"relaxed_schema_validation"`
}

func (o *CollectorSourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *CollectorSourceType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	return o.FromString(s)
}

// GetAllCollectors gets all the Collectors for all services
func (o *Client) GetAllCollectors(ctx context.Context) ([]Collector, error) {
	var response struct {
		Items map[string]interface{} `json:"items"`
	}
	var collectors []Collector
	// First Reach out to /collectors , we are interested really only in the keys, so that we can pull the collectors
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlCollectors),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	for k := range response.Items {
		cs, err := o.GetCollectorsByServiceName(ctx, k)
		if err != nil {
			return nil, convertTtaeToAceWherePossible(err)
		}
		for _, v := range cs {
			v.ServiceName = k
			collectors = append(collectors, v)
		}
	}
	return collectors, nil
}

// GetCollectorsByServiceName gets all the Collectors that correspond to a particular service
func (o *Client) GetCollectorsByServiceName(ctx context.Context, name string) ([]Collector, error) {
	var ace ClientErr
	var Response struct {
		Items []Collector `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlCollectorsByServiceName, name),
		apiResponse: &Response,
	})

	if err != nil {
		err = convertTtaeToAceWherePossible(err)
		if errors.As(err, &ace) && ace.Type() == ErrNotfound {
			return nil, nil
		}
		return nil, err
	}

	for i := range Response.Items {
		Response.Items[i].ServiceName = name
	}
	return Response.Items, nil
}

// CreateCollector creates a collector
func (o *Client) CreateCollector(ctx context.Context, in *Collector) error {
	// Check if this is the first collector for that service name
	//cs, err := o.GetCollectorsByServiceName(ctx, in.ServiceName)
	//if err != nil {
	//	return err
	//}
	var Request struct {
		ServiceName string      `json:"service_name"`
		Items       []Collector `json:"collectors"`
	}
	Request.ServiceName = in.ServiceName
	Request.Items = append(Request.Items, *in)
	// This is the first collector for this service name
	// So we POST
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlCollectors),
		apiInput: &Request,
	})
	err = convertTtaeToAceWherePossible(err)
	var ace ClientErr
	if !(errors.As(err, &ace) && ace.Type() == ErrConflict) {
		return err // fatal error
	}

	// There are other collectors, so this is a patch
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlCollectorsByServiceName, in.ServiceName),
		apiInput: &Request,
	}))
}

// UpdateCollector Updates a collector
func (o *Client) UpdateCollector(ctx context.Context, in *Collector) error {
	var Request struct {
		Collectors []Collector `json:"collectors"`
	}
	Request.Collectors = append(Request.Collectors, *in)
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlCollectorsByServiceName, in.ServiceName),
		apiInput: &Request,
	}))
}

// DeleteAllCollectorsInService deletes all the collectors under a service
func (o *Client) DeleteAllCollectorsInService(ctx context.Context, name string) error {
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlCollectorsByServiceName, name),
	}))
}

func (p1 *CollectorPlatform) Equals(p2 *CollectorPlatform) bool {
	if p1.OsType != p2.OsType {
		return false
	}
	if p1.OsVersion != p2.OsVersion {
		return false
	}
	if p1.Model != p2.Model {
		return false
	}
	if len(p1.OsFamily) != len(p2.OsFamily) {
		return false
	}

	m := make(map[CollectorOSFamily]bool, len(p1.OsFamily))
	for _, v := range p1.OsFamily {
		m[v] = true
	}
	for _, v := range p2.OsFamily {
		_, ok := m[v]
		if !ok {
			return false
		}
	}
	return true
}

// DeleteCollector deletes the collector described in the object
func (o *Client) DeleteCollector(ctx context.Context, in *Collector) error {
	var Request struct {
		ServiceName string      `json:"service_name"`
		Items       []Collector `json:"collectors"`
	}

	cs, err := o.GetCollectorsByServiceName(ctx, in.ServiceName)
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	// There are no collectors
	if len(cs) == 0 {
		return nil
	}

	// If there is only one collector, we need to call DELETE
	if len(cs) == 1 {
		return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
			method: http.MethodDelete,
			urlStr: fmt.Sprintf(apiUrlCollectorsByServiceName, in.ServiceName),
		}))
	}

	// There is more than one collector, so we need to drop this collector from the list and PUT it backsxxa
	Request.ServiceName = in.ServiceName
	for _, c := range cs {
		if !c.Platform.Equals(&in.Platform) {
			Request.Items = append(Request.Items, c)
		}
	}

	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlCollectorsByServiceName, in.ServiceName),
		apiInput: &Request,
	}))
}
