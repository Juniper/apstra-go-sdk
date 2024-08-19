package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlFfDp       = apiUrlBlueprintById + apiUrlPathDelim + "device-profiles"
	apiUrlFfDpImport = apiUrlFfDp + apiUrlPathDelim + "import"
	apiUrlFfDpById   = apiUrlFfDp + apiUrlPathDelim + "%s"
)

func (o *FreeformClient) ImportDeviceProfile(ctx context.Context, id ObjectId) (ObjectId, error) {
	var response struct {
		Ids []ObjectId `json:"ids"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		urlStr: fmt.Sprintf(apiUrlFfDpImport, o.blueprintId),
		apiInput: struct {
			DeviceProfileIds []ObjectId `json:"device_profile_ids"`
		}{
			DeviceProfileIds: []ObjectId{id},
		},
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	if len(response.Ids) != 1 {
		return "", fmt.Errorf("expected one ObjectId got %d", len(response.Ids))
	}

	return response.Ids[0], nil
}

func (o *FreeformClient) GetDeviceProfilesByName(ctx context.Context, desired string) ([]DeviceProfile, error) {
	deviceProfiles, err := o.getAllDeviceProfiles(ctx)
	if err != nil {
		return nil, err
	}
	var result []DeviceProfile
	for _, deviceProfile := range deviceProfiles {
		if deviceProfile.Label == desired {
			result = append(result, *deviceProfile.polish())
		}
	}
	return result, nil
}

func (o *FreeformClient) GetDeviceProfilesById(ctx context.Context, desired ObjectId) (*DeviceProfile, error) {
	deviceProfiles, err := o.getAllDeviceProfiles(ctx)
	if err != nil {
		return nil, err
	}
	var result DeviceProfile
	for _, deviceProfile := range deviceProfiles {
		if deviceProfile.Id == desired {
			result = *deviceProfile.polish()
		}
	}
	// todo I think we need to figure out how to return this as a DeviceProfile type instead of a []rawDeviceProfile... is that the Polish[]?
	// todo please have Chris review.
	return &result, nil
}

func (o *FreeformClient) getAllDeviceProfiles(ctx context.Context) ([]rawDeviceProfile, error) {
	response := &struct {
		Items []rawDeviceProfile
	}{}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfDp, o.blueprintId),
		apiResponse: response,
	})
	return response.Items, convertTtaeToAceWherePossible(err)
}

func (o *FreeformClient) UpdateDeviceProfile(ctx context.Context, id ObjectId) error {
	var response struct {
		Ids []ObjectId `json:"ids"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		urlStr: fmt.Sprintf(apiUrlFfDpImport, o.blueprintId),
		apiInput: struct {
			DeviceProfileIds []ObjectId `json:"device_profile_ids"`
		}{
			DeviceProfileIds: []ObjectId{id},
		},
		apiResponse: &response,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	if len(response.Ids) != 1 {
		return fmt.Errorf("expected one ObjectId got %d", len(response.Ids))
	}

	return nil
}

func (o *FreeformClient) DeleteDeviceProfile(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfDpById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
