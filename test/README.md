# Code Controller Tests

This directory contains tests for the Code Controller (CC) project. The test suite is organized as follows:

## Test Structure

- `test/testutil`: Test utilities and helpers
- `test/integration`: Integration tests that verify multiple components working together
- `test/e2e`: End-to-end tests that test the entire system

Component tests are located alongside the code they test:
- `internal/*/package_test.go`: Tests for internal packages
- `pkg/*/package_test.go`: Tests for public packages
- `cmd/cc/main_test.go`: Tests for CLI commands

## Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run end-to-end tests only
make test-e2e

# Run tests with coverage
make test-cover
```

## Test Philosophy

The Code Controller project follows a test-driven development (TDD) approach:

1. Write tests first
2. Implement functionality to satisfy the tests
3. Refactor while keeping tests passing

This ensures:
- Features are well-defined before implementation
- Edge cases are considered early
- Code is more maintainable and less prone to regressions

## Mock Generation

Mocks for interfaces are generated using `go.uber.org/mock/mockgen`:

```bash
# Generate all mocks
make generate-mocks
```

Generated mocks are stored in `*/mocks` directories alongside the interfaces they mock.

## Test Dependencies

- `github.com/stretchr/testify`: Assertion library
- `go.uber.org/mock`: Mock generation tool