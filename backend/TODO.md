## The following functionalities to be improved/ implemented

### Docker
- currently we have renamed debian:buster-slim to buster-slim in the ecr. the name needs to change to original `debian:buster-slim` and the same needs to be updated in both the docker files - `Dockerfile` and `Dockerfile-migration`
- move the docker related files into a directory

### CI
- move the integration and deployment related files and scripts to a directory so that the root path looks cleaner

### Documentation
- review Readme.md to ensure all steps mentioned are correct and one can set up a working app with all tests running successfully
- verify the swagger documentation

### Tests
- the integration tests use `TestMain` in setup_test.go as the entry point. while this is a prescribed approach and works just fine, it is good to look for alternate options and choose the best one
- mocked interfaces are generated through mockery. these have not been updated after small changes in the application code. while the tests are currently passing, it will be good idea to regenerate those mocks and verify once again

### App
- it is a good idea to move the main.go and migration.go to inside cmd/server and cmd/migration respectively. That will give an interpretation that the app has two entry point commands - server and migration
- bookings/database - we should check whether base_db.go and connection.go can be combined into same package
- bookings/database - base_db.go has a method `GormDB()` that is specifically exposed for tests. Review how this can be avoided by making necessary changes in the testDB.go.
- use zap logger instance for the DB in connection.go
- connection.go uses `panic` to signal when database connection could not be established. get rid of panic by adding another return type error and handling the error in the caller. main should throw `Os.Exit(1)` instead of panic