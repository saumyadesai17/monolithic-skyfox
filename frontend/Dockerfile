FROM node:lts-bullseye as build

WORKDIR /app

COPY package.json ./
RUN yarn install

COPY . .

ARG ENVIRONMENT
ARG PUBLIC_HOSTNAME_AND_PORT
ARG REACT_APP_ON_EC2

ENV REACT_APP_ENVIRONMENT=$ENVIRONMENT
ENV REACT_APP_PUBLIC_HOSTNAME_AND_PORT=$PUBLIC_HOSTNAME_AND_PORT
ENV REACT_APP_ON_EC2=$REACT_APP_ON_EC2

RUN yarn build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]