PROJECT_ID = zing-dev-197522
CLUSTER_ID = jpl-gke
ZONE = us-central1-f
TAG = poc-test
FULL_CLUSTER = gke_$(PROJECT_ID)_$(ZONE)_$(CLUSTER_ID)
IMAGE_NAME = gcr.io/zing-registry-188222/grpctest:$(TAG)

SRC_DIR = /go/src/github.com/zenoss/grpctest

default: build

auth:
	@gcloud container clusters get-credentials --project ${PROJECT_ID} ${CLUSTER_ID}

build:
	@docker build -t $(IMAGE_NAME) .
#	@docker push $(IMAGE_NAME)

deploy:
	@kubectl apply -f grpctest.yaml --cluster $(FULL_CLUSTER)

protoc-image:
	@cd pb && docker build -f Dockerfile -t protoc-image .

protoc: protoc-image
	@docker run -v $(CURDIR):${SRC_DIR} --rm protoc-image protoc --proto_path=${SRC_DIR}/googleapis -I=${SRC_DIR} ${SRC_DIR}/pb/grpc_test.proto --go_out=plugins=grpc:${SRC_DIR}
#	@docker run --rm protoc-image protoc --include_imports --include_source_info --proto_path=${GOOGLEAPIS_DIR} --proto_path=$(pwd) --descriptor_set_out=pb/api_descriptor.pb pb/grpc_test.proto

.PHONY: client

client:
	@docker run -w ${SRC_DIR}/client -v $(CURDIR):${SRC_DIR} --rm golang:latest go build .
