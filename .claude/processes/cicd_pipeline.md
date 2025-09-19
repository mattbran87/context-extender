# CI/CD Pipeline Process

## Overview
This document defines the Continuous Integration and Continuous Deployment pipeline for the context-extender Go project, integrated with the 4-phase cyclical development framework.

## Pipeline Architecture

### Pipeline Stages Overview
```
Code Commit → Build → Test → Security Scan → Quality Check → Package → Deploy → Monitor
     ↓          ↓       ↓           ↓              ↓           ↓         ↓         ↓
   Notify     Success  Report    Vulnerabilities  Metrics   Artifacts  Status  Alerts
```

## GitHub Actions Implementation

### Main CI/CD Workflow

```yaml
# .github/workflows/main.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop, 'feature/**', 'release/**']
  pull_request:
    branches: [main, develop]

env:
  GO_VERSION: '1.21'
  GOLANGCI_LINT_VERSION: 'v1.55'

jobs:
  # Stage 1: Build
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: ${{ runner.os }}-go-
      
      - name: Download dependencies
        run: go mod download
      
      - name: Build
        run: go build -v ./...
      
      - name: Upload build artifacts
        uses: actions/upload-artifact@v3
        with:
          name: build-artifacts
          path: |
            ./context-extender
            ./go.mod
            ./go.sum

  # Stage 2: Test
  test:
    name: Test
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Run unit tests
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
          go tool cover -html=coverage.out -o coverage.html
      
      - name: Upload coverage reports
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
          flags: unittests
          name: codecov-umbrella
      
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage is below 80%"
            exit 1
          fi
      
      - name: Run integration tests
        run: go test -v -tags=integration ./...
      
      - name: Run benchmarks
        run: go test -bench=. -benchmem ./...

  # Stage 3: Security Scan
  security:
    name: Security Scan
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run gosec security scanner
        uses: securego/gosec@master
        with:
          args: '-fmt sarif -out results.sarif ./...'
      
      - name: Upload SARIF file
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: results.sarif
      
      - name: Check for vulnerabilities
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

  # Stage 4: Quality Check
  quality:
    name: Quality Check
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: ${{ env.GOLANGCI_LINT_VERSION }}
          args: --timeout=5m
      
      - name: Check go fmt
        run: |
          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
            echo "Please run 'gofmt -s -w .' to format your code"
            gofmt -s -l .
            exit 1
          fi
      
      - name: Check go vet
        run: go vet ./...
      
      - name: Check go mod tidy
        run: |
          go mod tidy
          if [ -n "$(git status --porcelain)" ]; then
            echo "Please run 'go mod tidy'"
            exit 1
          fi

  # Stage 5: Package
  package:
    name: Package
    needs: [test, security, quality]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' || startsWith(github.ref, 'refs/heads/release/')
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Build binaries
        run: |
          # Build for multiple platforms
          GOOS=linux GOARCH=amd64 go build -o dist/context-extender-linux-amd64 ./cmd/context-extender
          GOOS=darwin GOARCH=amd64 go build -o dist/context-extender-darwin-amd64 ./cmd/context-extender
          GOOS=windows GOARCH=amd64 go build -o dist/context-extender-windows-amd64.exe ./cmd/context-extender
      
      - name: Create Docker image
        run: |
          docker build -t context-extender:${{ github.sha }} .
          docker tag context-extender:${{ github.sha }} context-extender:latest
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: release-artifacts
          path: dist/

  # Stage 6: Deploy (Staging)
  deploy-staging:
    name: Deploy to Staging
    needs: package
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    environment: staging
    steps:
      - name: Deploy to staging
        run: |
          echo "Deploying to staging environment"
          # Add actual deployment commands here

  # Stage 7: Deploy (Production)
  deploy-production:
    name: Deploy to Production
    needs: package
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - name: Deploy to production
        run: |
          echo "Deploying to production environment"
          # Add actual deployment commands here
```

### Feature Flag Workflow

```yaml
# .github/workflows/feature-flags.yml
name: Feature Flag Management

on:
  workflow_dispatch:
    inputs:
      feature:
        description: 'Feature flag name'
        required: true
      action:
        description: 'Action to perform'
        required: true
        type: choice
        options:
          - enable
          - disable
          - rollout
      percentage:
        description: 'Rollout percentage (if applicable)'
        required: false
        default: '0'

jobs:
  manage-feature:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Update feature flag
        run: |
          echo "Updating feature flag: ${{ github.event.inputs.feature }}"
          echo "Action: ${{ github.event.inputs.action }}"
          echo "Percentage: ${{ github.event.inputs.percentage }}"
          # Add actual feature flag management here
```

## Pipeline Integration with 4-Phase Cycle

### Phase-Specific Pipeline Behavior

