all: build verify unit-tests integration-tests

fast-check: verify unit-tests

unit-tests:
	go test -v .

integration-tests:
	go test -tags integration -v .

build:
	go build .

verify: lint-revive lint-staticcheck fmt-check vet

fmt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofmtcheck.sh"

lint-revive:
	revive -set_exit_status -config revive.toml .

lint-staticcheck:
	staticcheck -tags integration .

vet:
	go vet


.PHONY: all build fmt-check lint lint-revive lint-staticcheck unit-tests verify vet
