# Contributing to GoPixoo

Thanks for your interest in contributing! Here's how to get started.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<your-username>/GoPixoo.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit and push
7. Open a pull request

## Development Setup

- Go 1.24+
- A Pixoo64 device on your local network (for integration testing)

```bash
# Build
go build -o gopixoo .

# Run tests
go test ./...

# Vet
go vet ./...
```

## Guidelines

- Keep PRs focused on a single change
- Add tests for new functionality
- Follow existing code style (`gofmt`)
- Update documentation if your change affects CLI usage

## Reporting Bugs

Open an issue with:
- GoPixoo version (`gopixoo --version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior

## Feature Requests

Open an issue describing the feature and why it would be useful. If it relates to a Pixoo64 API endpoint, include a link to the relevant documentation if available.
