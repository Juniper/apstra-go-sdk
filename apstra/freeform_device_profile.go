package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlFFDp       = apiUrlBlueprintById + apiUrlPathDelim + "device-profiles"
	apiUrlFFDpImport = apiUrlFFDp + apiUrlPathDelim + "import"
	apiUrlFFDpById   = apiUrlFFDp + apiUrlPathDelim + "%s"
)

func (o *FreeformClient) ImportDeviceProfile(ctx context.Context, id ObjectId) (ObjectId, error) {
	var response struct {
		Ids []ObjectId `json:"ids"`
	}
	input := struct {
		DeviceProfileIds []ObjectId `json:"device_profile_ids"`
	}{
		DeviceProfileIds: []ObjectId{id},
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFDpImport, o.blueprintId),
		apiInput:    input,
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

func (o *FreeformClient) DeleteDeviceProfile(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFDpById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
