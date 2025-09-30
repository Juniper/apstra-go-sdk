// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"fmt"
	sdk "github.com/Juniper/apstra-go-sdk"
)

// FeatureVersion details whether a feature is enabled on the given NOS version
type FeatureVersion struct {
	Version string `json:"version"`
	Enabled bool   `json:"value"`
}

type FeatureVersions []FeatureVersion

// Validate ensures that there are no Version string collisions within a FeatureVersions
func (f FeatureVersions) Validate() error {
	versionMap := make(map[string]struct{}, len(f))
	for _, v := range f {
		if _, ok := versionMap[v.Version]; ok {
			return sdk.ErrMultipleMatch(fmt.Sprintf("duplicate feature version: %s", v.Version))
		}
		versionMap[v.Version] = struct{}{}
	}
	return nil
}
