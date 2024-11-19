DOCKER_REGISTRY := us-central1-docker.pkg.dev/suite-189022/image-repo
IMAGE_NAME := test-api
IMAGE_TAG := latest

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: dep-tidy
dep-tidy:
	go mod tidy

.PHONY: clean
clean:
	@rd /s /q vendor
	# rm -rf vendor

.PHONY: build
build:
	docker build -f deploy/Dockerfile -t $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: push
push: build
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: deploy-local
deploy-local:
	skaffold run --no-prune=false --cache-artifacts=false --port-forward -f deploy/skaffold/skaffold-local.yaml

.PHONY: deploy-local-verbose
deploy-local-verbose:
	skaffold run --no-prune=false --cache-artifacts=false --port-forward -f deploy/skaffold/skaffold-local.yaml -v debug

.PHONY: deploy-local-clean
deploy-local-clean:
	skaffold delete -f deploy/skaffold/skaffold-local.yaml

.PHONY: deploy-gcp
deploy-gcp:
	skaffold run --no-prune=false --cache-artifacts=false -f deploy/skaffold/skaffold-gcp.yaml

.PHONY: deploy-gcp-verbose
deploy-gcp-verbose:
	skaffold run --no-prune=false --cache-artifacts=false -f deploy/skaffold/skaffold-gcp.yaml -v debug

.PHONY: deploy-gcp-clean
deploy-gcp-clean:
	skaffold delete -f deploy/skaffold/skaffold-gcp.yaml