# Contributing to go-auth

Thank you for your interest in contributing to go-auth! We welcome contributions from the community and are grateful for your support.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of Go and authentication concepts

### Development Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/go-auth.git
   cd go-auth
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run tests to ensure everything works**
   ```bash
   go test ./...
   ```

## 🛠️ Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, well-documented code
   - Follow Go conventions and best practices
   - Add tests for new functionality

3. **Test your changes**
   ```bash
   # Run all tests
   go test ./...
   
   # Run with race detection
   go test ./... -race
   
   # Check coverage
   go test ./... -cover
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

Examples:
```
feat: add password reset functionality
fix: resolve token validation issue
docs: update README with new examples
test: add integration tests for middleware
```

## 🧪 Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./pkg/auth

# With verbose output
go test ./... -v

# With coverage
go test ./... -cover

# Race condition detection
go test ./... -race

# Benchmarks
go test ./... -bench=.
```

### Writing Tests

- Write tests for all new functionality
- Maintain or improve test coverage
- Include both unit and integration tests
- Test error conditions and edge cases

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    t.Run("successful_case", func(t *testing.T) {
        // Test successful operation
    })
    
    t.Run("error_case", func(t *testing.T) {
        // Test error handling
    })
}
```

## 📝 Code Style

### Go Standards

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Write clear, concise comments
- Keep functions focused and small
- Handle errors appropriately

### Documentation

- Document all public functions and types
- Include examples in documentation
- Update README.md for significant changes


## 🐛 Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Go version** (`go version`)
2. **Operating system** and version
3. **go-auth version**
4. **Minimal reproduction case**
5. **Expected vs actual behavior**
6. **Error messages** (if any)

### Feature Requests

For feature requests, please provide:

1. **Use case** - Why is this feature needed?
2. **Proposed solution** - How should it work?
3. **Alternatives considered** - What other approaches did you consider?
4. **Additional context** - Any other relevant information

## 🔒 Security

### Reporting Security Issues

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please email us directly at: pragneshbagary1699@gmail.com

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond as quickly as possible and work with you to resolve the issue.

## 📋 Pull Request Process

### Before Submitting

1. **Ensure tests pass**
   ```bash
   go test ./...
   go test ./... -race
   ```

2. **Check code formatting**
   ```bash
   go fmt ./...
   ```

3. **Update documentation** if needed

4. **Add examples** for new features

### Submitting

1. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request**
   - Use a clear, descriptive title
   - Reference any related issues
   - Describe what your PR does
   - Include testing information

3. **Respond to feedback**
   - Address review comments promptly
   - Make requested changes
   - Keep the conversation constructive

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] Updated documentation

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
```

## 🎯 Areas for Contribution

We welcome contributions in these areas:

### High Priority
- **Bug fixes** - Help us maintain stability
- **Performance improvements** - Optimize critical paths
- **Security enhancements** - Strengthen security features
- **Documentation** - Improve guides and examples

### Medium Priority
- **New storage backends** - PostgreSQL, SQLite, Redis, etc.
- **Testing improvements** - Better test coverage

### Low Priority
- **Code cleanup** - Refactoring and optimization
- **Developer tools** - Improve development experience
- **Examples** - More real-world examples

## 🏆 Recognition

Contributors will be:
- Listed in the project README
- Mentioned in release notes
- Given credit in commit messages
- Invited to join the maintainer team (for significant contributions)

## 📞 Getting Help

If you need help with contributing:

- 💬 [GitHub Discussions](https://github.com/pragneshbagary/go-auth/discussions)
- 📧 Email: pragneshbagary1699@gmail.com
- 🐛 [Issues](https://github.com/pragneshbagary/go-auth/issues) for bugs and feature requests

## 📄 License

By contributing to go-auth, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

**Thank you for contributing to go-auth!** 🙏

Your contributions help make authentication in Go better for everyone.