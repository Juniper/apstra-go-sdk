package apstra

import (
	"context"
	"fmt"
	"github.com/orsinium-labs/enum"
	"net/http"
)

const (
	apiUrlPolicies   = apiUrlBlueprintById + apiUrlPathDelim + "policies"
	apiUrlPolicyById = apiUrlPolicies + apiUrlPathDelim + "%s"
)

type PolicyApplicationPointType enum.Member[string]

var (
	PolicyApplicationPointTypeGroup          = PolicyApplicationPointType{Value: "group"}
	PolicyApplicationPointTypeInternal       = PolicyApplicationPointType{Value: "internal"}
	PolicyApplicationPointTypeExternal       = PolicyApplicationPointType{Value: "external"}
	PolicyApplicationPointTypeSecurityZone   = PolicyApplicationPointType{Value: "security_zone"}
	PolicyApplicationPointTypeVirtualNetwork = PolicyApplicationPointType{Value: "virtual_network"}
	PolicyApplicationPointTypes              = enum.New(
		PolicyApplicationPointTypeGroup,
		PolicyApplicationPointTypeInternal,
		PolicyApplicationPointTypeExternal,
		PolicyApplicationPointTypeSecurityZone,
		PolicyApplicationPointTypeVirtualNetwork,
	)
)

type Policy struct {
	Id   ObjectId `json:"object_id,omitempty"`
	Data *PolicyData
}

type PolicyData struct {
	Enabled             bool                        `json:"enabled"`
	Label               string                      `json:"label"`
	Description         string                      `json:"description"`
	SrcApplicationPoint *PolicyApplicationPointData `json:"src_application_point"`
	DstApplicationPoint *PolicyApplicationPointData `json:"dst_application_point"`
	Rules               []PolicyRule                `json:"rules"`
	Tags                []string                    `json:"tags"`
}

func (o PolicyData) request() *policyRequest {
	rules := make([]rawPolicyRule, len(o.Rules))
	for i, rule := range o.Rules {
		rules[i] = *rule.raw()
	}
	return &policyRequest{
		Enabled:             o.Enabled,
		Label:               o.Label,
		Description:         o.Description,
		SrcApplicationPoint: o.SrcApplicationPoint.ObjectId(),
		DstApplicationPoint: o.DstApplicationPoint.ObjectId(),
		Rules:               rules,
		Tags:                o.Tags,
	}
}

type policyRequest struct {
	Enabled             bool            `json:"enabled"`
	Label               string          `json:"label"`
	Description         string          `json:"description"`
	SrcApplicationPoint ObjectId        `json:"src_application_point,omitempty"`
	DstApplicationPoint ObjectId        `json:"dst_application_point,omitempty"`
	Rules               []rawPolicyRule `json:"rules"`
	Tags                []string        `json:"tags"`
}

type rawPolicy struct {
	Enabled             bool                        `json:"enabled"`
	Label               string                      `json:"label"`
	Description         string                      `json:"description"`
	SrcApplicationPoint *PolicyApplicationPointData `json:"src_application_point,omitempty"`
	DstApplicationPoint *PolicyApplicationPointData `json:"dst_application_point,omitempty"`
	Rules               []rawPolicyRule             `json:"rules"`
	Tags                []string                    `json:"tags"`
	Id                  ObjectId                    `json:"id"`
}

func (o rawPolicy) polish() (*Policy, error) {
	rules := make([]PolicyRule, len(o.Rules))
	for i, raw := range o.Rules {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		rules[i] = *polished
	}
	return &Policy{
		Id: o.Id,
		Data: &PolicyData{
			Enabled:             o.Enabled,
			Label:               o.Label,
			Description:         o.Description,
			SrcApplicationPoint: o.SrcApplicationPoint,
			DstApplicationPoint: o.DstApplicationPoint,
			Rules:               rules,
			Tags:                o.Tags,
		},
	}, nil
}

func (o rawPolicy) request() *policyRequest {
	return &policyRequest{
		Enabled:             o.Enabled,
		Label:               o.Label,
		Description:         o.Description,
		SrcApplicationPoint: o.SrcApplicationPoint.Id,
		DstApplicationPoint: o.DstApplicationPoint.Id,
		Rules:               o.Rules,
		Tags:                o.Tags,
	}
}

type PolicyApplicationPointData struct {
	Id    ObjectId `json:"id"`
	Label string   `json:"label"`
	Type  string   `json:"type"` // group, internal, external, security_zone, virtual_network
}

func (o PolicyApplicationPointData) ObjectId() ObjectId {
	return o.Id
}

func (o *TwoStageL3ClosClient) getAllPolicies(ctx context.Context) ([]rawPolicy, error) {
	response := &struct {
		Policies []rawPolicy `json:"policies"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlPolicies, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Policies, nil
}

func (o *TwoStageL3ClosClient) getPolicy(ctx context.Context, id ObjectId) (*rawPolicy, error) {
	response := &rawPolicy{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) createPolicy(ctx context.Context, data *policyRequest) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlPolicies, o.blueprintId),
		apiInput:    data,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageL3ClosClient) deletePolicy(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) updatePolicy(ctx context.Context, id ObjectId, data *policyRequest) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
		apiInput: data,
	})
	return convertTtaeToAceWherePossible(err)
}
