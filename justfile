# Project task runner. Run `just` or `just --list` to see available tasks.

default:
  @just --list

# Run all unit tests.
test:
  go test ./...

# Run static analysis and configured linters.
lint:
  golangci-lint run

# Vendor Go dependencies for offline Homebrew Formula builds.
vendor:
  go mod vendor

# Format all Go source files.
fmt:
  go fmt ./...

# Build the dibo binary into the ignored dist directory.
build:
  mkdir -p dist
  go build -o dist/dibo .

# Run formatting, linting, tests, and a build.
check: fmt lint test build

# Run dibo with optional arguments, e.g. `just run detect`.
run *args:
  go run . {{args}}

# Install dibo into the current Go bin directory.
install:
  go install .
