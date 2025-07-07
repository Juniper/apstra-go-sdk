package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"net/http"
	"sort"
	"time"
)

var _ Template = &TemplatePodBased{}

type TemplatePodBased struct {
	Id             ObjectId
	templateType   TemplateType
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *TemplatePodBasedData
}

func (o *TemplatePodBased) Type() TemplateType {
	return o.templateType
}

func (o *TemplatePodBased) ID() ObjectId {
	return o.Id
}

func (o *TemplatePodBased) OverlayControlProtocol() OverlayControlProtocol {
	if o == nil || o.Data == nil || len(o.Data.PodInfo) == 0 {
		return OverlayControlProtocolNone
	}

	// return the first record
	for _, v := range o.Data.PodInfo {
		if v.TemplateRackBasedData != nil {
			return v.TemplateRackBasedData.VirtualNetworkPolicy.OverlayControlProtocol
		}
	}

	return OverlayControlProtocolNone
}

type rawTemplatePodBased struct {
	Id                      ObjectId                 `json:"id"`
	Type                    templateType             `json:"type"`
	DisplayName             string                   `json:"display_name"`
	AntiAffinityPolicy      *rawAntiAffinityPolicy   `json:"anti_affinity_policy,omitempty"`
	Superspine              rawSuperspine            `json:"superspine"`
	CreatedAt               time.Time                `json:"created_at"`
	LastModifiedAt          time.Time                `json:"last_modified_at"`
	Capability              templateCapability       `json:"capability,omitempty"`
	RackBasedTemplates      []rawTemplateRackBased   `json:"rack_based_templates"`
	RackBasedTemplateCounts []RackBasedTemplateCount `json:"rack_based_template_counts"`
}

func (o rawTemplatePodBased) polish() (*TemplatePodBased, error) {
	superspine, err := o.Superspine.polish()
	if err != nil {
		return nil, err
	}
	capability, err := o.Capability.parse()
	if err != nil {
		return nil, err
	}
	var _, rbt []TemplateRackBased
	for _, rrbt := range o.RackBasedTemplates {
		if rrbt.Type == templateTypeNone {
			// because sometimes Apstra doesn't fill this in, but we know based on context
			rrbt.Type = templateTypeRackBased
		}
		polished, err := rrbt.polish()
		if err != nil {
			return nil, err
		}
		rbt = append(rbt, *polished)
	}
	tType, err := o.Type.parse()
	if err != nil {
		return nil, err
	}
	var aap *AntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		aap, err = o.AntiAffinityPolicy.polish()
		if err != nil {
			return nil, err
		}
	}

	if len(o.RackBasedTemplates) != len(o.RackBasedTemplateCounts) {
		return nil, fmt.Errorf("template '%s' has %d rack_based_templates and %d rack_based_template_counts - "+
			"these should match", o.Id, len(o.RackBasedTemplates), len(o.RackBasedTemplateCounts))
	}

	podTypeInfos := make(map[ObjectId]TemplatePodBasedInfo, len(o.RackBasedTemplates))
OUTER:
	for _, rrbt := range o.RackBasedTemplates {
		prbt, err := rrbt.polish()
		if err != nil {
			return nil, err
		}
		for _, rbtc := range o.RackBasedTemplateCounts { // loop over rack based template counts looking for matching ID
			if prbt.Id == rbtc.RackBasedTemplateId {
				podTypeInfos[rbtc.RackBasedTemplateId] = TemplatePodBasedInfo{
					Count:                 rbtc.Count,
					TemplateRackBasedData: prbt.Data,
				}
				continue OUTER
			}
		}
		return nil, fmt.Errorf("template contains rack_based_template '%s' which does not appear among rack_based_tempalte_counts", rrbt.Id)
	}

	return &TemplatePodBased{
		Id:             o.Id,
		templateType:   TemplateType(tType),
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &TemplatePodBasedData{
			DisplayName:        o.DisplayName,
			AntiAffinityPolicy: aap,
			Superspine:         *superspine,
			Capability:         TemplateCapability(capability),
			PodInfo:            podTypeInfos,
		},
	}, nil
}

