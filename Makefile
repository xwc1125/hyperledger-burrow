# ----------------------------------------------------------
# REQUIREMENTS

# - Go 1.12
# - Make
# - jq
# - find
# - bash
# - protoc (for rebuilding protobuf files)

# ----------------------------------------------------------

SHELL := /bin/bash
REPO := $(shell pwd)

# Our own Go files containing the compiled bytecode of solidity files as a constant

CI_IMAGE="hyperledger/burrow:ci"

VERSION := $(shell scripts/version.sh)
# Gets implicit default GOPATH if not set
GOPATH?=$(shell go env GOPATH)
BIN_PATH?=$(GOPATH)/bin
HELM_PATH?=helm/package
HELM_PACKAGE=$(HELM_PATH)/burrow-$(VERSION).tgz
ARCH?=linux-amd64

export GO111MODULE=on

### Formatting, linting and vetting

# check the code for style standards; currently enforces go formatting.
# display output first, then check for success
.PHONY: check
check:
	@echo "Checking code for formatting style compliance."
	@gofmt -l -d $(shell go list -f "{{.Dir}}" ./...)
	@gofmt -l $(shell go list -f "{{.Dir}}" ./...) | read && echo && echo "Your marmot has found a problem with the formatting style of the code." 1>&2 && exit 1 || true

# Just fix it
.PHONY: fix
fix:
	@goimports -l -w $(shell go list -f "{{.Dir}}" ./...)

# fmt runs gofmt -w on the code, modifying any files that do not match
# the style guide.
.PHONY: fmt
fmt:
	@echo "Correcting any formatting style corrections."
	@gofmt -l -w $(shell go list -f "{{.Dir}}" ./...)

# lint installs golint and prints recommendations for coding style.
lint:
	@echo "Running lint checks."
	go get -u github.com/golang/lint/golint
	@for file in $(shell go list -f "{{.Dir}}" ./...); do \
		echo; \
		golint --set_exit_status $${file}; \
	done

# vet runs extended compilation checks to find recommendations for
# suspicious code constructs.
.PHONY: vet
vet:
	@echo "Running go vet."
	@go vet $(shell go list ./... )

# run the megacheck tool for code compliance
.PHONY: megacheck
megacheck:
	@go get honnef.co/go/tools/cmd/megacheck
	@for pkg in $(shell go list ./... ); do megacheck "$$pkg"; done

# Protobuffing

BURROW_TS_PATH = ./js
PROTO_GEN_TS_PATH = ${BURROW_TS_PATH}/proto
NODE_BIN = ${BURROW_TS_PATH}/node_modules/.bin

PROTO_FILES = $(shell find . -path $(BURROW_TS_PATH) -prune -o -path ./node_modules -prune -o -type f -name '*.proto' -print)
PROTO_GO_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
PROTO_GO_FILES_REAL = $(shell find . -type f -name '*.pb.go' -print)
PROTO_TS_FILES = $(patsubst %.proto, %.pb.ts, $(PROTO_FILES))

.PHONY: protobuf
protobuf: $(PROTO_GO_FILES) $(PROTO_TS_FILES) fix

# Implicit compile rule for GRPC/proto files (note since pb.go files no longer generated
# in same directory as proto file this just regenerates everything
%.pb.go: %.proto
	protoc -I ./protobuf $< --gogo_out=plugins=grpc:${GOPATH}/src

# Using this: https://github.com/agreatfool/grpc_tools_node_protoc_ts
%.pb.ts: %.proto
	$(NODE_BIN)/grpc_tools_node_protoc -I protobuf \
		--plugin="protoc-gen-ts=$(NODE_BIN)/protoc-gen-ts" \
		--js_out="import_style=commonjs,binary:${PROTO_GEN_TS_PATH}" \
		--ts_out="generate_package_definition:${PROTO_GEN_TS_PATH}" \
		--grpc_out="generate_package_definition:${PROTO_GEN_TS_PATH}" \
		$<

