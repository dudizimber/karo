# Write Documentation Prompt Template

## Prompt
```
I'm working on Karo (Kubernetes Alert Reaction Operator) and need to create or update documentation. The operator is a Kubernetes controller that creates Jobs in response to Prometheus alerts via AlertManager webhooks.

## Documentation Task
**Type of documentation:**
- [ ] README updates
- [ ] API documentation
- [ ] User guide/tutorial
- [ ] Developer documentation
- [ ] Architecture documentation
- [ ] Troubleshooting guide
- [ ] Installation guide
- [ ] Configuration reference
- [ ] Migration guide
- [ ] Release notes/changelog

**Target audience:**
- [ ] End users (platform teams, DevOps engineers)
- [ ] Developers (contributors to the project)
- [ ] Operators (cluster administrators)
- [ ] Security teams (compliance and security review)

## Current Documentation State
**Existing documentation:**
- `README.md` - {current_state_and_what_needs_updating}
- `docs/` directory - {what_exists_currently}
- Code comments - {documentation_quality_in_code}
- Examples - {current_examples_and_their_quality}

**Documentation gaps:**
1. {gap_1_description}
2. {gap_2_description}
3. {gap_3_description}

## Content Requirements
**What needs to be documented:**
{detailed_description_of_what_you_want_documented}

**Specific topics to cover:**
- [ ] {topic_1} - {why_this_is_important}
- [ ] {topic_2} - {why_this_is_important}
- [ ] {topic_3} - {why_this_is_important}

**Technical depth needed:**
- [ ] High-level overview (marketing/executive summary)
- [ ] User-focused tutorial (step-by-step guides)
- [ ] Technical reference (detailed API/configuration)
- [ ] Developer deep-dive (architecture and internals)

## Project Context for Documentation
**Key concepts to explain:**
- **AlertReaction CRD** - Custom resource that defines alert-to-job mappings
- **Prometheus-style matchers** - Alert filtering using =, !=, =~, !~ operators
- **AlertManager integration** - Webhook-based alert reception
- **Job creation** - Kubernetes Jobs triggered by matching alerts
- **Multi-action support** - Multiple responses per alert type

**Architecture overview:**
```
Prometheus → AlertManager → Webhook → Controller → AlertReaction → Jobs
```

**Current features:**
- v1alpha1 API with Prometheus-style matchers
- Multiple AlertReactions per alert support
- Comprehensive matcher evaluation (regex, exact match)
- Volume mounting and service account support
- Environment variable substitution from alert data
- Git hooks for code quality automation

## Examples and Use Cases
**Real-world scenarios to document:**
1. {use_case_1} - {brief_description}
2. {use_case_2} - {brief_description}
3. {use_case_3} - {brief_description}

**Example configurations needed:**
```yaml
{paste_example_alertreaction_configurations_if_relevant}
```

**Expected outcomes:**
{describe_what_users_should_achieve_after_following_the_documentation}

## Format and Style
**Preferred documentation format:**
- [ ] Markdown files
- [ ] OpenAPI/Swagger specs
- [ ] Code comments (GoDoc)
- [ ] Wiki pages
- [ ] Interactive tutorials
- [ ] Video content (scripts)

**Style requirements:**
- [ ] Beginner-friendly language
- [ ] Technical precision
- [ ] Step-by-step instructions
- [ ] Troubleshooting sections
- [ ] Code examples with explanations
- [ ] Visual diagrams (ASCII art, mermaid)

## Integration Requirements
**Documentation should integrate with:**
- [ ] GitHub repository structure
- [ ] CI/CD pipeline (automated checks)
- [ ] Code generation (API docs)
- [ ] Version control (versioned docs)
- [ ] External documentation sites

**Links and references:**
- [ ] Link to related Kubernetes concepts
- [ ] Reference external tools (Prometheus, AlertManager)
- [ ] Point to relevant code sections
- [ ] Include troubleshooting flowcharts

## Quality Standards
**Documentation should be:**
- [ ] **Accurate** - Reflects current codebase and behavior
- [ ] **Complete** - Covers all necessary topics
- [ ] **Clear** - Easy to understand for the target audience
- [ ] **Actionable** - Provides concrete steps and examples
- [ ] **Maintainable** - Easy to keep up-to-date
- [ ] **Searchable** - Well-organized with good headings

**Validation criteria:**
- [ ] Technical reviewer approval
- [ ] User testing (can someone follow the docs successfully?)
- [ ] Link checking (no broken references)
- [ ] Code example testing (examples actually work)

## Request for Copilot
Please help me create comprehensive documentation that:

1. **Explains the concepts clearly** for the target audience
2. **Provides practical examples** that users can follow
3. **Covers edge cases and troubleshooting** common issues
4. **Follows documentation best practices** for structure and style
5. **Integrates well** with existing project documentation
6. **Is maintainable** and easy to keep current
7. **Includes visual aids** where helpful (diagrams, flowcharts)

**Specific assistance needed:**
1. How should I structure the documentation for different audiences?
2. What examples would be most helpful for users?
3. How can I make complex concepts more accessible?
4. What troubleshooting scenarios should I cover?
5. How should I organize the information architecture?

## Success Metrics
**Good documentation will result in:**
- [ ] Reduced support questions and issues
- [ ] Faster user onboarding and adoption
- [ ] Fewer documentation-related bug reports
- [ ] Positive feedback from users and contributors
- [ ] Self-service capability for common tasks

---

## Usage Instructions
1. Replace `{placeholders}` with your specific documentation needs
2. Describe the current state and gaps clearly
3. Specify the target audience and technical depth
4. Include examples of what you want documented
5. Copy and paste into Copilot chat

## Example Usage
```
**Type of documentation:** User guide/tutorial
**Target audience:** Platform teams setting up alert automation

**What needs to be documented:**
Complete end-to-end tutorial for setting up automated disk cleanup when storage alerts fire.

**Specific topics to cover:**
- Setting up AlertManager webhook integration
- Creating AlertReaction with disk cleanup job
- Configuring RBAC permissions for cleanup operations
- Testing the complete workflow
- Troubleshooting common issues

**Real-world scenarios:**
1. Disk space alerts trigger automated log cleanup
2. High memory alerts trigger application restart
3. Database connection issues trigger backup creation
```