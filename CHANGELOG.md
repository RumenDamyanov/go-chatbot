# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-10-28

### Security

- Require Go 1.26.10 or newer. Go 1.25 and older are end of life and do not include the October 2026 toolchain fixes.
- Require `github.com/andybalholm/brotli` v1.2.7. Earlier versions panic on crafted Brotli input used by Gin and Fiber.
- Require `golang.org/x/net` v0.60.0 and `github.com/quic-go/quic-go` v0.59.1. Older copies are pulled in by Gin and are reachable from HTTP handlers.

### Changed

- Default models are now OpenAI `gpt-5.4-mini`, Anthropic `claude-sonnet-5`, Google `gemini-3.8-flash`, xAI `grok-4.6`, Meta `llama-3.3-70b-instruct`, and Ollama `llama3.3`.
- Model names must match `^[A-Za-z0-9._-]+$` before they are sent to a provider. Gemini places the name in the request URL.

## [Unreleased]

### Added

- Initial project setup and scaffolding
- Core package architecture planning

## [1.0.0] - 2025-01-XX

### Added

- Initial release of go-chatbot package
- Framework-agnostic Go chat implementation
- Gin adapter with middleware support
- Echo adapter with middleware support
- Fiber adapter with middleware support
- Chi adapter with middleware support
- Vanilla net/http support
- Support for multiple AI providers:
  - OpenAI (GPT-4, GPT-3.5, etc.)
  - Anthropic (Claude models)
  - Google Gemini
  - Meta Llama models
  - xAI Grok models
  - Ollama (local models)
  - Default fallback model
- Model factory for easy provider switching
- Configurable chat prompts and behavior
- Frontend components (HTML/CSS/JS)
- Framework-specific components (Vue, React, Angular)
- TypeScript examples
- SCSS support for styling
- Comprehensive test suite with Go testing and testify
- Static analysis with golangci-lint
- Go coding standards (gofmt, goimports, go vet)
- High test coverage (80%+)
- Documentation and examples
- CI/CD workflows (GitHub Actions)
- Code quality tools integration
- Context-aware operations with cancellation support
- Structured logging with slog
- Message filtering middleware
- Rate limiting capabilities
- Security safeguards against abuse

### Features

- Plug-and-play chat API and UI components
- AI model abstraction layer
- Customizable prompts and configuration
- Multi-language and emoji support
- Security safeguards and rate limiting
- Framework adapters for easy integration
- Extensible architecture
- Comprehensive documentation
- Context-aware operations
- Graceful error handling

### Requirements

- Go 1.21 or higher
- Go modules for dependency management
- Optional: Popular Go web frameworks (Gin, Echo, Fiber, Chi)

---

*For more details about releases, see [GitHub Releases](https://go.rumenx.com/chatbot/releases).*
