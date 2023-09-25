git package apstra

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

type rawIbaDashboard struct {
	Description         string      `json:"description"`
	Default             bool        `json:"default"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	Label               string      `json:"label"`
	Grid                [][]string  `json:"grid"`
	PredefinedDashboard string      `json:"predefined_dashboard"`
	Id                  string      `json:"id"`
	UpdatedBy           interface{} `json:"updated_by"`
}

type IBADashboard struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *IBADashboardData
}

type IBADashboardData struct {

}

type rawConfigletData struct {
	RefArchs    []refDesign             `json:"ref_archs"`
	Generators  []rawConfigletGenerator `json:"generators"`
	DisplayName string                  `json:"display_name"`
}

type rawConfiglet struct {
	RefArchs       []refDesign             `json:"ref_archs"`
	Generators     []rawConfigletGenerator `json:"generators"`
	CreatedAt      time.Time               `json:"created_at"`
	Id             ObjectId                `json:"id,omitempty"`
	LastModifiedAt time.Time               `json:"last_modified_at"`
	DisplayName    string                  `json:"display_name"`
}

func (o *ConfigletData) raw() *rawConfigletData {
	refArchs := make([]refDesign, len(o.RefArchs))
	for i, j := range o.RefArchs {
		refArchs[i] = refDesign(j.String())
	}

	generators := make([]rawConfigletGenerator, len(o.Generators))
	for i, j := range o.Generators {
		generators[i] = *j.raw()
	}

	return &rawConfigletData{
		DisplayName: o.DisplayName,
		RefArchs:    refArchs,
		Generators:  generators,
	}
}

func (o *rawConfigletData) polish() (*ConfigletData, error) {
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
	return &ConfigletData{
		RefArchs:    refArchs,
		Generators:  generators,
		DisplayName: o.DisplayName,
	}, nil
}

func (o *rawConfigletGenerator) polish() (*ConfigletGenerator, error) {
	platform, err := o.ConfigStyle.parse()
	if err != nil {
		return nil, err
	}
	section, err := o.Section.parse()
	if err != nil {
		return nil, err
	}
	return &ConfigletGenerator{
		ConfigStyle:          PlatformOS(platform),
		Section:              ConfigletSection(section),
		TemplateText:         o.TemplateText,
		NegationTemplateText: o.NegationTemplateText,
		Filename:             o.Filename,
	}, nil
}

func (o *ConfigletGenerator) raw() *rawConfigletGenerator {
	return &rawConfigletGenerator{
		TemplateText:         o.TemplateText,
		Filename:             o.Filename,
		NegationTemplateText: o.NegationTemplateText,
		ConfigStyle:          o.ConfigStyle.raw(),
		Section:              o.Section.raw(),
	}
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
	configlets, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	foundIdx := -1
	for i, configlet := range configlets {
		if configlet.DisplayName == name {
			if foundIdx >= 0 {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple Configlets have name %q", name),
				}
			}
			foundIdx = i
		}
	}

	if foundIdx >= 0 {
		return &configlets[foundIdx], nil
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Configlet with name '%s' found", name),
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

func (o *Client) createConfiglet(ctx context.Context, in *rawConfigletData) (ObjectId, error) {
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignConfiglets,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateConfiglet(ctx context.Context, id ObjectId, in *rawConfigletData) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignConfigletsById, id),
		apiInput: in,
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
