# Repository Guidelines

Welcome to the **Despite** project. This guide outlines the structure, conventions, and workflows for contributing to this repository. Please follow these guidelines to ensure smooth collaboration.

## Project Structure & Module Organization

The repository is organized as a Go project with the following key directories and files:
- **Source Code**: Core logic resides in `.go` files at the root, such as `main.go`, `player.go`, and `dsmap.go`.
- **Tests**: Integration tests are in `integration_test.go` at the root level.
- **Assets**: Game or data files like `lev01.dsmap` are stored in the root directory.

## Build, Test, and Development Commands

Key commands for working on this project:
- `go build`: Compiles the project into an executable.
- `go test`: Runs all tests, including those in `integration_test.go`.
- `go run main.go`: Executes the application locally for development.

## Coding Style & Naming Conventions

We adhere to standard Go conventions:
- **Indentation**: Use tabs or 2 spaces (follow existing code style).
- **Formatting**: Run `go fmt` to ensure consistent style.
- **Naming**: Use camelCase for variables and functions, and PascalCase for types (e.g., `PlayerData`, `handleInput`).

## Testing Guidelines

- **Framework**: Tests use the built-in Go `testing` package.
- **Coverage**: Aim for meaningful coverage of core functionality.
- **Naming**: Test functions should start with `Test` and describe intent (e.g., `TestPlayerMovement`).
- **Running Tests**: Use `go test` to execute all tests.

## Commit & Pull Request Guidelines

- **Commit Messages**: Use descriptive, concise messages with a summary line (e.g., "Fix player movement bug in level 1").
- **Pull Requests**: Include a clear description, link to related issues, and mention any specific testing steps or screenshots if applicable.

## Agent-Specific Instructions

For AI agents or automated tools contributing to this repo:
- Ensure changes align with Go conventions using `go fmt`.
- Validate builds with `go test` before proposing changes.

Thank you for contributing to Despite!
