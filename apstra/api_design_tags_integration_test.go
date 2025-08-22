// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetAllTags(t *testing.T) {
	ctx := context.Background()
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			idList, err := client.Client.ListAllTags(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			tagList, err := client.Client.GetAllTags(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			if len(idList) != len(tagList) {
				t.Fatalf("got %d tag IDs but %d tags", len(idList), len(tagList))
			}

			for _, id := range idList {
				tag, err := client.Client.GetTag(context.TODO(), id)
				if err != nil {
					t.Fatal(err)
				}
				log.Println(tag)
			}
		})
	}
}

func TestCreateGetDeleteTag(t *testing.T) {
	ctx := context.Background()
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			label := testutils.RandString(10, "hex")
			description := testutils.RandString(10, "hex")
			id, err := client.Client.CreateTag(ctx, &apstra.DesignTagRequest{
				Label:       label,
				Description: description,
			})
			require.NoError(t, err)

			tag, err := client.Client.GetTag(context.TODO(), id)
			require.NoError(t, err)
			require.Equal(t, label, tag.Data.Label)
			require.Equal(t, description, tag.Data.Description)

			err = client.Client.DeleteTag(context.TODO(), id)
			require.NoError(t, err)
		})
	}
}

func TestCreateTagCollision(t *testing.T) {
	ctx := context.Background()
	clients := testclient.GetTestClients(t, ctx)
	label := testutils.RandString(10, "hex")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			id1, err := client.Client.CreateTag(ctx, &apstra.DesignTagRequest{Label: label})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Client.DeleteTag(ctx, id1)) })

			_, err = client.Client.CreateTag(ctx, &apstra.DesignTagRequest{Label: label})
			require.Error(t, err)

			var ace apstra.ClientErr
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrExists)
		})
	}
}

func TestGetTagsByLabels(t *testing.T) {
	ctx := context.Background()
	clients := testclient.GetTestClients(t, ctx)

	labels := make([]string, 2)
	for i := range labels {
		labels[i] = testutils.RandString(5, "hex")
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			labelIds := make([]apstra.ObjectId, len(labels))
			for i := range labels {
				id, err := client.Client.CreateTag(ctx, &apstra.DesignTagRequest{Label: labels[i]})
				require.NoError(t, err)
				t.Cleanup(func() { require.NoError(t, client.Client.DeleteTag(ctx, id)) })
				labelIds[i] = id
			}

			tags, err := client.Client.GetTagsByLabels(context.Background(), labels)
			require.NoError(t, err)
			require.Equal(t, len(labels), len(tags))
		})
	}
}
