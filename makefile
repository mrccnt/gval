.PHONY: snapshot release tag lint --tag --msg
.DEFAULT_GOAL := snapshot
MAKEFLAGS += --no-print-directory

snapshot:
	@goreleaser release --clean --snapshot

release:
	@goreleaser release --clean --skip=publish

tag: --tag --msg
	@git tag -a $(TAG) -m "$(MSG)"
	@git push origin (TAG)

lint:
	@goreleaser check
	@golangci-lint run

--tag:
	@if [ "$(TAG)" = "" ]; then echo "TAG not set" && exit 1; else exit 0; fi

--msg:
	@if [ "$(MSG)" = "" ]; then echo "MSG not set" && exit 1; else exit 0; fi
