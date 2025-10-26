# Controller Logic Prompt Template

## Context
# Controller Logic Prompt Template

## Prompt
```
I'm working on the Karo (Kubernetes Alert Reaction Operator) controller logic. The controller uses controller-runtime to reconcile AlertReaction resources and creates Jobs based on incoming Prometheus alerts via webhooks.

## Current Controller Architecture
**Main reconciliation function:** `Reconcile(ctx context.Context, req ctrl.Request)`
**Key methods:**
- `ProcessAlert(alert Alert)` - Processes incoming webhook alerts
- `alertMatches(alert Alert, matchers []AlertMatcher)` - Evaluates alert matching
- `createJobForAction(alertReaction *AlertReaction, action Action, alert Alert)` - Creates Jobs

**Current workflow:**
1. Webhook receives alert from AlertManager
2. Controller finds matching AlertReaction resources
3. Evaluates matchers against alert labels/annotations
4. Creates Jobs for matching actions
5. Updates AlertReaction status

## Change Requirements
**What needs to be modified:**
{detailed_description_of_controller_changes}

**Specific areas:**
- [ ] Reconciliation logic
- [ ] Alert processing workflow
- [ ] Job creation and management
- [ ] Status updates and conditions
- [ ] Error handling and retries
- [ ] Finalizer handling
- [ ] Watch filters and predicates

## Functional Requirements
**New behavior needed:**
{describe_the_new_behavior_the_controller_should_have}

**Integration points:**
- [ ] Kubernetes API interactions
- [ ] AlertReaction CRD handling
- [ ] Job lifecycle management
- [ ] Status condition updates
- [ ] Event recording
- [ ] Metrics collection

## Technical Details
**Controller-runtime components involved:**
- [ ] Manager setup and configuration
- [ ] Reconciler implementation
- [ ] Client operations (Get, List, Create, Update, Patch)
- [ ] Status subresource updates
- [ ] Owner references and garbage collection
- [ ] Watches and event filtering

**Kubernetes resources to manage:**
- [ ] AlertReaction (primary resource)
- [ ] Jobs (created/managed by controller)
- [ ] ConfigMaps/Secrets (if needed)
- [ ] Events (for user feedback)

**Error scenarios to handle:**
- [ ] Resource creation failures
- [ ] Permission/RBAC errors
- [ ] Resource conflicts and retries
- [ ] Invalid AlertReaction configurations
- [ ] Webhook processing errors

## Current Code Context
**Existing reconciliation logic:**
```go
{paste_current_reconcile_function_or_relevant_code}
```

**Related helper functions:**
```go
{paste_relevant_helper_functions}
```

## Expected Changes
**New/modified functions needed:**
- `{function_name}` - {purpose_and_behavior}
- `{function_name}` - {purpose_and_behavior}

**Controller setup changes:**
{describe_any_changes_to_controller_setup_or_watches}

**Status updates:**
{describe_how_the_status_should_be_updated}

## Testing Requirements
**Unit tests needed:**
- [ ] Reconciliation logic with various AlertReaction configurations
- [ ] Error handling scenarios
- [ ] Status update correctness
- [ ] Job creation logic

**Integration tests needed:**
- [ ] Full reconciliation cycles with fake client
- [ ] Multi-resource scenarios
- [ ] Garbage collection behavior
- [ ] Controller restart scenarios

## Performance Considerations
**Scalability requirements:**
- [ ] Handle multiple AlertReaction resources efficiently
- [ ] Process high-volume alert streams
- [ ] Minimize API server requests
- [ ] Proper resource cleanup

**Resource usage:**
- [ ] Memory usage for large numbers of resources
- [ ] CPU usage during high alert activity
- [ ] API server load considerations

## Request for Copilot
Please help me implement these controller changes by:

1. **Analyzing the current controller structure** and identifying optimal modification points
2. **Implementing the new reconciliation logic** following controller-runtime best practices
3. **Adding proper error handling** with appropriate retries and backoff
4. **Updating status conditions** to reflect the new functionality
5. **Creating comprehensive tests** for all scenarios
6. **Ensuring thread safety** and proper resource cleanup
7. **Following Kubernetes controller patterns** and conventions

**Specific questions:**
1. What's the best way to structure the new reconciliation logic?
2. How should I handle race conditions and concurrent reconciliations?
3. What status conditions should I add for the new functionality?
4. How can I make this change backward compatible?
5. What are the performance implications of this change?

## Error Handling Strategy
**Retry scenarios:**
- [ ] Transient API server errors
- [ ] Resource creation conflicts
- [ ] Network timeouts

**Non-retry scenarios:**
- [ ] Invalid AlertReaction configurations
- [ ] Permission denied errors
- [ ] Resource quota exceeded

**Status reporting:**
- [ ] Success conditions
- [ ] Error conditions with helpful messages
- [ ] Progress indicators for long-running operations

---

## Usage Instructions
1. Replace `{placeholders}` with your specific controller change requirements
2. Include relevant existing code for context
3. Describe the expected behavior and integration points
4. Specify error handling and testing requirements
5. Copy and paste into Copilot chat

## Example Usage
```
**What needs to be modified:**
Add support for alert prioritization where high-priority alerts are processed before low-priority ones.

**New behavior needed:**
The controller should maintain a priority queue of alerts and process them in order of priority (highest first). When multiple alerts arrive simultaneously, they should be queued and processed based on their priority value from the AlertReaction configuration.

**Integration points:**
- Webhook handler needs to queue alerts by priority instead of processing immediately
- Reconciler needs to process alerts from the priority queue
- Status should show queue length and processing statistics
```