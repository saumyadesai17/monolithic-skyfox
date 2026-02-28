
module.exports = {
  "name": "task-definition",
  "config" : {
    "requiresCompatibilities": [
      "EC2"
    ],
    "executionRoleArn": "catalystSecretsExecutionRoleCI",
    "containerDefinitions": [
      {
        "name": "movie-service",
        "image": process.env.IMAGE,
        "repositoryCredentials": {
          "credentialsParameter": process.env.REGISTRY_SECRET_ARN
        },
        "memory": 50,
        "cpu": 256,
        "essential": true,
        "portMappings": [
          {
            "containerPort": 4567,
            "hostPort": parseInt(process.env.PORT),
            "protocol": "tcp"
          }
        ],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-create-group": "true",
            "awslogs-group": "catalyst-log-group",
            "awslogs-region": "ap-south-1",
            "awslogs-stream-prefix": "movie-service-"+process.env.CI_ENVIRONMENT_SLUG
          }
        }
      }
    ],
    "volumes": [],
    "networkMode": "bridge",
    "placementConstraints": [],
    "family": "movie-service-"+process.env.CI_ENVIRONMENT_SLUG
  }
}

