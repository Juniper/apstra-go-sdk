package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlConfigTemplateAssignments = apiUrlBlueprintById + apiUrlPathDelim + "config-templates-assignments"
)

// ListConfigTemplateAssignments returns map of [ObjectId]*ObjectId where keys are system IDs and values are Config Template IDs
func (o *FreeformClient) ListConfigTemplateAssignments(ctx context.Context) (map[ObjectId]*ObjectId, error) {
	var response struct {
		Assignments map[ObjectId]*ObjectId `json:"assignments"`
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
		if *v == in {
			result = append(result, k)
		}
	}
	return result, nil
}

// UpdateConfigTemplateAssignments returns map [ObjectId]*ObjectId where keys are system IDs and values are Config Template IDs.
// A nil value will clear the assignment of the Config Template to the System.
func (o *FreeformClient) UpdateConfigTemplateAssignments(ctx context.Context, assignments map[ObjectId]*ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(apiUrlConfigTemplateAssignments, o.blueprintId),
		apiInput: struct {
			Assignments map[ObjectId]*ObjectId `json:"assignments"`
		}{
			Assignments: assignments,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) UpdateConfigTemplateAssignmentsByTemplate(ctx context.Context, ctID ObjectId, sysIDs []ObjectId) error {
	// set a lock and defer the unlock until return
	o.client.lock("configTemplateAssignments")
	defer o.client.unlock("configTemplateAssignments")

	// read the current set of CT assignments from the api and load it into current.
	current, err := o.GetConfigTemplateAssignments(ctx, ctID)
	if err != nil {
		return err
	}

	// turn current into a map [ObjectID]struct{}
	currentMap := make(map[ObjectId]bool)
	for _, sysID := range current {
		currentMap[sysID] = true
	}

	request := make(map[ObjectId]*ObjectId)

	// desired map
	desiredMap := make(map[ObjectId]bool)
	for _, sysID := range sysIDs {
		desiredMap[sysID] = true

		if !currentMap[sysID] {
			// lookup failed, add
			request[sysID] = &ctID
		}
	}

	// fill the request map with the ids to be deleted
	for _, sysID := range current {
		if !desiredMap[sysID] {
			request[sysID] = nil
		}
	}

	if len(request) == 0 {
		return nil
	}

	// Update the assignments with the newly formed request.
	err = o.UpdateConfigTemplateAssignments(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
