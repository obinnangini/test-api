DOCKER_REGISTRY := us-central1-docker.pkg.dev/suite-189022/image-repo
IMAGE_NAME := test-api
IMAGE_TAG := latest

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: dep-tidy
dep-tidy:
	 go clean -modcache && go mod tidy

.PHONY: clean
clean:
	@rd /s /q vendor
	# rm -rf vendor

.PHONY: compose-clean
compose-clean:
	@rd /s /q deploy\.compose-data

.PHONY: build
build:
	docker build -f deploy/Dockerfile -t $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: compose-up
compose-up:
	docker compose -f deploy/compose.yaml up --build --abort-on-container-failure --renew-anon-volumes --remove-orphans

.PHONY: compose-down
compose-down:
	docker compose -f deploy/compose.yaml down

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