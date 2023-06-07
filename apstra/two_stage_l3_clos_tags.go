package apstra

import "context"

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
