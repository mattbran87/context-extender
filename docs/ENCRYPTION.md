# Context Extender Encryption System

## Overview

Context Extender implements a robust AES-256-GCM encryption system that provides SQLCipher-equivalent security without requiring CGO compilation. The system uses application-level encryption with field-level granularity, allowing selective encryption of sensitive data.

## Key Features

- **AES-256-GCM Encryption**: Industry-standard authenticated encryption with 256-bit keys
- **PBKDF2 Key Derivation**: 100,000 iterations (configurable) for strong key generation from passwords
- **Field-Level Encryption**: Selective encryption of sensitive fields only
- **Pure Go Implementation**: No CGO dependencies required
- **High Performance**: Sub-microsecond encryption operations for small data
- **Key Rotation Support**: Built-in support for rotating encryption keys
- **Multiple Provider Support**: Extensible architecture for different encryption methods

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                     Application                      │
├─────────────────────────────────────────────────────┤
│                  Field Middleware                    │
│            (Selective Field Encryption)              │
├─────────────────────────────────────────────────────┤
│               Encryption Provider Factory            │
│         (AES-GCM, Plaintext, Future Methods)         │
├─────────────────────────────────────────────────────┤
│                 Encryption Interface                 │
│              (Standardized Operations)               │
├─────────────────────────────────────────────────────┤
│                   Database Layer                     │
│               (Pure Go SQLite Driver)                │
└─────────────────────────────────────────────────────┘
```

## Quick Start

### 1. Initialize Encryption

```bash
# Initialize encryption with interactive setup
context-extender encrypt init

# Or with quick setup using defaults
context-extender encrypt init --quick

# Check encryption status
context-extender encrypt status
```

### 2. Configuration

The encryption system can be configured via:
- Command-line flags
- Configuration file (`~/.context-extender/config.json`)
- Environment variables

Example configuration:
```json
{
  "encryption": {
    "method": "aes-gcm",
    "key_derivation": "pbkdf2",
    "pbkdf2_iterations": 100000,
    "key_path": "~/.context-extender/encryption.key"
  }
}
```

### 3. Programmatic Usage

```go
import "context-extender/internal/encryption"

// Create configuration
config := encryption.DefaultEncryptionConfig()
config.KeyPath = "./my-key.key"

// Create and initialize provider
provider, err := encryption.NewAESGCMProvider(config)
if err != nil {
    return err
}
defer provider.Close()

if err := provider.Initialize(); err != nil {
    return err
}

// Encrypt data
encrypted, err := provider.EncryptField([]byte("sensitive data"))
if err != nil {
    return err
}

// Decrypt data
decrypted, err := provider.DecryptField(encrypted)
if err != nil {
    return err
}
```

## Performance Characteristics

Based on benchmarks (Intel i7-1185G7):

| Operation | Data Size | Time | Memory | Allocations |
|-----------|-----------|------|--------|-------------|
| Encrypt | 100 bytes | 623ns | 240B | 3 |
| Encrypt | 1KB | 1.2μs | 1136B | 3 |
| Encrypt | 10KB | 3.5μs | 10352B | 3 |
| Decrypt | 100 bytes | 267ns | 112B | 1 |
| Decrypt | 1KB | 619ns | 1024B | 1 |
| Decrypt | 10KB | 3.0μs | 10240B | 1 |

## Security Considerations

### Key Management

- **Key Storage**: Encryption keys are stored separately from the database
- **Key Derivation**: Uses PBKDF2 with SHA-256 and configurable iterations (default: 100,000)
- **Key Rotation**: Support for rotating keys while maintaining data access
- **Key Backup**: Always backup encryption keys in a secure location

### Best Practices

1. **Use Strong Passwords**: When prompted for a password, use a strong, unique password
2. **Backup Keys**: Always maintain secure backups of encryption keys
3. **Regular Key Rotation**: Rotate keys periodically for enhanced security
4. **Selective Encryption**: Only encrypt sensitive fields to maintain performance
5. **Test Recovery**: Regularly test backup and recovery procedures

### Field-Level Encryption Rules

The middleware supports configurable rules for selective encryption:

```go
// Default rules encrypt sensitive field patterns
rules := &FieldRules{
    Tables: map[string]*TableRule{
        "sessions": {
            Fields: map[string]bool{
                "metadata": true,  // Encrypt metadata field
                "context": true,   // Encrypt context field
            },
        },
    },
}
```

## CLI Commands

### Encryption Management

```bash
# Initialize encryption
context-extender encrypt init [flags]
  --quick           Quick setup with defaults
  --key-path PATH   Path to store encryption key
  --iterations N    PBKDF2 iterations (default: 100000)

# View key information
context-extender encrypt key-info

# Rotate encryption key
context-extender encrypt rotate

# Migrate unencrypted data
context-extender encrypt migrate

# Verify encryption
context-extender encrypt verify

# Check encryption status
context-extender encrypt status
```

### Configuration Management

```bash
# Show current configuration
context-extender config show

# Generate default configuration
context-extender config generate

# Initialize configuration
context-extender config init

# Validate configuration
context-extender config validate
```

## Migration Guide

### From SQLCipher to AES-GCM

1. Export data from SQLCipher database
2. Initialize new encryption system
3. Import data with encryption enabled
4. Verify data integrity

```bash
# Export from SQLCipher
sqlcipher old.db "PRAGMA key='password'; .dump" > export.sql

# Initialize new encryption
context-extender encrypt init

# Import with encryption
context-extender import export.sql

# Verify
context-extender encrypt verify
```

### From Unencrypted to Encrypted

```bash
# Backup existing database
cp ~/.context-extender/data.db ~/.context-extender/data.db.backup

# Initialize encryption
context-extender encrypt init

# Migrate data
context-extender encrypt migrate

# Verify
context-extender encrypt verify
```

## Troubleshooting

### Common Issues

**Q: "Key ID mismatch" error**
A: This indicates the data was encrypted with a different key. Ensure you're using the correct key file.

**Q: Performance degradation after encryption**
A: Consider:
- Reducing PBKDF2 iterations if key derivation is slow
- Encrypting only sensitive fields
- Using batch operations for bulk data

**Q: Lost encryption key**
A: Without the encryption key, encrypted data cannot be recovered. Always maintain secure backups.

### Debug Mode

Enable debug logging for troubleshooting:
```bash
context-extender --debug encrypt status
```

## Advanced Topics

### Custom Encryption Providers

Implement the `EncryptionProvider` interface:

```go
type CustomProvider struct {
    // Your implementation
}

func (p *CustomProvider) Initialize() error { /* ... */ }
func (p *CustomProvider) EncryptField(data []byte) (*FieldEncryption, error) { /* ... */ }
func (p *CustomProvider) DecryptField(encrypted *FieldEncryption) ([]byte, error) { /* ... */ }
// ... other interface methods
```

### Performance Tuning

1. **Batch Operations**: Use `BulkEncrypt` and `BulkDecrypt` for multiple fields
2. **Connection Pooling**: Configure appropriate pool size for concurrent operations
3. **Caching**: Consider caching decrypted data in memory (with appropriate security measures)
4. **Selective Encryption**: Only encrypt truly sensitive data

### Integration with CI/CD

```yaml
# Example GitHub Actions workflow
- name: Setup encryption
  run: |
    echo "${{ secrets.ENCRYPTION_KEY }}" | base64 -d > encryption.key
    context-extender encrypt init --key-path encryption.key

- name: Run tests
  run: context-extender test
```

## API Reference

See the [API Documentation](./API.md) for detailed interface specifications.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines on contributing to the encryption system.

## License

The encryption system is part of the Context Extender project and follows the same license terms.