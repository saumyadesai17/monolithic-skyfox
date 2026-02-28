#!/bin/bash
set -e

function check_status() {
  ecs-cli ps --cluster-config "$CLUSTER_CONFIG_NAME" --ecs-profile "$CLUSTER_PROFILE_NAME" --desired-status RUNNING
  status=$(ecs-cli ps --cluster-config "$CLUSTER_CONFIG_NAME" --ecs-profile "$CLUSTER_PROFILE_NAME" --desired-status RUNNING | grep "$ENVIRONMENT" | grep "/web" | awk '{print $NF}')
  if [ "$status" = "HEALTHY" ]; then
    echo "cluster returned $status status for booking service."
    echo "deployment for version $VERSION is successful."
    exit 0
  else
    echo "booking service status : $status"
    echo "This could be due to either the deployment for $VERSION version is not successful or there is is a poblem with the cluster. Try accessing the http://host:8080/version to verify further"
    return 1
  fi
}

for i in {1..4};
do
  echo "healthcheck retry count $i"
  check_status && break || sleep 30;
done

exit 1
