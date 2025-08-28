# Copyright (c) Juniper Networks, Inc., 2022-2025.
# All rights reserved.
# SPDX-License-Identifier: Apache-2.0

all: compliance-check license-header-check verify unit-tests integration-tests

check-repo-clean:
	git update-index --refresh && git diff-index --quiet HEAD --

compliance:
	@sh -c "$(CURDIR)/.ci/scripts/compliance.sh"

compliance-check: compliance check-repo-clean

license-header-check:
	@sh -c "$(CURDIR)/.ci/scripts/license_header_check.sh"

fast-check: verify unit-tests

unit-tests:
	go test -v .

integration-tests:
	go test -tags integration -v .

verify: fmt-check vet

fmt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofmt_check.sh"

fumpt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofumpt_check.sh"

vet:
	go vet -v ./apstra/...

.PHONY: all fmt-check license-header-check unit-tests verify vet
