# CI/CD Updates Prompt Template

## Context
# CI/CD Updates Prompt Template

## Prompt
```
I'm working on the Karo (Kubernetes Alert Reaction Operator) CI/CD pipeline using GitHub Actions. The current pipeline includes building, testing, security scanning, and multi-architecture container image builds.

## Current Pipeline Overview
**Existing workflows:**
- `.github/workflows/ci-cd.yml` - Main CI/CD pipeline
- `.github/workflows/security.yml` - Security scanning (if separate)
- `.github/workflows/release.yml` - Release automation (if exists)

**Current pipeline stages:**
1. **Code Quality** - linting, formatting, vet
2. **Testing** - unit tests, integration tests
3. **Security** - Trivy vulnerability scanning, gosec
4. **Build** - multi-arch container images (amd64, arm64)
5. **Deploy** - push to ghcr.io registry

## Proposed Changes
**What I want to add/modify:**
{detailed_description_of_pipeline_changes}

**New stages/jobs needed:**
- [ ] {new_job_1} - {purpose_and_description}
- [ ] {new_job_2} - {purpose_and_description}
- [ ] {new_job_3} - {purpose_and_description}

**Existing stages to modify:**
- [ ] {existing_stage} - {what_changes_are_needed}

**New triggers/conditions:**
- [ ] {trigger_condition} - {when_it_should_run}

## Technical Requirements
**GitHub Actions features needed:**
- [ ] Matrix builds (multiple versions/architectures)
- [ ] Conditional job execution
- [ ] Artifact sharing between jobs
- [ ] Secrets management
- [ ] Environment protection rules
- [ ] Manual approval gates
- [ ] Scheduled runs (cron)
- [ ] External service integration

**Tools/Actions to integrate:**
- [ ] {tool_name} - {purpose_and_integration_details}
- [ ] {action_name} - {what_it_does}

**Performance considerations:**
- [ ] Parallel job execution
- [ ] Caching strategies (Go modules, Docker layers)
- [ ] Resource usage optimization
- [ ] Build time reduction

## Security Requirements
**Security scanning:**
- [ ] Container image vulnerability scanning
- [ ] Code security analysis (gosec, CodeQL)
- [ ] Dependency vulnerability checking
- [ ] License compliance checking
- [ ] Supply chain security (SLSA, provenance)

**Secrets and permissions:**
- [ ] Registry authentication
- [ ] Signing keys management
- [ ] Minimal required permissions
- [ ] Environment-specific secrets

## Integration Points
**External services:**
- [ ] Container registry (ghcr.io, Docker Hub)
- [ ] Security scanning services
- [ ] Kubernetes cluster (for deployment)
- [ ] Monitoring/alerting systems
- [ ] Issue tracking integration

**Notifications:**
- [ ] Slack/Teams notifications
- [ ] Email alerts
- [ ] GitHub status checks
- [ ] PR comments with results

## Deployment Strategy
**Environments:**
- [ ] Development (automatic deployment)
- [ ] Staging (automatic with tests)
- [ ] Production (manual approval)

**Deployment methods:**
- [ ] Helm chart deployment
- [ ] kubectl apply
- [ ] GitOps (ArgoCD/Flux)
- [ ] Container registry only

## Quality Gates
**Required checks before deployment:**
- [ ] All tests pass (unit, integration, e2e)
- [ ] Security scans pass
- [ ] Code quality metrics meet thresholds
- [ ] Manual approval (for production)
- [ ] Documentation is up to date

**Failure handling:**
- [ ] Automatic rollback on deployment failure
- [ ] Notification on pipeline failures
- [ ] Retry logic for transient failures
- [ ] Emergency bypass procedures

## Current Pipeline Content
**Existing workflow file:**
```yaml
{paste_current_workflow_yaml_or_relevant_sections}
```

**Areas that work well:**
{describe_what_currently_works_well}

**Pain points:**
{describe_current_issues_or_limitations}

## Request for Copilot
Please help me implement these CI/CD improvements by:

1. **Analyzing the current pipeline** and identifying optimization opportunities
2. **Creating/updating GitHub Actions workflows** with the new requirements
3. **Implementing security best practices** for CI/CD pipelines
4. **Adding proper error handling** and retry logic
5. **Optimizing build performance** with caching and parallelization
6. **Setting up proper testing stages** with comprehensive coverage
7. **Implementing deployment automation** with appropriate safeguards

**Specific questions:**
1. What's the best way to structure multi-environment deployments?
2. How should I implement approval gates for production deployments?
3. What caching strategies would improve build performance?
4. How can I make the pipeline more secure and compliant?
5. What monitoring should I add to track pipeline health?

## Success Criteria
**Pipeline should achieve:**
- [ ] Fast feedback (< 10 minutes for basic checks)
- [ ] High reliability (> 95% success rate)
- [ ] Security compliance (all scans pass)
- [ ] Easy maintenance and updates
- [ ] Clear visibility into pipeline status
- [ ] Automated recovery from common failures

---

## Usage Instructions
1. Replace `{placeholders}` with your specific CI/CD requirements
2. Include current pipeline configuration for context
3. Describe the integration points and external dependencies
4. Specify security and quality requirements
5. Copy and paste into Copilot chat

## Example Usage
```
**What I want to add/modify:**
Add automated end-to-end testing against a real Kubernetes cluster and implement progressive deployment with canary releases.

**New stages/jobs needed:**
- e2e-testing - Deploy operator to kind cluster and run end-to-end tests with real AlertManager
- canary-deployment - Deploy to 10% of production traffic and monitor for 30 minutes
- full-deployment - Complete rollout if canary succeeds

**Tools/Actions to integrate:**
- kind-action - Create Kubernetes cluster for testing
- helm/chart-testing-action - Test Helm charts
- prometheus/alertmanager - Set up real alert testing environment
```