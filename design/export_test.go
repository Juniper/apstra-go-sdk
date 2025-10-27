// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"fmt"
)

// The `SetID()` function on each design object is not available to end users
// due to the `_test.go` filename, but we use it in tests.

func (c *ConfigTemplate) SetID(id string) {
	if c.id != "" {
		panic(fmt.Sprintf("id already has value %q", c.id))
	}

	c.id = id
	return
}

func (c *Configlet) SetID(id string) {
	if c.id != "" {
		panic(fmt.Sprintf("id already has value %q", c.id))
	}

	c.id = id
	return
}

func (i *InterfaceMap) SetID(id string) {
	if i.id != "" {
		panic(fmt.Sprintf("id already has value %q", i.id))
	}

	i.id = id
	return
}

func (l *LogicalDevice) SetID(id string) {
	if l.id != "" {
		panic(fmt.Sprintf("id already has value %q", l.id))
	}

	l.id = id
	return
}

func (r *RackType) SetID(id string) {
	if r.id != "" {
		panic(fmt.Sprintf("id already has value %q", r.id))
	}

	r.id = id
	return
}

func (t *Tag) SetID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}

func (t *TemplateL3Collapsed) SetID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}

func (t *TemplatePodBased) SetID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}

func (t *TemplateRackBased) SetID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}

func (t *TemplateRailCollapsed) SetID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}
