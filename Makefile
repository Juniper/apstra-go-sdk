all: compliance-check verify unit-tests integration-tests

check-repo-clean:
	git update-index --refresh && git diff-index --quiet HEAD --

compliance:
	go run github.com/chrismarget-j/go-licenses save   --ignore github.com/Juniper --save_path Third_Party_Code --force ./... || exit 1 ;\
	go run github.com/chrismarget-j/go-licenses report --ignore github.com/Juniper --template .notices.tpl ./... > Third_Party_Code/NOTICES.md || exit 1 ;\

compliance-check: compliance check-repo-clean


fast-check: verify unit-tests

unit-tests:
	go test -v .

integration-tests:
	go test -tags integration -v .

verify: lint-revive lint-staticcheck fmt-check vet

fmt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofmt_check.sh"

fumpt-check:
	@sh -c "$(CURDIR)/.ci/scripts/gofumpt_check.sh"

lint-revive:
	revive -set_exit_status -config revive.toml .

lint-staticcheck:
	staticcheck -tags integration .

vet:
	go vet -v ./apstra/...

.PHONY: all fmt-check lint lint-revive lint-staticcheck unit-tests verify vet
