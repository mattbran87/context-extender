# Conversation Continuity Management

## 🎯 Purpose
Prevent context and progress loss during conversation compression by maintaining persistent state documentation.

## 📋 Pre-Compression Checklist

### **1. Update Current Status** (CRITICAL)
```bash
# Always update before conversation ends
.claude/current_status
```
- Current cycle and phase
- Day number in current phase
- Major achievements completed
- Work in progress
- Next planned activities

### **2. Commit Code Changes** (CRITICAL)
```bash
# Ensure all work is saved
git add .
git commit -m "Progress checkpoint: [describe current state]"
```

### **3. Update Progress Documentation**
- Todo list status
- Implementation completion status
- Any architectural decisions made
- Test results and validations

### **4. Create Conversation Checkpoint**
```markdown
# Conversation Checkpoint - [Date]
## What We Just Completed:
- [List major accomplishments]

## Current Implementation State:
- [What's working/tested]
- [What's partially complete]
- [What's not started]

## Next Session Should Focus On:
- [Immediate next steps]
- [Priority items]

## Critical Context:
- [Important decisions made]
- [Architecture changes]
- [Integration points completed]
```

## 🔄 Post-Compression Recovery

### **1. Read Current Status First**
```bash
# Always start new session by reading:
.claude/current_status
.claude/project_guide.md
```

### **2. Verify Implementation State**
```bash
# Test what's actually working
go build .
go test ./...
```

### **3. Compare Documentation vs Reality**
- Check git commits for actual progress
- Test functionality described as "complete"
- Verify architectural decisions were implemented

## 🚨 Red Flags to Watch For

### **Documentation Says "Complete" But:**
- Code doesn't compile
- Tests don't pass
- Features don't work as described
- Dependencies are missing

### **Conversation Suggests Work Done That:**
- Isn't reflected in git history
- Doesn't match file timestamps
- Contradicts current code state
- Was planned but not implemented

## 🛠️ Recovery Actions

### **When Context Loss Detected:**
1. **Stop and assess** - Don't continue work based on unclear state
2. **Read all status files** - Get ground truth from documentation
3. **Test current functionality** - Verify what actually works
4. **Update status documents** - Correct any inaccuracies
5. **Clarify with user** - Confirm understanding before proceeding

### **When Implementation State Unclear:**
1. **Run tests** - See what passes/fails
2. **Check git log** - Review actual commit history
3. **Verify dependencies** - Ensure all required packages present
4. **Test core functionality** - Validate basic operations work
5. **Document findings** - Update status to reflect reality

## 📊 Continuous State Tracking

### **Real-Time Documentation**
- Update `.claude/current_status` frequently during work
- Use TodoWrite tool to track progress
- Commit code changes regularly with descriptive messages
- Test functionality as it's implemented

### **Session Boundaries**
- End each session with status update
- Start each session with status verification
- Use git commits as source of truth
- Cross-reference documentation with actual code

## 🎯 Best Practices

### **For Development Sessions:**
1. **Start**: Read current status, verify code state
2. **During**: Update todos, commit changes frequently
3. **End**: Update status, document progress, commit final state

### **For Planning Sessions:**
1. **Clearly separate** what's planned vs implemented
2. **Use future tense** for planned work
3. **Use past tense** for completed work
4. **Mark implementation status** clearly in documents

### **For Conversation Continuity:**
1. **Assume nothing** about previous context
2. **Verify claims** about implementation status
3. **Test before building** on supposed foundations
4. **Update docs to match reality**

## 🔍 Quality Assurance

### **Before Major Work:**
- Verify current phase and cycle
- Confirm what's actually implemented
- Test existing functionality
- Clear understanding of next steps

### **After Major Milestones:**
- Update all status documents
- Commit all code changes
- Document architectural decisions
- Verify tests pass

### **Session Handoffs:**
- Clear status documentation
- Working code state
- Tested functionality
- Documented next steps

This process ensures that conversation compression doesn't lead to lost context or confusion about implementation state.