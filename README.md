# Context-Extender

**Automatically capture and manage your Claude Code conversations**

Context-Extender seamlessly integrates with Claude Code to capture all your conversations, enabling powerful features like conversation history, search, analytics, and data export.

## 🚀 Quick Start

### 1. Install (First Time)
```bash
# Interactive installation wizard (recommended)
context-extender install

# Or quick setup for experienced users
context-extender configure
```

### 2. Start Using
- Open Claude Code and have conversations as normal
- Context-Extender automatically captures everything
- No changes needed to your workflow!

### 3. View Your Data
```bash
# List all conversations
context-extender query list

# Export to Excel with charts
context-extender export --format xlsx --output my-conversations.xlsx
```

## ✨ Features

- **🔄 Automatic Capture**: Zero-effort conversation logging
- **🔍 Search & Filter**: Find any conversation instantly
- **📊 Rich Exports**: CSV, JSON, Excel with charts and analytics
- **📈 Usage Analytics**: Track your Claude Code productivity
- **🔒 Privacy First**: All data stays local on your computer
- **⚡ High Performance**: Zero impact on Claude Code speed

## 🎯 Perfect For

- **Developers** tracking Claude-assisted coding sessions
- **Researchers** analyzing AI conversation patterns
- **Teams** sharing Claude Code insights and best practices
- **Power Users** optimizing their AI workflow

## 📋 Requirements

- Claude Code installed and working
- Windows, macOS, or Linux
- No additional dependencies

## 🛠 Installation Options

### Interactive Wizard (Best for First-Time Users)
```bash
context-extender install
```
Step-by-step guidance with explanations at every step.

### Quick Setup (For Experienced Users)
```bash
context-extender configure
```
Immediate installation without interactive prompts.

### Check Installation
```bash
context-extender configure --status
```

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
context-extender export --format xlsx --from 2024-01-01 --to 2024-01-07 --output week.xlsx

# Specific columns only
context-extender export --format csv --columns session_id,user_prompts,claude_replies --output summary.csv

# Preview before exporting
context-extender export --format csv --preview
```

## 🔧 Common Commands

```bash
# Check what's being captured
context-extender query list

# View specific conversation
context-extender query show <session-id>

# Database status and statistics
context-extender database status

# Complete uninstall
context-extender uninstall

# Remove hooks only (keep data)
context-extender uninstall --keep-data

# See all options
context-extender --help
```

## 🛡 Privacy & Security

- **Local Storage Only**: All data stays on your computer
- **No Cloud Sync**: Unless you explicitly export data
- **No Telemetry**: Zero data collection or tracking
- **Open Source**: Full transparency in what the tool does

## 🚨 Troubleshooting

**No conversations appearing?**
1. Check hook status: `context-extender configure --status`
2. Restart Claude Code after installation
3. Ensure you're having actual conversations (not just opening Claude Code)

**Installation issues?**
1. Verify Claude Code is installed and working
2. Run with administrator privileges if needed
3. Check that Claude Code settings file isn't read-only

**Need help?** See [GETTING_STARTED.md](GETTING_STARTED.md) for detailed instructions.

## 🤝 Contributing

Found a bug or want a feature? [File an issue](https://github.com/your-repo/context-extender/issues) or submit a pull request.

## 📄 License

Open source under MIT License. See LICENSE file for details.

---

**Made with ❤️ for the Claude Code community**