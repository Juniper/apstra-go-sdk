// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/policy"
)

var (
	_ Template          = (*TemplatePodBased)(nil)
	_ internal.IDSetter = (*TemplatePodBased)(nil)
	_ json.Marshaler    = (*TemplatePodBased)(nil)
	_ json.Unmarshaler  = (*TemplatePodBased)(nil)
)

type TemplatePodBased struct {
	Label              string
	AntiAffinityPolicy *policy.AntiAffinity // required for 4.2.0 only
	Superspine         Superspine
	Pods               []PodWithCount

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t TemplatePodBased) TemplateType() enum.TemplateType {
	return enum.TemplateTypePodBased
}

func (t TemplatePodBased) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t *TemplatePodBased) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *TemplatePodBased) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (t TemplatePodBased) MarshalJSON() ([]byte, error) {
	type rawRackBasedTemplateCount struct {
		Count               int    `json:"count"`
		RackBasedTemplateId string `json:"rack_based_template_id"`
	}

	raw := struct {
		Capability              enum.TemplateCapability     `json:"capability"`
		DisplayName             string                      `json:"display_name"`
		AntiAffinityPolicy      *policy.AntiAffinity        `json:"anti_affinity_policy,omitempty"`
		RackBasedTemplateCounts []rawRackBasedTemplateCount `json:"rack_based_template_counts"`
		RackBasedTemplates      []TemplateRackBased         `json:"rack_based_templates"`
		Superspine              Superspine                  `json:"superspine"`
		Type                    enum.TemplateType           `json:"type"`
	}{
		Capability:         enum.TemplateCapabilityBlueprint,
		DisplayName:        t.Label,
		AntiAffinityPolicy: t.AntiAffinityPolicy,
		Superspine:         t.Superspine,
		Type:               t.TemplateType(),
	}

	// used to generate IDs within the template
	hasher := md5.New()

	// set the superspine logical device ID if necessary
	if raw.Superspine.LogicalDevice.ID() == nil {
		raw.Superspine.LogicalDevice.mustSetHashID(hasher)
	}

	// keep track of rack type IDs (hashes of rack data). if two rack types are
	// identical twins (have the same contents) we don't want to add them to
	// raw.RackTypes twice. we will add them to raw.RackTypeCounts twice, and
	// the Apstra API will amend the totals as needed.
	podIDs := make(map[string]struct{}, len(t.Pods))

	// loop over pods, calculate a fresh ID, count the type of each
	for _, podWithCount := range t.Pods {
		pod := podWithCount.Pod.Replicate()  // fresh copy without metadata
		pod.mustSetHashID(hasher)            // assign the ID
		pod.skipTypeDuringMarshalJSON = true // don't marshal the nested template's type

		// add an entry to raw.RackTypeCounts without regard to twins
		raw.RackBasedTemplateCounts = append(raw.RackBasedTemplateCounts, rawRackBasedTemplateCount{Count: podWithCount.Count, RackBasedTemplateId: pod.id})

		// add an entry to raw.RackTypes only if it's not a twin
		if _, ok := podIDs[pod.id]; !ok {
			podIDs[pod.id] = struct{}{}
			raw.RackBasedTemplates = append(raw.RackBasedTemplates, pod)
		}
	}

	//// used to keep track of pod type quantity by ID in case we have identical racks
	//podIDToCount := make(map[string]int, len(t.Pods))
	//
	//// initialize raw.RackBasedTemplateCounts so we can append to it without shuffling memory around
	//raw.RackBasedTemplateCounts = make([]rawRackBasedTemplateCount, 0, len(t.Pods))
	//
	//// loop over pods, calculate a fresh ID, count the type of each
	//for _, pod := range t.Pods {
	//	rackBasedTemplate := pod.Pod.Replicate() // fresh copy without metadata
	//	rackBasedTemplate.mustSetHashID(hash)    // assign the ID
	//	if _, ok := podIDToCount[rackBasedTemplate.id]; !ok {
	//		raw.RackBasedTemplates = append(raw.RackBasedTemplates, rackBasedTemplate) // previously unseen pod type - append it to the slice
	//	}
	//	podIDToCount[rackBasedTemplate.id] += pod.Count // adjust the quantity for this pod type
	//}
	//
	//// prepare raw.RackTypeCounts from rackTypeIDToCount
	//raw.RackBasedTemplateCounts = make([]rawRackBasedTemplateCount, 0, len(podIDToCount))
	//for id, count := range podIDToCount {
	//	raw.RackBasedTemplateCounts = append(raw.RackBasedTemplateCounts, rawRackBasedTemplateCount{RackBasedTemplateId: id, Count: count})
	//}

	return json.Marshal(&raw)
}

func (t *TemplatePodBased) UnmarshalJSON(bytes []byte) error {
	type rawRackBasedTemplateCount struct {
		Count               int    `json:"count"`
		RackBasedTemplateId string `json:"rack_based_template_id"`
	}

	var raw struct {
		DisplayName             string                      `json:"display_name"`
		AntiAffinityPolicy      *policy.AntiAffinity        `json:"anti_affinity_policy"`
		RackBasedTemplateCounts []rawRackBasedTemplateCount `json:"rack_based_template_counts"`
		RackBasedTemplates      []TemplateRackBased         `json:"rack_based_templates"`
		Superspine              Superspine                  `json:"superspine"`
		Type                    enum.TemplateType           `json:"type"`

		ID             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshaling pod based template: %w", err)
	}

	t.Label = raw.DisplayName
	t.AntiAffinityPolicy = raw.AntiAffinityPolicy
	t.Superspine = raw.Superspine
	t.id = raw.ID
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	idToPodType := make(map[string]TemplateRackBased, len(raw.RackBasedTemplates))
	for _, podType := range raw.RackBasedTemplates {
		idToPodType[podType.id] = podType
	}

	t.Pods = make([]PodWithCount, len(raw.RackBasedTemplates))
	for i, rackBasedTemplateCount := range raw.RackBasedTemplateCounts {
		if rackBasedTemplate, ok := idToPodType[rackBasedTemplateCount.RackBasedTemplateId]; ok {
			t.Pods[i] = PodWithCount{Pod: rackBasedTemplate, Count: rackBasedTemplateCount.Count}
			continue
		}

		// we should not get here
		return sdk.ErrAPIResponseInvalid(fmt.Sprintf("payload specifies %d instances of pod type %q which does not exist", rackBasedTemplateCount.Count, rackBasedTemplateCount.RackBasedTemplateId))
	}

	return nil
}

func (t TemplatePodBased) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplatePodBased) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func NewPodBasedTemplate(id string) TemplatePodBased {
	return TemplatePodBased{id: id}
}
