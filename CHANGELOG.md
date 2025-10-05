# Changelog

All notable changes to the Alert Reaction Operator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Prepare for next release

## [0.1.4] - 2025-10-05

### Added
- AlertReaction Custom Resource Definition (CRD)
- AlertReaction controller with reconciliation logic
- Webhook server for AlertManager integration
- Job creation and management based on alert triggers
- Environment variable substitution in job templates
- Comprehensive test suite with unit and integration tests
- GitHub Actions CI/CD pipeline
- Complete Helm chart for deployment
- Support for multiple deployment environments (dev, prod)
- Monitoring and observability features
- Security policies and RBAC configuration
- Documentation and contribution guidelines
- Draft-first release process with comprehensive release automation
- CHANGELOG automation and validation tooling
- Enhanced Helm chart release pipeline with CRD bundling
- Release preparation and management scripts
- Comprehensive release documentation and automation guides

### Changed
- Restructured CI/CD pipeline to separate concerns (testing vs releasing)
- Moved release process to draft-first approach for better control
- Enhanced Helm chart automation with multi-format publishing
- Improved release workflow with CHANGELOG integration and validation
- Updated chart documentation to reflect bundled CRDs

### Improved
- Release process now requires manual approval before publishing
- Helm charts automatically bundle latest CRDs from config/crd/
- Enhanced release notes generation with CHANGELOG integration
- Better separation between CI/CD and release pipelines
- Streamlined user experience with automatic CRD installation

### Features
- **Alert Matching**: Flexible alert selection using labels and annotations
- **Job Execution**: Automatic Kubernetes Job creation in response to alerts
- **Template Support**: Environment variable substitution in job specifications
- **Cooldown Period**: Configurable delays to prevent rapid re-execution
- **Status Tracking**: Comprehensive status reporting for AlertReaction resources
- **Webhook Integration**: HTTP endpoint for AlertManager webhook notifications
- **Health Checks**: Readiness and liveness probes for operator reliability

### Infrastructure
- **Kubernetes Support**: Compatible with Kubernetes 1.19+
- **Go Runtime**: Built with Go 1.24 and controller-runtime v0.21.0
- **Container Images**: Multi-architecture container images (amd64/arm64)
- **Helm Deployment**: Production-ready Helm chart with configurable values
- **CI/CD Pipeline**: Automated testing, building, and security scanning
- **Monitoring**: Prometheus metrics and ServiceMonitor support
- **Release Automation**: Draft-first release workflow with manual approval gates
- **CHANGELOG Management**: Automated CHANGELOG processing and validation
- **Helm Publishing**: Multi-format chart publishing (OCI registry, GitHub Pages, release assets)
- **CRD Bundling**: Automatic inclusion of latest CRDs in Helm charts during release
- **Documentation**: Enhanced release process documentation and user guides

### Security
- **RBAC**: Minimal required permissions with ClusterRole/ClusterRoleBinding
- **Security Context**: Non-root user, read-only filesystem, dropped capabilities
- **Network Policies**: Optional traffic restriction policies
- **Image Scanning**: Automated vulnerability scanning in CI/CD
- **Secret Management**: Secure handling of configuration data

### Documentation
- **User Guide**: Comprehensive README with examples and troubleshooting
- **API Reference**: Detailed CRD specification and field documentation
- **Helm Chart Guide**: Complete Helm chart documentation
- **Contributing Guide**: Development setup and contribution guidelines
- **Operations Guide**: Monitoring, troubleshooting, and maintenance procedures

### Removed
- Duplicate CRD files from charts directory (now bundled during release)
- Manual CRD installation requirement for users
- Release job from CI/CD pipeline (moved to separate workflow)

[Unreleased]: https://github.com/dudizimber/k8s-alert-reaction-operator/compare/v0.1.4...HEAD
[0.1.4]: https://github.com/dudizimber/k8s-alert-reaction-operator/releases/tag/v0.1.4
