package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintConfiglets       = apiUrlBlueprintById + apiUrlPathDelim + "configlets"
	apiUrlBlueprintConfigletsPrefix = apiUrlBlueprintConfiglets + apiUrlPathDelim
	apiUrlBlueprintConfigletsById   = apiUrlBlueprintConfigletsPrefix + "%s"
)

type rawTwoStageL3ClosConfiglet struct {
	Data      rawConfigletData `json:"configlet"`
	Id        string           `json:"id"`
	Condition string           `json:"condition"`
	Label     string           `json:"label"`
}

type TwoStageL3ClosConfiglet struct {
	Data      ConfigletData
	Id        string
	Condition string
	Label     string
}

func (o *TwoStageL3ClosConfiglet) raw() *rawTwoStageL3ClosConfiglet {
	rawc := rawTwoStageL3ClosConfiglet{}
	rawc.Data = *o.Data.raw()
	rawc.Id = o.Id
	rawc.Condition = o.Condition
	rawc.Label = o.Label
	return &rawc
}

func (o *rawTwoStageL3ClosConfiglet) polish() (*TwoStageL3ClosConfiglet, error) {
	c := TwoStageL3ClosConfiglet{}
	c.Id = o.Id
	c.Condition = o.Condition
	c.Label = o.Label
	d, err := o.Data.polish()
	c.Data = *d
	return &c, err
}
func (o *TwoStageL3ClosClient) getAllConfiglets(ctx context.Context) ([]TwoStageL3ClosConfiglet, error) {
	response := &struct {
		Items []TwoStageL3ClosConfiglet `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *TwoStageL3ClosClient) getAllConfigletIds(ctx context.Context) ([]ObjectId, error) {
	configlets, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	ids := make([]ObjectId, len(configlets))
	for i, c := range configlets {
		ids[i] = ObjectId(c.Id)
	}
	return ids, nil
}

func (o *TwoStageL3ClosClient) getConfiglet(ctx context.Context, id ObjectId) (*TwoStageL3ClosConfiglet, error) {
	response := &rawTwoStageL3ClosConfiglet{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.polish()
}

func (o *TwoStageL3ClosClient) getConfigletByName(ctx context.Context, name string) (*TwoStageL3ClosConfiglet, error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	var c *TwoStageL3ClosConfiglet
	for _, t := range cgs {
		if t.Label == name {
			if c == nil {
				c = &t
			} else { // This is clearly the second occurrence
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify a configlet", name),
				}
			}
		}
	}
	if c != nil {
		return c, nil
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf(" Configlet with name '%s' not found", name),
	}
}

func (o *TwoStageL3ClosClient) importConfiglet(ctx context.Context, c ConfigletData, condition string, label string) (ObjectId, error) {
	response := &objectIdResponse{}

	in := TwoStageL3ClosConfiglet{
		Data:      c,
		Condition: condition,
		Label:     label,
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}
func (o *TwoStageL3ClosClient) importConfigletById(ctx context.Context, id ObjectId, condition string,
	label string) (ObjectId, error) {
	response := &objectIdResponse{}
	cfglet, err := o.client.GetConfiglet(ctx, id)
	if len(label) == 0 {
		label = cfglet.Data.DisplayName
	}
	in := TwoStageL3ClosConfiglet{
		Data:      *cfglet.Data,
		Condition: condition,
		Label:     label,
	}
	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateConfiglet(ctx context.Context, in *TwoStageL3ClosConfiglet) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), in.Id),
		apiInput: in.raw(),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *TwoStageL3ClosClient) deleteConfiglet(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
	})
}
