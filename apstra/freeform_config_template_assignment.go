package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlConfigTemplateAssignments = apiUrlBlueprintById + apiUrlPathDelim + "config-templates-assignments"
)

func (o *FreeformClient) ListConfigTemplateAssignments(ctx context.Context) (map[ObjectId]ObjectId, error) {
	var response struct {
		Assignments map[ObjectId]ObjectId `json:"assignments"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplateAssignments, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Assignments, nil
}

func (o *FreeformClient) GetConfigTemplateAssignments(ctx context.Context, in ObjectId) ([]ObjectId, error) {
	assignments, err := o.ListConfigTemplateAssignments(ctx)
	if err != nil {
		return nil, err
	}

	var result []ObjectId
	for k, v := range assignments {
		if v == in {
			result = append(result, k)
		}
	}
	return result, nil
}

func (o *FreeformClient) UpdateConfigTemplateAssignments(ctx context.Context, ctId ObjectId, sysIds []ObjectId) error {
	var apiInput struct {
		Assignments map[ObjectId]ObjectId `json:"assignments"`
	}

	apiInput.Assignments = make(map[ObjectId]ObjectId, len(sysIds))
	for _, sysId := range sysIds {
		apiInput.Assignments[sysId] = ctId
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlConfigTemplateAssignments, o.blueprintId),
		apiInput: apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
