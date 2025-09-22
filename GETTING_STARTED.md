# Getting Started with Context-Extender

Welcome to Context-Extender! This guide will help you get up and running quickly.

## 🚀 Quick Installation

### Step 1: Download for Your Platform

#### 🪟 **Windows**
```bash
# Download directly or visit GitHub releases
curl -L -o context-extender.exe https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-windows-amd64.exe
```

#### 🍎 **macOS**
```bash
# Intel Mac
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-darwin-amd64

# Apple Silicon (M1/M2/M3)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-darwin-arm64

# Make executable
chmod +x context-extender
```

#### 🐧 **Linux**
```bash
# x86_64 (Intel/AMD)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-linux-amd64

# ARM64 (Raspberry Pi, etc.)
curl -L -o context-extender https://github.com/mattbran87/context-extender/releases/download/v1.2.0/context-extender-1.2.0-linux-arm64

# Make executable
chmod +x context-extender
```

### Step 2: Choose Installation Method

#### Option 1: Interactive Wizard (Recommended for first-time users)
```bash
# Windows
./context-extender.exe install

# macOS/Linux
./context-extender install
```
The wizard will guide you through every step with explanations.

#### Option 2: Quick Setup (For experienced users)
```bash
# Windows
./context-extender.exe configure

# macOS/Linux
./context-extender configure
```
This installs hooks immediately without interactive prompts.

## ✅ Verify Installation

Check that everything is working:
```bash
context-extender configure --status
```

You should see all hooks marked as "✅ Installed".

## 📝 Start Capturing Conversations

1. **Start a new Claude Code session**
2. **Have conversations with Claude as normal**
3. **Context-Extender automatically captures everything!**

## 🔍 View Your Data

### List all conversations
```bash
context-extender query list
```

### View specific conversation
```bash
context-extender query show <session-id>
```

### Check database status
```bash
context-extender database status
```

## 📊 Export Your Data

### Export to Excel (with charts and formatting)
```bash
context-extender export --format xlsx --output my-conversations.xlsx
```

### Export to CSV (for analysis in Excel/Google Sheets)
```bash
context-extender export --format csv --output my-conversations.csv
```

### Export to JSON (for programmatic use)
```bash
context-extender export --format json --output my-conversations.json
```

### Export with specific columns
```bash
context-extender export --format csv --columns session_id,user_prompts,claude_replies,total_words --output summary.csv
```

## 🔧 Common Commands

### Check what's being captured
```bash
context-extender query list --recent 5
```

### Export last week's conversations
```bash
context-extender export --format xlsx --from $(date -d '7 days ago' +%Y-%m-%d) --output last-week.xlsx
```

### See all available commands
```bash
context-extender --help
```

## 🛠 Troubleshooting

### No conversations appearing?
1. Check hook status: `context-extender configure --status`
2. Restart Claude Code after installation
3. Make sure you're having actual conversations (prompts + responses)

### Installation issues?
1. Make sure Claude Code is installed and working
2. Try running as administrator (Windows) or with sudo (Mac/Linux)
3. Check Claude Code settings file isn't read-only

### Need to uninstall?
```bash
# Complete removal (including data)
context-extender uninstall

# Remove hooks but keep conversations
context-extender uninstall --keep-data

# Force uninstall without prompts
context-extender uninstall --force
```

## 🎯 Pro Tips

1. **Export regularly** to avoid losing data
2. **Use meaningful session names** in Claude Code for better organization
3. **Check database status** occasionally to monitor storage usage
4. **Use filters** when exporting to focus on specific time periods or projects

## 📚 Learn More

- Run `context-extender [command] --help` for detailed command information
- Check the export options: `context-extender export --help`
- Explore query capabilities: `context-extender query --help`

## 🔒 Privacy & Security

- **All data stays local** on your computer
- **No cloud sync** unless you explicitly export data
- **No telemetry** or data collection
- **Open source** - you can verify what the tool does

---

**Need help?** File an issue on GitHub or check the documentation.
**Love the tool?** Star the repository and share with other Claude Code users!