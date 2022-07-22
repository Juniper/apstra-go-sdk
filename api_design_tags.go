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

type TagLabel string

type optionsTagsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getAllTagsResponse struct {
	Items []DesignTag `json:"items"`
}

type DesignTag struct {
	Id             ObjectId  `json:"id,omitempty"`
	Label          TagLabel  `json:"label"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

func (o *Client) listAllTags(ctx context.Context) ([]ObjectId, error) {
	response := &optionsTagsResponse{}
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

func (o *Client) getTag(ctx context.Context, id ObjectId) (*DesignTag, error) {
	response := &DesignTag{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTagById, id),
		apiResponse: response,
	})
}

func (o *Client) getTagByLabel(ctx context.Context, label TagLabel) (*DesignTag, error) {
	tags, err := o.getAllTags(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	labelNoCase := strings.ToLower(string(label))

	for _, t := range tags {
		if strings.ToLower(string(t.Label)) == labelNoCase {
			return &t, nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("tag with label '%s' not found", label),
	}
}

func (o *Client) getAllTags(ctx context.Context) ([]DesignTag, error) {
	response := &getAllTagsResponse{}
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

func (o *Client) createTag(ctx context.Context, in *DesignTag) (ObjectId, error) {
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

func (o *Client) updateTag(ctx context.Context, id ObjectId, in *DesignTag) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPut,
		urlStr:      fmt.Sprintf(apiUrlDesignTagById, id),
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) deleteTag(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignTagById, id),
	})
}
