// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestPeekParseResponseBodyAsTaskId(t *testing.T) {
	testData := "{\"id\": \"9e299b2f-7a24-4358-a10c-d93a71ed120d\", \"task_id\": \"48f06e19-7177-4f9d-bf65-0bb8949e1cce\"}"
	httpResp := http.Response{Body: io.NopCloser(bytes.NewReader([]byte(testData)))}
	var taskResp taskIdResponse

	ok, err := peekParseResponseBodyAsTaskId(&httpResp, &taskResp)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(ok, len(testData))
}
