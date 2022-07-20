package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintInterfaceMapAssignment = apiUrlBlueprintById + apiUrlPathDelim + "interface-map-assignments"
)

// SystemIdToInterfaceMapAssignment maps graph db 'system' nodes (their id is
// the string value) to graph db 'interface_map' nodes. interface{} is used for
// the interface_map nodes because apstra expects 'null' in the JSON fields
// where no map is assigned.
type SystemIdToInterfaceMapAssignment map[string]interface{}

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

func (o *TwoStageLThreeClosClient) setInterfaceMapAssignments(ctx context.Context, assignments SystemIdToInterfaceMapAssignment) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintInterfaceMapAssignment, o.blueprintId),
		apiInput: &interfaceMapAssignment{Assignments: assignments},
	})
}