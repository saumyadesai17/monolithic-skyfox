#!/bin/sh
set -e
set -v

export PROMOTE_FROM=$1
export IMAGE=$BOOKING_WEB_IMAGE_NAME:$CI_COMMIT_SHORT_SHA-$PROMOTE_FROM

echo "Promoting image - $IMAGE"

docker -v login -u gitlab-ci-token -p "$GITLAB_REGISTRY_PASSWORD" registry.gitlab.com

docker pull "$IMAGE"
docker tag "$IMAGE" "$BOOKING_WEB_IMAGE_NAME:$CI_COMMIT_SHORT_SHA-$CI_ENVIRONMENT_SLUG"
docker push "$BOOKING_WEB_IMAGE_NAME:$CI_COMMIT_SHORT_SHA-$CI_ENVIRONMENT_SLUG"

set +v

