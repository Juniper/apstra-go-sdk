// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprints              = "/api/blueprints"
	apiUrlPathDelim               = "/"
	apiUrlBlueprintsPrefix        = apiUrlBlueprints + apiUrlPathDelim
	apiUrlBlueprintById           = apiUrlBlueprintsPrefix + "%s"
	apiUrlBlueprintByIdPrefix     = apiUrlBlueprintById + apiUrlPathDelim
	apiUrlBlueprintNodes          = apiUrlBlueprintById + apiUrlPathDelim + "nodes"
	apiUrlBlueprintNodeById       = apiUrlBlueprintNodes + apiUrlPathDelim + "%s"
	apiUrlBlueprintNodeByIdPrefix = apiUrlBlueprintNodeById + apiUrlPathDelim

	initTypeFromTemplate   = "template_reference"
	initTypeTemplateInline = "rack_based_template_inline"

	nodeQueryNodeTypeUrlParam = "node_type"

	cablingMapMaxWaitSec = 30
)

type BlueprintRequestFabricAddressingPolicy struct {
	SpineSuperspineLinks AddressingScheme
	SpineLeafLinks       AddressingScheme
	FabricL3Mtu          *uint16
}

func (o *BlueprintRequestFabricAddressingPolicy) raw() *rawBlueprintRequestFabricAddressingPolicy {
	var fabricL3Mtu *uint16
	if o.FabricL3Mtu != nil {
		t := *o.FabricL3Mtu // copy the pointed-to value
		fabricL3Mtu = &t
	}

	return &rawBlueprintRequestFabricAddressingPolicy{
		SpineSuperspineLinks: o.SpineSuperspineLinks.raw(),
		SpineLeafLinks:       o.SpineLeafLinks.raw(),
		FabricL3Mtu:          fabricL3Mtu,
	}
}

type rawBlueprintRequestFabricAddressingPolicy struct {
	SpineSuperspineLinks addressingScheme `json:"spine_superspine_links,omitempty"`
	SpineLeafLinks       addressingScheme `json:"spine_leaf_links,omitempty"`
	FabricL3Mtu          *uint16          `json:"fabric_l3_mtu,omitempty"`
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
	Design         enum.RefDesign
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
	Design         enum.RefDesign             `json:"design"`
	LastModifiedAt time.Time                  `json:"last_modified_at"`
	Label          string                     `json:"label"`
	Relationships  map[string]json.RawMessage `json:"relationships"`
	Nodes          map[string]json.RawMessage `json:"nodes"`
	SourceVersions struct {
		ConfigBlueprint int `json:"config_blueprint"`
	} `json:"source_versions"`
}

