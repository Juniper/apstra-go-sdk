// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const apiUrlFeatures = "/api/features"

func (o *Client) getFeatures(ctx context.Context) error {
	var response map[string]struct {
		Status string `json:"status"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlFeatures,
		apiResponse: &response,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	o.lock(apiUrlFeatures)
	defer o.unlock(apiUrlFeatures)

	o.features = make(map[enum.ApiFeature]bool, len(response))
	for f, s := range response {
		// Parse the map key (feature name)
		var feature enum.ApiFeature
		err = feature.FromString(f)
		if err != nil {
			var eErr enum.Error
			if errors.As(err, &eErr) && eErr.Type() == enum.ErrorTypeParsingFailed {
				// permit unknown values by shoving them into the enum directly
				feature = enum.ApiFeature{Value: f}
			} else {
				return err
			}
		}

		// Parse the map value (feature status). This should always be "enabled" or "disabled".
		var status enum.FeatureSwitch
		err = status.FromString(s.Status)
		if err != nil {
			return err
		}

		// Add the feature status to the o.features map
		switch status {
		case enum.FeatureSwitchEnabled:
			o.features[feature] = true
		case enum.FeatureSwitchDisabled:
			o.features[feature] = false
		default:
			return fmt.Errorf("status for feature %q has unknown value %q", feature, status)
		}
	}

	return nil
}

func (o *Client) FeatureEnabled(f enum.ApiFeature) bool {
	o.lock(apiUrlFeatures)
	defer o.unlock(apiUrlFeatures)

	return o.features[f]
}

func (o *Client) FeatureExists(f enum.ApiFeature) bool {
	o.lock(apiUrlFeatures)
	defer o.unlock(apiUrlFeatures)

	_, exists := o.features[f]

	return exists
}
