# Context Extender Deployment Guide

## Overview

This guide covers deployment strategies, build processes, and production configuration for Context Extender.

## Build Requirements

### Pure Go Build (Recommended)
- Go 1.21 or later
- No CGO dependencies required
- Cross-platform compilation supported

### Build Commands

```bash
# Build for current platform
go build -o context-extender main.go

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o context-extender-linux main.go
GOOS=darwin GOARCH=amd64 go build -o context-extender-mac main.go
GOOS=windows GOARCH=amd64 go build -o context-extender.exe main.go

# Build with version information
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o context-extender main.go

# Optimized production build (smaller binary)
go build -ldflags="-s -w" -o context-extender main.go
```

## Deployment Strategies

### 1. Standalone Binary

The simplest deployment method - distribute a single executable.

```bash
# Build
go build -ldflags="-s -w" -o context-extender main.go

# Deploy
scp context-extender user@server:/usr/local/bin/
ssh user@server chmod +x /usr/local/bin/context-extender
```

### 2. Docker Container

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o context-extender main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/context-extender .
COPY --from=builder /app/config/default.json ./config/
EXPOSE 8080
CMD ["./context-extender"]
```

Build and run:
```bash
docker build -t context-extender:latest .
docker run -v ~/.context-extender:/root/.context-extender context-extender:latest
```

### 3. Systemd Service

```ini
# /etc/systemd/system/context-extender.service
[Unit]
Description=Context Extender Service
After=network.target

[Service]
Type=simple
User=context-extender
Group=context-extender
WorkingDirectory=/opt/context-extender
ExecStart=/opt/context-extender/context-extender serve
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=context-extender

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/context-extender

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable context-extender
sudo systemctl start context-extender
```

### 4. Kubernetes Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: context-extender
  labels:
    app: context-extender
spec:
  replicas: 3
  selector:
    matchLabels:
      app: context-extender
  template:
    metadata:
      labels:
        app: context-extender
    spec:
      containers:
      - name: context-extender
        image: context-extender:latest
        ports:
        - containerPort: 8080
        env:
        - name: CE_DATABASE_PATH
          value: "/data/context-extender.db"
        - name: CE_ENCRYPTION_ENABLED
          value: "true"
        volumeMounts:
        - name: data
          mountPath: /data
        - name: config
          mountPath: /config
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: context-extender-pvc
      - name: config
        configMap:
          name: context-extender-config
```

## Configuration Management

### Environment Variables

All configuration options can be set via environment variables:

```bash
export CE_DATABASE_PATH=/var/lib/context-extender/data.db
export CE_ENCRYPTION_ENABLED=true
export CE_ENCRYPTION_KEY_PATH=/etc/context-extender/encryption.key
export CE_LOG_LEVEL=info
export CE_MAX_CONNECTIONS=100
export CE_CACHE_SIZE=1000000
```

### Configuration File

Production configuration template:

```json
{
  "database": {
    "path": "/var/lib/context-extender/data.db",
    "max_connections": 100,
    "cache_size": 1000000,
    "wal_mode": true,
    "synchronous": "NORMAL"
  },
  "encryption": {
    "enabled": true,
    "method": "aes-gcm",
    "key_path": "/etc/context-extender/encryption.key",
    "pbkdf2_iterations": 100000
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "/var/log/context-extender/app.log",
    "rotation": {
      "enabled": true,
      "max_size": "100MB",
      "max_age": "30d",
      "max_backups": 10
    }
  },
  "metrics": {
    "enabled": true,
    "port": 9090,
    "path": "/metrics"
  }
}
```

## Production Checklist

### Pre-Deployment

- [ ] Run all tests: `go test ./...`
- [ ] Run benchmarks: `go test -bench=. ./...`
- [ ] Static analysis: `go vet ./...`
- [ ] Security scan: `gosec ./...`
- [ ] Update dependencies: `go mod tidy`
- [ ] Generate documentation
- [ ] Review configuration

### Security

