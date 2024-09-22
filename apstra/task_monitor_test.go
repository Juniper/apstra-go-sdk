// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"fmt"
	"net/url"
	"testing"
)

func TestBlueprintIdFromUrl(t *testing.T) {
	testBpId := ObjectId(randString(10, "hex"))
	test := "https://host:443" + fmt.Sprintf(apiUrlBlueprintById, testBpId)
	parsed, err := url.Parse(test)
	if err != nil {
		t.Fatal(err)
	}
	resultBpId := blueprintIdFromUrl(parsed)
	if testBpId != resultBpId {
		t.Fatalf("expected '%s', got '%s'", testBpId, resultBpId)
	}
}
