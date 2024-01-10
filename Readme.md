### Summery

- This is an async request sender.
- The project adheres to the Domain-Driven Design (DDD) principles in its structure.
- A message broker is employed to enhance scalability and prevent data loss.
- There are 3 entities that are stored to the database
- Server service handles two HTTP endpoints to register the request and expose the result
- Tasker service produce and consume the task(request) and sends the request to the third-party servers and store the result in db

![Flow](./images/request_tasker_flow.drawio.svg)

![Schema](./images/request_tasker_DB.drawio.svg)

### Dependencies

- [oapi-codegen](https://github.com/deepmap/oapi-codegen)     
- [taskfile](https://taskfile.dev/)                           
- [mockery](https://github.com/vektra/mockery)                
- [golang-migrate](https://github.com/golang-migrate/migrate) 

``` bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

``` bash
go get -u github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen
```

``` bash
go install github.com/vektra/mockery/v2@v2.38.0
```

``` bash
brew install go-task
```

``` bash
brew install mockery
```

``` bash
brew install golang-migrate
```

### Taskfile

Instead of MakeFile there is Taskfile in this project. Check this [Taskfile.yml](./Taskfile.yml) for all of the commands and descriptions.


### Run

``` bash
task dep-up
```

``` bash
task run
```

### Configs

- All of the configs are in [.env](./.env) file.
- Docker conainer overwrites a few configs in [./dep/.env](./dep/.env)

### Testing

There are both mocked and integration tests in this project, implemented by mockery and [testcontainer](https://testcontainers.com/).
Run all tests by this commad: `task test`
- Run only a single test:
  - `task test-usecase`
  - `task test-api`
  - `task test-tasker`
  - `task test-db`
  - `task test-kafka`


### Future
Right now there isn't any retry mechanisem for failed requests so if a dead letter queue get added in this project then the failed requests can be retried.