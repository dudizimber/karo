# Refactor Code Prompt Template

## Context
# Refactor Code Prompt Template

## Prompt
```
I'm working on refactoring code in Karo (Kubernetes Alert Reaction Operator) to improve maintainability, readability, and code quality. The operator is built with Go 1.24+ and controller-runtime.

## Refactoring Target
**Code to refactor:**
{describe_the_specific_code_area_or_function_to_refactor}

**File(s) involved:**
- `{primary_file_path}` - {what_needs_refactoring_here}
- `{secondary_file_path}` - {related_changes_needed}

**Current implementation:**
```go
{paste_the_current_code_that_needs_refactoring}
```

## Refactoring Goals
**Primary objectives:**
- [ ] **Improve readability** - Make code easier to understand
- [ ] **Reduce complexity** - Simplify complex functions or logic
- [ ] **Eliminate duplication** - Remove repeated code patterns
- [ ] **Improve testability** - Make code easier to unit test
- [ ] **Better error handling** - More robust error management
- [ ] **Performance improvement** - Optimize without changing behavior
- [ ] **Follow Go idioms** - Align with Go best practices
- [ ] **Enhance maintainability** - Easier to modify and extend

**Specific improvements wanted:**
1. {specific_improvement_1}
2. {specific_improvement_2}
3. {specific_improvement_3}

## Current Issues
**Problems with the existing code:**
- [ ] **Function too long** - {function_name} is {X} lines, should be broken down
- [ ] **Too many responsibilities** - Function does multiple unrelated things
- [ ] **Hard to test** - Difficult to write unit tests due to tight coupling
- [ ] **Unclear naming** - Variable/function names don't reflect their purpose
- [ ] **Complex conditionals** - Nested if statements or complex boolean logic
- [ ] **Error handling scattered** - Inconsistent error handling patterns
- [ ] **Code duplication** - Similar logic repeated in multiple places
- [ ] **Poor abstraction** - Missing interfaces or appropriate abstractions

**Technical debt:**
{describe_any_technical_debt_that_should_be_addressed}

**Code smells:**
- [ ] Long method
- [ ] Large class/struct
- [ ] Duplicate code
- [ ] Long parameter list
- [ ] Feature envy
- [ ] Data clumps
- [ ] Primitive obsession
- [ ] Switch statements (in Go context)

## Constraints
**Must maintain:**
- [ ] **Exact same functionality** - No behavior changes
- [ ] **Backward compatibility** - Existing interfaces unchanged
- [ ] **Performance characteristics** - No performance regression
- [ ] **Error handling behavior** - Same error conditions and responses
- [ ] **Logging output** - Consistent with existing log messages

**Cannot change:**
- [ ] Public API interfaces
- [ ] Configuration file formats
- [ ] Database schema (if applicable)
- [ ] External integration contracts

## Design Patterns and Principles
**Apply these patterns:**
- [ ] **Single Responsibility Principle** - Each function/struct has one job
- [ ] **Dependency Injection** - Make dependencies explicit and testable
- [ ] **Interface Segregation** - Small, focused interfaces
- [ ] **Factory Pattern** - For complex object creation
- [ ] **Strategy Pattern** - For interchangeable algorithms
- [ ] **Template Method** - For common workflows with variations

**Go-specific patterns:**
- [ ] **Error wrapping** - Use `fmt.Errorf` with `%w` verb
- [ ] **Context propagation** - Pass context through call chains
- [ ] **Channel patterns** - For concurrent communication
- [ ] **Interface satisfaction** - Small, focused interfaces
- [ ] **Functional options** - For optional parameters

## Proposed Structure
**New structure/organization:**
```go
{describe_or_sketch_the_refactored_structure}
```

**New functions/methods to extract:**
- `{function_name}` - {responsibility_and_signature}
- `{function_name}` - {responsibility_and_signature}

**New types/interfaces to introduce:**
- `{type_name}` - {purpose_and_methods}
- `{interface_name}` - {contract_definition}

## Testing Strategy
**Current test coverage:**
{describe_existing_tests_and_coverage}

**Testing improvements needed:**
- [ ] **Unit tests for extracted functions** - Test smaller, focused units
- [ ] **Mock interfaces** - Enable testing with mocked dependencies  
- [ ] **Table-driven tests** - Cover multiple scenarios efficiently
- [ ] **Error path testing** - Ensure error conditions are well tested
- [ ] **Integration tests** - Verify components work together

**Test refactoring needed:**
{describe_how_tests_should_change_to_match_refactored_code}

## Related Code
**Similar patterns in codebase:**
{reference_other_code_that_follows_good_patterns_or_needs_similar_refactoring}

**Dependencies that might be affected:**
- {dependency_1} - {how_it_might_be_impacted}
- {dependency_2} - {how_it_might_be_impacted}

## Request for Copilot
Please help me refactor this code by:

1. **Analyzing the current implementation** and identifying specific refactoring opportunities
2. **Breaking down large functions** into smaller, focused units
3. **Extracting common patterns** into reusable functions or types
4. **Improving error handling** with consistent patterns and proper wrapping
5. **Adding appropriate abstractions** through interfaces where beneficial
6. **Optimizing for testability** by reducing coupling and dependencies
7. **Following Go best practices** and idiomatic patterns
8. **Maintaining backward compatibility** while improving internal structure

**Specific questions:**
1. What's the best way to break down the large function while maintaining readability?
2. What interfaces should I introduce to make this more testable?
3. How can I eliminate the code duplication without over-engineering?
4. What's the most Go-idiomatic way to structure this functionality?
5. How should I handle the error cases more consistently?

## Success Criteria
**Refactoring will be successful if:**
- [ ] **Code is more readable** - Easier for new developers to understand
- [ ] **Functions are focused** - Each function has a single, clear responsibility
- [ ] **Tests are comprehensive** - Better coverage with cleaner test code
- [ ] **No functionality changes** - All existing behavior preserved
- [ ] **Performance maintained** - No regression in performance metrics
- [ ] **Technical debt reduced** - Fewer code smells and maintainability issues

---

## Usage Instructions
1. Replace `{placeholders}` with your specific refactoring context
2. Paste the actual code that needs refactoring
3. Clearly describe the problems and desired improvements
4. Specify constraints and requirements
5. Copy and paste into Copilot chat

## Example Usage
```
**Code to refactor:**
The ProcessAlert function in controllers/alertreaction_controller.go is 150 lines long and handles alert matching, job creation, status updates, and error handling all in one function.

**Current implementation:**
```go
func (r *AlertReactionReconciler) ProcessAlert(alert Alert) error {
    // Find all AlertReaction resources
    var alertReactions AlertReactionList
    if err := r.List(ctx, &alertReactions); err != nil {
        return err
    }
    
    // Process each AlertReaction (50+ lines of complex logic)
    for _, ar := range alertReactions.Items {
        // Complex matching logic (20+ lines)
        // Job creation logic (30+ lines)  
        // Status update logic (20+ lines)
        // Error handling scattered throughout
    }
    
    return nil
}
```

**Primary objectives:**
- Break down into smaller, testable functions
- Separate alert matching from job creation
- Improve error handling consistency
- Make each piece independently testable
```