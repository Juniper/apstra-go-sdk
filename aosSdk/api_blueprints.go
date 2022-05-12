package aosSdk

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlBlueprints        = "/api/blueprints"
	apiUrlPathDelim         = "/"
	apiUrlBlueprintsPrefix  = apiUrlBlueprints + apiUrlPathDelim
	apiUrlRoutingZonePrefix = apiUrlBlueprintsPrefix
	apiUrlRoutingZoneSuffix = "/security-zones/"
)

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

// getAllBlueprintIds returns the Ids of all blueprints
func (o Client) getAllBlueprintIds() ([]ObjectId, error) {
	response, err := o.getBluePrints()
	if err != nil {
		return nil, fmt.Errorf("error calling getBluePrints - %w", err)
	}

	var result []ObjectId
	for _, item := range response.Items {
		result = append(result, item.Id)
	}
	return result, nil
}

func (o Client) getBluePrints() (*getBlueprintsResponse, error) {
	aosUrl, err := url.Parse(apiUrlBlueprints)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlBlueprints, err)
	}
	response := &getBlueprintsResponse{}
	return response, o.talkToAos(&talkToAosIn{
		method:      http.MethodGet,
		url:         aosUrl,
		apiResponse: response,
	})
}

func (o Client) getBlueprint(in ObjectId) (*GetBlueprintResponse, error) {
	aosUrl, err := url.Parse(apiUrlBlueprintsPrefix + string(in))
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlBlueprints+string(in), err)
	}
	response := &GetBlueprintResponse{}
	return response, o.talkToAos(&talkToAosIn{
		method:      http.MethodGet,
		url:         aosUrl,
		apiResponse: response,
	})
}

type RtPolicy struct {
	// todo: what's an RtPolicy?
	//ImportRTs interface{} `json:"import_RTs"`
	//ExportRTs interface{} `json:"export_RTs"`
}

// createRoutingZoneRequest doesn't appear in swagger, but shows up in Ryan Booth's
// postman collection:
//   https://www.getpostman.com/collections/6ad3bd003d83e4cba47b
// It's sent to {{aos_server_api}}/blueprints/{{blueprint_id}}/security-zones/
// via POST
type createRoutingZoneRequest struct {
	SzType          string   `json:"sz_type,omitempty"`
	RoutingPolicyId string   `json:"routing_policy_id,omitempty"`
	RtPolicy        RtPolicy `json:"rt_policy,omitempty"`
	VrfName         string   `json:"vrf_name,omitempty"`
	Label           string   `json:"label,omitempty"`
}

// CreateRoutingZoneCfg is the public version of createRoutingZoneRequest. The
// difference being that it includes the BlueprintId, which is required in the
// URL path when calling the API
type CreateRoutingZoneCfg struct {
	SzType          string
	RoutingPolicyId string
	RtPolicy        RtPolicy
	VrfName         string
	Label           string
	BlueprintId     ObjectId
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

func (o Client) createRoutingZone(cfg *CreateRoutingZoneCfg) (*objectIdResponse, error) {
	aosUrl, err := url.Parse(apiUrlRoutingZonePrefix + string(cfg.BlueprintId) + apiUrlRoutingZoneSuffix)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlRoutingZonePrefix+string(cfg.BlueprintId)+apiUrlRoutingZoneSuffix, err)
	}
	toServer := &createRoutingZoneRequest{
		SzType:          cfg.SzType,
		RoutingPolicyId: cfg.RoutingPolicyId,
		RtPolicy:        cfg.RtPolicy,
		VrfName:         cfg.VrfName,
		Label:           cfg.Label,
	}
	result := &objectIdResponse{}
	return result, o.talkToAos(&talkToAosIn{
		method:      http.MethodPost,
		url:         aosUrl,
		apiInput:    toServer,
		apiResponse: result,
	})
}
func (o Client) getRoutingZone(blueprintId ObjectId, zone ObjectId) (*SecurityZone, error) {
	urlString := apiUrlRoutingZonePrefix + string(blueprintId) + apiUrlRoutingZoneSuffix + string(zone)
	aosUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlString, err)
	}
	result := &SecurityZone{}
	return result, o.talkToAos(&talkToAosIn{
		method:      http.MethodGet,
		url:         aosUrl,
		apiInput:    nil,
		apiResponse: result,
		doNotLogin:  false,
	})
}

func (o Client) getAllRoutingZones(blueprintId ObjectId) ([]SecurityZone, error) {
	urlString := apiUrlRoutingZonePrefix + string(blueprintId) + apiUrlRoutingZoneSuffix
	aosUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlString, err)
	}
	response := &getAllSecurityZonesResponse{}
	err = o.talkToAos(&talkToAosIn{
		method:      http.MethodGet,
		url:         aosUrl,
		apiInput:    nil,
		apiResponse: response,
		doNotLogin:  false,
	})
	if err != nil {
		return nil, err
	}
	var result []SecurityZone
	for _, v := range response.Items {
		result = append(result, v)
	}
	return result, nil
}

func (o Client) deleteRoutingZone(blueprintId ObjectId, zoneId ObjectId) error {
	aosUrl, err := url.Parse(apiUrlRoutingZonePrefix + string(blueprintId) + apiUrlRoutingZoneSuffix + string(zoneId))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlRoutingZonePrefix+string(blueprintId)+apiUrlRoutingZoneSuffix+string(zoneId), err)
	}
	return o.talkToAos(&talkToAosIn{
		method: http.MethodDelete,
		url:    aosUrl,
	})
}
