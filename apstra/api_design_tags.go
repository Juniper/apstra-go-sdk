package apstra

import (
	"context"
	"fmt"
	"net/http"
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

	for _, t := range tags {
		if t.Label == label {
			return &t, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("tag with label '%s' not found", label),
	}
}

func (o *Client) getTagsByLabels(ctx context.Context, labels []string) ([]rawDesignTag, error) {
	requestedLabels := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		requestedLabels[label] = struct{}{}
	}

	if len(requestedLabels) != len(labels) {
		return nil, fmt.Errorf("slice of requested labels contains duplicate entries")
	}

	tags, err := o.getAllTags(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	allTagsMap := make(map[string]rawDesignTag, len(tags))
	for _, tag := range tags {
		allTagsMap[tag.Label] = tag
	}

	var result []rawDesignTag
	for requestedLabel := range requestedLabels {
		if tag, ok := allTagsMap[requestedLabel]; ok {
			result = append(result, tag)
		} else {
			return nil, ClientErr{
				errType: ErrNotfound,
				err:     fmt.Errorf("tag with label '%s' not found", requestedLabel),
			}
		}
	}
	return result, nil
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
