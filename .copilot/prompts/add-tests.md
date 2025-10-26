# Add Tests Prompt Template

# Add Tests Prompt Template

## Prompt
```
I'm working on Karo (Kubernetes Alert Reaction Operator), a Kubernetes operator built with Go 1.24+ and controller-runtime. I need to add comprehensive tests following our project's testing patterns.

## Testing Target
**Function/Feature to test:** {function_or_feature_name}
**File location:** {file_path}
**Test file location:** {test_file_path}

### Code to Test
```go
{paste_the_code_that_needs_testing}
```

## Test Requirements
**Test type needed:**
- [ ] Unit tests (isolated function testing)
- [ ] Integration tests (component interaction testing)
- [ ] Controller tests (with fake client)
- [ ] Webhook tests (HTTP endpoint testing)
- [ ] End-to-end tests (full workflow testing)

**Test coverage should include:**
- [ ] Happy path scenarios
- [ ] Error conditions and edge cases
- [ ] Input validation
- [ ] Boundary conditions
- [ ] Concurrent access (if applicable)
- [ ] Resource cleanup
- [ ] Timeout handling
- [ ] Permission/RBAC scenarios

## Specific Test Scenarios
**Happy path scenarios:**
1. {scenario_1_description}
2. {scenario_2_description}
3. {scenario_3_description}

**Error scenarios:**
1. {error_scenario_1}
2. {error_scenario_2}
3. {error_scenario_3}

**Edge cases:**
1. {edge_case_1}
2. {edge_case_2}
3. {edge_case_3}

## Existing Test Patterns
**Similar tests in the codebase:**
- `{existing_test_file_1}` - {what_pattern_it_demonstrates}
- `{existing_test_file_2}` - {what_pattern_it_demonstrates}

**Test utilities available:**
- Fake Kubernetes client from controller-runtime
- Test logger setup
- Mock webhook server utilities
- Alert/Job test fixtures

## Dependencies and Mocks
**External dependencies to mock/fake:**
- [ ] Kubernetes API client
- [ ] HTTP requests/responses
- [ ] Time-based operations
- [ ] File system operations
- [ ] Environment variables

**Test data needed:**
- Sample AlertReaction resources
- Sample Alert payloads
- Expected Job specifications
- Error response examples

## Test Structure Preferences
**Preferred test style:**
- [ ] Table-driven tests (for multiple scenarios)
- [ ] Individual test functions (for complex scenarios)
- [ ] Testify suite (for setup/teardown needs)
- [ ] Ginkgo/Gomega (BDD style)

**Assertions framework:**
- [ ] Standard Go testing
- [ ] Testify assertions
- [ ] Gomega matchers

## Request for Copilot
Please help me create comprehensive tests that:

1. **Follow Go testing best practices**
2. **Use table-driven tests** where appropriate
3. **Include proper setup and teardown**
4. **Have clear test names** that describe what they test
5. **Cover all the scenarios** listed above
6. **Use appropriate mocks/fakes** for dependencies
7. **Include helpful assertions** with good error messages
8. **Are maintainable and readable**

**Specific requirements:**
- Use existing test utilities and patterns from the codebase
- Ensure tests are deterministic and can run in parallel
- Include both positive and negative test cases
- Add performance tests if relevant
- Document complex test scenarios with comments

## Expected Test Structure
```go
func Test{FunctionName}(t *testing.T) {
    tests := []struct {
        name           string
        input          {InputType}
        expectedOutput {OutputType}
        expectedError  string
        setupFunc      func()
        cleanupFunc    func()
    }{
        // Test cases here
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## Usage Instructions
1. Replace `{placeholders}` with your specific testing needs
2. Paste the actual code that needs testing
3. Specify the test scenarios you want to cover
4. Reference existing test patterns in the codebase
5. Copy and paste into Copilot chat

## Example Usage
```
**Function/Feature to test:** alertMatches function
**File location:** controllers/alertreaction_controller.go
**Test file location:** controllers/matcher_test.go

### Code to Test
```go
func (r *AlertReactionReconciler) alertMatches(alert Alert, matchers []AlertMatcher) bool {
    for _, matcher := range matchers {
        if !r.evaluateMatcher(alert, matcher) {
            return false
        }
    }
    return true
}
```

**Happy path scenarios:**
1. Alert matches all matchers with exact string matching
2. Alert matches all matchers with regex patterns
3. Empty matchers list (should match all alerts)

**Error scenarios:**
1. Alert missing required labels
2. Invalid regex patterns in matchers
3. Nil alert or matchers input
```