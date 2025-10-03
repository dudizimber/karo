# Copilot Prompt System

## 🚀 Quick Start

### For New Features
1. Copy `.copilot/prompts/add-feature.md`
2. Replace `{feature_name}` and `{requirements}`
3. Paste into Copilot chat

### For Bug Fixes
1. Copy `.copilot/prompts/fix-bug.md`  
2. Fill in `{issue_description}` and error logs
3. Get debugging assistance

### For Commit Messages
1. Copy `.copilot/prompts/commit-messages.md`
2. Describe your changes
3. Get properly formatted conventional commit

## 📁 File Structure

```
.copilot/
├── README.md                    # Complete usage guide
├── quick-prompts.md            # Copy-paste ready prompts
├── context-templates.md        # Project context for complex queries
└── prompts/
    ├── add-feature.md          # Adding new operator features
    ├── fix-bug.md              # Debugging and troubleshooting
    ├── add-tests.md            # Creating comprehensive tests
    ├── commit-messages.md      # Conventional commit formatting
    ├── code-review.md          # Code review assistance
    ├── crd-changes.md          # AlertReaction CRD modifications
    ├── controller-logic.md     # Controller development
    ├── ci-cd-updates.md        # GitHub Actions improvements
    ├── write-docs.md           # Documentation creation
    ├── optimize-performance.md # Performance improvements
    └── refactor-code.md        # Code refactoring
```

## 🎯 Common Use Cases

| Task | Template | Example |
|------|----------|---------|
| Add feature | `add-feature.md` | "Add alert prioritization" |
| Fix issue | `fix-bug.md` | "Controller creates duplicate Jobs" |
| Write tests | `add-tests.md` | "Test matcher evaluation logic" |
| Code review | `code-review.md` | "Review new webhook handler" |
| Commit msg | `commit-messages.md` | "feat(controller): add priority queue" |
| Update docs | `write-docs.md` | "Document new matcher operators" |

## ✨ Tips for Better Results

1. **Be Specific**: Include file paths, function names, exact requirements
2. **Provide Context**: Use the context templates for complex questions
3. **Include Examples**: Show current code and expected outcomes
4. **Iterate**: Refine prompts based on Copilot's initial responses
5. **Test Results**: Always validate generated code with tests

## 🔧 Customization

### Adding New Prompts
1. Create new `.md` file in `prompts/` directory
2. Follow existing template structure
3. Include placeholders `{like_this}` for customization
4. Add usage examples

### Modifying Existing Prompts
1. Update template with new requirements
2. Test with actual use cases
3. Update documentation and examples

## 🚀 Integration with Development Workflow

### Pre-Development
```bash
# Use feature planning prompt
cp .copilot/prompts/add-feature.md /tmp/feature-prompt.md
# Fill in details and get implementation plan
```

### During Development
```bash
# Quick debugging
grep -A 10 "Bug Fix" .copilot/quick-prompts.md
# Copy, customize, and paste into Copilot
```

### Pre-Commit
```bash
# Generate commit message
grep -A 15 "Commit Message" .copilot/quick-prompts.md
# Fill in changes and get conventional commit format
```

## 📚 Project-Specific Context

When using any prompt, Copilot knows about:
- **Kubernetes Operator** built with Go 1.24+ and controller-runtime
- **AlertReaction CRD** with Prometheus-style matchers (=, !=, =~, !~)
- **Alert Processing** workflow: AlertManager → Webhook → Controller → Jobs
- **Testing Patterns** with table-driven tests and fake clients
- **Git Hooks** for automated code quality enforcement
- **CI/CD Pipeline** with GitHub Actions and security scanning

## 🎯 Success Metrics

Good prompts result in:
- ✅ Code that follows project patterns and conventions
- ✅ Comprehensive test coverage for new functionality
- ✅ Proper error handling and edge case coverage
- ✅ Documentation that matches implementation
- ✅ Performance-conscious solutions
- ✅ Security-aware implementations

## 🤝 Contributing

To improve the prompt system:
1. Use prompts and note what works/doesn't work
2. Create issues for prompt improvements needed
3. Submit PRs with enhanced templates
4. Share successful prompt patterns with the team

---

**Made with ❤️ to streamline k8s-alert-reaction-operator development**