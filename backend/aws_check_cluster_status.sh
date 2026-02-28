#!/bin/bash

aws configure set aws_access_key_id $AWS_ACCESS_KEY
aws configure set aws_secret_access_key $AWS_SECRET_KEY

while true; do
    eval status=`aws ecs describe-clusters --cluster "$CLUSTER_NAME" --region ap-south-1 | jq .clusters[].status`
    aws ecs describe-clusters --cluster "$CLUSTER_NAME" --region ap-south-1
    if [ "$status" == "ACTIVE" ]; then
        break
    fi
    sleep 10
done