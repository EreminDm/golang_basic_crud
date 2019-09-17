# Changelog

All notable changes to this project will be documented in this file.

## [0.0.2] - 2019-09-17

### Added

- GRPC layer, which represent ways for communication with service
- Nets layer, which includes implementation beetween transport layer and other layers

### Changed

- HTTP layer migrated to nets

## [0.0.1] - 2019-09-13

### Added

- Mongo database methods for CRUD operations
- Controller layer, which represents communication between HTTP and Database layers
- HTTP layer, which include mux handler for CRUD operations
- Entity module, which describes objects struct
- Tests for database methods with 96.5% covarage
- Tests for controller methods with 100% covarage
- Tests for HTTP methods with 75% covarage
- Swagger documentation
- Makefile
- Automation builds using Travis, including tests and linters runing
- Golangci linters with extended linters, which are described in the .golangci.yaml file
