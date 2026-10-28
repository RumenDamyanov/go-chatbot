# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

Build with **Go 1.26.10+** or **1.27.3+**. Go 1.25 and older are end of life and do not receive the October 2026 security releases. Those fixes are in the toolchain (`crypto/x509`, `net/http`, `html/template`), not in this module. The `go` line asks the `go` command for at least Go 1.26.10 when `GOTOOLCHAIN=auto`, which is the default. `GOTOOLCHAIN=local` fails the build until that toolchain is installed.

- [CVE-2026-111943](https://pkg.go.dev/vuln/GO-2026-6809) (`crypto/x509`): untrustworthy eMudhra intermediates are rejected. Fixed in Go 1.26.10 and 1.27.3 ([announcement](https://groups.google.com/g/golang-announce/c/66pojueaQek/m/tmWX2Pq9AgAJ)).
- Go 1.26.9 and 1.27.2 also fixed HTTP/2 server crashes, trailer memory use, and `html/template` context tracking ([announcement](https://groups.google.com/g/golang-announce/c/U2fTuyDJznI)).
- Gin and Fiber decompress Brotli with `github.com/andybalholm/brotli`. Versions before v1.2.7 panic on crafted input ([CVE-2026-111497](https://www.wiz.io/blog/cve-2026-111497-brotli-rce-node-js)). This module requires v1.2.7.

## Reporting a Vulnerability

The go-chatbot team and community take all security bugs seriously. Thank you for improving the security of go-chatbot. We appreciate your efforts and responsible disclosure and will make every effort to acknowledge your contributions.

### How to Report a Security Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to **<contact@rumenx.com>**.

If the security vulnerability is accepted, we will acknowledge receipt of your vulnerability report, work with you to understand the scope and severity, and provide you with an estimated timeline for a fix.

### What to Include in Your Report

To help us understand the nature and scope of the possible issue, please include as much of the following information as possible:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

### Our Commitment

- We will respond to your report within 72 hours with our evaluation of the report and an expected resolution date.
- If you have followed the instructions above, we will not take any legal action against you in regard to the report.
- We will handle your report with strict confidentiality, and not pass on your personal details to third parties without your permission.
- We will keep you informed of the progress towards resolving the problem.
- In the public information concerning the problem reported, we will give your name as the discoverer of the problem (unless you desire otherwise).

## Security Features

go-chatbot includes several built-in security features:

### Input Validation and Sanitization

- All user inputs are validated and sanitized before processing
- Message filtering middleware to prevent abuse and malicious content
- Rate limiting to prevent spam and DoS attacks
- Content filtering for profanity and harmful content

### API Security

- Secure API key handling with environment variable support
- Request/response validation
- Context-aware operations with timeout support
- Proper error handling without information leakage

### Configuration Security

- No hardcoded credentials or secrets
- Secure configuration loading from environment variables
- Validation of configuration parameters
- Safe defaults for all security-related settings

### Dependencies

- Regular dependency updates
- Security scanning of dependencies
- Minimal dependency footprint
- Use of well-maintained, trusted libraries

## Best Practices for Users

When using go-chatbot in your application, please follow these security best practices:

### Environment Variables

- Store all API keys and secrets in environment variables or secure secret management systems
- Never commit credentials to version control
- Use different API keys for different environments (development, staging, production)

### Rate Limiting

- Implement appropriate rate limiting for your use case
- Monitor for unusual usage patterns
- Set reasonable timeouts for API calls

### Input Validation

- Validate all user inputs before passing to the chatbot
- Implement additional content filtering if needed for your specific use case
- Log suspicious activities for monitoring

### Network Security

- Use HTTPS for all communications
- Implement proper CORS policies
- Consider using API gateways for additional security layers

### Monitoring and Logging

- Monitor chatbot usage for anomalies
- Log security-related events
- Set up alerts for suspicious activities
- Regularly review logs for security incidents

## Known Security Considerations

### AI Model Interactions

- AI responses may occasionally contain unexpected content
- Implement additional filtering for sensitive applications
- Monitor AI responses for quality and appropriateness
- Consider implementing human review for critical applications

### Third-Party AI Services

- API keys provide access to third-party AI services
- Monitor usage to prevent unexpected charges
- Implement usage limits and monitoring
- Review third-party service security policies

## Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed. We recommend:

- Keeping go-chatbot updated to the latest version
- Subscribing to GitHub notifications for security advisories
- Following our release notes for security-related changes

## Contact

For security-related questions or concerns, please contact:

**Email**: <contact@rumenx.com>

For general questions, please use our [GitHub Issues](https://go.rumenx.com/chatbot/issues) page.

---

Thank you for helping keep go-chatbot and our users safe!
