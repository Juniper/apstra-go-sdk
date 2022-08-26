package goapstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlBlueprints          = "/api/blueprints"
	apiUrlPathDelim           = "/"
	apiUrlBlueprintsPrefix    = apiUrlBlueprints + apiUrlPathDelim
	apiUrlBlueprintById       = apiUrlBlueprintsPrefix + "%s"
	apiUrlBlueprintNodes      = apiUrlBlueprintById + apiUrlPathDelim + "nodes"
	apiUrlBlueprintNodeById   = apiUrlBlueprintNodes + apiUrlPathDelim + "%s"
	apiUrlBlueprintLockStatus = apiUrlBlueprintById + apiUrlPathDelim + "lock-status"

	initTypeFromTemplate      = "template_reference"
	nodeQueryNodeTypeUrlParam = "node_type"
)

type LockStatus int
type lockStatus string

const (
	LockStatusUnlocked = LockStatus(iota)
	LockStatusLockedByRestrictedUser
	LockStatusLockedByAdmin
	LockStatusLockedByDeletedUser
	LockStatusUnknown = "unknown lock status %s"

	lockStatusUnlocked               = lockStatus("unlocked")
	lockStatusLockedByRestrictedUser = lockStatus("locked_by_restricted_user")
	lockStatusLockedByAdmin          = lockStatus("locked_by_admin")
	lockStatusLockedByDeletedUser    = lockStatus("locked_by_deleted_user")
	lockStatusUnknown                = "unknown lock status %d"
)

func (o LockStatus) String() string {
	switch o {
	case LockStatusUnlocked:
		return string(lockStatusUnlocked)
	case LockStatusLockedByRestrictedUser:
		return string(lockStatusLockedByRestrictedUser)
	case LockStatusLockedByAdmin:
		return string(lockStatusLockedByAdmin)
	case LockStatusLockedByDeletedUser:
		return string(lockStatusLockedByDeletedUser)
	default:
		return fmt.Sprintf(lockStatusUnknown, o)
	}
}

func (o lockStatus) parse() (LockStatus, error) {
	switch o {
	case lockStatusUnlocked:
		return LockStatusUnlocked, nil
	case lockStatusLockedByRestrictedUser:
		return LockStatusLockedByRestrictedUser, nil
	case lockStatusLockedByAdmin:
		return LockStatusLockedByAdmin, nil
	case lockStatusLockedByDeletedUser:
		return LockStatusLockedByDeletedUser, nil
	default:
		return 0, fmt.Errorf(LockStatusUnknown, o)
	}
}

const (
	NodeTypeNone = NodeType(iota)
	NodeTypeMetadata
	NodeTypeSystem
	NodeTypeUnknown = "unknown node type %s"

	nodeTypeNone     = nodeType("")
	nodeTypeMetadata = nodeType("metadata")
	nodeTypeSystem   = nodeType("system")
	nodeTypeUnknown  = "unknown node type %d"
)

type NodeType int
type nodeType string

func (o NodeType) String() string {
	switch o {
	case NodeTypeNone:
		return string(nodeTypeNone)
	case NodeTypeMetadata:
		return string(nodeTypeMetadata)
	case NodeTypeSystem:
		return string(nodeTypeSystem)
	default:
		return fmt.Sprintf(nodeTypeUnknown, o)
	}
}

func (o nodeType) parse() (NodeType, error) {
	switch o {
	case nodeTypeNone:
		return NodeTypeNone, nil
	case nodeTypeMetadata:
		return NodeTypeMetadata, nil
	case nodeTypeSystem:
		return NodeTypeSystem, nil
	default:
		return 0, fmt.Errorf(NodeTypeUnknown, o)
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

type CreateBlueprintFromTemplate struct {
	RefDesign  RefDesign
	Label      string
	TemplateId ObjectId
}

func (o *CreateBlueprintFromTemplate) raw() *rawCreateBluePrintFromTemplate {
	return &rawCreateBluePrintFromTemplate{
		RefDesign:  o.RefDesign.String(),
		Label:      o.Label,
		InitType:   initTypeFromTemplate,
		TemplateId: o.TemplateId,
	}
}

type rawCreateBluePrintFromTemplate struct {
	RefDesign  string   `json:"design"`
	Label      string   `json:"label"`
	InitType   string   `json:"init_type"`
	TemplateId ObjectId `json:"template_id"`
}

type rawLockInfo struct {
	UserName         string     `json:"username"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	UserId           string     `json:"user_id"`
	PossibleOverride bool       `json:"possible_override"`
	LockStatus       lockStatus `json:"lock_status"`
}

func (o *rawLockInfo) polish() (*LockInfo, error) {
	ls, err := o.LockStatus.parse()
	if err != nil {
		return nil, err
	}
	return &LockInfo{
		UserName:         o.UserName,
		FirstName:        o.FirstName,
		LastName:         o.LastName,
		UserId:           o.UserId,
		PossibleOverride: o.PossibleOverride,
		LockStatus:       ls,
	}, nil
}

type LockInfo struct {
	UserName         string
	FirstName        string
	LastName         string
	UserId           string
	PossibleOverride bool
	LockStatus       LockStatus
}

func (o *LockInfo) raw() *rawLockInfo {
	ls := lockStatus(o.LockStatus.String())
	return &rawLockInfo{
		UserName:         o.UserName,
		FirstName:        o.FirstName,
		LastName:         o.LastName,
		UserId:           o.UserId,
		PossibleOverride: o.PossibleOverride,
		LockStatus:       ls,
	}
}

func (o *Client) getBlueprintLockInfo(ctx context.Context, id ObjectId) (*LockInfo, error) {
	response := &rawLockInfo{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintLockStatus, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.polish()
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
	if err != nil {
		return nil, err
	}
	result := make([]BlueprintStatus, len(response.Items))
	for i, item := range response.Items {
		p, err := item.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *p
	}
	return result, nil
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
		err:     fmt.Errorf("found %d blueprints but one with id '%s' wasn't among them", len(blueprintStatuses), id),
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
			err:     fmt.Errorf("found %d blueprints but one named '%s' wasn't among them", len(blueprintStatuses), name),
		}
	}
}

func (o *Client) createBlueprintFromTemplate(ctx context.Context, cfg *CreateBlueprintFromTemplate) (ObjectId, error) {
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

func (o *Client) runQuery(ctx context.Context, blueprint ObjectId, query *QEQuery, response interface{}) error {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintQueryEngine, blueprint))
	if err != nil {
		return err
	}

	if query.blueprintType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, query.blueprintType.string())
		apstraUrl.RawQuery = params.Encode()
	}

	apiInput := &struct {
		Query string `json:"query"`
	}{Query: query.string()}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPost,
		url:            apstraUrl,
		apiInput:       apiInput,
		apiResponse:    response,
		unsynchronized: true,
	})
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
