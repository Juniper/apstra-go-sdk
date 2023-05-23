package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlBlueprints        = "/api/blueprints"
	apiUrlPathDelim         = "/"
	apiUrlBlueprintsPrefix  = apiUrlBlueprints + apiUrlPathDelim
	apiUrlBlueprintById     = apiUrlBlueprintsPrefix + "%s"
	apiUrlBlueprintNodes    = apiUrlBlueprintById + apiUrlPathDelim + "nodes"
	apiUrlBlueprintNodeById = apiUrlBlueprintNodes + apiUrlPathDelim + "%s"

	initTypeFromTemplate      = "template_reference"
	nodeQueryNodeTypeUrlParam = "node_type"
)

const (
	NodeTypeNone = NodeType(iota)
	NodeTypeMetadata
	NodeTypePolicy
	NodeTypeRedundancyGroup
	NodeTypeRoutingPolicy
	NodeTypeSecurityZone
	NodeTypeSystem
	NodeTypeVirtualNetwork
	NodeTypeUnknown = "unknown node type %s"

	nodeTypeNone            = nodeType("")
	nodeTypeMetadata        = nodeType("metadata")
	nodeTypePolicy          = nodeType("policy")
	nodeTypeRedundancyGroup = nodeType("redundancy_group")
	nodeTypeRoutingPolicy   = nodeType("routing_policy")
	nodeTypeSecurityZone    = nodeType("security_zone")
	nodeTypeSystem          = nodeType("system")
	nodeTypeVirtualNetwork  = nodeType("virtual_network")
	nodeTypeUnknown         = "unknown node type %d"
)

type NodeType int
type nodeType string

func (o NodeType) String() string {
	switch o {
	case NodeTypeNone:
		return string(nodeTypeNone)
	case NodeTypeMetadata:
		return string(nodeTypeMetadata)
	case NodeTypePolicy:
		return string(nodeTypePolicy)
	case NodeTypeRedundancyGroup:
		return string(nodeTypeRedundancyGroup)
	case NodeTypeRoutingPolicy:
		return string(nodeTypeRoutingPolicy)
	case NodeTypeSecurityZone:
		return string(nodeTypeSecurityZone)
	case NodeTypeSystem:
		return string(nodeTypeSystem)
	case NodeTypeVirtualNetwork:
		return string(nodeTypeVirtualNetwork)
	default:
		return fmt.Sprintf(nodeTypeUnknown, o)
	}
}

const (
	RefDesignTwoStageL3Clos = RefDesign(iota)
	RefDesignFreeform
	RefDesignDatacenter = RefDesignTwoStageL3Clos
	RefDesignUnknown    = "unknown reference design '%s'"

	refDesignDatacenter = refDesign("two_stage_l3clos")
	refDesignFreeform   = refDesign("freeform")
	refDesignUnknown    = refDesign("unknown reference design %d")
)

type RefDesign int
type refDesign string

func (o RefDesign) String() string {
	switch o {
	case RefDesignDatacenter:
		return string(refDesignDatacenter)
	case RefDesignFreeform:
		return string(refDesignFreeform)
	default:
		return fmt.Sprintf(string(refDesignUnknown), o)
	}
}

func (o *RefDesign) FromString(s string) error {
	i, err := refDesign(s).parse()
	if err != nil {
		return err
	}
	*o = i
	return nil
}

func (o refDesign) parse() (RefDesign, error) {
	switch o {
	case refDesignDatacenter:
		return RefDesignDatacenter, nil
	case refDesignFreeform:
		return RefDesignFreeform, nil
	default:
		return 0, fmt.Errorf(RefDesignUnknown, o)
	}
}

type getBluePrintsResponse struct {
	Items []rawBlueprintStatus `json:"items"`
}

type optionsBlueprintsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type postBlueprintsResponse struct {
	Id     ObjectId `json:"id"`
	TaskId TaskId   `json:"task_id"`
}

type Blueprint struct {
	client         *Client
	Id             ObjectId
	Version        int
	Design         RefDesign
	LastModifiedAt time.Time
	Label          string
	Relationships  map[string]json.RawMessage
	Nodes          map[string]json.RawMessage
	SourceVersions struct {
		ConfigBlueprint int
	}
}

type rawBlueprint struct {
	Id             ObjectId                   `json:"id"`
	Version        int                        `json:"version"`
	Design         refDesign                  `json:"design"`
	LastModifiedAt time.Time                  `json:"last_modified_at"`
	Label          string                     `json:"label"`
	Relationships  map[string]json.RawMessage `json:"relationships"`
	Nodes          map[string]json.RawMessage `json:"nodes"`
	SourceVersions struct {
		ConfigBlueprint int `json:"config_blueprint"`
	} `json:"source_versions"`
}

