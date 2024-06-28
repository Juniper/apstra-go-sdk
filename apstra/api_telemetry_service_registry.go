package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlTelemetryServiceRegistry            = "/api/telemetry-service-registry"
	apiUrlTelemetryServiceRegistryEntryByName = apiUrlTelemetryServiceRegistry + apiUrlPathDelim + "%s"
)

var (
	_ json.Marshaler   = (*TelemetryServiceRegistryEntry)(nil)
	_ json.Unmarshaler = (*TelemetryServiceRegistryEntry)(nil)
)

type TelemetryServiceRegistryEntry struct {
	ServiceName       string
	ApplicationSchema json.RawMessage
	StorageSchemaPath StorageSchemaPath
	Builtin           bool
	Description       string
	Version           string
}

func (o *TelemetryServiceRegistryEntry) UnmarshalJSON(data []byte) error {
	var raw struct {
		ServiceName       string          `json:"service_name"`
		ApplicationSchema json.RawMessage `json:"application_schema"`
		StorageSchemaPath string          `json:"storage_schema_path"`
		Builtin           bool            `json:"builtin"`
		Description       string          `json:"description"`
		Version           string          `json:"version"`
	}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	o.ServiceName = raw.ServiceName
	o.ApplicationSchema = raw.ApplicationSchema
	o.Builtin = raw.Builtin
	o.Description = raw.Description
	o.Version = raw.Version
	err = o.StorageSchemaPath.FromString(raw.StorageSchemaPath)
	if err != nil {
		return err
	}

	return nil
}

func (o *TelemetryServiceRegistryEntry) MarshalJSON() ([]byte, error) {
	var raw struct {
		ServiceName       string          `json:"service_name"`
		ApplicationSchema json.RawMessage `json:"application_schema"`
		StorageSchemaPath string          `json:"storage_schema_path"`
		Builtin           bool            `json:"builtin"`
		Description       string          `json:"description"`
		Version           string          `json:"version"`
	}

	raw.ServiceName = o.ServiceName
	raw.StorageSchemaPath = o.StorageSchemaPath.String()
	raw.Builtin = o.Builtin
	raw.Description = o.Description
	raw.Version = o.Version
	raw.ApplicationSchema = o.ApplicationSchema

	return json.Marshal(raw)
}

// GetAllTelemetryServiceRegistryEntries gets all the Telemetry Service Registry Entries
func (o *Client) GetAllTelemetryServiceRegistryEntries(ctx context.Context) ([]TelemetryServiceRegistryEntry, error) {
	var response struct {
		Items []TelemetryServiceRegistryEntry `json:"items"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistry),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

// GetTelemetryServiceRegistryEntry gets all the Telemetry Service Registry Entries
func (o *Client) GetTelemetryServiceRegistryEntry(ctx context.Context, name string) (*TelemetryServiceRegistryEntry, error) {
	var response TelemetryServiceRegistryEntry

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

// CreateTelemetryServiceRegistryEntry creates a telemetry service registry entry
func (o *Client) CreateTelemetryServiceRegistryEntry(ctx context.Context, in TelemetryServiceRegistryEntry) (string, error) {
	var response struct {
		Name string `json:"service_name"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistry),
		apiInput:    &in,
		apiResponse: &response,
	})
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

// UpdateTelemetryServiceRegistryEntry updates a telemetry service registry entry
func (o *Client) UpdateTelemetryServiceRegistryEntry(ctx context.Context, name string, in *TelemetryServiceRegistryEntry) error {
	var response struct {
		Name string `json:"service_name"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPut,
		urlStr:      fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// DeleteTelemetryServiceRegistryEntry deletes a telemetry service registry entry
func (o *Client) DeleteTelemetryServiceRegistryEntry(ctx context.Context, name string) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
