package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlTelemetryServiceRegistry            = "/api/telemetry-service-registry"
	apiUrlTelemetryServiceRegistryPrefix      = apiUrlTelemetryServiceRegistry + apiUrlPathDelim
	apiUrlTelemetryServiceRegistryEntryByName = apiUrlTelemetryServiceRegistryPrefix + "%s"
)

type rawTelemetryServiceRegistryEntry struct {
	ServiceName       string          `json:"service_name"`
	ApplicationSchema json.RawMessage `json:"application_schema"`
	StorageSchemaPath string          `json:"storage_schema_path"`
	Builtin           bool            `json:"builtin"`
	Description       string          `json:"description"`
	Version           string          `json:"version"`
}

type TelemetryServiceRegistryEntry struct {
	ServiceName       string
	StorageSchemaPath StorageSchemaPath
	ApplicationSchema json.RawMessage
	Builtin           bool
	Description       string
	Version           string
}

func (o *TelemetryServiceRegistryEntry) UnmarshalJSON(data []byte) error {
	var sspath StorageSchemaPath
	var r rawTelemetryServiceRegistryEntry
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}
	err = sspath.FromString(r.StorageSchemaPath)
	if err != nil {
		return err
	}
	*o = TelemetryServiceRegistryEntry{
		ServiceName:       r.ServiceName,
		StorageSchemaPath: sspath,
		ApplicationSchema: r.ApplicationSchema,
		Builtin:           r.Builtin,
		Description:       r.Description,
		Version:           r.Version,
	}
	return nil
}

func (o *TelemetryServiceRegistryEntry) MarshalJSON() ([]byte, error) {

	return json.Marshal(rawTelemetryServiceRegistryEntry{
		ServiceName:       o.ServiceName,
		StorageSchemaPath: o.StorageSchemaPath.String(),
		Builtin:           o.Builtin,
		Description:       o.Description,
		Version:           o.Version,
		ApplicationSchema: o.ApplicationSchema,
	})
}

// GetAllTelemetryServiceRegistryEntries gets all the Telemetry Service Registry Entries
func (o *Client) GetAllTelemetryServiceRegistryEntries(ctx context.Context) ([]TelemetryServiceRegistryEntry, error) {
	response := &struct {
		Items []TelemetryServiceRegistryEntry `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistry),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

// GetTelemetryServiceRegistryEntry gets all the Telemetry Service Registry Entries
func (o *Client) GetTelemetryServiceRegistryEntry(ctx context.Context, name string) (*TelemetryServiceRegistryEntry, error) {
	response := &TelemetryServiceRegistryEntry{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

// CreateTelemetryServiceRegistryEntry creates a telemetry service registry entry
func (o *Client) CreateTelemetryServiceRegistryEntry(ctx context.Context, r TelemetryServiceRegistryEntry) (string, error) {
	response := &struct {
		Name string `json:"service_name"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistry),
		apiInput:    &r,
		apiResponse: &response,
	})
	if err != nil {
		return "", err
	}
	return response.Name, nil
}

// UpdateTelemetryServiceRegistryEntry updates a telemetry service registry entry
func (o *Client) UpdateTelemetryServiceRegistryEntry(ctx context.Context, name string, r *TelemetryServiceRegistryEntry) error {
	response := &struct {
		Name string `json:"service_name"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPut,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
		apiInput:    &r,
		apiResponse: &response,
	})

	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

// DeleteTelemetryServiceRegistryEntry deletes a telemetry service registry entry
func (o *Client) DeleteTelemetryServiceRegistryEntry(ctx context.Context, name string) error {
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
	}))
}