func (o *rawBlueprint) polish() (*Blueprint, error) {
	design, err := o.Design.parse()
	if err != nil {
		return nil, err
	}
	return &Blueprint{
		client:         nil,
		Id:             o.Id,
		Version:        o.Version,
		Design:         design,
		LastModifiedAt: o.LastModifiedAt,
		Label:          o.Label,
		Relationships:  o.Relationships,
		Nodes:          o.Nodes,
		SourceVersions: struct {
			ConfigBlueprint int
		}{ConfigBlueprint: o.SourceVersions.ConfigBlueprint},
	}, nil
}

type BlueprintDeploymentStatus struct {
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
}

type BlueprintAnomalyCounts struct {
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
}

type BlueprintStatus struct {
	Id                     ObjectId                  `json:"id"`
	Label                  string                    `json:"label"`
	Status                 string                    `json:"status"`
	Design                 RefDesign                 `json:"design"`
	HasUncommittedChanges  bool                      `json:"has_uncommitted_changes"`
	Version                int                       `json:"version"`
	LastModifiedAt         time.Time                 `json:"last_modified_at"`
	SuperspineCount        int                       `json:"superspine_count"`
	SpineCount             int                       `json:"spine_count"`
	LeafCount              int                       `json:"leaf_count"`
	AccessCount            int                       `json:"access_count"`
	GenericCount           int                       `json:"generic_count"`
	ExternalRouterCount    int                       `json:"external_router_count"`
	L2ServerCount          int                       `json:"l2_server_count"`
	L3ServerCount          int                       `json:"l3_server_count"`
	RemoteGatewayCount     int                       `json:"remote_gateway_count"`
	BuildWarningsCount     int                       `json:"build_warnings_count"`
	RootCauseCount         int                       `json:"root_cause_count"`
	TopLevelRootCauseCount int                       `json:"top_level_root_cause_count"`
	BuildErrorsCount       int                       `json:"build_errors_count"`
	DeploymentStatus       BlueprintDeploymentStatus `json:"deployment_status"`
	AnomalyCounts          BlueprintAnomalyCounts    `json:"anomaly_counts"`
}

type rawBlueprintStatus struct {
	Id                     ObjectId                  `json:"id"`
	Label                  string                    `json:"label"`
	Status                 string                    `json:"status"`
	Design                 refDesign                 `json:"design"`
	HasUncommittedChanges  bool                      `json:"has_uncommitted_changes"`
	Version                int                       `json:"version"`
	LastModifiedAt         time.Time                 `json:"last_modified_at"`
	SuperspineCount        int                       `json:"superspine_count"`
	SpineCount             int                       `json:"spine_count"`
	LeafCount              int                       `json:"leaf_count"`
	AccessCount            int                       `json:"access_count"`
	GenericCount           int                       `json:"generic_count"`
	ExternalRouterCount    int                       `json:"external_router_count"`
	L2ServerCount          int                       `json:"l2_server_count"`
	L3ServerCount          int                       `json:"l3_server_count"`
	RemoteGatewayCount     int                       `json:"remote_gateway_count"`
	BuildWarningsCount     int                       `json:"build_warnings_count"`
	BuildErrorsCount       int                       `json:"build_errors_count"`
	RootCauseCount         int                       `json:"root_cause_count"`
	TopLevelRootCauseCount int                       `json:"top_level_root_cause_count"`
	DeploymentStatus       BlueprintDeploymentStatus `json:"deployment_status"`
	AnomalyCounts          BlueprintAnomalyCounts    `json:"anomaly_counts"`
	// todo 4.1.1 introduced (?) the following:
	//   "deploy_modes_summary": {
	//     "ready": 0,
	//     "undeploy": 0,
	//     "drain": 0,
	//     "deploy": 0
	//   },

}

func (o *rawBlueprintStatus) polish() (*BlueprintStatus, error) {
	design, err := o.Design.parse()
	if err != nil {
		return nil, err
	}
	return &BlueprintStatus{
		Id:                     o.Id,
		Label:                  o.Label,
		Status:                 o.Status,
		Design:                 design,
		HasUncommittedChanges:  o.HasUncommittedChanges,
		Version:                o.Version,
		LastModifiedAt:         o.LastModifiedAt,
		SuperspineCount:        o.SuperspineCount,
		SpineCount:             o.SpineCount,
		LeafCount:              o.LeafCount,
		AccessCount:            o.AccessCount,
		GenericCount:           o.GenericCount,
		ExternalRouterCount:    o.ExternalRouterCount,
		L2ServerCount:          o.L2ServerCount,
		L3ServerCount:          o.L3ServerCount,
		RemoteGatewayCount:     o.RemoteGatewayCount,
		BuildWarningsCount:     o.BuildWarningsCount,
		RootCauseCount:         o.RootCauseCount,
		TopLevelRootCauseCount: o.TopLevelRootCauseCount,
		BuildErrorsCount:       o.BuildErrorsCount,
		DeploymentStatus:       BlueprintDeploymentStatus{},
		AnomalyCounts:          BlueprintAnomalyCounts{},
	}, nil
}

