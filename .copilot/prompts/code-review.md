# Code Review Prompt Template

## Context
# Code Review Prompt Template

## Prompt
```
I'm working on Karo (Kubernetes Alert Reaction Operator) and need a comprehensive code review. The operator is built with Go 1.24+, controller-runtime, and follows Kubernetes controller patterns.

## Code Review Request
**Review type:**
- [ ] Pre-commit review (before committing changes)
- [ ] Pull request review (comprehensive review)
- [ ] Refactoring review (code quality improvement)
- [ ] Security review (security-focused analysis)
- [ ] Performance review (optimization opportunities)

**Files to review:**
```
{list_of_files_to_review}
```

## Changes Made
**Summary of changes:**
{brief_summary_of_what_was_changed}

**Type of changes:**
- [ ] New feature implementation
- [ ] Bug fix
- [ ] Performance optimization
- [ ] Security improvement
- [ ] Code refactoring
- [ ] Test additions
- [ ] Documentation updates

## Code to Review
**Primary changes:**
```go
{paste_the_main_code_changes_here}
```

**Supporting changes:**
```go
{paste_any_related_helper_functions_or_utilities}
```

**Test changes:**
```go
{paste_relevant_test_changes}
```

## Review Focus Areas
**Please pay particular attention to:**
- [ ] **Correctness** - Does the code do what it's supposed to do?
- [ ] **Error handling** - Are errors properly handled and propagated?
- [ ] **Resource management** - Are resources properly cleaned up?
- [ ] **Thread safety** - Is the code safe for concurrent access?
- [ ] **Performance** - Are there any performance bottlenecks?
- [ ] **Security** - Are there any security vulnerabilities?
- [ ] **Maintainability** - Is the code readable and maintainable?
- [ ] **Testing** - Is there adequate test coverage?

## Specific Concerns
**Areas I'm unsure about:**
1. {specific_concern_1}
2. {specific_concern_2}
3. {specific_concern_3}

**Questions for the reviewer:**
1. {question_1}
2. {question_2}
3. {question_3}

## Context and Constraints
**Business requirements:**
{describe_what_this_code_needs_to_accomplish}

**Technical constraints:**
- Must maintain backward compatibility
- Should follow existing code patterns
- Must be performant under high alert load
- Should integrate with existing monitoring/observability

**Dependencies and integrations:**
- Kubernetes API server interactions
- AlertManager webhook integration
- Job creation and lifecycle management
- Custom resource management

## Testing Coverage
**Tests included:**
- [ ] Unit tests for core logic
- [ ] Integration tests with fake Kubernetes client
- [ ] Error scenario testing
- [ ] Edge case coverage
- [ ] Performance/load testing

**Test coverage gaps (if any):**
{describe_any_areas_that_need_more_testing}

## Review Checklist
Please review for:

### Code Quality
- [ ] **Readability** - Clear variable names, good function structure
- [ ] **Documentation** - Adequate comments for complex logic
- [ ] **Consistency** - Follows project conventions and patterns
- [ ] **Simplicity** - Avoids unnecessary complexity

### Kubernetes Best Practices
- [ ] **Controller patterns** - Follows controller-runtime conventions
- [ ] **Resource management** - Proper use of owner references, finalizers
- [ ] **Status updates** - Appropriate status conditions and reporting
- [ ] **Event recording** - User-facing events for important operations

### Go Best Practices
- [ ] **Error handling** - Proper error wrapping and context
- [ ] **Memory management** - Avoids memory leaks and excessive allocations
- [ ] **Concurrency** - Safe concurrent access patterns
- [ ] **Interface usage** - Appropriate abstraction levels

### Security
- [ ] **Input validation** - Proper validation of user inputs
- [ ] **RBAC compliance** - Uses minimal required permissions
- [ ] **Secrets handling** - Secure handling of sensitive data
- [ ] **Attack surface** - Minimal exposure to potential attacks

### Performance
- [ ] **Efficiency** - Optimal algorithms and data structures
- [ ] **Resource usage** - Reasonable CPU and memory consumption
- [ ] **API usage** - Efficient Kubernetes API interactions
- [ ] **Scalability** - Handles growth in resources/load

## Request for Copilot
Please provide a comprehensive code review that covers:

1. **Overall architecture assessment** - Is the approach sound?
2. **Code quality analysis** - Readability, maintainability, patterns
3. **Bug identification** - Potential issues or edge cases
4. **Performance optimization** - Areas for improvement
5. **Security vulnerabilities** - Potential security issues
6. **Best practices compliance** - Go, Kubernetes, and project conventions
7. **Testing adequacy** - Coverage and quality of tests
8. **Improvement suggestions** - Specific recommendations

**Please be specific with:**
- Line-by-line feedback where appropriate
- Concrete suggestions for improvements
- Examples of better patterns or approaches
- Explanation of why changes are recommended

---

## Usage Instructions
1. Replace `{placeholders}` with your specific code and context
2. Paste the actual code you want reviewed
3. Specify your areas of concern and questions
4. Check the relevant review focus areas
5. Copy and paste into Copilot chat

## Example Usage
```
**Files to review:**
- controllers/alertreaction_controller.go (modified ProcessAlert method)
- controllers/priority_queue.go (new file for alert prioritization)

**Summary of changes:**
Added priority-based alert processing to handle high-priority alerts before low-priority ones.

**Primary changes:**
```go
func (r *AlertReactionReconciler) ProcessAlert(alert Alert) error {
    // Find matching AlertReactions
    reactions, err := r.findMatchingReactions(alert)
    if err != nil {
        return err
    }
    
    // Queue alerts by priority
    for _, reaction := range reactions {
        priority := reaction.Spec.Priority
        r.alertQueue.Enqueue(alert, priority)
    }
    
    return nil
}
```

**Areas I'm unsure about:**
1. Is the priority queue implementation thread-safe for concurrent webhook calls?
2. Should I persist the queue state or is in-memory sufficient?
3. Are there potential memory leaks if alerts aren't processed quickly enough?
```