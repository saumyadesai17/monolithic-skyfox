const API_HOST_LOCAL = "http://localhost:8080";
const API_HOST_CI = "https://falcon-api-dev.neev.idfcfirstbank.com";
const API_HOST_STAGING = "http://ec2-13-232-102-189.ap-south-1.compute.amazonaws.com:8080";
const API_HOST_PROD = "http://ec2-13-232-102-189.ap-south-1.compute.amazonaws.com:8080";

const UI_HOST_LOCAL = "https://localhost:3000";
const UI_HOST_CI = "https://localhost:3000";
const UI_HOST_STAGING = "https://localhost:3000";
const UI_HOST_PROD = "https://localhost:3000";

const ENV_LOCAL = "local";

const HOSTS = {
    local: {
        "API": API_HOST_LOCAL,
        "UI": UI_HOST_LOCAL
    },
    dev: {
        "API": API_HOST_CI,
        "UI": UI_HOST_CI
    },
    qa: {
        "API": API_HOST_STAGING,
        "UI": UI_HOST_STAGING
    },
    prod: {
        "API": API_HOST_PROD,
        "UI": UI_HOST_PROD
    }
};

export const serviceUrl = () => {
    if(import.meta.env.VITE_ON_EC2 === "true") {
        console.log("Using public hostname and port");
        console.log(import.meta.env.VITE_PUBLIC_HOSTNAME_AND_PORT);
        return import.meta.env.VITE_PUBLIC_HOSTNAME_AND_PORT || ENV_LOCAL;
    }

    console.log("Environment is ");
    console.log(import.meta.env.VITE_ENVIRONMENT);

    const environment = import.meta.env.VITE_ENVIRONMENT || ENV_LOCAL;
    return HOSTS[environment].API
};

export const urls = {
    service: serviceUrl()
};

export const featureToggles = {
    dummy: true
};
