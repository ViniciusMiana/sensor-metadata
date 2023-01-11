PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
GO_EXEC_PATH := .tmp_exec_path
VERSION ?= $(shell git describe --tags)
COMMIT ?= $(shell git rev-parse --short HEAD)
BASE_TAG := viniciusmiana/sensor
AUTH_TAG := viniciusmiana/auth
TAG ?= $(BASE_TAG):$(VERSION)
SCRIPTS_DIR=./scripts
REGISTRY ?= $(shell docker ps -q -f name=registry)

all: build

# TODO Include memory limit, besides GOGC
# TODO Add CI/CD support
# TODO Run mocks and add mocked tests
# TODO Separate into to makefiles and/or separate targets

TEST_CLIENT := client
SOURCES := $(shell find ./ -name '*.go' | grep -v test/$(TEST_CLIENT) )

all: lint test todo

## init: Install tools used by the build
init:
	@echo "[STATUS] Installing tools used by the build"
	mkdir -p  $(GO_EXEC_PATH) out test/$(TEST_CLIENT) mocks
	# install all pre-reqs
	$(SCRIPTS_DIR)/installTools.sh -e swagger -e linter -e mockery
	$(GO_EXEC_PATH)/bin/golangci-lint cache clean

## clean: Delete all generated files on this project
clean:
	@echo "[STATUS] Cleaning everything to start fresh"
	go clean
	rm -rf out
	rm -rf test/$(TEST_CLIENT)
	rm -rf mocks
	rm -rf $(GO_EXEC_PATH)


## mocks: Generate mocks used for testing
mocks: init
	$(GO_EXEC_PATH)/bin/mockery --all --dir cmd/sensor/db --output mocks/sensor/db --case underscore --disable-version-string --exported

## format: Formats the code
format: init swagger mocks
	@echo "[STATUS] Formatting code"
	gofmt -s -w .
	goimports -w .
	go mod tidy

## lint: Run linter
lint: format
	@echo "[STATUS] Running linter"
	GOGC=20 $(GO_EXEC_PATH)/bin/golangci-lint run --config .golangci.yml -v --timeout=10m



## build: Compile this project and add the generated file to out directory
build: $(SOURCES) init
	@echo "[STATUS] Building app build $(VERSION)"
	mkdir -p ./out
	CGO_ENABLED=0 go build -ldflags "-w -s -X main.version=$(VERSION)" -installsuffix 'static' -tags timetzdata -o ./out/sensor ./cmd/sensor/main.go
	CGO_ENABLED=0 go build -ldflags "-w -s -X main.version=$(VERSION)" -installsuffix 'static' -tags timetzdata -o ./out/authenticator ./cmd/authenticator/main.go

## test: Run all repository tests
test: build swagger | start_mongo
	GOGC=20 go test -v -p 1 -timeout 900s -covermode=count -coverprofile=./out/coverage.out  ./test ./cmd/sensor/...
	make stop_mongo

## todo: Generate the TODO.md file containing all repository TODOS
todo: TODO.md

TODO.md:  $(SOURCES)
	@echo "[STATUS] Updating TODO.md"
	@grep -rn --exclude=Makefile --exclude=TODO.md --exclude=\*.{pack,bin,md} --exclude-dir=./test/client '\/\/\s*\(TODO\|todo\)' . | grep -v 'Binary' > 'TODO.md'


## deploy: Deploy the docker image to the local cluster
deploy: docker
	@echo "[STATUS] Deploying to local cluster"
	-helm uninstall sensor
	helm install sensor ./deployments


## docker: Create the docker image
docker: dockerSensor dockerAuth

dockerAuth:  cmd/authenticator/Dockerfile
	@echo "[STATUS] Creating auth image: ${AUTH_TAG}:${VERSION}-${COMMIT}"
	docker build -t  ${AUTH_TAG}:${VERSION}-${COMMIT} -f cmd/authenticator/Dockerfile . --progress=plain --platform linux/amd64
	docker tag ${AUTH_TAG}:${VERSION}-${COMMIT} ${AUTH_TAG}:latest
	docker push ${AUTH_TAG}:${VERSION}-${COMMIT}
	docker push ${AUTH_TAG}:latest
	@echo "[STATUS] Image created:  ${AUTH_TAG}:${VERSION}-${COMMIT}"

dockerSensor: cmd/sensor/Dockerfile
	@echo "[STATUS] Creating image: ${TAG}-${COMMIT}"
	docker build -t  ${TAG}-${COMMIT} -f cmd/sensor/Dockerfile . --progress=plain --platform linux/amd64
	docker tag ${TAG}-${COMMIT} $(BASE_TAG):latest
	docker push ${TAG}-${COMMIT}
	docker push $(BASE_TAG):latest
	@echo "[STATUS] Image created:  ${TAG}-${COMMIT}"

## start_images: starts mongo
start_mongo:
	$(SCRIPTS_DIR)/mongo.sh START

## stop_images: stops mongo
stop_mongo:
	$(SCRIPTS_DIR)/mongo.sh STOP

## swagger: Generate the swagger client
swagger: api/swagger.yml
	@echo "[STATUS] Generating swagger client"
	mkdir -p test/$(TEST_CLIENT)
	$(GO_EXEC_PATH)/bin/swagger generate client --spec=api/swagger.yml --target=test/$(TEST_CLIENT)

## cover: Analyzes test coverage
cover: test
	go tool cover -func ./out/coverage.out


help: Makefile
	@echo
	@echo " Make rules for "$(ORG)/$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/\t/'
	@echo
