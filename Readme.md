### Summary

- This app serves as an asynchronous request sender.
- The project follows the principles of Domain-Driven Design (DDD) in its structure.
- A message broker is utilized to improve scalability and prevent data loss.
- Three entities are stored in the database.
- The Server service manages two HTTP endpoints for registering requests and exposing results.
- The Tasker service produces and consumes tasks (requests), sending them to third-party servers and storing the results in the database.
- The API development approach prioritizes documentation, employing a document-first methodology, with code generation facilitated by oapi-codegen.

![Request Tasker Flow](./images/request_tasker_flow.drawio.svg)

![Request Tasker Database Schema](./images/request_tasker_DB.drawio.svg)

### Dependencies

- [oapi-codegen](https://github.com/deepmap/oapi-codegen)     
- [Taskfile](https://taskfile.dev/)                           
- [mockery](https://github.com/vektra/mockery)                
- [golang-migrate](https://github.com/golang-migrate/migrate) 

```bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

```bash
go get -u github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen
```

```bash
go install github.com/vektra/mockery/v2@v2.38.0
```

```bash
brew install go-task
```

```bash
brew install mockery
```

```bash
brew install golang-migrate
```

### Taskfile

Instead of a Makefile, this project uses a Taskfile. Check [Taskfile.yml](./Taskfile.yml) for a list of commands and descriptions.

### Run

```bash
task dep-up
```

```bash
task run
```

### API Documentation

Chcek this [api.yml](./openapi/api.yml) openapi file.

### Configs

- All configurations are in the [.env](./.env) file.
- The Docker container overwrites a few configurations in [./dep/.env](./dep/.env).

### Testing

This project includes both mocked and integration tests implemented with Mockery and [testcontainer](https://testcontainers.com/). Run all tests with the command: `task test`.

To run specific tests:
- `task test-usecase`
- `task test-api`
- `task test-tasker`
- `task test-db`
- `task test-kafka`

In the case of any changes in interfaces, it's requried to run the `task generate-mock` command.

### Additional Considerations for the Future

- [ ] Currently, there is no retry mechanism for failed requests. If a dead-letter queue is added to this project, failed requests can be retried.
- [ ] Optimize test environment by leveraging local infrastructure instead of creating and deleting containers for each test.
- [ ] Add health_check.
- [ ] Use Transaction for sql queries.
- [ ] Separate api and worker to different build binary files or use the component separation in the same binary way.