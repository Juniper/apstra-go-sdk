package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	oenum "github.com/orsinium-labs/enum"
	"net/http"
)

const (
	apiUrlTelemetryServiceRegistry            = "/api/telemetry-service-registry"
	apiUrlTelemetryServiceRegistryPrefix      = apiUrlTelemetryServiceRegistry + apiUrlPathDelim
	apiUrlTelemetryServiceRegistryEntryByName = apiUrlTelemetryServiceRegistryPrefix + "%s"
)

type StorageSchemaPath oenum.Member[string]
type SchemaType oenum.Member[string]

func (o SchemaType) String() string {
	return o.Value
}

func (o *SchemaType) FromString(s string) error {
	t := StorageSchemaPaths.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse SchemaType %q", s)
	}
	o.Value = t.Value
	return nil
}

func (o StorageSchemaPath) String() string {
	return o.Value
}

func (o *StorageSchemaPath) FromString(s string) error {
	t := StorageSchemaPaths.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse StorageSchemaPath %q", s)
	}
	o.Value = t.Value
	return nil
}

var (
	StorageSchemaPathXCVR               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.xcvr"}
	StorageSchemaPathGRAPH              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.graph"}
	StorageSchemaPathROUTE              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route"}
	StorageSchemaPathMAC                = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mac"}
	StorageSchemaPathOPTICAL_XCVR       = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.optical_xcvr"}
	StorageSchemaPathHOSTNAME           = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.hostname"}
	StorageSchemaPathGENERIC            = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.generic"}
	StorageSchemaPathLAG                = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lag"}
	StorageSchemaPathBGP                = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.bgp"}
	StorageSchemaPathINTERFACE          = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface"}
	StorageSchemaPathMLAG               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mlag"}
	StorageSchemaPathIBA_STRING_DATA    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_string_data"}
	StorageSchemaPathIBA_INTEGER_DATA   = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_integer_data"}
	StorageSchemaPathROUTE_LOOKUP       = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route_lookup"}
	StorageSchemaPathINTERFACE_COUNTERS = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface_counters"}
	StorageSchemaPathARP                = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.arp"}
	StorageSchemaPathCPP_GRAPH          = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.cpp_graph"}
	StorageSchemaPathNSXT               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.nsxt"}
	StorageSchemaPathENVIRONMENT        = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.environment"}
	StorageSchemaPathLLDP               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lldp"}
	StorageSchemaPaths                  = oenum.New(StorageSchemaPathXCVR, StorageSchemaPathGRAPH, StorageSchemaPathROUTE, StorageSchemaPathMAC, StorageSchemaPathOPTICAL_XCVR, StorageSchemaPathHOSTNAME, StorageSchemaPathGENERIC, StorageSchemaPathLAG, StorageSchemaPathBGP, StorageSchemaPathINTERFACE, StorageSchemaPathMLAG, StorageSchemaPathIBA_STRING_DATA, StorageSchemaPathIBA_INTEGER_DATA, StorageSchemaPathROUTE_LOOKUP, StorageSchemaPathINTERFACE_COUNTERS, StorageSchemaPathARP, StorageSchemaPathCPP_GRAPH, StorageSchemaPathNSXT, StorageSchemaPathENVIRONMENT, StorageSchemaPathLLDP)

	//SchemaTypeInteger = SchemaType{Value: "integer"}
	//SchemaTypeString  = SchemaType{Value: "string"}
	//SchemaTypeObject  = SchemaType{Value: "object"}
	//SchemaTypes       = oenum.New(SchemaTypeInteger, SchemaTypeString, SchemaTypeObject)
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

func (o *rawTelemetryServiceRegistryEntry) polish() (*TelemetryServiceRegistryEntry, error) {
	var sspath StorageSchemaPath
	err := sspath.FromString(o.StorageSchemaPath)

	if err != nil {
		return nil, err
	}
	return &TelemetryServiceRegistryEntry{
		ServiceName:       o.ServiceName,
		StorageSchemaPath: sspath,
		ApplicationSchema: o.ApplicationSchema,
		Builtin:           o.Builtin,
		Description:       o.Description,
		Version:           o.Version,
	}, nil
}

func (o *TelemetryServiceRegistryEntry) raw() *rawTelemetryServiceRegistryEntry {

	return &rawTelemetryServiceRegistryEntry{
		ServiceName:       o.ServiceName,
		StorageSchemaPath: o.StorageSchemaPath.String(),
		Builtin:           o.Builtin,
		Description:       o.Description,
		Version:           o.Version,
		ApplicationSchema: o.ApplicationSchema,
	}
}

func (o *Client) getAllTelemetryServiceRegistryEntries(ctx context.Context) ([]rawTelemetryServiceRegistryEntry, error) {
	response := &struct {
		Items []rawTelemetryServiceRegistryEntry `json:"items"`
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

func (o *Client) getTelemetryServiceRegistryEntry(ctx context.Context, name string) (*rawTelemetryServiceRegistryEntry, error) {
	response := &rawTelemetryServiceRegistryEntry{}
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

func (o *Client) createTelemetryServiceRegistryEntry(ctx context.Context, r rawTelemetryServiceRegistryEntry) (string, error) {
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

func (o *Client) updateTelemetryServiceRegistryEntry(ctx context.Context, name string, r *rawTelemetryServiceRegistryEntry) error {
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

func (o *Client) deleteTelemetryServiceRegistryEntry(ctx context.Context, name string) error {
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlTelemetryServiceRegistryEntryByName, name),
	}))
}
