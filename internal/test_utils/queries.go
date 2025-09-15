// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testutils

import (
	"context"

	"github.com/Juniper/apstra-go-sdk/apstra"
)

func GetSystemIdsByRole(ctx context.Context, bp *apstra.TwoStageL3ClosClient, role string) ([]apstra.ObjectId, error) {
	leafQuery := new(apstra.PathQuery).
		SetClient(bp.Client()).
		SetBlueprintId(bp.Id()).
		SetBlueprintType(apstra.BlueprintTypeStaging).
		Node([]apstra.QEEAttribute{
			apstra.NodeTypeSystem.QEEAttribute(),
			{Key: "role", Value: apstra.QEStringVal(role)},
			{Key: "name", Value: apstra.QEStringVal("n_system")},
		})

	var leafQueryResult struct {
		Items []struct {
			System struct {
				Id apstra.ObjectId `json:"id"`
			} `json:"n_system"`
		} `json:"items"`
	}

	err := leafQuery.Do(ctx, &leafQueryResult)
	if err != nil {
		return nil, err
	}

	result := make([]apstra.ObjectId, len(leafQueryResult.Items))
	for i, item := range leafQueryResult.Items {
		result[i] = item.System.Id
	}

	return result, nil
}
