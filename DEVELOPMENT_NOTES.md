# Development Notes

Last updated: 2026-03-13

## NEXT PRIORITY: Event Deduplication During Workflow Execution

### Goal
Demonstrate event deduplication when the reconcile loop (workflow) is still running. When multiple UPDATE events arrive for the same resource while the workflow is processing (sleeping), intermediate updates should be skipped and only the latest state should be processed.

### Demo Setup
- **Backend**: `cmd/example_two/main.go` - uses `NewApplictaionTwo()` workflow with a 5-second sleep
- **Event Generator**: `poetry run example_two` - sends rapid UPDATE events
- **Expected Behavior**:
  - Workflow starts processing event #1 (replicas: 2)
  - While sleeping (5s), events #2, #3, #4 arrive (replicas: 3, 4, 5)
  - Events #2, #3, #4 should be deduplicated
  - After sleep completes, workflow should process event #4 (latest state, replicas: 5)
  - Intermediate events #2 and #3 should be skipped

### Current Implementation Gap
The `EventProcessor` in `internal/runtime/http/main.go` does NOT currently deduplicate events. It processes all events sequentially from the queue without checking if newer events for the same resource have arrived.

### Implementation Approach
Need to implement deduplication logic in `EventProcessor`:
1. When processing events from the queue, group them by resource (kind + name)
2. For each resource, only keep the latest event (highest resourceVersion)
3. Discard intermediate UPDATE events
4. Add logging/metrics to show which events were skipped

### Test Scenario
```bash
# Terminal 1: Run backend with sleep workflow
LOG_LEVEL=DEBUG go run cmd/example_two/main.go

# Terminal 2: Send rapid events (1 second apart, workflow takes 5 seconds)
cd demo
poetry run example_two --replicas 1-5 --sleep 1

# Expected logs should show:
# - Event #1 (replicas: 2) starts processing
# - Events #2, #3, #4 arrive during 5s sleep
# - After sleep: Events #2, #3 skipped, Event #4 (replicas: 5) processed
```

---

## Recent Changes

### Node Update Fix (COMPLETED - 2026-03-13)
Fixed bug where UPDATE events for existing nodes were not updating the node's `Data` field. The `EventProcessor` in `internal/runtime/http/main.go` now correctly updates `existing.Data` when a node already exists, and logs `"updating existing node"` with kind and name.

### Replica Pointer Dereference Fix (COMPLETED)
Fixed a separate bug in `internal/applicaton/actions.go:21` where `deployment.Spec.Replicas` was being logged as a pointer address instead of the actual value. Now correctly dereferences the `*int32` pointer.

### Demo Tool Refactoring (COMPLETED)
- Moved `demo/templates/` → `demo/src/crystal/templates/`
- Moved `demo/mock_deployment_updates.py` → `demo/src/crystal/scripts/mock_deployment_updates.py`
- Added Poetry script entry point: `example_two`
- Updated all README files to reflect new structure and usage

### Mock Event Generator (COMPLETED)
Created a Python-based mock Kubernetes admission webhook event generator with:
- Realistic Deployment UPDATE events
- Auto-incrementing replica counts, resourceVersion, and generation
- Extensible template system for future ConfigMap/Secret support
- CLI with multiple configuration options

**Usage:**
```bash
cd demo
poetry run example_two --replicas 1-5 --sleep 2
```

## Next Steps

### Phase 1: Fix Node Updates (COMPLETED)
Node updates now work correctly. See "Node Update Fix" in Recent Changes.

### Phase 2: Implement Event Deduplication (Main Goal)
1. **Modify EventProcessor** to deduplicate events by resource (kind + name)
2. Keep only the latest event per resource based on resourceVersion
3. Add metrics for:
   - Total events received
   - Events skipped due to deduplication
   - Events processed
4. Add logging to show which intermediate events were skipped
5. Test deduplication demo:
   - Run `cmd/example_two/main.go` (5s sleep workflow)
   - Send rapid events with `poetry run example_two --replicas 1-10 --sleep 1`
   - Verify intermediate updates are skipped in logs

### Phase 3: Visualization (Optional)
1. Add frontend visualization showing:
   - Events received vs events processed
   - Deduplication statistics
   - Timeline of event arrivals vs workflow execution

## Project Structure

```
crystal/
├── cmd/
│   ├── example_one/main.go    # Main backend server (port 8000)
│   └── example_two/main.go    # Alternative backend example
├── internal/
│   ├── applicaton/
│   │   ├── main.go           # Workflow definitions
│   │   └── actions.go        # Workflow action functions (alan, tom, sleepy, etc.)
│   ├── runtime/http/
│   │   └── main.go           # HTTP handlers and EventProcessor
│   ├── workflow/
│   │   └── nodes.go          # Node and NodeList implementation
│   └── api/
│       └── deployment.go     # Deployment struct definitions
└── demo/
    ├── src/crystal/
    │   ├── templates/         # K8s resource templates
    │   └── scripts/
    │       └── mock_deployment_updates.py  # Mock event generator
    └── app.py                # Streamlit frontend

```

## Useful Commands

**Run backend (example_one - no sleep):**
```bash
LOG_LEVEL=DEBUG go run cmd/example_one/main.go
```

**Run backend (example_two - with 5s sleep for deduplication demo):**
```bash
LOG_LEVEL=DEBUG go run cmd/example_two/main.go
```

**Run frontend:**
```bash
cd demo
poetry run streamlit run app.py
```

**Generate mock events:**
```bash
cd demo
poetry run example_two --replicas 1-5 --sleep 2 --verbose
```

**Check event list:**
```bash
curl http://localhost:8000/list | jq
```
