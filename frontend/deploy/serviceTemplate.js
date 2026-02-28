module.exports = {
  "name": "service",
  "config" : {
    "cluster": "team" + process.env.TEAM_ID,
    "taskDefinition": "booking-web"+process.env.TEAM_ID+"-"+process.env.CI_ENVIRONMENT_SLUG,
    "serviceName": "booking-web-"+process.env.CI_ENVIRONMENT_SLUG,
    "desiredCount": 1,
    "clientToken": process.env.CLIENT_TOKEN,
    "launchType": "EC2",
    "deploymentConfiguration": {
      "minimumHealthyPercent": 0
    }
  }
}