### Request flow

TODO: draw.io
Call create task api --> produce to Kafka --> consume from kafka --> send request --> save result is db

| Dependencies                                                |
| ----------------------------------------------------------- |
| [oapi-codegen](https://github.com/deepmap/oapi-codegen)     |
| [taskfile](https://taskfile.dev/)                           |
| [mockery](https://github.com/vektra/mockery)                |
| [golang-migrate](https://github.com/golang-migrate/migrate) |

### Install Dependencies

``` bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

``` bash
go get -u github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen
```

``` bash
brew install go-task
```

``` bash
go install github.com/vektra/mockery/v2@v2.38.0
```

``` bash
brew install mockery
```

``` bash
brew install golang-migrate
```

### Run

``` bash
task dep-up
```

``` bash
task run
```

### future
dead letter queue for missed requests