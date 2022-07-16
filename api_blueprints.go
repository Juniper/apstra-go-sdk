package goapstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	apiUrlBlueprints                  = "/api/blueprints"
	apiUrlPathDelim                   = "/"
	apiUrlBlueprintsPrefix            = apiUrlBlueprints + apiUrlPathDelim
	apiUrlBlueprintById               = apiUrlBlueprintsPrefix + "%s"
	apiUrlBlueprintRoutingZones       = apiUrlBlueprintById + apiUrlPathDelim + "security_zones"
	apiUrlBlueprintRoutingZonesPrefix = apiUrlBlueprintRoutingZones + apiUrlPathDelim
	apiUrlBlueprintRoutingZonesById   = apiUrlBlueprintRoutingZonesPrefix + "%s"

	initTypeFromTemplate = "template_reference"
)

const (
	RefDesignDatacenter = RefDesign(iota)

	refDesignDatacenter = "two_stage_l3clos"
	refDesignUnknown    = "unknown reference design %d"
)

type RefDesign int

func (o RefDesign) String() string {
	switch o {
	case RefDesignDatacenter:
		return refDesignDatacenter
	default:
		return fmt.Sprintf(refDesignUnknown, o)
	}
}

type getBluePrintsResponse struct {
	Items []BlueprintStatus `json:"items"`
}

type optionsBlueprintsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type optionsBlueprintsRoutingzonesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type postBlueprintsResponse struct {
	Id     ObjectId `json:"id"`
	TaskId TaskId   `json:"task_id"`
}

type Blueprint struct {
	client         *Client
	Id             ObjectId                   `json:"id"`
	Version        int                        `json:"version"`
	Design         string                     `json:"design"`
	LastModifiedAt time.Time                  `json:"last_modified_at"`
	Label          string                     `json:"label"`
	Relationships  map[string]json.RawMessage `json:"relationships"`
	Nodes          map[string]json.RawMessage `json:"nodes"`
	SourceVersions struct {
		ConfigBlueprint int `json:"config_blueprint"`
	} `json:"source_versions"`
}

type BlueprintStatus struct {
	Id                     ObjectId  `json:"id"`
	Label                  string    `json:"label"`
	Status                 string    `json:"status"`
	Design                 string    `json:"design"`
	HasUncommittedChanges  bool      `json:"has_uncommitted_changes"`
	Version                int       `json:"version"`
	LastModifiedAt         time.Time `json:"last_modified_at"`
	SuperspineCount        int       `json:"superspine_count"`
	SpineCount             int       `json:"spine_count"`
	LeafCount              int       `json:"leaf_count"`
	AccessCount            int       `json:"access_count"`
	GenericCount           int       `json:"generic_count"`
	ExternalRouterCount    int       `json:"external_router_count"`
	L2ServerCount          int       `json:"l2_server_count"`
	L3ServerCount          int       `json:"l3_server_count"`
	RemoteGatewayCount     int       `json:"remote_gateway_count"`
	BuildWarningsCount     int       `json:"build_warnings_count"`
	RootCauseCount         int       `json:"root_cause_count"`
	TopLevelRootCauseCount int       `json:"top_level_root_cause_count"`
	BuildErrorsCount       int       `json:"build_errors_count"`
	DeploymentStatus       struct {
		ServiceConfig struct {
			NumSucceeded int `json:"num_succeeded"`
			NumFailed    int `json:"num_failed"`
			NumPending   int `json:"num_pending"`
		} `json:"service_config"`
		DrainConfig struct {
			NumSucceeded int `json:"num_succeeded"`
			NumFailed    int `json:"num_failed"`
			NumPending   int `json:"num_pending"`
		} `json:"drain_config"`
		Discovery2Config struct {
			NumSucceeded int `json:"num_succeeded"`
			NumFailed    int `json:"num_failed"`
			NumPending   int `json:"num_pending"`
		} `json:"discovery2_config"`
	} `json:"deployment_status"`
	AnomalyCounts struct {
		Arp                int `json:"arp"`
		Probe              int `json:"probe"`
		Hostname           int `json:"hostname"`
		Streaming          int `json:"streaming"`
		Series             int `json:"series"`
		Cabling            int `json:"cabling"`
		Route              int `json:"route"`
		Counter            int `json:"counter"`
		All                int `json:"all"`
		Bgp                int `json:"bgp"`
		BlueprintRendering int `json:"blueprint_rendering"`
		Mac                int `json:"mac"`
		Mlag               int `json:"mlag"`
		Deployment         int `json:"deployment"`
		Interface          int `json:"interface"`
		Liveness           int `json:"liveness"`
		Config             int `json:"config"`
		Lag                int `json:"lag"`
	} `json:"anomaly_counts"`
}

type CreateBluePrintFromTemplate struct {
	RefDesign  RefDesign
	Label      string
	TemplateId string
}

func (o *CreateBluePrintFromTemplate) raw() *rawCreateBluePrintFromTemplate {
	return &rawCreateBluePrintFromTemplate{
		RefDesign:  o.RefDesign.String(),
		Label:      o.Label,
		InitType:   initTypeFromTemplate,
		TemplateId: o.TemplateId,
	}
}

type rawCreateBluePrintFromTemplate struct {
	RefDesign  string `json:"design"`
	Label      string `json:"label"`
	InitType   string `json:"init_type"`
	TemplateId string `json:"template_id"`
}

