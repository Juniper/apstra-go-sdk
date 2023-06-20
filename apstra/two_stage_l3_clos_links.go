package apstra

import (
	"context"
	"fmt"
)

// SystemNodesFromLinkIds performs a graph query to determine the 'system' nodes
// which are connected to either end of the specified links. When NodeRole has
// any value other than SystemNodeRoleNone, that role is used as an additional
// filter to select nodes of a specified type. For example, given a slice of
// spine/leaf link IDs, setting SystemNodeRoleLeaf will cause the returned slice
// to contain only leaf switch IDs.
func (o *TwoStageL3ClosClient) SystemNodesFromLinkIds(ctx context.Context, linkIds []ObjectId, nodeRole SystemNodeRole) ([]ObjectId, error) {
	_, result, err := o.systemNodesFromLinkIds(ctx, linkIds, nodeRole)
	return result, err
}

func (o *TwoStageL3ClosClient) systemNodesFromLinkIds(ctx context.Context, linkIds []ObjectId, nodeRole SystemNodeRole) (QEQuery, []ObjectId, error) {
	systemQueryAttributes := []QEEAttribute{
		NodeTypeSystem.QEEAttribute(),
		{"name", QEStringVal("n_system")},
	}
	if nodeRole != SystemNodeRoleNone {
		systemQueryAttributes = append(systemQueryAttributes, nodeRole.QEEAttribute())
	}

	linkQuery := new(MatchQuery).
		SetBlueprintId(o.blueprintId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(o.Client())
	for _, linkId := range linkIds {
		linkQuery.Match(
			new(PathQuery).
				Node([]QEEAttribute{NodeTypeLink.QEEAttribute(),
					{"id", QEStringVal(linkId.String())},
				}).
				In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node(systemQueryAttributes),
		)
	}

	linkQueryResult := struct {
		Items []struct {
			System struct {
				Id ObjectId `json:"id"`
			} `json:"n_system"`
		} `json:"items"`
	}{}

	err := linkQuery.Do(ctx, &linkQueryResult)
	if err != nil {
		return nil, nil, err
	}

	response := make([]ObjectId, len(linkQueryResult.Items))
	for i, item := range linkQueryResult.Items {
		response[i] = item.System.Id
	}

	return linkQuery, response, nil
}

// SystemNodeFromLinkIds performs a graph query to determine the 'system' node
// which is connected to either end of the specified links. When NodeRole has
// any value other than SystemNodeRoleNone, that role is used as an additional
// filter to select nodes of a specified type. For example, given a slice of
// spine/leaf link IDs, setting SystemNodeRoleLeaf will cause the returned ID
// to represent the leaf system.
// If more than no systems or more than one system match the criteria an error
// is returned.
func (o *TwoStageL3ClosClient) SystemNodeFromLinkIds(ctx context.Context, linkIds []ObjectId, nodeRole SystemNodeRole) (ObjectId, error) {
	query, nodeIds, err := o.systemNodesFromLinkIds(ctx, linkIds, nodeRole)
	if err != nil {
		return "", err
	}

	switch len(nodeIds) {
	case 0:
		return "", ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no node matches query: %q", query.String()),
		}
	case 1:
		return nodeIds[0], nil
	default:
		return "", ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("multiple nodes match query: %q", query.String()),
		}
	}
}