.PHONY: protobuf_deps
protobuf_deps:
	@go get -u github.com/gogo/protobuf/protoc-gen-gogo
	@cd ${BURROW_TS_PATH} && yarn install --only=dev

.PHONY: clean_protobuf
clean_protobuf:
	@rm -f $(PROTO_GO_FILES_REAL)

### PEG query grammar

# This allows us to filter tagged objects with things like (EventID = 'foo' OR Height > 10) AND EventName CONTAINS 'frog'

.PHONY: peg_deps
peg_deps:
	go get -u github.com/pointlander/peg

# regenerate the parser
.PHONY: peg
peg:
	peg event/query/query.peg

### Building github.com/hyperledger/burrow

# Output commit_hash but only if we have the git repo (e.g. not in docker build
.PHONY: commit_hash
commit_hash:
	@git status &> /dev/null && scripts/commit_hash.sh > commit_hash.txt || true

# build all targets in github.com/hyperledger/burrow
.PHONY: build
build:	check build_burrow build_burrow_debug

# build all targets in github.com/hyperledger/burrow with checks for race conditions
.PHONY: build_race
build_race:	check build_race_db

# build burrow and vent
.PHONY: build_burrow
build_burrow: commit_hash
	go build $(BURROW_BUILD_FLAGS) -ldflags "-extldflags '-static' \
	-X github.com/hyperledger/burrow/project.commit=$(shell cat commit_hash.txt) \
	-X github.com/hyperledger/burrow/project.date=$(shell date '+%Y-%m-%d')" \
	-o ${REPO}/bin/burrow$(BURROW_BUILD_SUFFIX) ./cmd/burrow

# With the sqlite tag - enabling Vent sqlite adapter support, but building a CGO binary
.PHONY: build_burrow_sqlite
build_burrow_sqlite: export BURROW_BUILD_SUFFIX=-vent-sqlite
build_burrow_sqlite: export BURROW_BUILD_FLAGS=-tags sqlite
build_burrow_sqlite:
	$(MAKE) build_burrow

# Builds a binary suitable for delve line-by-line debugging through CGO with optimisations (-N) and inling (-l) disabled
.PHONY: build_burrow_debug
build_burrow_debug: export BURROW_BUILD_SUFFIX=-debug
build_burrow_debug: export BURROW_BUILD_FLAGS=-gcflags "all=-N -l"
build_burrow_debug:
	$(MAKE) build_burrow

.PHONY: install
install: build_burrow
	mkdir -p ${BIN_PATH}
	install -T ${REPO}/bin/burrow ${BIN_PATH}/burrow

# build burrow with checks for race conditions
.PHONY: build_race_db
build_race_db:
	go build -race -o ${REPO}/bin/burrow ./cmd/burrow

### Build docker images for github.com/hyperledger/burrow

# build docker image for burrow
.PHONY: docker_build
docker_build: check commit_hash
	@scripts/build_tool.sh

### Testing github.com/hyperledger/burrow

