# Add Feature Prompt Template

## Context
I'm working on the k8s-alert-reaction-operator, a Kubernetes operator that creates Jobs in response to Prometheus alerts. The operator uses:
- Go 1.24+ with controller-runtime framework
- AlertReaction CRD with Prometheus-style matchers
- Webhook endpoint for receiving AlertManager notifications
- Job creation based on alert conditions

## Feature Request
I need to add **{feature_name}** to the operator.

### Requirements
**Primary Functionality:**
{detailed_feature_description}

**Acceptance Criteria:**
- [ ] {acceptance_criteria_1}
- [ ] {acceptance_criteria_2}
- [ ] {acceptance_criteria_3}

**Technical Constraints:**
- Must maintain backward compatibility with existing AlertReaction resources
- Should follow existing code patterns and conventions
- Must include comprehensive tests
- Should update relevant documentation

### Files Likely Needing Changes
- [ ] `api/v1alpha1/alertreaction_types.go` - API types and CRD structure
- [ ] `controllers/alertreaction_controller.go` - Controller reconciliation logic
- [ ] `internal/webhook/webhook.go` - Webhook processing logic
- [ ] `config/crd/` - Generated CRD manifests
- [ ] Test files - Comprehensive test coverage
- [ ] Documentation - README, examples, etc.

## Implementation Plan
Please help me implement this feature by:

1. **API Changes**: Update the AlertReaction CRD with new fields/structures
2. **Controller Logic**: Modify the reconciliation logic to handle the new feature
3. **Webhook Updates**: Update webhook processing if needed
4. **Validation**: Add proper field validation and webhooks
5. **Testing**: Create comprehensive unit and integration tests
6. **Examples**: Provide usage examples and documentation
7. **Migration**: Consider any migration needs for existing resources

## Similar Patterns
{reference_to_existing_similar_code_if_any}

## Expected Behavior
**Before:** {current_behavior}
**After:** {expected_new_behavior}

## Questions for Copilot
1. What's the best approach to implement this feature while maintaining backward compatibility?
2. Are there any security implications I should consider?
3. What edge cases should I test for?
4. How should this integrate with existing Prometheus matcher functionality?

---

## Usage Instructions
1. Replace `{feature_name}` with your specific feature name
2. Fill in `{detailed_feature_description}` with comprehensive requirements
3. List specific `{acceptance_criteria_*}` for the feature
4. Reference any `{existing_similar_code}` if applicable
5. Describe `{current_behavior}` and `{expected_new_behavior}`
6. Copy and paste into Copilot chat

## Example Usage
```
I need to add **alert priority handling** to the operator.

### Requirements
**Primary Functionality:**
Add support for priority-based alert processing where high-priority alerts are processed before low-priority ones.

**Acceptance Criteria:**
- [ ] AlertReaction CRD supports priority field (integer 1-10, default 5)
- [ ] Controller processes alerts in priority order (higher numbers first)
- [ ] Webhook queues alerts by priority
- [ ] Backward compatibility maintained for existing AlertReactions
```