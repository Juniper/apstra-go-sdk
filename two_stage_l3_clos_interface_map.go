package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintInterfaceMapAssignment = apiUrlBlueprintById + apiUrlPathDelim + "interface-map-assignments"
)

type SystemIdToInterfaceMapAssignment map[string]string

type interfaceMapAssignment struct {
	Assignments SystemIdToInterfaceMapAssignment `json:"assignments"`
}

func (o *TwoStageLThreeClosClient) getInterfaceMapAssignments(ctx context.Context) (SystemIdToInterfaceMapAssignment, error) {
	response := &interfaceMapAssignment{}
	return response.Assignments, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintInterfaceMapAssignment, o.blueprintId),
		apiResponse: response,
	})
}

func (o *TwoStageLThreeClosClient) setInterfaceMapAssignments(ctx context.Context, assignments SystemIdToInterfaceMapAssignment) (int, error) {
	response := &struct {
		ConfigBlueprintVersion int `json:"config_blueprint_version"`
	}{}
	return response.ConfigBlueprintVersion, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPatch,
		urlStr:      fmt.Sprintf(apiUrlBlueprintInterfaceMapAssignment, o.blueprintId),
		apiInput:    &interfaceMapAssignment{Assignments: assignments},
		apiResponse: response,
	})
}
