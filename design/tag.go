// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDer            = (*Tag)(nil)
	_ internal.Replicator[Tag] = (*Tag)(nil)
	_ json.Marshaler           = (*Tag)(nil)
	_ json.Unmarshaler         = (*Tag)(nil)
	_ timeutils.Stamper        = (*Tag)(nil)
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

// SetID sets a the value returned by ID only if it was previously un-set. An
// error is returned If the value was previously set. Presence of an existing
// value is the only reason SetID will return an error. If the value is known to
// be empty, use MustSetID.
func (t *Tag) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *Tag) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

// Replicate returns a copy of itself with zero values for metadata fields
func (t Tag) Replicate() Tag {
	return Tag{
		Label:       t.Label,
		Description: t.Description,
	}
}

func (t Tag) CreatedAt() *time.Time {
	return t.createdAt
}

func (t Tag) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func (t Tag) MarshalJSON() ([]byte, error) {
	raw := struct {
		Label       string `json:"label"`
		Description string `json:"description"`
	}{
		Label:       t.Label,
		Description: t.Description,
	}
	return json.Marshal(&raw)
}

func (t *Tag) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string     `json:"id"`
		Label          string     `json:"label"`
		Description    string     `json:"description"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
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

func populateTagsByLabel(parentTags, childTags []Tag) error {
CHILDTAG:
	for i, childTag := range childTags {
		for _, parentTag := range parentTags {
			if childTag.Label == parentTag.Label {
				childTags[i] = parentTag.Replicate()
				continue CHILDTAG
			}
		}
		return fmt.Errorf("tag with label %q not found", childTag.Label)
	}
	sort.Slice(childTags, func(i, j int) bool {
		return childTags[i].Label < childTags[j].Label
	})
	return nil
}
