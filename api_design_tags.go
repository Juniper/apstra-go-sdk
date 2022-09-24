package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiUrlDesignTags       = apiUrlDesignPrefix + "tags"
	apiUrlDesignTagsPrefix = apiUrlDesignTags + apiUrlPathDelim
	apiUrlDesignTagById    = apiUrlDesignTagsPrefix + "%s"
)

type DesignTagRequest DesignTagData

type DesignTagData struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type DesignTag struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *DesignTagData
}

func (o *DesignTag) raw() *rawDesignTag {
	var label, description string
	if o.Data != nil {
		label = o.Data.Label
		description = o.Data.Description
	}
	return &rawDesignTag{
		Id:             o.Id,
		Label:          label,
		Description:    description,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
	}
}

type rawDesignTag struct {
	Id             ObjectId  `json:"id,omitempty"`
	Label          string    `json:"label"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

func (o *rawDesignTag) polish() *DesignTag {
	return &DesignTag{
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &DesignTagData{
			Label:       o.Label,
			Description: o.Description,
		},
	}
}

func (o *Client) listAllTags(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignTags,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getTag(ctx context.Context, id ObjectId) (*rawDesignTag, error) {
	response := &rawDesignTag{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTagById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getTagByLabel(ctx context.Context, label string) (*rawDesignTag, error) {
	tags, err := o.getAllTags(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	labelNoCase := strings.ToLower(label)

	for _, t := range tags {
		if strings.ToLower(t.Label) == labelNoCase {
			return &t, nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("tag with label '%s' not found", label),
	}
}

func (o *Client) getAllTags(ctx context.Context) ([]rawDesignTag, error) {
	response := &struct {
		Items []rawDesignTag `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTags,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) createTag(ctx context.Context, in *DesignTagRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignTags,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateTag(ctx context.Context, id ObjectId, in *DesignTagRequest) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignTagById, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteTag(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignTagById, id),
	})
}
