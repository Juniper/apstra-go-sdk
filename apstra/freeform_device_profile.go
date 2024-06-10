package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlFFDp       = apiUrlBlueprintById + apiUrlPathDelim + "device-profiles"
	apiUrlFFDpImport = apiUrlFFDp + apiUrlPathDelim + "import"
)

func (o *FreeformClient) ImportDeviceProfile(ctx context.Context, dpName string) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFDpImport, o.blueprintId),
		apiInput:    dpName,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}