type TemplatePodBasedData struct {
	DisplayName        string
	AntiAffinityPolicy *AntiAffinityPolicy
	Superspine         Superspine
	Capability         TemplateCapability
	PodInfo            map[ObjectId]TemplatePodBasedInfo
}

type RackBasedTemplateCount struct {
	RackBasedTemplateId ObjectId `json:"rack_based_template_id"`
	Count               int      `json:"count"`
}

type Superspine struct {
	PlaneCount         int
	Tags               []DesignTagData
	SuperspinePerPlane int
	LogicalDevice      LogicalDeviceData
}

type TemplateElementSuperspineRequest struct {
	PlaneCount         int
	Tags               []ObjectId
	SuperspinePerPlane int
	LogicalDeviceId    ObjectId
}

func (o *TemplateElementSuperspineRequest) raw(ctx context.Context, client *Client) (*rawSuperspine, error) {
	tags := make([]DesignTagData, len(o.Tags))
	for i, tagId := range o.Tags {
		rawTag, err := client.getTag(ctx, tagId)
		if err != nil {
			return nil, err
		}
		tags[i] = *rawTag.polish().Data
	}

	logicalDevice, err := client.getLogicalDevice(ctx, o.LogicalDeviceId)
	if err != nil {
		return nil, err
	}

	return &rawSuperspine{
		PlaneCount:         o.PlaneCount,
		Tags:               tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice:      *logicalDevice,
	}, nil
}

type rawSuperspine struct {
	PlaneCount         int              `json:"plane_count"`
	Tags               []DesignTagData  `json:"tags"`
	SuperspinePerPlane int              `json:"superspine_per_plane"`
	LogicalDevice      rawLogicalDevice `json:"logical_device"`
}

func (o rawSuperspine) polish() (*Superspine, error) {
	ld, err := o.LogicalDevice.polish()
	if err != nil {
		return nil, err
	}
	return &Superspine{
		PlaneCount:         o.PlaneCount,
		Tags:               o.Tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice: LogicalDeviceData{
			DisplayName: ld.Data.DisplayName,
			Panels:      ld.Data.Panels,
		},
	}, nil
}

type TemplatePodBasedInfo struct {
	Count                 int
	TemplateRackBasedData *TemplateRackBasedData
}

func (o *Client) getPodBasedTemplate(ctx context.Context, id ObjectId) (*rawTemplatePodBased, error) {
	rawTemplate, err := o.getTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	tType, err := rawTemplate.templateType()
	if err != nil {
		return nil, err
	}

	if tType != templateTypePodBased {
		return nil, ClientErr{
			errType: ErrWrongType,
			err:     fmt.Errorf("template '%s' is of type '%s', not '%s'", id, tType, templateTypePodBased),
		}
	}

	result := &rawTemplatePodBased{}
	err = json.Unmarshal(rawTemplate, result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling raw pod-based template - %w", err)
	}

	// force 'type' field of included rack-based templates to "rack_based" b/c Apstra rejects empty string.
	for i, rbt := range result.RackBasedTemplates {
		switch rbt.Type {
		case "":
			result.RackBasedTemplates[i].Type = templateTypeRackBased
		case templateTypeRackBased: // fallthrough
		default:
			return nil, fmt.Errorf("rack-based template '%s' within pod-based template '%s' claims to be type '%s', expected '%s'",
				rbt.DisplayName, result.Id, rbt.Type, templateTypeRackBased)
		}
	}
	return result, nil
}

func (o *Client) getAllPodBasedTemplates(ctx context.Context) ([]rawTemplatePodBased, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var result []rawTemplatePodBased
	for _, t := range templates {
		tType, err := t.templateType()
		if err != nil {
			return nil, err
		}
		if tType != templateTypePodBased {
			continue
		}
		var raw rawTemplatePodBased
		err = json.Unmarshal(t, &raw)
		if err != nil {
			return nil, err
		}
		result = append(result, raw)
	}

	return result, nil
}

