# Building a TO-DO List app ![Build status](https://github.com/rRupeshRanjan/my-todo-app/actions/workflows/go.yml/badge.svg)

In this project, we aim to build a set of REST APIs, where we can track our todo lists. As a user, we would have access to create,
update and delete tasks to list, as well as we should be able to see all the tasks in list. A task will have due
date, respective status and created date with it.

This REST API can further be integrated with a UI (coming up) for better visualizations.

#### Technology and libraries used
1. Go v1.15.2
2. [Fiber](https://github.com/gofiber/fiber/v2) v2.3.0 (for Http requests)
3. [Viper](https://github.com/spf13/viper) v1.7.1 (for config management)
4. [Zap](https://go.uber.org/zap) v1.16.0 (for logging)
5. [mysql](https://github.com/go-sql-driver/mysql) v1.5.0 (for sql driver)
6. [squirrel](https://github.com/Masterminds/squirrel) v1.5.0 (for sql query building)
7. [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) v1.5.0 (for sql tests)

#### Project Structure
- config
    - appConfig.go
- domain
    - task.go
    - constants.go
    - scenario.go
- services
    - taskService.go
    - taskRepositoryInterface.go
    - taskService_test.go
    - taskServiceBenchmark_test.go
- repository
    - taskRepository.go
    - taskRepository_test.go
    - taskRepositoryBenchmark_test.go
- main.go
- config.yaml

#### Testing mechanisms:
1. **Run all tests**: _go test ./..._
2. **Run all tests
