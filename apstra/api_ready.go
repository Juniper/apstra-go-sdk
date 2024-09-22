// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
)

// ready tests whether an Apstra server is ready to handle client requests.
// It's based on a code sample shared by JP Senior
// https://apstra-eng.slack.com/archives/C2DFCFHJR/p1674246878325579?thread_ts=1674244723.467159&cid=C2DFCFHJR
func (o *Client) ready(ctx context.Context) error {
	var err error

	_, err = o.getVersion(ctx)
	if err != nil {
		return fmt.Errorf("error getting Version - %w", err)
	}

	_, err = o.listSystems(ctx)
	if err != nil {
		return fmt.Errorf("error listing Systems - %w", err)
	}

	_, err = o.getAllBlueprintStatus(ctx)
	if err != nil {
		return fmt.Errorf("error listing Blueprints - %w", err)
	}

	_, err = o.getAllSystemAgents(ctx)
	if err != nil {
		return fmt.Errorf("error getting System Agents - %w", err)
	}

	_, err = o.getAuditConfig(ctx)
	if err != nil {
		return fmt.Errorf("error getting Audit Config- %w", err)
	}

	_, err = o.getAsnPools(ctx)
	if err != nil {
		return fmt.Errorf("error getting ASN Pools - %w", err)
	}

	_, err = o.getTelemetryQuery(ctx)
	if err != nil {
		return fmt.Errorf("error getting Telemetry Query - %w", err)
	}

	return nil
}

// waitUntilReady calls ready() until it returns without error or the supplied
// context is cancelled.
func (o *Client) waitUntilReady(ctx context.Context) error {
	for {
		err := o.ready(ctx)
		if err != nil {
			break
		}
	}

	return nil
}
