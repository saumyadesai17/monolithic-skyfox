
module.exports = {
  "name": "service",
  "config" : {
    "cluster": "team-trainers",
    "taskDefinition": "movie-service-"+process.env.CI_ENVIRONMENT_SLUG,
    "serviceName": "movie-service-"+process.env.CI_ENVIRONMENT_SLUG,
    "desiredCount": 1,
    "clientToken": process.env.CLIENT_TOKEN,
    "launchType": "EC2",
    "deploymentConfiguration": {
      "minimumHealthyPercent": 0
    }
  }
}

