package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlBlueprints                  = "/api/blueprints"
	apiUrlPathDelim                   = "/"
	apiUrlBlueprintsPrefix            = apiUrlBlueprints + apiUrlPathDelim
	apiUrlBlueprintById               = apiUrlBlueprintsPrefix + "%s"
	apiUrlBluePrintRoutingZones       = apiUrlBlueprintById + "/security_zones"
	apiUrlBluePrintRoutingZonesPrefix = apiUrlBluePrintRoutingZones + apiUrlPathDelim
	apiUrlBluePrintRoutingZonesById   = apiUrlBluePrintRoutingZonesPrefix + "%s"
)

type optionsBlueprintsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type optionsBlueprintsRoutingzonesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

// getBlueprintsResponse is returned by Apstra in response to
// 'GET apiUrlBlueprints'
type getBlueprintsResponse struct {
	Items []struct {
		Status           string `json:"status"`
		Version          int    `json:"version"`
		Design           string `json:"design"`
		DeploymentStatus struct {
			AdditionalProp1 struct {
				NumSucceeded int `json:"num_succeeded"`
				NumFailed    int `json:"num_failed"`
				NumPending   int `json:"num_pending"`
			} `json:"additionalProp1"`
			AdditionalProp2 struct {
				NumSucceeded int `json:"num_succeeded"`
				NumFailed    int `json:"num_failed"`
				NumPending   int `json:"num_pending"`
			} `json:"additionalProp2"`
			AdditionalProp3 struct {
				NumSucceeded int `json:"num_succeeded"`
				NumFailed    int `json:"num_failed"`
				NumPending   int `json:"num_pending"`
			} `json:"additionalProp3"`
		} `json:"deployment_status"`
		AnomalyCounts struct {
			AdditionalProp1 int `json:"additionalProp1"`
			AdditionalProp2 int `json:"additionalProp2"`
			AdditionalProp3 int `json:"additionalProp3"`
		} `json:"anomaly_counts"`
		Id             ObjectId  `json:"id"`
		LastModifiedAt time.Time `json:"last_modified_at"`
		Label          string    `json:"label"`
	} `json:"items"`
}

// GetBlueprintResponse is returned by Apstra in response to
// 'GET apiUrlBlueprintsPrefix + <id>'
type GetBlueprintResponse struct {
	Relationships struct {
		AdditionalProp1 struct {
			SourceId string `json:"source_id"`
			TargetId string `json:"target_id"`
			Type     string `json:"type"`
			Id       string `json:"id"`
		} `json:"additionalProp1"`
		AdditionalProp2 struct {
			SourceId string `json:"source_id"`
			TargetId string `json:"target_id"`
			Type     string `json:"type"`
			Id       string `json:"id"`
		} `json:"additionalProp2"`
		AdditionalProp3 struct {
			SourceId string `json:"source_id"`
			TargetId string `json:"target_id"`
			Type     string `json:"type"`
			Id       string `json:"id"`
		} `json:"additionalProp3"`
	} `json:"relationships"`
	Version        int       `json:"version"`
	Design         string    `json:"design"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	Nodes          struct {
		AdditionalProp1 struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"additionalProp1"`
		AdditionalProp2 struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"additionalProp2"`
		AdditionalProp3 struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"additionalProp3"`
	} `json:"nodes"`
	Id             string `json:"id"`
	SourceVersions struct {
		AdditionalProp1 int `json:"additionalProp1"`
		AdditionalProp2 int `json:"additionalProp2"`
		AdditionalProp3 int `json:"additionalProp3"`
	} `json:"source_versions"`
	Label string `json:"label"`
}

func (o Client) getBluePrints(ctx context.Context) (*getBlueprintsResponse, error) {
	response := &getBlueprintsResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlBlueprints,
		apiResponse: response,
	})
}

func (o Client) getBlueprint(ctx context.Context, id ObjectId) (*GetBlueprintResponse, error) {
	response := &GetBlueprintResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintById, id),
		apiResponse: response,
	})
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

func (o *Client) createRoutingZone(ctx context.Context, Id ObjectId, cfg *CreateRoutingZoneCfg) (*objectIdResponse, error) {
	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBluePrintRoutingZones, Id),
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
		urlStr:      fmt.Sprintf(apiUrlBluePrintRoutingZonesById, blueprintId, zoneId),
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
		urlStr: fmt.Sprintf(apiUrlBluePrintRoutingZonesById, blueprintId, zoneId),
	})
}
