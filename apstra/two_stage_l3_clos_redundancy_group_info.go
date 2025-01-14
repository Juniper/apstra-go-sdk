// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

type RedundancyGroupInfo struct {
	Id         ObjectId
	Type       enum.RedundancyGroupType
	SystemType enum.SystemType
	SystemRole enum.NodeRole
	SystemIds  [2]ObjectId
}

func (o *TwoStageL3ClosClient) GetRedundancyGroupInfo(ctx context.Context, id ObjectId) (*RedundancyGroupInfo, error) {
	resultMap, err := o.getRedundancyGroupInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	result, ok := resultMap[id]
	if !ok {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("redundancy group %q not found", id),
		}
	}

	return &result, nil
}

func (o *TwoStageL3ClosClient) GetAllRedundancyGroupInfo(ctx context.Context) (map[ObjectId]RedundancyGroupInfo, error) {
	return o.getRedundancyGroupInfo(ctx, "")
}

func (o *TwoStageL3ClosClient) getRedundancyGroupInfo(ctx context.Context, id ObjectId) (map[ObjectId]RedundancyGroupInfo, error) {
	rgNodeAttrs := []QEEAttribute{
		NodeTypeRedundancyGroup.QEEAttribute(),
		{Key: "name", Value: QEStringVal("n_redundancy_group")},
	}
	if id != "" {
		rgNodeAttrs = append(rgNodeAttrs, QEEAttribute{Key: "id", Value: QEStringVal(id)})
	}

	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node(rgNodeAttrs).
		Out([]QEEAttribute{RelationshipTypeComposedOfSystems.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeSystem.QEEAttribute(), {Key: "name", Value: QEStringVal("n_system")}})

	var queryResult struct {
		Items []struct {
			RedundancyGroup struct {
				Id   ObjectId                 `json:"id"`
				Type enum.RedundancyGroupType `json:"rg_type"`
			} `json:"n_redundancy_group"`
			System struct {
				Id   ObjectId        `json:"id"`
				Role enum.NodeRole   `json:"role"`
				Type enum.SystemType `json:"system_type"`
			} `json:"n_system"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResult)
	if err != nil {
		return nil, fmt.Errorf("graph query %q failed - %w", query, err)
	}

	result := make(map[ObjectId]RedundancyGroupInfo, len(queryResult.Items)/2)
	for _, item := range queryResult.Items {
		rgInfo, ok := result[item.RedundancyGroup.Id]
		if !ok {
			// create the map entry
			result[item.RedundancyGroup.Id] = RedundancyGroupInfo{
				Id:         item.RedundancyGroup.Id,
				Type:       item.RedundancyGroup.Type,
				SystemType: item.System.Type,
				SystemRole: item.System.Role,
				SystemIds:  [2]ObjectId{item.System.Id, ""},
			}
			continue
		}

		// validate the existing map entry
		if rgInfo.Type != item.RedundancyGroup.Type {
			return nil, fmt.Errorf("graph query %q returned inconsistent redundancy group types for group %q", query, item.RedundancyGroup.Id)
		}
		if rgInfo.SystemType != item.System.Type {
			return nil, fmt.Errorf("graph query %q returned inconsistent system types for group %q", query, item.RedundancyGroup.Id)
		}
		if rgInfo.SystemRole != item.System.Role {
			return nil, fmt.Errorf("graph query %q returned inconsistent system roles for group %q", query, item.RedundancyGroup.Id)
		}
		if rgInfo.SystemIds[1] != "" {
			return nil, fmt.Errorf("graph query %q returned too many system nodes for redundancy group %q", query, item.RedundancyGroup.Id)
		}

		// add the second system ID to the existing map entry
		rgInfo.SystemIds[1] = item.System.Id
		result[item.RedundancyGroup.Id] = rgInfo
	}

	// ensure that each redundancy group has both system IDs
	for k, v := range result {
		if v.SystemIds[0] == "" || v.SystemIds[1] == "" {
			return nil, fmt.Errorf("graph query %q didn't find system pairs for redundancy group %q, got: %q", query, k, v)
		}
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) GetRedundancyGroupInfoBySystemId(ctx context.Context, id ObjectId) (*RedundancyGroupInfo, error) {
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{NodeTypeSystem.QEEAttribute(), {Key: "id", Value: QEStringVal(id)}}).
		Out([]QEEAttribute{RelationshipTypePartOfRedundancyGroup.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeRedundancyGroup.QEEAttribute(), {Key: "name", Value: QEStringVal("n_redundancy_group")}}).
		Out([]QEEAttribute{RelationshipTypeComposedOfSystems.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeSystem.QEEAttribute(), {Key: "name", Value: QEStringVal("n_system")}})

	var queryResult struct {
		Items []struct {
			RedundancyGroup struct {
				Id   ObjectId                 `json:"id"`
				Type enum.RedundancyGroupType `json:"rg_type"`
			} `json:"n_redundancy_group"`
			System struct {
				Id   ObjectId        `json:"id"`
				Role enum.NodeRole   `json:"role"`
				Type enum.SystemType `json:"system_type"`
			} `json:"n_system"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResult)
	if err != nil {
		return nil, fmt.Errorf("graph query %q failed - %w", query, err)
	}

	switch len(queryResult.Items) {
	case 0:
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("redundancy group associated with system %q not found", id),
		}
	case 2:
	default:
		return nil, fmt.Errorf("graph query %q returned an unexpected number of results. Expected 0 or 2, got %d", query, len(queryResult.Items))
	}

	var result RedundancyGroupInfo
	for i, item := range queryResult.Items {
		if i == 0 {
			result.Id = item.RedundancyGroup.Id
			result.Type = item.RedundancyGroup.Type
			result.SystemType = item.System.Type
			result.SystemRole = item.System.Role
			result.SystemIds[i] = item.System.Id
		} else {
			if result.Id != item.RedundancyGroup.Id {
				return nil, fmt.Errorf("graph query %q returned inconsistent redundancy group IDs for system %q: %q and %q", query, id, result.Id, item.RedundancyGroup.Id)
			}
			if result.Type != item.RedundancyGroup.Type {
				return nil, fmt.Errorf("graph query %q returned inconsistent redundancy group types for system %q: %q and %q", query, id, result.Type, item.RedundancyGroup.Type)
			}
			if result.SystemType != item.System.Type {
				return nil, fmt.Errorf("graph query %q returned inconsistent system types for system %q: %q and %q", query, id, result.SystemType, item.System.Type)
			}
			if result.SystemRole != item.System.Role {
				return nil, fmt.Errorf("graph query %q returned inconsistent system roles for system %q: %q and %q", query, id, result.SystemRole, item.System.Role)
			}
			result.SystemIds[1] = item.System.Id
		}
	}

	return &result, nil
}
