# Contributing to NexusFlow

Thank you for your interest in contributing to NexusFlow! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/yourusername/nexusflow/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)
   - Relevant logs or screenshots

### Suggesting Features

1. Check [Discussions](https://github.com/yourusername/nexusflow/discussions) for existing feature requests
2. Create a new discussion with:
   - Clear use case
   - Proposed solution
   - Alternative approaches considered
   - Impact on existing functionality

### Pull Requests

1. **Fork the repository** and create a branch from `main`
2. **Follow coding conventions** (see below)
3. **Write tests** for new functionality
4. **Update documentation** as needed
5. **Ensure all tests pass** (`make test-all`)
6. **Run linters** (`make lint`)
7. **Submit a pull request** with a clear description

## Development Setup

See [Development Setup Guide](docs/development/setup.md) for detailed instructions.

## Coding Conventions

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write clear, self-documenting code
- Add comments for complex logic

### Protobuf Conventions

- Use `snake_case` for field names
- Include comprehensive comments
- Version all packages (e.g., `v1`, `v2`)
- Group related messages together

### Git Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:

```
feat(issue-service): add custom field support

Implement custom fields with 10+ field types including text,
number, date, select, multi-select, user, checkbox, URL, email,
and textarea.

Closes #123
```

### Testing

- Write unit tests for all business logic
- Write integration tests for API endpoints
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

### Documentation

- Update README.md for user-facing changes
- Update API documentation for API changes
- Add inline comments for complex code
- Update architecture docs for structural changes

## Project Structure

```
nexusflow/
├── services/           # Microservices
├── pkg/               # Shared libraries
├── proto/             # Protobuf definitions
├── deployments/       # Infrastructure as Code
├── docs/              # Documentation
└── scripts/           # Build and utility scripts
```

## Pull Request Process

1. **Create a feature branch**: `git checkout -b feature/amazing-feature`
2. **Make your changes** following conventions above
3. **Test thoroughly**: `make test-all && make lint`
4. **Commit with clear messages**: Follow conventional commits
5. **Push to your fork**: `git push origin feature/amazing-feature`
6. **Open a pull request** with:
   - Clear title and description
   - Link to related issues
   - Screenshots for UI changes
   - Test results

### PR Review Process

- Maintainers will review within 2-3 business days
- Address review comments promptly
- Keep PR scope focused and manageable
- Squash commits before merging (if requested)

## Community

- 💬 [GitHub Discussions](https://github.com/yourusername/nexusflow/discussions) - General questions and ideas
- 🐛 [Issue Tracker](https://github.com/yourusername/nexusflow/issues) - Bug reports and feature requests
- 📧 [Mailing List](mailto:dev@nexusflow.io) - Development discussions

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

## Recognition

Contributors will be recognized in:

- README.md contributors section
- Release notes
- Annual contributor highlights

Thank you for contributing to NexusFlow! 🎉
