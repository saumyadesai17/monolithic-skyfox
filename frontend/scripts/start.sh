#!/bin/bash
set -e
set -v

echo "Running this as `whoami`"
echo "Currently running this in $PWD"

IFS='-' read -ra identifiers <<< $DEPLOYMENT_GROUP_NAME  # DEPLOYMENT_GROUP_NAME is of form neev-xx-xx-[backend|frontend]-[deployment|seed]-[integration|staging|production]
export BATCH_ID=${identifiers[1]}
export TEAM_ID=${identifiers[2]}
export ENVIRONMENT=${identifiers[5]}
export PREFIX="/neev-$BATCH_ID/team-$TEAM_ID/$ENVIRONMENT"

cd /home/ec2-user/deployment/
mkdir -p frontend_$ENVIRONMENT
mv frontend-outputs.sh frontend_$ENVIRONMENT/
cd frontend_$ENVIRONMENT
echo "Existing contents of the directory are"
ls

. ./frontend-outputs.sh # exports UI_IMAGE

export UI_PORT=`aws ssm get-parameters --name "$PREFIX/UI_PORT" | jq ".Parameters[0].Value" | tr -d \"`
export REGISTRY_ID=`aws ssm get-parameters --name "$PREFIX/REGISTRY_ID" | jq ".Parameters[0].Value" | tr -d \"`

env > /home/ec2-user/environment_variables_of_frontend_for_$ENVIRONMENT

echo "Logging into ECR"
$(aws ecr get-login --no-include-email --registry-ids $REGISTRY_ID)

echo "Removing last container if it was there..."
docker rm -f frontend_${ENVIRONMENT} || true

echo "Starting new container../"
docker container run -p $UI_PORT:80 --rm --name frontend_${ENVIRONMENT} -d $UI_IMAGE