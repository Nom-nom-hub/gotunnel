# 🤝 Contributing to GoTunnel

Thank you for your interest in contributing to GoTunnel! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork the repository**
2. **Clone your fork:**
   ```bash
   git clone https://github.com/your-username/gotunnel.git
   cd gotunnel
   ```

3. **Set up development environment:**
   ```bash
   # Install Go dependencies
   go mod download
   
   # Install dashboard dependencies
   cd dashboard
   npm install
   cd ..
   ```

4. **Make your changes**
5. **Test your changes**
6. **Submit a pull request**

## 🛠️ Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Building from Source
```bash
# Build Go binaries
make build

# Build dashboard
cd dashboard
npm run build
cd ..
```

### Running Tests
```bash
# Go tests
go test ./...

# Dashboard tests
cd dashboard
npm test
cd ..
```

## 📝 Code Style

### Go Code
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused

### JavaScript/TypeScript
- Use ESLint and Prettier
- Follow Airbnb style guide
- Add JSDoc comments for functions
- Use meaningful variable names

### Git Commits
- Use conventional commit format
- Keep commits atomic and focused
- Write clear commit messages

## 🧪 Testing

### Go Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/tunnel
```

### Dashboard Tests
```bash
cd dashboard
npm test
npm run test:e2e
```

### Integration Tests
```bash
# Test tunnel functionality
make test-tunnel

# Test GUI functionality
make test-dashboard
```

## 📋 Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write clean, well-documented code
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes:**
   ```bash
   make test
   make build
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a pull request:**
   - Use the PR template
   - Describe your changes clearly
   - Link any related issues

## 🎯 Areas for Contribution

### High Priority
- **Bug fixes** - Fix reported issues
- **Documentation** - Improve docs and examples
- **Tests** - Add test coverage
- **Performance** - Optimize tunnel performance

### Medium Priority
- **New features** - Add requested functionality
- **UI improvements** - Enhance dashboard interface
- **Security** - Security audits and improvements
- **Platform support** - Add new OS support

### Low Priority
- **Examples** - Add usage examples
- **Tools** - Development and deployment tools
- **Monitoring** - Add monitoring and metrics
- **Integrations** - Third-party integrations

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment details:**
   - OS and version
   - Go version
   - Node.js version

2. **Steps to reproduce:**
   - Clear, step-by-step instructions
   - Expected vs actual behavior

3. **Logs and error messages:**
   - Full error output
   - Relevant log files

4. **Additional context:**
   - Screenshots if applicable
   - Related issues or discussions

## 💡 Feature Requests

When requesting features, please:

1. **Describe the problem** - What issue does this solve?
2. **Propose a solution** - How should it work?
3. **Consider alternatives** - Are there other approaches?
4. **Provide examples** - Show how it would be used

## 📚 Documentation

### Code Documentation
- Add comments for exported functions
- Include examples in comments
- Update README for new features

### User Documentation
- Update README.md for new features
- Add examples to examples/ directory
- Update website documentation

## 🔒 Security

### Reporting Security Issues
- **DO NOT** create public issues for security problems
- Email security@gotunnel.dev instead
- Include detailed reproduction steps
- Allow time for response before disclosure

### Security Guidelines
- Never commit secrets or credentials
- Use environment variables for sensitive data
- Follow security best practices
- Validate all user inputs

## 🏷️ Release Process

### Versioning
- Follow [Semantic Versioning](https://semver.org/)
- Use conventional commit messages
- Update CHANGELOG.md

### Release Steps
1. Update version numbers
2. Update CHANGELOG.md
3. Create release tag
4. Build and upload binaries
5. Update documentation

## 🤝 Community Guidelines

### Be Respectful
- Treat everyone with respect
- Be constructive in feedback
- Help newcomers learn

### Be Inclusive
- Welcome contributors of all backgrounds
- Use inclusive language
- Consider accessibility

### Be Professional
- Keep discussions on-topic
- Be patient with questions
- Focus on the code, not the person

## 📞 Getting Help

- **Discord:** [Join our community](https://discord.gg/gotunnel)
- **Email:** support@gotunnel.dev
- **Issues:** [GitHub Issues](https://github.com/gotunnel/gotunnel/issues)
- **Discussions:** [GitHub Discussions](https://github.com/gotunnel/gotunnel/discussions)

## 🙏 Acknowledgments

Thank you for contributing to GoTunnel! Your contributions help make tunneling accessible to developers worldwide.

---

**Made with ❤️ by the GoTunnel community** 