# Solidity fixtures
.PHONY: solidity
solidity: $(patsubst %.sol, %.sol.go, $(wildcard ./execution/solidity/*.sol)) build_burrow

%.sol.go: %.sol
	@burrow compile $^

# Solang fixtures
.PHONY: solang
solang: $(patsubst %.solang, %.solang.go, $(wildcard ./execution/wasm/*.solang)) build_burrow

%.solang.go: %.solang
	@burrow compile --wasm $^

# node/js
.PHONY: yarn_install
yarn_install:
	@cd ${BURROW_TS_PATH} && yarn install

# Test

.PHONY: test_js
test_js:
	@cd ${BURROW_TS_PATH} && yarn test

.PHONY: test
test: check bin/solc
	@tests/scripts/bin_wrapper.sh go test ./... ${GO_TEST_ARGS}

.PHONY: test_keys
test_keys:
	burrow_bin="${REPO}/bin/burrow" tests/keys_server/test.sh

.PHONY:	test_truffle
test_truffle:
	burrow_bin="${REPO}/bin/burrow" tests/web3/truffle.sh

.PHONY:	test_integration_vent
test_integration_vent:
	# Include sqlite adapter with tests - will build with CGO but that's probably fine
	go test -count=1 -v -tags 'integration sqlite' ./vent/...

.PHONY:	test_integration_vent_postgres
test_integration_vent_postgres:
	docker-compose run burrow make test_integration_vent

.PHONY: test_restore
test_restore:
	@tests/scripts/bin_wrapper.sh tests/dump/test.sh

# Go will attempt to run separate packages in parallel

.PHONY: test_integration
test_integration:
	@go test -count=1 -v -tags integration ./integration/...

.PHONY: test_integration_all
test_integration_all: test_keys test_deploy test_integration_vent_postgres test_restore test_truffle test_integration

.PHONY: test_integration_all_no_postgres
test_integration_all_no_postgres: test_keys test_deploy test_integration_vent test_restore test_truffle test_integration

.PHONY: test_deploy
test_deploy:
	@tests/scripts/bin_wrapper.sh tests/deploy.sh

bin/solc: ./tests/scripts/deps/solc.sh
	@mkdir -p bin
	@tests/scripts/deps/solc.sh bin/solc
	@touch bin/solc

# test burrow with checks for race conditions
.PHONY: test_race
test_race: build_race
	@go test -race $(shell go list ./... )

### Clean up

# clean removes the target folder containing build artefacts
.PHONY: clean
clean:
	-rm -r ./bin

### Release and versioning

# Print version
.PHONY: version
version:
	@echo $(VERSION)

# Generate full changelog of all release notes
CHANGELOG.md: project/history.go project/cmd/changelog/main.go
	@go run ./project/cmd/changelog/main.go > CHANGELOG.md

# Generated release note for this version
NOTES.md: project/history.go project/cmd/notes/main.go
	@go run ./project/cmd/notes/main.go > NOTES.md

.PHONY: docs
docs: CHANGELOG.md NOTES.md

# Tag the current HEAD commit with the current release defined in
# ./project/history.go
.PHONY: tag_release
tag_release: test check docs build
	@scripts/tag_release.sh

.PHONY: build_ci_image
build_ci_image:
	docker build -t ${CI_IMAGE} -f ./.github/Dockerfile .

.PHONY: push_ci_image
push_ci_image: build_ci_image
	docker push ${CI_IMAGE}

.PHONY: ready_for_pull_request
ready_for_pull_request: docs fix

.PHONY: staticcheck
staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

# Note --set flag currently needs helm 3 version < 3.0.3 https://github.com/helm/helm/issues/3141 - but hopefully they will reintroduce support
bin/helm:
	@echo Downloading helm...
	mkdir -p bin
	curl https://get.helm.sh/helm-v3.0.2-$(ARCH).tar.gz | tar xvzO $(ARCH)/helm > bin/helm && chmod +x bin/helm

.PHONY: helm_deps
helm_deps: bin/helm
	@bin/helm repo add --username "$(CM_USERNAME)" --password "$(CM_PASSWORD)" chartmuseum $(CM_URL)

.PHONY: helm_test
helm_test: bin/helm
	bin/helm dep up helm/burrow
	bin/helm lint helm/burrow

helm_package: $(HELM_PACKAGE)

$(HELM_PACKAGE): helm_test bin/helm
	bin/helm package helm/burrow \
		--version "$(VERSION)" \
		--app-version "$(VERSION)" \
		--set "image.tag=$(VERSION)" \
		--dependency-update \
		--destination helm/package

.PHONY: helm_push
helm_push: helm_package
	@echo pushing helm chart...
	@curl -u ${CM_USERNAME}:${CM_PASSWORD} \
		--data-binary "@$(HELM_PACKAGE)" $(CM_URL)/api/charts