func (o *rawBlueprint) polish() (*Blueprint, error) {
	return &Blueprint{
		client:         nil,
		Id:             o.Id,
		Version:        o.Version,
		Design:         o.Design,
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
	Design                 enum.RefDesign            `json:"design"`
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
	Design                 enum.RefDesign            `json:"design"`
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
	return &BlueprintStatus{
		Id:                     o.Id,
		Label:                  o.Label,
		Status:                 o.Status,
		Design:                 o.Design,
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
	RefDesign                 enum.RefDesign
	Label                     string
	TemplateId                ObjectId
	FabricSettings            *FabricSettings
	SkipCablingReadinessCheck bool
}

func (o *CreateBlueprintFromTemplateRequest) raw() *rawCreateBlueprintFromTemplateRequest {
	var fabricSettings *rawFabricSettings
	if o.FabricSettings != nil {
		fabricSettings = o.FabricSettings.raw()
	}
	return &rawCreateBlueprintFromTemplateRequest{
		RefDesign:      o.RefDesign.String(),
		Label:          o.Label,
		InitType:       initTypeFromTemplate,
		TemplateId:     o.TemplateId,
		FabricSettings: fabricSettings,
	}
}

func (o *CreateBlueprintFromTemplateRequest) raw420() *rawCreateBlueprintFromTemplateRequest420 {
	var fabricAddressingPolicy *rawBlueprintRequestFabricAddressingPolicy
	if o.FabricSettings != nil {
		fabricAddressingPolicy = o.FabricSettings.rawBlueprintRequestFabricAddressingPolicy()
	}
	return &rawCreateBlueprintFromTemplateRequest420{
		RefDesign:              o.RefDesign.String(),
		Label:                  o.Label,
		InitType:               initTypeFromTemplate,
		TemplateId:             o.TemplateId,
		FabricAddressingPolicy: fabricAddressingPolicy,
	}
}

type rawCreateBlueprintFromTemplateRequest struct {
	RefDesign      string             `json:"design"`
	Label          string             `json:"label"`
	InitType       string             `json:"init_type"`
	TemplateId     ObjectId           `json:"template_id"`
	FabricSettings *rawFabricSettings `json:"fabric_policy,omitempty"`
}

type rawCreateBlueprintFromTemplateRequest420 struct {
	RefDesign              string                                     `json:"design"`
	Label                  string                                     `json:"label"`
	InitType               string                                     `json:"init_type"`
	TemplateId             ObjectId                                   `json:"template_id"`
	FabricAddressingPolicy *rawBlueprintRequestFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
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
				return "", ClientErr{
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
	return "", ClientErr{
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
	var response getBluePrintsResponse
	var errs []error

	for i := range o.GetTuningParam("BlueprintStatusMaxRetries") {
		err := o.talkToApstra(ctx, &talkToApstraIn{
			method:      http.MethodGet,
			urlStr:      apiUrlBlueprints,
			apiResponse: &response,
		})
		if err == nil { // success!
			return response.Items, nil
		}

		// we got an error
		err = convertTtaeToAceWherePossible(err)
		var ace ClientErr
		if errors.As(err, &ace) && ace.IsRetryable() {
			// AOS-45313 issue?
			errs = append(errs, fmt.Errorf("retryable error at attempt %d while fetching blueprint status - %w", i, err))
			time.Sleep(time.Duration(o.GetTuningParam("BlueprintStatusRetryIntervalMs")) * time.Millisecond)
			continue
		} else {
			errs = append(errs, fmt.Errorf("non-retryable error at attempt %d while fetching blueprint status - %w", i, err))
			break
		}
	}

	return nil, errors.Join(errs...)
}

func (o *Client) getBlueprintStatus(ctx context.Context, id ObjectId) (*rawBlueprintStatus, error) {
	blueprintStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return nil, err
	}

	// try to find the requested blueprint
	for _, bps := range blueprintStatuses {
		if bps.Id == id {
			return &bps, nil
		}
	}
	return nil, ClientErr{
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
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("blueprint with name '%s' not found", desired),
		}
	case 1:
		return &blueprintStatuses[0], nil
	default:
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("multiple blueprints with name '%s' found", desired),
		}
	}
}

func (o *Client) waitForBlueprintCabling(ctx context.Context, bpId ObjectId) error {
	var ace ClientErr
	deadline := time.Now().Add(cablingMapMaxWaitSec * time.Second)

	for {
		// create a new bp client with every loop iteration so we can be sure the blueprint exists
		bp, err := o.NewTwoStageL3ClosClient(ctx, bpId)
		if err != nil {
			return err
		}

		// try to fetch the cabling map
		_, err = bp.GetCablingMapLinks(ctx)
		if err == nil {
			return nil // success! we're done here.
		}
		if errors.As(err, &ace) && ace.Type() == ErrNotfound {
			// Have we been trying too long?
			if time.Now().After(deadline) {
				return ClientErr{
					errType: ErrTimeout,
					err:     fmt.Errorf("timed out waiting for blueprint %q cabling map to become available", bpId),
				}
			}

			// 404 is likely a transient error. try again after delay.
			time.Sleep(clientPollingIntervalMs * time.Millisecond)
			continue
		}

		return err // any other error is fatal
	}
}

func (o *Client) createBlueprintFromTemplate(ctx context.Context, req *rawCreateBlueprintFromTemplateRequest) (ObjectId, error) {
	response := &postBlueprintsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlBlueprints,
		apiInput:    req,
		apiResponse: response,
	})
	return response.Id, convertTtaeToAceWherePossible(err)
}

func (o *Client) createBlueprintFromTemplate420(ctx context.Context, req *rawCreateBlueprintFromTemplateRequest420) (ObjectId, error) {
	response := &postBlueprintsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlBlueprints,
		apiInput:    req,
		apiResponse: response,
	})
	return response.Id, convertTtaeToAceWherePossible(err)
}

func (o *Client) deleteBlueprint(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintById, id),
	})
	return convertTtaeToAceWherePossible(err)
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

	// talkToApstra will copy the http response into this buffer
	httpBody := new(bytes.Buffer)

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPost,
		url:            apstraUrl,
		apiInput:       apiInput,
		httpBodyWriter: httpBody,
		unsynchronized: true,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	// save the raw http response into the query
	query.setRawResult(httpBody.Bytes())

	// unpack the http response into the caller-supplied struct, if any
	if response != nil {
		err = json.NewDecoder(httpBody).Decode(response)
		if err != nil {
			return fmt.Errorf("error while decoding API query response %q - %w", httpBody.String(), err)
		}
	}

	return nil
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

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *Client) getNode(ctx context.Context, blueprintId ObjectId, nodeId ObjectId, target interface{}) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeById, blueprintId, nodeId),
		apiResponse: target,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *Client) patchNode(ctx context.Context, blueprint ObjectId, node ObjectId, request interface{}, response interface{}, unsafe bool) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPatch,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeById, blueprint, node),
		apiInput:    request,
		apiResponse: response,
		unsafe:      unsafe,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *Client) patchNodes(ctx context.Context, blueprint ObjectId, request []interface{}) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintNodes, blueprint),
		apiInput: request,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