- [ ] Enable encryption for sensitive data
- [ ] Secure key storage (use key management service in production)
- [ ] Set appropriate file permissions (600 for keys, 644 for configs)
- [ ] Configure TLS for network communication
- [ ] Implement rate limiting
- [ ] Set up audit logging
- [ ] Configure firewall rules
- [ ] Use least-privilege service account

### Performance

- [ ] Configure connection pooling
- [ ] Enable WAL mode for SQLite
- [ ] Set appropriate cache sizes
- [ ] Configure GOMAXPROCS for container environments
- [ ] Enable pprof for profiling (development only)
- [ ] Set up monitoring and alerting

### Backup & Recovery

- [ ] Automated database backups
- [ ] Encryption key backups (secure storage)
- [ ] Configuration backups
- [ ] Test restore procedures
- [ ] Document recovery processes

## Monitoring

### Health Checks

```go
// Health endpoint implementation
func healthCheck(w http.ResponseWriter, r *http.Request) {
    checks := map[string]string{
        "database": checkDatabase(),
        "encryption": checkEncryption(),
        "disk_space": checkDiskSpace(),
    }

    status := http.StatusOK
    for _, check := range checks {
        if check != "ok" {
            status = http.StatusServiceUnavailable
            break
        }
    }

    w.WriteHeader(status)
    json.NewEncoder(w).Encode(checks)
}
```

### Metrics

Prometheus metrics endpoint:

```go
// Metrics to track
- context_extender_requests_total
- context_extender_request_duration_seconds
- context_extender_database_connections_active
- context_extender_encryption_operations_total
- context_extender_cache_hits_total
- context_extender_cache_misses_total
```

### Logging

Structured logging configuration:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))
```

## Scaling Strategies

### Horizontal Scaling

For read-heavy workloads:
1. Deploy multiple read replicas
2. Use load balancer for distribution
3. Implement caching layer (Redis/Memcached)
4. Consider CDN for static assets

### Vertical Scaling

For write-heavy workloads:
1. Increase CPU/memory resources
2. Optimize database settings
3. Use SSD storage
4. Increase connection pool size

## Troubleshooting

### Common Issues

**High Memory Usage**
```bash
# Check memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Adjust GOGC for more aggressive garbage collection
export GOGC=50
```

**Database Lock Errors**
```bash
# Enable WAL mode
context-extender config set database.wal_mode=true

# Increase busy timeout
context-extender config set database.busy_timeout=5000
```

**Slow Startup**
```bash
# Check initialization logs
journalctl -u context-extender -f

# Profile startup
go test -cpuprofile=cpu.prof -memprofile=mem.prof
```

### Debug Mode

Enable debug logging:
```bash
context-extender --debug serve
```

## CI/CD Pipeline

### GitHub Actions Example

```yaml
name: Deploy

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Test
      run: go test -v ./...

    - name: Build
      run: |
        VERSION=${GITHUB_REF#refs/tags/}
        go build -ldflags="-s -w -X main.Version=$VERSION" -o context-extender

    - name: Build Docker image
      run: |
        docker build -t context-extender:${GITHUB_REF#refs/tags/} .
        docker push context-extender:${GITHUB_REF#refs/tags/}

    - name: Deploy
      run: |
        kubectl set image deployment/context-extender context-extender=context-extender:${GITHUB_REF#refs/tags/}
```

## Migration Strategies

### Zero-Downtime Deployment

1. **Blue-Green Deployment**
   - Deploy new version to green environment
   - Test thoroughly
   - Switch traffic from blue to green
   - Keep blue as rollback option

2. **Rolling Updates**
   - Update instances one at a time
   - Health checks between updates
   - Automatic rollback on failures

3. **Canary Deployment**
   - Deploy to small percentage of users
   - Monitor metrics and errors
   - Gradually increase traffic
   - Full rollout or rollback based on results

## Support

For deployment assistance:
- Documentation: [docs/](.)
- Issues: [GitHub Issues](https://github.com/yourusername/context-extender/issues)
- Community: [Discussions](https://github.com/yourusername/context-extender/discussions)