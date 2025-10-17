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
	"github.com/stretchr/testify/require"
)

var testTemplatesPodBased = map[string]design.TemplatePodBased{}

func TestTemplatePodBased_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	//template, err := clients[0].Client.GetTemplatePodBased2(ctx, "L2_superspine_multi_plane")
	//require.NoError(t, err)
	//log.Printf("\n%#v\n", template)
	//return

	type testCase struct {
		create design.TemplatePodBased
		update design.TemplatePodBased
	}

	testCases := map[string]testCase{}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					var id string
					var err error
					var obj design.TemplatePodBased

					// create the object (by type)
					id, err = client.Client.CreateTemplatePodBased2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID then validate
					template, err := client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok := template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr := objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, *objPtr)

					// retrieve the object by ID (by type) then validate
					obj, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, obj)

					// retrieve the object by label then validate
					template, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, *objPtr)

					// retrieve the object by label (by type) then validate
					obj, err = client.Client.GetTemplatePodBasedByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must be in there)
					ids, err = client.Client.ListTemplatesPodBased2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) then validate
					templates, err := client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr := slice.MustFindByID(templates, id)
					require.NotNil(t, templatePtr)
					objPtr, ok = (*templatePtr).(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, *objPtr)

					// retrieve the list of objects (by type) (ours must be in there) then validate
					objs, err := client.Client.GetTemplatesPodBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, obj)

					// update the object then validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.update, *objPtr)

					// retrieve the updated object by ID (by type) type then validate
					update, err := client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.update, update)

					// restore the object (by type)
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplatePodBased2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the restored object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, *objPtr)

					// retrieve the restored object by ID (by type) then validate
					obj, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteTemplate2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by ID (by type)
					_, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label (by type)
					_, err = client.Client.GetTemplatePodBasedByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must *not* be in there)
					ids, err = client.Client.ListTemplatesPodBased2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					templates, err = client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr = slice.MustFindByID(templates, id)
					require.Nil(t, templatePtr)

					// retrieve the list of objects (by type) (ours must *not* be in there)
					objs, err = client.Client.GetTemplatesPodBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// update the object (by type)
					err = client.Client.UpdateTemplatePodBased2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