type CreatePodBasedTemplateRequest struct {
	DisplayName        string
	Superspine         *TemplateElementSuperspineRequest
	PodInfos           map[ObjectId]TemplatePodBasedInfo
	AntiAffinityPolicy *AntiAffinityPolicy
}

func (o *CreatePodBasedTemplateRequest) raw(ctx context.Context, client *Client) (*rawCreatePodBasedTemplateRequest, error) {
	templatesRackBased := make([]rawTemplateRackBased, len(o.PodInfos))
	rackBasedTemplatesCounts := make([]RackBasedTemplateCount, len(o.PodInfos))
	var i int
	for k, pi := range o.PodInfos {
		if pi.TemplateRackBasedData != nil {
			return nil, fmt.Errorf("the TemplateRackBasedData (pod info) field must be nil when creating a pod-based template")
		}

		// grab the rack-based template (pod) from the API using the caller's map key (ObjectId) and stash it in templatesRackBased
		rbt, err := client.getRackBasedTemplate(ctx, k)
		if err != nil {
			return nil, err
		}
		templatesRackBased[i] = *rbt

		// prep the RackBasedTemplateCount object using the caller's map key (ObjectId) as
		// the link between the rawTemplateRackBased data copy and the RackBasedTemplateCount
		rackBasedTemplatesCounts[i].RackBasedTemplateId = k
		rackBasedTemplatesCounts[i].Count = pi.Count
		i++
	}

	sort.Slice(templatesRackBased, func(i, j int) bool {
		return templatesRackBased[i].DisplayName < templatesRackBased[j].DisplayName
	})

	switch {
	case o.Superspine == nil:
		return nil, errors.New("super spine cannot be <nil> when creating a pod-based template")
	case o.AntiAffinityPolicy == nil && compatibility.TemplateRequestRequiresAntiAffinityPolicy.Check(client.apiVersion):
		return nil, fmt.Errorf("anti-affinity policy cannot be <nil> when creating a pod-based template with Apstra %s", compatibility.TemplateRequestRequiresAntiAffinityPolicy)
	}

	var err error
	var superspine *rawSuperspine
	if o.Superspine != nil {
		superspine, err = o.Superspine.raw(ctx, client)
		if err != nil {
			return nil, err
		}
	}

	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	return &rawCreatePodBasedTemplateRequest{
		Type:                    templateTypePodBased,
		DisplayName:             o.DisplayName,
		Superspine:              *superspine,
		RackBasedTemplates:      templatesRackBased,
		RackBasedTemplateCounts: rackBasedTemplatesCounts,
		AntiAffinityPolicy:      antiAffinityPolicy,
	}, nil
}

type rawCreatePodBasedTemplateRequest struct {
	Type                    templateType             `json:"type"`
	DisplayName             string                   `json:"display_name"`
	Superspine              rawSuperspine            `json:"superspine"`
	RackBasedTemplates      []rawTemplateRackBased   `json:"rack_based_templates"`
	RackBasedTemplateCounts []RackBasedTemplateCount `json:"rack_based_template_counts"`
	AntiAffinityPolicy      *rawAntiAffinityPolicy   `json:"anti_affinity_policy,omitempty"`
}

func (o *Client) createPodBasedTemplate(ctx context.Context, in *rawCreatePodBasedTemplateRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignTemplates,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", err
	}

	return response.Id, nil
}

func (o *Client) updatePodBasedTemplate(ctx context.Context, id ObjectId, in *CreatePodBasedTemplateRequest) error {
	apiInput, err := in.raw(ctx, o)
	if err != nil {
		return err
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiInput: apiInput,
	})
	if err != nil {
		return err
	}

	return nil
}
