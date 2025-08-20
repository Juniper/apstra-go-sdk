# Copyright (c) Juniper Networks, Inc., 2022-2024.
# All rights reserved.
# SPDX-License-Identifier: Apache-2.0

all: compliance-check license-header-check verify unit-tests integration-tests

check-repo-clean:
	git update-index --refresh && git diff-index --quiet HEAD --

compliance:
	go run github.com/chrismarget-j/go-licenses save   --ignore github.com/Juniper --save_path Third_Party_Code --force ./... || exit 1 ;\
	go run github.com/chrismarget-j/go-licenses report --ignore github.com/Juniper --template .notices.tpl ./... > Third_Party_Code/NOTICES.md || exit 1 ;\

compliance-check: compliance check-repo-clean

license-header-check:
	@sh -c "$(CURDIR)/.ci/scripts/license_header_check.sh"

fast-check: verify unit-tests

unit-tests:
	go test -v .

integration-tests:
	go test -tags integration -v .

verify: lint-staticcheck fmt-check vet

fmt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofmt_check.sh"

fumpt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofumpt_check.sh"

lint-staticcheck:
	staticcheck -tags integration .

vet:
	go vet -v ./apstra/...

.PHONY: all fmt-check license-header-check lint lint-staticcheck unit-tests verify vet
