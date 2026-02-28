<!-- ## Install colima to run containers
- Set up colima for docker:
  - `brew install colima`
  - start colima by executing `colima start`
    - Colima has a dependency on Docker CLI hence you should not face any errors while running the next command
    - In the terminal type `docker ps`
    - Possibly a line indicating headers like "CONTAINER ID, IMAGE etc" will be displayed -->

# Set up Docker to run containers

- Install Docker Desktop from the IDFCFirst Bank Self Service portal.
- Run the Docker engine from inside Docker Desktop
- Install `docker` CLI:
  - `brew install docker`
- Log in to IDFC Repository from your Terminal:
  - `docker login https://artifactory.idfcfirstbank.com`
  - Enter your username eg. `jane.doe_int` if your IDFC ID is `jane.doe_int@idfcbank.com`
  - Enter your password
  - Login should be successful now
  - [Optional] Test a Docker pull: `docker pull artifactory.idfcfirstbank.com/neev-docker/golang:1.22`

## Install postgres
- Setup Postgres using docker
  - Start db command
    `docker run --name postgresdb -e POSTGRES_DB=bookingengine -e POSTGRES_USER=bookingengine -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d artifactory.idfcfirstbank.com/infra-common-docker/postgres`
  - Install below tools to connect to the database
   ```
      brew install --cask pgadmin4   
      brew install psqlodbc
   ```
  - Connect to db command
    `psql -h localhost -U bookingengine -d bookingengine`

<!-- ## Run movie service

`docker run --name movie-service -p 4567:4567 artifactory.idfcfirstbank.com/neev-docker/movie-service` -->

## Seed data
- set DB_PASSWORD, POSTGRES_USERNAME, POSTGRES_DB as environment variables before running the script
  - To see values set for environment variable use
    echo $POSTGRES_DB
  - Alternatively add the export commands to an env file and run ‘source env’
   ```
   export POSTGRES_DB=bookingengine
   export POSTGRES_USER=bookingengine
   export DB_PASSWORD=postgres
   ```
- params
  - starting date
  - number of weeks to seed data from given starting date

`sh scripts/seedShowData.sh 2026-02-20 3`

## Run the application
- Run the migrations locally
```shell script
make run-migrations
```
- Build the application and run the test along with coverage
```shell script
make build
```
- Run the server locally on localhost:8080
```shell script
make run
```
- To build & run the application
```shell script
make deploy
```

## Integration Testing
### Change docker host to run test containers
- set TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE,  DOCKER_HOST as environment variables
   ```
   export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock
   export DOCKER_HOST="unix://${HOME}/.colima/docker.sock"
   ```

### Mockery to create mocks
- install mockery:
  `brew install mockery`
- create mocks using following command:
  `mockery --dir=directory/forwhich/mocks/needtobe/created --output=path/where/mocks/arecreated --outpkg={package name} --all`

## Add api documentation
- Run these commands in terminal
  `export GOBIN=$(go env GOPATH)/bin`
  `export PATH="/Users/<username>/go/bin:$PATH"`
  `source ~/.zshrc `
- Install swag 
  `go get -u github.com/swaggo/swag/cmd/swag`
- Run 
  `swag init`
- For more info on annotations
  https://github.com/swaggo/swag
  https://swaggo.github.io/swaggo.io/declarative_comments_format/general_api_info.html
