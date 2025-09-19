// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestTag_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.Tag
		update design.Tag
	}

	testCases := map[string]testCase{
		"a": {
			create: design.Tag{
				Label:       testutils.RandString(6, "hex"),
				Description: testutils.RandString(10, "hex"),
			},
			update: design.Tag{Description: testutils.RandString(20, "hex")},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					var id string
					var err error
					var create design.Tag

					t.Run("create_test_obj", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						id, err = client.Client.CreateTag2(ctx, tCase.create)
						require.NoError(t, err)
					})

					t.Run("get_test_obj_by_id_and_compare", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						create, err = client.Client.GetTag2(ctx, id)
						require.NoError(t, err)
						idPtr := create.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						compare.Tag(t, tCase.create, create)
					})

					t.Run("get_test_obj_by_id_and_compare", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						create, err = client.Client.GetTagByLabel2(ctx, tCase.create.Label)
						require.NoError(t, err)
						idPtr := create.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						compare.Tag(t, tCase.create, create)
					})

					t.Run("find_id_in_list", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						ids, err := client.Client.ListTags2(ctx)
						require.NoError(t, err)
						require.Contains(t, ids, id)
					})

					t.Run("find_obj_in_list", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						objs, err := client.Client.GetTags2(ctx)
						require.NoError(t, err)
						objPtr := slice.ObjectWithID(objs, id)
						require.NotNil(t, objPtr)
						idPtr := objPtr.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						compare.Tag(t, tCase.create, *objPtr)
					})

					var update design.Tag

					t.Run("prepare_obj_update_payload", func(t *testing.T) {
						update = design.NewTag(id)
						update.Label = create.Label // label is immutable on tags
						update.Description = tCase.update.Description
						require.NotNil(t, update.ID())
						require.Equal(t, id, *update.ID())
					})

					t.Run("update_test_obj", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						err = client.Client.UpdateTag2(ctx, update)
						require.NoError(t, err)
					})

					t.Run("get_updated_obj_by_id_and_compare", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						obj, err := client.Client.GetTag2(ctx, id)
						require.NoError(t, err)
						idPtr := obj.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						compare.Tag(t, update, obj)
					})

					t.Run("delete_obj", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						err = client.Client.DeleteTag2(ctx, id)
						require.NoError(t, err)
					})

					var ace apstra.ClientErr

					t.Run("get_deleted_obj_by_id", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						_, err = client.Client.GetTag2(ctx, id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					})

					t.Run("get_deleted_obj_by_label", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						_, err = client.Client.GetTagByLabel2(ctx, tCase.create.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					})

					t.Run("fail_to_find_id_in_list", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						ids, err := client.Client.ListTags2(ctx)
						require.NoError(t, err)
						require.NotContains(t, ids, id)
					})

					t.Run("fail_to_find_obj_in_list", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						objs, err := client.Client.GetTags2(ctx)
						require.NoError(t, err)
						objPtr := slice.ObjectWithID(objs, id)
						require.Nil(t, objPtr)
					})

					t.Run("fail_to_update_obj", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						err = client.Client.UpdateTag2(ctx, update)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					})

					t.Run("fail_to_delete_obj", func(t *testing.T) {
						ctx := testutils.ContextWithTestID(ctx, t)
						err = client.Client.DeleteTag2(ctx, id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					})
				})
			}
		})
	}
}
