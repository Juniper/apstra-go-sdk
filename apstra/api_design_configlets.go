package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignConfiglets       = apiUrlDesignPrefix + "configlets"
	apiUrlDesignConfigletsPrefix = apiUrlDesignConfiglets + apiUrlPathDelim
	apiUrlDesignConfigletsById   = apiUrlDesignConfigletsPrefix + "%s"
)

type ConfigletGenerator struct {
	ConfigStyle          PlatformOS
	Section              ConfigletSection
	SectionCondition     *string
	TemplateText         string
	NegationTemplateText string
	Filename             string
}

type rawConfigletGenerator struct {
	ConfigStyle          string `json:"config_style"`
	Section              string `json:"section"`
	SectionCondition     string `json:"section_condition,omitempty"`
	TemplateText         string `json:"template_text"`
	NegationTemplateText string `json:"negation_template_text"`
	Filename             string `json:"filename"`
}

type Configlet struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *ConfigletData
}

type ConfigletData struct {
	RefArchs    []RefDesign
	Generators  []ConfigletGenerator
	DisplayName string
}

type rawConfigletData struct {
	RefArchs    []string                `json:"ref_archs"`
	Generators  []rawConfigletGenerator `json:"generators"`
	DisplayName string                  `json:"display_name"`
}

type rawConfiglet struct {
	RefArchs       []string                `json:"ref_archs"`
	Generators     []rawConfigletGenerator `json:"generators"`
	CreatedAt      time.Time               `json:"created_at"`
	Id             ObjectId                `json:"id,omitempty"`
	LastModifiedAt time.Time               `json:"last_modified_at"`
	DisplayName    string                  `json:"display_name"`
}

type ConfigletRequest ConfigletData
type rawConfigletRequest rawConfigletData

func (o *ConfigletRequest) raw() *rawConfigletRequest {
	rawcr := rawConfigletRequest{}
	rawcr.DisplayName = o.DisplayName
	rawcr.RefArchs = make([]string, len(o.RefArchs))
	rawcr.Generators = make([]rawConfigletGenerator, len(o.Generators))
	for i, j := range o.RefArchs {
		rawcr.RefArchs[i] = j.String()
	}
	for i, j := range o.Generators {
		rawcr.Generators[i] = *j.raw()
	}

	return &rawcr
}

func (o *rawConfigletGenerator) polish() (*ConfigletGenerator, error) {
	platform, err := platformOS(o.ConfigStyle).parse()
	if err != nil {
		return nil, err
	}

	section, err := configletSection(o.Section).parse()
	if err != nil {
		return nil, err
	}

	return &ConfigletGenerator{
		ConfigStyle:          PlatformOS(platform),
		Section:              ConfigletSection(section),
		SectionCondition:     &o.SectionCondition,
		TemplateText:         o.TemplateText,
		NegationTemplateText: o.NegationTemplateText,
		Filename:             o.Filename,
	}, nil
}

func (o *ConfigletGenerator) raw() *rawConfigletGenerator {
	cg := rawConfigletGenerator{}
	cg.TemplateText = o.TemplateText
	cg.Filename = o.Filename
	cg.NegationTemplateText = o.NegationTemplateText
	cg.ConfigStyle = o.ConfigStyle.raw().string()
	cg.Section = string(o.Section.raw())
	if o.SectionCondition != nil {
		cg.SectionCondition = *o.SectionCondition
	}
	return &cg
}

func (o *rawConfiglet) polish() (*Configlet, error) {
	var err error

	refArchs := make([]RefDesign, len(o.RefArchs))
	for i, refArch := range o.RefArchs {
		refArchs[i], err = refDesign(refArch).parse()
		if err != nil {
			return nil, err
		}
	}

	generators := make([]ConfigletGenerator, len(o.Generators))
	for i, generator := range o.Generators {
		polished, err := generator.polish()
		if err != nil {
			return nil, err
		}
		generators[i] = *polished
	}

	return &Configlet{
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &ConfigletData{
			RefArchs:    refArchs,
			Generators:  generators,
			DisplayName: o.DisplayName,
		},
	}, nil
}

func (o *Client) listAllConfiglets(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignConfiglets,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getConfiglet(ctx context.Context, id ObjectId) (*rawConfiglet, error) {
	response := &rawConfiglet{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignConfigletsById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getConfigletByName(ctx context.Context, name string) (*rawConfiglet, error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	for _, t := range cgs {
		if t.DisplayName == name {
			return &t, nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf(" Configlet with name '%s' not found", name),
	}
}

func (o *Client) getAllConfiglets(ctx context.Context) ([]rawConfiglet, error) {
	response := &struct {
		Items []rawConfiglet `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignConfiglets,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) createConfiglet(ctx context.Context, in *ConfigletRequest) (ObjectId, error) {

	cr := in.raw()
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignConfiglets,
		apiInput:    cr,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateConfiglet(ctx context.Context, id ObjectId, in *ConfigletRequest) error {
	cr := in.raw()

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignConfigletsById, id),
		apiInput: cr,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteConfiglet(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignConfigletsById, id),
	})
}