#### Research Phase
- **Branch**: `research/cycle-XXX`
- **Pipeline**: Build and test only
- **Deployment**: None
- **Purpose**: Validate technical feasibility

#### Planning Phase
- **Branch**: `planning/cycle-XXX`
- **Pipeline**: Full pipeline except deployment
- **Deployment**: None
- **Purpose**: Validate design decisions

#### Implementation Phase
- **Branch**: `feature/story-XXX`
- **Pipeline**: Full pipeline
- **Deployment**: Staging on merge to develop
- **Purpose**: Continuous integration and testing

#### Review Phase
- **Branch**: `release/cycle-XXX`
- **Pipeline**: Full pipeline with additional checks
- **Deployment**: Production on approval
- **Purpose**: Final validation and release

## Quality Gates

### Pipeline Quality Gates

| Gate | Criteria | Action on Failure |
|------|----------|-------------------|
| **Build** | Successful compilation | Stop pipeline |
| **Test Coverage** | ≥ 80% | Stop pipeline |
| **Security Scan** | No high/critical vulnerabilities | Stop pipeline |
| **Lint** | No errors | Stop pipeline |
| **Format** | Code formatted | Stop pipeline |
| **Dependencies** | go mod tidy clean | Stop pipeline |

### Branch Protection Rules

```yaml
# Branch protection for main
main:
  required_reviews: 1
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
  required_status_checks:
    - build
    - test
    - security
    - quality
  enforce_admins: false
  restrictions:
    users: []
    teams: []
```

## Monitoring and Notifications

### Pipeline Notifications

```yaml
# .github/workflows/notifications.yml
name: Pipeline Notifications

on:
  workflow_run:
    workflows: ["CI/CD Pipeline"]
    types: [completed]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - name: Send notification
        run: |
          if [ "${{ github.event.workflow_run.conclusion }}" == "failure" ]; then
            echo "Pipeline failed - sending notification"
            # Add Slack/Email notification here
          fi
```

### Metrics Collection

```yaml
# Pipeline metrics to track
metrics:
  - build_duration
  - test_duration
  - test_coverage
  - security_vulnerabilities
  - quality_issues
  - deployment_frequency
  - lead_time
  - mean_time_to_recovery
```

## Local Development Setup

### Pre-commit Hooks

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: gofmt -s -w
        language: system
        files: '\.go$'
      
      - id: go-vet
        name: go vet
        entry: go vet
        language: system
        files: '\.go$'
      
      - id: go-test
        name: go test
        entry: go test
        language: system
        files: '\.go$'
      
      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run
        language: system
        files: '\.go$'
```

### Local Pipeline Testing

```bash
# Install act for local GitHub Actions testing
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash  # Linux

# Test pipeline locally
act -j build
act -j test
act -j security
```

## Pipeline Maintenance

### Regular Maintenance Tasks

| Task | Frequency | Responsibility |
|------|-----------|----------------|
| Update Go version | Quarterly | Claude |
| Update action versions | Monthly | Claude |
| Review security policies | Monthly | Claude |
| Clean up old artifacts | Weekly | Automated |
| Monitor pipeline performance | Daily | Claude |

### Pipeline Optimization

1. **Cache Optimization**
   - Cache dependencies
   - Cache build artifacts
   - Cache test results

2. **Parallel Execution**
   - Run independent jobs in parallel
   - Use matrix builds for multiple versions

3. **Conditional Execution**
   - Skip steps based on file changes
   - Run heavy tests only on main branch

## Troubleshooting Guide

### Common Issues and Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| Build fails | Dependency issues | Clear cache, run go mod tidy |
| Tests timeout | Long-running tests | Increase timeout, optimize tests |
| Security scan fails | New vulnerability | Update dependencies, patch code |
| Coverage drops | New untested code | Add tests before merging |
| Deployment fails | Configuration issue | Check environment variables |

## Integration with Existing Processes

### Risk Management
- Security scanning identifies vulnerabilities
- Failed pipelines trigger risk assessment
- Deployment rollback plans in place

### Quality Assurance
- Automated testing ensures quality
- Coverage thresholds enforce standards
- Quality gates prevent bad code

### Stakeholder Communication
- Pipeline status visible in GitHub
- Notifications on failures
- Deployment announcements

## Success Metrics

### Pipeline Metrics
- **Build Success Rate**: > 95%
- **Mean Time to Feedback**: < 10 minutes
- **Deployment Frequency**: Daily
- **Lead Time**: < 1 day
- **Mean Time to Recovery**: < 1 hour

### Quality Metrics
- **Test Coverage**: > 80%
- **Security Vulnerabilities**: 0 critical/high
- **Code Quality Score**: > 95%
- **Failed Deployments**: < 5%