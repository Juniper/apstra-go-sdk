package aosSdk

import "time"

const (
	apiUrlBlueprints        = "/api/blueprints"
	apiUrlRoutingZonePrefix = apiUrlBlueprints + "/"
	apiUrlRoutingZoneSuffix = "/security-zones/"
	apiUrlBlueprintsPrefix  = apiUrlBlueprints + "/"
)

// getBlueprintsResult is returned by Apstra in response to
// 'GET apiUrlBlueprints'
type getBlueprintsResult struct {
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
	var response getBlueprintsResult
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlBlueprints,
		fromServerPtr: &response,
	})
	if err != nil {
		return nil, err
	}
	var result []ObjectId
	for _, item := range response.Items {
		result = append(result, item.Id)
	}
	return result, nil
}

func (o Client) getBlueprint(in ObjectId) (*GetBlueprintResponse, error) {
	var response GetBlueprintResponse
	return &response, o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlBlueprintsPrefix + string(in),
		fromServerPtr: &response,
	})
}

type RtPolicy struct {
	// todo
	//ImportRTs interface{} `json:"import_RTs"`
	//ExportRTs interface{} `json:"export_RTs"`
}

type CreateRoutingZoneCfg struct {
	SzType          string   `json:"sz_type,omitempty"`
	RoutingPolicyId string   `json:"routing_policy_id,omitempty"`
	RtPolicy        RtPolicy `json:"rt_policy,omitempty"`
	VrfName         string   `json:"vrf_name,omitempty"`
	Label           string   `json:"label,omitempty"`
}

type GetRoutingZoneResult struct {
	VniId           int      `json:"vni_id"`
	SzType          string   `json:"sz_type"`
	RoutingPolicyId string   `json:"routing_policy_id"`
	Label           string   `json:"label"`
	VrfName         string   `json:"vrf_name"`
	RtPolicy        RtPolicy `json:"rt_policy"`
	RouteTarget     string   `json:"route_target"`
	Id              string   `json:"id"`
	VlanId          int      `json:"vlan_id"`
}

type getAllRoutingZonesResult struct {
	Items map[string]GetRoutingZoneResult
}

func (o Client) CreateRoutingZone(blueprintId string, cfg *CreateRoutingZoneCfg) (*GetRoutingZoneResult, error) {
	result := &GetRoutingZoneResult{}
	return result, o.talkToAos(&talkToAosIn{
		method:        httpMethodPost,
		url:           apiUrlRoutingZonePrefix + blueprintId + apiUrlRoutingZoneSuffix,
		toServerPtr:   cfg,
		fromServerPtr: result,
	})
}

func (o Client) DeleteRoutingZone(blueprintId string, zoneId string) error {
	return o.talkToAos(&talkToAosIn{
		method: httpMethodDelete,
		url:    apiUrlRoutingZonePrefix + blueprintId + apiUrlRoutingZonePrefix + zoneId,
	})
}
