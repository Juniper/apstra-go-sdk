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
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/stretchr/testify/require"
)

var testTags = map[string]design.Tag{
	"label_6_description_0": {
		Label: testutils.RandString(6, "hex"),
	},
	"label_4_description_20": {
		Label:       testutils.RandString(4, "hex"),
		Description: testutils.RandString(20, "hex"),
	},
}

func TestTag_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.Tag
		update design.Tag
	}

	testCases := map[string]testCase{
		"a": {
			create: testTags["label_6_description_0"],
			update: testTags["label_4_description_20"],
		},
	}

	// Tag-only quirk: "update" tag labels must match "create" tag labels
	for k, v := range testCases {
		v.update.Label = v.create.Label
		testCases[k] = v
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					require.NotEqual(t, tCase.create, zero.Of(tCase.create)) // make sure we didn't use a bogus map key
					require.NotEqual(t, tCase.update, zero.Of(tCase.update)) // make sure we didn't use a bogus map key

					var id string
					var err error
					var obj design.Tag

					// create the object
					id, err = client.Client.CreateTag2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTag2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetTag2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Tag(t, tCase.create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetTagByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Tag(t, tCase.create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListTags2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetTags2(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Tag(t, tCase.create, *objPtr)

					// update the object and validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTag2(ctx, tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetTag2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Tag(t, tCase.update, obj)

					// restore the object to the original state
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTag2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetTag2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Tag(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteTag2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetTag2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTagByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTags2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetTags2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTag2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteTag2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
