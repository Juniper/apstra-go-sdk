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
	apiUrlBlueprintTagging = apiUrlBlueprintById + apiUrlPathDelim + "tagging"
)

func (o *TwoStageL3ClosClient) GetNodeTags(ctx context.Context, nodeId ObjectId) ([]string, error) {
	query := new(PathQuery).
		SetBlueprintType(BlueprintTypeStaging).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{{"id", QEStringVal(nodeId.String())}}).
		In([]QEEAttribute{RelationshipTypeTag.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeTag.QEEAttribute(),
			{"name", QEStringVal("n_tag")},
		})

	var response struct {
		Items []struct {
			Tag struct {
				Label string `json:"label"`
			} `json:"n_tag"`
		} `json:"items"`
	}

	err := query.Do(ctx, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		// no tags found in the graph query. does the node even exist?
		var trash struct{}
		return nil, o.Client().GetNode(ctx, o.blueprintId, nodeId, &trash)
	}

	result := make([]string, len(response.Items))
	for i, item := range response.Items {
		result[i] = item.Tag.Label
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) SetNodeTags(ctx context.Context, nodeId ObjectId, tags []string) error {
	desiredTags := make(map[string]bool, len(tags))
	for _, tag := range tags {
		desiredTags[tag] = true
	}

	tags, err := o.GetNodeTags(ctx, nodeId)
	if err != nil {
		return err
	}
	currentTags := make(map[string]bool, len(tags))
	for _, tag := range tags {
		currentTags[tag] = true
	}

	var addTags, removeTags []string
	for k := range desiredTags {
		if currentTags[k] {
			delete(currentTags, k)
			delete(desiredTags, k)
		}
	}

	if len(currentTags) == 0 && len(desiredTags) == 0 {
		// nothing to add, nothing to remove - our job is done
		return nil
	}

	for k := range desiredTags {
		addTags = append(addTags, k)
	}

	for k := range currentTags {
		removeTags = append(removeTags, k)
	}

	apiInput := struct {
		Nodes  []string `json:"nodes"`
		Add    []string `json:"add"`
		Remove []string `json:"remove"`
	}{
		Nodes:  []string{nodeId.String()},
		Add:    addTags,
		Remove: removeTags,
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlBlueprintTagging, o.blueprintId),
		apiInput: &apiInput,
	})
	return convertTtaeToAceWherePossible(err)
}
