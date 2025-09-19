// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/slice"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

const (
	TagUrl     = urlPrefix + "tags"
	TagUrlByID = TagUrl + "/%s"
)

var (
	_ json.Marshaler    = (*Tag)(nil)
	_ json.Unmarshaler  = (*Tag)(nil)
	_ slice.IDer        = (*Tag)(nil)
	_ timeutils.Stamper = (*Tag)(nil)
)

type Tag struct {
	Label       string
	Description string

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t Tag) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

func (t Tag) CreatedAt() *time.Time {
	return t.createdAt
}

func (t Tag) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func (t Tag) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID          string `json:"id,omitempty"`
		Label       string `json:"label,omitempty"`
		Description string `json:"description"`
	}{
		ID:          t.id,
		Label:       t.Label,
		Description: t.Description,
	}
	return json.Marshal(&raw)
}

func (t *Tag) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string     `json:"id,omitempty"`
		Label          string     `json:"label"`
		Description    string     `json:"description"`
		CreatedAt      *time.Time `json:"created_at,omitempty"`
		LastModifiedAt *time.Time `json:"last_modified_at,omitempty"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling tag: %w", err)
	}

	t.id = raw.ID
	t.Label = raw.Label
	t.Description = raw.Description
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func NewTag(id string) Tag {
	return Tag{id: id}
}
