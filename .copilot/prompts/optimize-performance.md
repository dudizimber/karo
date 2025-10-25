# Optimize Performance Prompt Template

## Prompt
```
I'm working on performance optimization for Karo (Kubernetes Alert Reaction Operator), a Kubernetes controller that processes Prometheus alerts and creates Jobs. I need to identify and resolve performance bottlenecks.

## Current Performance Profile
**Observed performance issues:**
{describe_the_performance_problems_youre_experiencing}

**Performance metrics (if available):**
- **Alert processing latency:** {current_latency_p50_p95_p99}
- **Job creation time:** {time_from_alert_to_job_creation}
- **Memory usage:** {current_memory_consumption}
- **CPU usage:** {current_cpu_utilization}
- **API server requests:** {requests_per_second_or_per_alert}
- **Throughput:** {alerts_processed_per_second}

**Current load characteristics:**
- **Alert volume:** {alerts_per_minute_hour_day}
- **AlertReaction resources:** {number_of_alertreaction_resources}
- **Concurrent jobs:** {typical_number_of_running_jobs}
- **Cluster size:** {number_of_nodes_pods_resources}

## Performance Requirements
**Target performance goals:**
- **Latency:** {target_p95_alert_processing_time}
- **Throughput:** {target_alerts_per_second}
- **Resource usage:** {max_cpu_memory_constraints}
- **Scalability:** {must_handle_X_alertreactions_Y_alerts}

**SLA requirements:**
- [ ] Alert processing under {X} seconds 95% of the time
- [ ] System remains responsive under {Y} alerts/minute
- [ ] Memory usage stays under {Z} MB
- [ ] CPU usage stays under {W}% average

## Suspected Bottlenecks
**Areas of concern:**
- [ ] **Controller reconciliation loop** - {specific_performance_issues}
- [ ] **Alert matching logic** - {matcher_evaluation_performance}
- [ ] **Kubernetes API calls** - {excessive_api_requests}
- [ ] **Job creation overhead** - {job_creation_latency}
- [ ] **Memory allocations** - {gc_pressure_object_churn}
- [ ] **Webhook processing** - {http_request_handling}
- [ ] **Status updates** - {frequent_status_writes}

**Specific code areas:**
```go
{paste_suspected_performance_critical_code_sections}
```

## Profiling Data
**Available profiling information:**
- [ ] Go pprof CPU profile
- [ ] Go pprof memory profile
- [ ] Go pprof goroutine profile
- [ ] Kubernetes resource usage metrics
- [ ] Custom application metrics
- [ ] Distributed tracing data

**Profiling results (if available):**
```
{paste_pprof_output_or_performance_analysis_results}
```

## Optimization Areas
**Categories to investigate:**

### Algorithm and Logic Optimization
- [ ] **Matcher evaluation** - Optimize Prometheus-style matcher logic
- [ ] **Alert deduplication** - Reduce duplicate processing
- [ ] **Resource lookup** - Efficient AlertReaction finding
- [ ] **Batch processing** - Group operations where possible

### Kubernetes API Optimization
- [ ] **Client-side caching** - Reduce API server round trips
- [ ] **List/Watch efficiency** - Optimize resource watching
- [ ] **Batch operations** - Group API calls
- [ ] **Informer optimization** - Efficient cache usage

### Memory and GC Optimization
- [ ] **Object pooling** - Reuse frequently allocated objects
- [ ] **Memory allocation patterns** - Reduce GC pressure
- [ ] **Data structure efficiency** - Use optimal data structures
- [ ] **Memory leaks** - Identify and fix resource leaks

### Concurrency and Parallelization
- [ ] **Worker pool patterns** - Parallel alert processing
- [ ] **Channel optimization** - Efficient message passing
- [ ] **Lock contention** - Reduce synchronization overhead
- [ ] **Goroutine management** - Optimal concurrency levels

### I/O and Network Optimization
- [ ] **HTTP server tuning** - Webhook endpoint optimization
- [ ] **Connection pooling** - Efficient client connections
- [ ] **Request batching** - Reduce network overhead
- [ ] **Timeout optimization** - Appropriate timeout values

## Load Testing Strategy
**Test scenarios to create:**
1. **High-volume alert burst** - {X} alerts in {Y} seconds
2. **Sustained load** - {Z} alerts/minute for {W} hours
3. **Many AlertReactions** - {A} different AlertReaction resources
4. **Complex matchers** - AlertReactions with {B} complex regex matchers
5. **Resource constraint** - Limited CPU/memory environment

**Performance metrics to track:**
- [ ] Alert-to-Job latency percentiles (p50, p95, p99)
- [ ] CPU and memory usage over time
- [ ] API server request rate and latency
- [ ] Goroutine count and growth
- [ ] GC frequency and pause times
- [ ] Error rates and retry patterns

## Code Context
**Current implementation to optimize:**
```go
{paste_the_code_sections_that_need_optimization}
```

**Related configuration:**
```yaml
{paste_relevant_controller_configuration_deployment_specs}
```

## Request for Copilot
Please help me optimize the performance by:

1. **Analyzing the current implementation** for bottlenecks and inefficiencies
2. **Identifying specific optimization opportunities** in the code
3. **Suggesting algorithmic improvements** for better time/space complexity
4. **Recommending Kubernetes API usage patterns** for better performance
5. **Providing memory optimization strategies** to reduce GC pressure
6. **Implementing efficient concurrency patterns** for parallel processing
7. **Creating performance benchmarks** to measure improvements
8. **Suggesting monitoring and alerting** for performance regression detection

**Specific questions:**
1. What are the most impactful optimizations I should focus on first?
2. How can I reduce Kubernetes API server load while maintaining functionality?
3. What concurrency patterns would work best for alert processing?
4. How should I implement client-side caching for AlertReaction resources?
5. What profiling and monitoring should I add to track performance?

## Success Criteria
**Optimization will be successful if:**
- [ ] Alert processing latency reduced by {X}%
- [ ] System handles {Y}x more load without degradation
- [ ] Memory usage reduced by {Z}%
- [ ] API server requests reduced by {W}%
- [ ] No functional regressions introduced
- [ ] Performance improvements are measurable and sustained

---

## Usage Instructions
1. Replace `{placeholders}` with your specific performance data and requirements
2. Include actual profiling data and measurements where available
3. Paste relevant code sections that need optimization
4. Specify your performance targets and constraints
5. Copy and paste into Copilot chat

## Example Usage
```
**Observed performance issues:**
Alert processing takes 5-10 seconds during high load, causing delayed job creation and user complaints.

**Performance metrics:**
- Alert processing latency: p95 = 8.5s (target: <2s)
- Memory usage: 500MB steady state, 1.2GB during spikes
- API server requests: ~50 requests per alert processed
- Throughput: 10 alerts/minute max (need 100 alerts/minute)

**Suspected bottlenecks:**
- Controller reconciliation loop making too many API calls
- Alert matching logic using inefficient regex compilation
- Status updates happening too frequently
```