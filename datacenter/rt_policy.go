// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

type RTPolicy struct {
	ImportRTs []string `json:"import_RTs"`
	ExportRTs []string `json:"export_RTs"`
}
