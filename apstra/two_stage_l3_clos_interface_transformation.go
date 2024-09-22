// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintInterfaceTransformation = apiUrlBlueprintById + apiUrlPathDelim + "interface-transformation"
)

type rawSetTransformationRequest struct {
	Force      bool `json:"force"` // not clear what this is for
	Interfaces []struct {
		TransformationId int      `json:"transformation_id"`
		SystemId         ObjectId `json:"system_id"`
		IfName           string   `json:"if_name"`
	} `json:"interfaces"`
}

// SetTransformIdByIfName attempts to update the transform ID of the named
// interface on the specified system. Note that it is not always possible to
// change the transform number, particularly when such a change would change
// the link speed.
func (o *TwoStageL3ClosClient) SetTransformIdByIfName(ctx context.Context, systemId ObjectId, ifName string, transformId int) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut,
		urlStr: fmt.Sprintf(apiUrlBlueprintInterfaceTransformation, o.blueprintId),
		apiInput: &rawSetTransformationRequest{
			Force: false,
			Interfaces: []struct {
				TransformationId int      `json:"transformation_id"`
				SystemId         ObjectId `json:"system_id"`
				IfName           string   `json:"if_name"`
			}{
				{
					TransformationId: transformId,
					SystemId:         systemId,
					IfName:           ifName,
				},
			},
		},
	})
	return convertTtaeToAceWherePossible(err)
}

// GetTransformationId returns the current transform number of
// the specified interface node
func (o *TwoStageL3ClosClient) GetTransformationId(ctx context.Context, interfaceNodeId ObjectId) (int, error) {
	query := new(PathQuery).
		SetBlueprintType(BlueprintTypeStaging).
		SetBlueprintId(o.Id()).
		SetClient(o.Client()).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "id", Value: QEStringVal(interfaceNodeId)},
			{Key: "name", Value: QEStringVal("n_interface")},
		}).
		In([]QEEAttribute{
			RelationshipTypeHostedInterfaces.QEEAttribute(),
		}).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
		}).
		Out([]QEEAttribute{RelationshipTypeInterfaceMap.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterfaceMap.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_interface_map")},
		})

	var queryResponse struct {
		Items []struct {
			Interface struct {
				IfName string `json:"if_name"`
			} `json:"n_interface"`
			InterfaceMap struct {
				Id         string `json:"id"`
				Interfaces []struct {
					Mapping []int  `json:"mapping"`
					Name    string `json:"name"`
				} `json:"interfaces"`
			} `json:"n_interface_map"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResponse)
	if err != nil {
		return -1, fmt.Errorf("error executing query %q - %w", query.String(), err)
	}

	switch len(queryResponse.Items) {
	case 0:
		return -1, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no interface map associated with node %q query: %q", interfaceNodeId, query.String()),
		}
	case 1:
		return transformByIfName(queryResponse.Items[0].Interface.IfName, queryResponse.Items[0].InterfaceMap.Interfaces)
	default:
		return -1, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("query %q found %d interface maps, expected 1", query.String(), len(queryResponse.Items)),
		}
	}
}

// GetTransformationIdByIfName returns the current transform number of the
// named interface on the specified system node
func (o *TwoStageL3ClosClient) GetTransformationIdByIfName(ctx context.Context, systemNodeId ObjectId, ifName string) (int, error) {
	query := new(PathQuery).
		SetBlueprintType(BlueprintTypeStaging).
		SetBlueprintId(o.Id()).
		SetClient(o.Client()).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{Key: "id", Value: QEStringVal(systemNodeId.String())},
		}).
		Out([]QEEAttribute{RelationshipTypeInterfaceMap.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterfaceMap.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_interface_map")},
		})

	var queryResponse struct {
		Items []struct {
			InterfaceMap struct {
				Id         string `json:"id"`
				Interfaces []struct {
					Mapping []int  `json:"mapping"`
					Name    string `json:"name"`
				} `json:"interfaces"`
			} `json:"n_interface_map"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResponse)
	if err != nil {
		return -1, fmt.Errorf("failed querying for node %s interface map - %w", systemNodeId, err)
	}

	switch len(queryResponse.Items) {
	case 0:
		return -1, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no interface map associated with node %q query: %q", systemNodeId, query.String()),
		}
	case 1:
		return transformByIfName(ifName, queryResponse.Items[0].InterfaceMap.Interfaces)
	default:
		return -1, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("query %q found %d interface maps, expected 1", query.String(), len(queryResponse.Items)),
		}
	}
}

func transformByIfName(ifName string, in []struct {
	Mapping []int  `json:"mapping"`
	Name    string `json:"name"`
},
) (int, error) {
	for _, iMapInterface := range in {
		if iMapInterface.Name != ifName {
			continue
		}
		return iMapInterface.Mapping[1], nil
	}

	// getting here means we failed to find a the interface within the mapping
	return -1, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no mapping for interface %q found in interface map", ifName),
	}
}
