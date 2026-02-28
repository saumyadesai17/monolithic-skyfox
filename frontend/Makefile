ARTIFACTORY_DOCKER_REGISTRY = artifactory.idfcfirstbank.com/neev-docker
IMAGE_NAME = falcon-booking-frontend

REVISION := $(GO_REVISION_FALCON_BOOKING_FRONTEND)

define METADATA
export GO_REVISION_FALCON_BOOKING_FRONTEND=$(REVISION)
export GIT_REVISION=$(GO_REVISION_FALCON_BOOKING_FRONTEND)
export GIT_REPO_NAME=falcon-booking-frontend
endef
export METADATA

generate_metadata:
	echo "$$METADATA" > metadata

docker/login:
	@docker login -u $(ARTIFACTORY_HA_USERNAME) -p $(ARTIFACTORY_HA_PASSWORD) $(ARTIFACTORY_DOCKER_REGISTRY)

docker/lint: docker/login
	docker pull $(ARTIFACTORY_DOCKER_REGISTRY)/hadolint/hadolint
	docker run --rm \
	-v $(PWD)/.hadolint.yaml:/.hadolint.yaml \
	-i $(ARTIFACTORY_DOCKER_REGISTRY)/hadolint/hadolint < Dockerfile

docker/build: docker/login
	docker build \
	--build-arg REGISTRY_URI="${ARTIFACTORY_DOCKER_REGISTRY}" \
	--build-arg ARTIFACTORY_HA_NPM_TOKEN="${ARTIFACTORY_HA_NPM_TOKEN}" \
	--build-arg ENVIRONMENT="${ENVIRONMENT}" \
	--tag ${IMAGE_NAME} .


docker/publish: docker/login
	docker tag ${IMAGE_NAME} $(ARTIFACTORY_DOCKER_REGISTRY)/${IMAGE_NAME}:${REVISION}
	docker push $(ARTIFACTORY_DOCKER_REGISTRY)/${IMAGE_NAME}:${REVISION}


.PHONY: test lint fmt compile deps