type CreateBlueprintFromTemplateRequest struct {
	RefDesign              RefDesign
	Label                  string
	TemplateId             ObjectId
	FabricAddressingPolicy *FabricAddressingPolicy
}

func (o *CreateBlueprintFromTemplateRequest) raw() *rawCreateBlueprintFromTemplateRequest {
	var fap *rawFabricAddressingPolicy
	if o.FabricAddressingPolicy != nil {
		fap = o.FabricAddressingPolicy.raw()
	}
	return &rawCreateBlueprintFromTemplateRequest{
		RefDesign:              o.RefDesign.String(),
		Label:                  o.Label,
		InitType:               initTypeFromTemplate,
		TemplateId:             o.TemplateId,
		FabricAddressingPolicy: fap,
	}
}

type rawCreateBlueprintFromTemplateRequest struct {
	RefDesign              string                     `json:"design"`
	Label                  string                     `json:"label"`
	InitType               string                     `json:"init_type"`
	TemplateId             ObjectId                   `json:"template_id"`
	FabricAddressingPolicy *rawFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
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
	}
	return "", ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("found %d blueprints but one named '%s' wasn't among them", len(blueprintStatuses), name),
	}
}

func (o *Client) getBlueprint(ctx context.Context, id ObjectId) (*Blueprint, error) {
	response := &rawBlueprint{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.polish()
}

//lint:ignore U1000 keep for future
func (o *Client) getBlueprintByName(ctx context.Context, name string) (*Blueprint, error) {
	id, err := o.getBlueprintIdByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return o.getBlueprint(ctx, id)
}

func (o *Client) getAllBlueprintStatus(ctx context.Context) ([]rawBlueprintStatus, error) {
	response := &getBluePrintsResponse{}
	return response.Items, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlBlueprints,
		apiResponse: response,
	})
}

func (o *Client) getBlueprintStatus(ctx context.Context, id ObjectId) (*rawBlueprintStatus, error) {
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
		err:     fmt.Errorf("found %d blueprints but one with id '%s' wasn't among them", len(blueprintStatuses), id),
	}
}

func (o *Client) getBlueprintStatusesByName(ctx context.Context, desired string) ([]rawBlueprintStatus, error) {
	blueprintStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	i := 0
	for i < len(blueprintStatuses) {
		if blueprintStatuses[i].Label != desired { // element not desired. delete element.
			// copy last element to current position
			blueprintStatuses[i] = blueprintStatuses[len(blueprintStatuses)-1]
			// delete last element
			blueprintStatuses = blueprintStatuses[:len(blueprintStatuses)-1]
		} else {
			i++
		}
	}
	return blueprintStatuses, nil
}

func (o *Client) getBlueprintStatusByName(ctx context.Context, desired string) (*rawBlueprintStatus, error) {
	blueprintStatuses, err := o.getBlueprintStatusesByName(ctx, desired)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	switch len(blueprintStatuses) {
	case 0:
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("blueprint with name '%s' not found", desired),
		}
	case 1:
		return &blueprintStatuses[0], nil
	default:
		return nil, ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("multiple blueprints with name '%s' found", desired),
		}
	}
}

func (o *Client) createBlueprintFromTemplate(ctx context.Context, req *rawCreateBlueprintFromTemplateRequest) (ObjectId, error) {
	response := &postBlueprintsResponse{}
	return response.Id, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlBlueprints,
		apiInput:    req,
		apiResponse: response,
	})
}

func (o *Client) deleteBlueprint(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintById, id),
	})
}

func (o *Client) runQuery(ctx context.Context, blueprint ObjectId, query QEQuery, response interface{}) error {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintQueryEngine, blueprint))
	if err != nil {
		return err
	}

	if query.getBlueprintType() != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, query.getBlueprintType().string())
		apstraUrl.RawQuery = params.Encode()
	}

	apiInput := &struct {
		Query string `json:"query"`
	}{Query: query.String()}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPost,
		url:            apstraUrl,
		apiInput:       apiInput,
		apiResponse:    response,
		unsynchronized: true,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *Client) getNodes(ctx context.Context, blueprint ObjectId, nodeType NodeType, response interface{}) error {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintNodes, blueprint))
	if err != nil {
		return err
	}

	if nodeType != NodeTypeNone {
		params := apstraUrl.Query()
		params.Set(nodeQueryNodeTypeUrlParam, nodeType.String())
		apstraUrl.RawQuery = params.Encode()
	}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

func (o *Client) patchNode(ctx context.Context, blueprint ObjectId, node ObjectId, request interface{}, response interface{}) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPatch,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeById, blueprint, node),
		apiInput:    request,
		apiResponse: response,
	})
}
