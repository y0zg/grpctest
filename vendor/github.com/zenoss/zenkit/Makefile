ZENKIT_VERSION       := 1.6.0
ROOTDIR              ?= $(CURDIR)
GINKGO               := $(shell command -v ginkgo 2> /dev/null)
PACKAGE              := github.com/zenoss/zenkit
LOCAL_USER_ID        := $(shell id -u)
BUILD_IMG            := zenoss/zenkit-build:$(ZENKIT_VERSION)
COVERAGE_DIR         := coverage
DOCKER_COMPOSE       := /usr/local/bin/docker-compose

DOCKER_PARAMS        := --rm -v /var/run/docker.sock:/var/run/docker.sock \
						     -v $(ROOTDIR):/go/src/$(PACKAGE):rw \
							 -e LOCAL_USER_ID=$(LOCAL_USER_ID) \
							 -w /go/src/$(PACKAGE)
DOCKER_CMD           := docker run -t $(DOCKER_PARAMS) $(BUILD_IMG)

.PHONY: default
default: test

.PHONY: test-containerized
test-containerized:
	@$(DOCKER_CMD) /bin/bash -c "go get ./... && make test"

.PHONY: test
ifndef GINKGO
test: test-containerized
else
test: COVERAGE_PROFILE := $(COVERAGE_DIR)/profile.out
test: COVERAGE_HTML    := $(COVERAGE_DIR)/index.html
test: COVERAGE_XML     := $(COVERAGE_DIR)/coverage.xml
test:
	@mkdir -p $(COVERAGE_DIR)
	@GODEBUG=netdns=go ginkgo -r \
		-cover \
		-covermode=count \
		-tags="integration" \
		--skipPackage vendor
	@gocovmerge $$(find . -name \*.coverprofile) > $(COVERAGE_PROFILE)
	@go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@gocov convert $(COVERAGE_PROFILE) | gocov-xml > $(COVERAGE_XML)
endif

.PHONY: clean
clean:
	rm -rf $(COVERAGE_DIR) **/*.coverprofile **/junit.xml

local-dev:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/wadey/gocovmerge
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/AlekSi/gocov-xml

.PHONY: docker-compose
docker-compose: $(DOCKER_COMPOSE)

$(DOCKER_COMPOSE): DOCKER_COMPOSE_VERSION := 1.16.1
$(DOCKER_COMPOSE):
	@if [ ! -w $(@D) ]; then echo 'No docker-compose found. Please run "sudo make docker-compose" to install it.'; exit 2; else true; fi
	@curl -L https://github.com/docker/compose/releases/download/$(DOCKER_COMPOSE_VERSION)/docker-compose-`uname -s`-`uname -m` > $@
	@chmod +x $@
