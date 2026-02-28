
#!/bin/bash
set -e
set -v

apt-get install nodejs -y
apt-get install jq -y

export CLUSTER_NAME="team-trainers"
export SERVICE="movie-service-$CI_ENVIRONMENT_SLUG"
export CLUSTER_CONFIG_NAME="team-trainers-config"
export CLUSTER_PROFILE_NAME="team-trainers-profile"

export IMAGE=$IMAGE_NAME:latest
echo  "Image URL: $IMAGE"

aws configure set aws_access_key_id "$AWS_ACCESS_KEY"
aws configure set aws_secret_access_key "$AWS_SECRET_KEY"
aws configure set default.region ap-south-1 --profile gitlabci

echo "Generating task definitions and service definitions"
node ecs.js serviceTemplate
node ecs.js taskDefinitionTemplate

echo "Registering task definition"
aws ecs register-task-definition --cli-input-json "file://task-definition.json"

status=$(aws ecs list-services --cluster "$CLUSTER_NAME" | grep "$SERVICE") || true

if [ -z "$status" ]; then
  echo "Registering service definition"
  aws ecs create-service --cli-input-json file://service.json
else
  echo "Service already exists. Updating service.."
  aws ecs update-service --force-new-deployment --service "$SERVICE" --task-definition "$SERVICE" --cluster "$CLUSTER_NAME"
fi

echo "running clean up"
rm -rf service.json
rm -rf task-definition.json

set +v

