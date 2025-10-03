# Commit Message Prompt Template

## Context
I'm working on the k8s-alert-reaction-operator and need to write a conventional commit message that follows our project standards.

## Changes Made
**Files modified:**
- {file_1}: {brief_description_of_changes}
- {file_2}: {brief_description_of_changes}
- {file_3}: {brief_description_of_changes}

**Type of changes:**
- [ ] New feature (feat)
- [ ] Bug fix (fix)
- [ ] Documentation (docs)
- [ ] Code style/formatting (style)
- [ ] Refactoring (refactor)
- [ ] Tests (test)
- [ ] Build/CI changes (ci)
- [ ] Maintenance tasks (chore)

**Scope (optional):**
- [ ] controller - Controller logic changes
- [ ] webhook - Webhook handling changes
- [ ] api - CRD/API type changes
- [ ] crd - CRD manifest changes
- [ ] tests - Test-related changes
- [ ] docs - Documentation changes
- [ ] ci - CI/CD pipeline changes
- [ ] deps - Dependency updates

## Detailed Description
**What was changed:**
{detailed_description_of_what_was_modified}

**Why it was changed:**
{reason_for_the_change}

**Impact:**
{what_this_change_affects_or_enables}

## Breaking Changes
- [ ] This is a breaking change
- [ ] This requires migration steps
- [ ] This changes public APIs

**If breaking, describe:**
{description_of_breaking_changes_and_migration_path}

## Related Issues
- Fixes #{issue_number}
- Closes #{issue_number}
- Relates to #{issue_number}

## Testing
**Tests added/modified:**
- {test_1_description}
- {test_2_description}

**Manual testing performed:**
- {manual_test_1}
- {manual_test_2}

## Request for Copilot
Please help me write a conventional commit message that:

1. **Follows the format:** `<type>[(scope)]: <description>`
2. **Has a clear, concise subject line** (≤50 characters)
3. **Includes a detailed body** explaining the what and why
4. **Lists any breaking changes** in the footer
5. **References related issues** appropriately

**Additional requirements:**
- Use imperative mood ("add" not "added")
- Capitalize the subject line
- Don't end subject with a period
- Separate subject, body, and footer with blank lines
- Wrap body at 72 characters

---

## Usage Instructions
1. Fill in all the `{placeholders}` with your specific change details
2. Check the relevant boxes for change type and scope
3. Provide detailed descriptions of what and why
4. Copy and paste into Copilot chat
5. Copilot will generate a properly formatted conventional commit message

## Example Usage
```
**Files modified:**
- api/v1alpha1/alertreaction_types.go: Added Priority field to AlertReactionSpec
- controllers/alertreaction_controller.go: Implemented priority-based alert processing
- controllers/alertreaction_controller_test.go: Added tests for priority handling

**Type of changes:**
- [x] New feature (feat)

**Scope:**
- [x] controller - Controller logic changes
- [x] api - CRD/API type changes

**What was changed:**
Added priority field to AlertReaction CRD and implemented priority-based processing in the controller where higher priority alerts (larger numbers) are processed before lower priority ones.

**Why it was changed:**
Users needed the ability to ensure critical alerts are processed immediately while less important alerts can wait in queue during high-load scenarios.

**Impact:**
This enables users to configure alert processing order based on business criticality, improving response times for critical issues.
```