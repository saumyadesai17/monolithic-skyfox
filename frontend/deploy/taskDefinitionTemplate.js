module.exports = {
  "name": "task-definition",
  "config" : {
    "requiresCompatibilities": [
      "EC2"
    ],
    "executionRoleArn": "catalystSecretsExecutionRoleCI",
    "containerDefinitions": [
      {
        "name": "booking-web",
        "image": process.env.IMAGE,
        "repositoryCredentials": {
          "credentialsParameter": process.env.REGISTRY_SECRET_ARN
        },
        "memory": 256,
        "cpu": 256,
        "essential": true,
        "portMappings": [
          {
            "containerPort": 80,
            "hostPort": parseInt(process.env.UI_PORT),
            "protocol": "tcp"
          }
        ],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-create-group": "true",
            "awslogs-group": "catalyst-log-group",
            "awslogs-region": "ap-south-1",
            "awslogs-stream-prefix": "booking-web"+process.env.TEAM_ID+"-"+process.env.CI_ENVIRONMENT_SLUG
          }
        }
      }
    ],
    "volumes": [],
    "networkMode": "bridge",
    "placementConstraints": [],
    "family": "booking-web"+process.env.TEAM_ID+"-"+process.env.CI_ENVIRONMENT_SLUG
  }
}