package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintSecurityZones       = apiUrlBlueprintById + apiUrlPathDelim + "security-zones"
	apiUrlBlueprintSecurityZonesPrefix = apiUrlBlueprintSecurityZones + apiUrlPathDelim
	apiUrlBlueprintSecurityZoneById    = apiUrlBlueprintSecurityZonesPrefix + "%s"
)

type RtPolicy struct {
	// todo: what's an RtPolicy?
	//ImportRTs interface{} `json:"import_RTs"`
	//ExportRTs interface{} `json:"export_RTs"`
}

type CreateSecurityZoneCfg struct {
	SzType          string   `json:"sz_type"`
	RoutingPolicyId string   `json:"routing_policy_id,omitempty"`
	RtPolicy        RtPolicy `json:"rt_policy"`
	VrfName         string   `json:"vrf_name"`
	Label           string   `json:"label"`
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
	VlanId          VLAN     `json:"vlan_id"`
}

func (o *TwoStageL3ClosClient) createSecurityZone(ctx context.Context, cfg *CreateSecurityZoneCfg) (*objectIdResponse, error) {
	response := &objectIdResponse{}
	return response, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})
}

func (o *TwoStageL3ClosClient) getSecurityZone(ctx context.Context, zoneId ObjectId) (*SecurityZone, error) {
	response := &SecurityZone{}
	return response, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
		apiResponse: response,
	})
}

func (o *TwoStageL3ClosClient) getSecurityZoneByName(ctx context.Context, vrfName string) (*SecurityZone, error) {
	zones, err := o.getAllSecurityZones(ctx)
	if err != nil {
		return nil, err
	}

	for _, zone := range zones {
		if zone.VrfName == vrfName {
			return o.getSecurityZone(ctx, zone.Id)
		}
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with vrf name '%s' in blueprint '%s' not found", vrfName, o.blueprintId),
	}
}

func (o *TwoStageL3ClosClient) getAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	response := &struct {
		Items map[string]SecurityZone `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
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

func (o *TwoStageL3ClosClient) updateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *CreateSecurityZoneCfg) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
		apiInput: cfg,
	})

}

func (o *TwoStageL3ClosClient) deleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
	})
}