type RtPolicy struct {
	// todo: what's an RtPolicy?
	//ImportRTs interface{} `json:"import_RTs"`
	//ExportRTs interface{} `json:"export_RTs"`
}

type CreateRoutingZoneCfg struct {
	SzType          string
	RoutingPolicyId string
	RtPolicy        RtPolicy
	VrfName         string
	Label           string
}

type getAllSecurityZonesResponse struct {
	Items map[string]SecurityZone
}

type SecurityZone struct {
	VniId           int      `json:"vni_id"`
	SzType          string   `json:"sz_type"`
	RoutingPolicyId ObjectId `json:"routing_policy_id"`
	Label           string   `json:"label"`
	VrfName         string   `json:"vrf_name"`
	RtPolicy        RtPolicy `json:"rt_policy"`
	RouteTarget     string   `json:"route_target"`
	Id              ObjectId `json:"id"`
	VlanId          int      `json:"vlan_id"`
}

func (o *Client) listAllBlueprintIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsBlueprintsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlBlueprints,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getBlueprintIdByName(ctx context.Context, name string) (ObjectId, error) {
	blueprintStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	// try to find the requested blueprint in the server's response
	found := -1
	for i, bps := range blueprintStatuses {
		if bps.Label == name {
			if found > 0 {
				return "", ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple blueprints have name '%s'", name),
				}
			}
			found = i
		}
	}

	// results
	if found >= 0 {
		return blueprintStatuses[found].Id, nil
	} else {
		return "", ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("found %d blueprints but one named '%s' wasn't among them", len(blueprintStatuses), name),
		}
	}
}

func (o *Client) getBlueprint(ctx context.Context, id ObjectId) (*Blueprint, error) {
	response := &Blueprint{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintById, id),
		apiResponse: response,
	})
}

func (o *Client) getBlueprintByName(ctx context.Context, name string) (*Blueprint, error) {
	id, err := o.getBlueprintIdByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return o.getBlueprint(ctx, id)
}

func (o *Client) getAllBlueprintStatus(ctx context.Context) ([]BlueprintStatus, error) {
	response := &getBluePrintsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlBlueprints,
		apiResponse: response,
	})
	return response.Items, err
}

func (o *Client) getBlueprintStatus(ctx context.Context, id ObjectId) (*BlueprintStatus, error) {
	blueprintStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	// try to find the requested blueprint
	for _, bps := range blueprintStatuses {
		if bps.Id == id {
			return &bps, nil
		}
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("found %d blueprints; '%s' wasn't among them", len(blueprintStatuses), id),
	}
}

func (o *Client) getBlueprintStatusByName(ctx context.Context, name string) (*BlueprintStatus, error) {
	blueprintStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	// try to find the requested blueprint
	found := -1
	for i, bps := range blueprintStatuses {
		if bps.Label == name {
			if found > 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple blueprints have name '%s'", name),
				}
			}
			found = i
		}
	}

	if found >= 0 {
		return &blueprintStatuses[found], nil
	} else {
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("found %d blueprints; '%s' wasn't among them", len(blueprintStatuses), name),
		}
	}
}

func (o *Client) createBlueprintFromTemplate(ctx context.Context, cfg *CreateBluePrintFromTemplate) (ObjectId, error) {
	response := &postBlueprintsResponse{}
	return response.Id, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlBlueprints,
		apiInput:    cfg.raw(),
		apiResponse: response,
	})
}

func (o *Client) deleteBlueprint(ctx context.Context, id ObjectId) error {
	log.Printf("delete blueprint id '%s'\n", id)
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintById, id),
	})
}

func (o *Client) createRoutingZone(ctx context.Context, Id ObjectId, cfg *CreateRoutingZoneCfg) (*objectIdResponse, error) {
	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRoutingZones, Id),
		apiInput:    cfg,
		apiResponse: response,
	})
}

func (o *Client) listAllRoutingZoneIds(ctx context.Context, blueprintId ObjectId) ([]ObjectId, error) {
	response := &optionsBlueprintsRoutingzonesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      fmt.Sprintf(apiUrlBlueprintById, blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getRoutingZone(ctx context.Context, blueprintId ObjectId, zoneId ObjectId) (*SecurityZone, error) {
	response := &SecurityZone{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRoutingZonesById, blueprintId, zoneId),
		apiResponse: response,
	})
}

func (o *Client) getAllRoutingZones(ctx context.Context, blueprintId ObjectId) ([]SecurityZone, error) {
	response := &getAllSecurityZonesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintById, blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	// This API endpoint returns a map. Convert to list for consistency with other 'getAll' functions.
	var result []SecurityZone
	for _, v := range response.Items {
		result = append(result, v)
	}
	return result, nil
}

func (o *Client) deleteRoutingZone(ctx context.Context, blueprintId ObjectId, zoneId ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintRoutingZonesById, blueprintId, zoneId),
	})
}

type QueryEngineQuery struct {
	Query string `json:"query"`
}

type QueryEngineResponse struct {
	Count int           `json:"count"`
	Items []interface{} `json:"items"`
}

func (o *Client) runQuery(ctx context.Context, blueprint ObjectId, in *QueryEngineQuery) (*QueryEngineResponse, error) {
	response := &QueryEngineResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintQueryEngine, blueprint),
		apiInput:    in,
		apiResponse: response,
	})

}
