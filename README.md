# Context-Extender

**Automatically capture and manage your Claude Code conversations**

[![Version](https://img.shields.io/badge/version-1.2.0-blue.svg)](https://github.com/mattbran87/context-extender/releases/latest)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)](#-requirements)

Context-Extender seamlessly integrates with Claude Code to capture all your conversations, enabling powerful features like conversation history, search, analytics, and data export.

## 🚀 Quick Start

### 1. Download & Install

#### 🪟 **Windows**
[**⬇️ Download for Windows**](https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-windows-amd64.exe)
```bash
# Rename downloaded file for easier use
mv context-extender-1.2.0-windows-amd64.exe context-extender.exe

# Interactive installation (recommended)
./context-extender.exe install
```

#### 🍎 **macOS**
```bash
# Intel Mac
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-darwin-amd64

# Apple Silicon (M1/M2/M3)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-darwin-arm64

# Make executable and install
chmod +x context-extender
./context-extender install
```

#### 🐧 **Linux**
```bash
# x86_64 (Intel/AMD)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-linux-amd64

# ARM64 (Raspberry Pi, etc.)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-linux-arm64

# Make executable and install
chmod +x context-extender
./context-extender install
```

**Installation Options:**
- `install` - Interactive wizard with explanations, perfect for first-time users
- `configure` - Direct installation without prompts, ideal for experienced users

### 2. Start Using
✅ Open Claude Code and have conversations as normal
✅ Context-Extender automatically captures everything
✅ No changes needed to your workflow!

### 3. View Your Data
```bash
# List all conversations
context-extender query list

# Export to Excel with charts and analytics
context-extender export --format xlsx --output my-conversations.xlsx
```

## ✨ Features

- **🔄 Automatic Capture**: Zero-effort conversation logging
- **🔍 Search & Filter**: Find any conversation instantly
- **📊 Rich Exports**: CSV, JSON, Excel with charts and analytics
- **📈 Usage Analytics**: Track your Claude Code productivity
- **🗑️ Safe Uninstall**: Professional uninstall with data protection (NEW in v1.2.0)
- **🔒 Privacy First**: All data stays local on your computer
- **⚡ High Performance**: Zero impact on Claude Code speed

## 🎯 Perfect For

- **Developers** tracking Claude-assisted coding sessions
- **Researchers** analyzing AI conversation patterns
- **Teams** sharing Claude Code insights and best practices
- **Power Users** optimizing their AI workflow

## 📋 Requirements

- **Claude Code** installed and working
- **Operating System**: Windows, macOS, or Linux
- **Architecture**: x86_64 (Intel/AMD) or ARM64 (Apple Silicon, Raspberry Pi)
- **Dependencies**: None! Pure Go binaries work out of the box

## 🛠 Installation Options

### 🌟 Interactive Wizard (Best for First-Time Users)
```bash
context-extender install
```
✅ Step-by-step guidance with explanations
✅ System requirements checking
✅ Installation testing and verification
✅ Built-in getting started tutorial

### ⚡ Quick Setup (For Experienced Users)
```bash
context-extender configure
```
✅ Immediate installation without interactive prompts
✅ Automatic database initialization
✅ Hook verification and testing

### 🔍 Check Installation Status
```bash
context-extender configure --status
```
Shows detailed hook status and Claude Code integration health.

## 📊 Export Your Data

### Excel with Charts and Analytics
```bash
context-extender export --format xlsx --output conversations.xlsx
```

### CSV for Spreadsheet Analysis
```bash
context-extender export --format csv --output conversations.csv
```

### JSON for Programmatic Access
```bash
context-extender export --format json --pretty --output conversations.json
```

### Advanced Filtering
```bash
# Last week's conversations
context-extender export --format xlsx --from 2025-01-15 --to 2025-01-22 --output week.xlsx

# Specific columns only
context-extender export --format csv --columns session_id,user_prompts,claude_replies --output summary.csv

# Preview before exporting (see what you'll get)
context-extender export --format csv --preview

# Compress large exports (90% size reduction)
context-extender export --format json --compress --output backup.json.gz
```

## 🔧 Essential Commands

```bash
# Check what's being captured
context-extender query list

# View specific conversation
context-extender query show <session-id>

# Database status and statistics
context-extender database status

# See all available commands
context-extender --help
```

## 🗑️ Uninstall Options

### Safe Complete Removal
```bash
context-extender uninstall
```
✅ Interactive prompts prevent accidental data loss
✅ Removes Claude Code hooks
✅ Deletes database and conversation data (with confirmation)
✅ Cleans up system PATH installation

### Keep Your Data
```bash
context-extender uninstall --keep-data
```
✅ Removes hooks but preserves all conversation data
✅ Perfect for temporary removal or reinstallation
✅ Data will be automatically detected on reinstall

### Automated Uninstall
```bash
context-extender uninstall --force
```
⚠️ Skips confirmation prompts (use with caution)
⚠️ Intended for scripts and automation

## 🛡 Privacy & Security

- **Local Storage Only**: All data stays on your computer
- **No Cloud Sync**: Unless you explicitly export data
- **No Telemetry**: Zero data collection or tracking
- **Open Source**: Full transparency in what the tool does

## 🚨 Troubleshooting

### No conversations appearing?
1. **Check hook status**: `context-extender configure --status`
2. **Restart Claude Code** after installation
3. **Have actual conversations** (prompts + responses, not just opening Claude Code)
4. **Verify installation**: Try `context-extender install` for guided diagnostics

### Installation issues?
1. **Verify Claude Code** is installed and working properly
2. **Run as administrator** (Windows) or with `sudo` (macOS/Linux) if needed
3. **Check file permissions** - Claude Code settings file isn't read-only
4. **Use interactive wizard**: `context-extender install` provides step-by-step diagnostics

### Database errors?
✅ **Fixed in v1.2.0**: Database tables now automatically created during installation
- Old installations: Run `context-extender configure` to fix database issues
- Fresh installations: Database setup is automatic

**Need more help?** See [GETTING_STARTED.md](GETTING_STARTED.md) for detailed instructions.

## 🤝 Contributing

Found a bug or want a feature? [File an issue](https://github.com/mattbran87/context-extender/issues) or submit a pull request.

## 📄 License

Open source under MIT License. See LICENSE file for details.

## 🆕 What's New in v1.2.0

- 🗑️ **Professional Uninstall System** - Safe removal with data protection options
- 🚀 **Interactive Installation Wizard** - Step-by-step setup for first-time users
- 🛠️ **Fixed Database Issues** - Automatic schema creation during installation
- 📚 **Complete Documentation** - Professional README and getting started guide
- 🔧 **Enhanced CLI Help** - Better command organization and examples

[**📝 Full Release Notes**](https://github.com/mattbran87/context-extender/releases/tag/v1.2.0)

---

**Made with ❤️ for the Claude Code community**