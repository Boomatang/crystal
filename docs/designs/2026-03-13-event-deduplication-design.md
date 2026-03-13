# Event Deduplication Design

**Date:** 2026-03-13
**Status:** Draft

## Problem Statement

When the reconcile workflow is running (e.g., during the 5-second sleep in `example_two`), multiple UPDATE events can arrive for the same resource. Currently, the `EventProcessor` processes all events sequentially, including intermediate states that are immediately superseded. This wastes processing time and can lead to inconsistent behavior.

## Goal

Deduplicate events so that only the latest state per resource is processed. When multiple events arrive for the same resource (identified by `kind/namespace/name`), keep only the event with the highest `resourceVersion`.

## Design

### Approach

Introduce a dedicated `EventQueue` type that encapsulates deduplication logic. The queue:
- Stores events in a map keyed by resource identity (`kind/namespace/name`)
- Deduplicates at ingestion time using `resourceVersion` comparison
- Provides thread-safe access via mutex
- Supports atomic drain-and-copy for batch processing

### New Type: EventQueue

**Location:** `internal/workflow/event_queue.go`

```go
type EventQueue struct {
    mu     sync.Mutex
    events map[string]AdmissionReview  // key: "kind/namespace/name"
}

func NewEventQueue() *EventQueue

func (q *EventQueue) Add(event AdmissionReview) (added bool, replaced bool)

func (q *EventQueue) DrainAndCopy() []AdmissionReview

func (q *EventQueue) Len() int

func (q *EventQueue) Snapshot() []AdmissionReview
```

**Key behavior:**
- Map key format: `fmt.Sprintf("%s/%s/%s", event.Request.Kind.Kind, event.Request.Namespace, event.Request.Name)`
- `Add()` compares resourceVersion of incoming vs existing event, keeps the higher one
- `Add()` returns `(added, replaced)` flags for logging:
  - `(true, false)` — new resource added
  - `(true, true)` — replaced older event
  - `(false, false)` — rejected (older or invalid resourceVersion)
- `DrainAndCopy()` locks, copies all values to a slice, clears the map, unlocks, returns slice
- `Len()` acquires mutex before reading map length (thread-safe)
- `Snapshot()` returns copy of current events without draining (for ListEventsHandler)

### resourceVersion Parsing

**Helper function:**

```go
func getResourceVersion(event AdmissionReview) (int64, error)
```

**Logic:**
1. Access `event.Request.Object["metadata"]` using comma-ok idiom
2. Type assert to `map[string]any` using comma-ok idiom
3. Extract `resourceVersion` field, handling both:
   - String type (expected): parse with `strconv.ParseInt`
   - Numeric type (`float64` from JSON): convert to int64
4. Return descriptive error if:
   - `Object` is nil
   - `metadata` key missing or wrong type
   - `resourceVersion` key missing or wrong type
   - String value unparseable

**Error handling in `Add()`:**
- Incoming event unparseable → reject (return `added=false, replaced=false`)
- Existing event unparseable → reject incoming, log warning (we can't safely determine ordering)
- Both parseable → keep higher resourceVersion

### Integration Changes

**`internal/runtime/http/main.go`:**

Remove:
- `var events []workflow.AdmissionReview`
- `var mu sync.Mutex`

Update `EventHandler`:
```go
func EventHandler(c chan bool, queue *workflow.EventQueue) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... decode event ...
        added, replaced := queue.Add(event)
        if added {
            c <- true
        }
        // Log add/replace/skip with kind, name, resourceVersion
    }
}
```

Update `EventProcessor`:
```go
func EventProcessor(world *workflow.World, nodes *workflow.NodeList, queue *workflow.EventQueue) {
    for range EventChan {  // EventChan remains a package-level variable
        load := queue.DrainAndCopy()
        // ... process events ...
    }
}
```

Update `ListEventsHandler`:
```go
func ListEventsHandler(queue *workflow.EventQueue) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(queue.Snapshot())
    }
}
```

**`cmd/example_two/main.go`:**
```go
queue := workflow.NewEventQueue()
go rhttp.EventProcessor(worldMain, nodes, queue)
mux.Handle("/event", rhttp.EventHandler(rhttp.EventChan, queue))
mux.Handle("/list", rhttp.ListEventsHandler(queue))
```

### Logging

Log the following with structured fields (`kind`, `namespace`, `name`, `resourceVersion`):
- `"event added"` — new resource, first event
- `"event replaced older version"` — newer event replaced existing (include `existing_resourceVersion`)
- `"event skipped (older than existing)"` — incoming event was older (include `existing_resourceVersion`)
- `"event rejected (invalid resourceVersion)"` — unparseable resourceVersion
- `"draining event queue"` — at start of DrainAndCopy (include `count`)

### Testing

**File:** `internal/workflow/event_queue_test.go`

| Test Case | Description |
|-----------|-------------|
| Add new event | Empty queue, add event, verify `Len() == 1`, `added=true, replaced=false` |
| Add newer event | Add resourceVersion "100", then "200", verify "200" kept, `added=true, replaced=true` |
| Reject older event | Add "200", then "100", verify "200" kept, `added=false, replaced=false` |
| Multiple resources | Add events for different kind/namespace/name, verify all kept separately |
| Same name different namespace | Add `Deployment/ns1/app` and `Deployment/ns2/app`, verify both kept |
| DrainAndCopy clears | Add events, drain, verify `Len() == 0`, slice has all events |
| Snapshot preserves | Add events, snapshot, verify `Len()` unchanged, slice has all events |
| Concurrent access | 10 goroutines adding 100 events each (mix of same/different keys), run with `-race`, verify no panics and final state consistent |
| Invalid resourceVersion string | Missing or unparseable string resourceVersion rejected |
| Numeric resourceVersion | resourceVersion as float64 (JSON number) handled correctly |
| Existing unparseable | When existing event has invalid resourceVersion, incoming is rejected with warning |

## Implementation Tasks

- [ ] [Create EventQueue type with deduplication logic](https://github.com/Boomatang/crystal/issues/1)
  - [ ] Implement NewEventQueue constructor
  - [ ] Implement Add method with resourceVersion comparison
  - [ ] Implement getResourceVersion helper function
  - [ ] Implement DrainAndCopy method
  - [ ] Implement Len method
  - [ ] Implement Snapshot method
  - [ ] Add structured logging

- [ ] [Create EventQueue unit tests](https://github.com/Boomatang/crystal/issues/2)
  - [ ] Test: Add new event
  - [ ] Test: Add newer event (replace)
  - [ ] Test: Reject older event
  - [ ] Test: Multiple resources
  - [ ] Test: Same name different namespace
  - [ ] Test: DrainAndCopy clears queue
  - [ ] Test: Snapshot preserves queue
  - [ ] Test: Concurrent access with -race
  - [ ] Test: Invalid resourceVersion string
  - [ ] Test: Numeric resourceVersion
  - [ ] Test: Existing unparseable resourceVersion

- [ ] [Update http handlers to use EventQueue](https://github.com/Boomatang/crystal/issues/3)
  - [ ] Remove global `events` slice and `mu` mutex
  - [ ] Update EventHandler to use EventQueue.Add
  - [ ] Update EventProcessor to use EventQueue.DrainAndCopy
  - [ ] Update ListEventsHandler to use EventQueue.Snapshot

- [ ] [Update example_two to use EventQueue](https://github.com/Boomatang/crystal/issues/4)
  - [ ] Create EventQueue instance
  - [ ] Pass queue to EventProcessor
  - [ ] Pass queue to EventHandler
  - [ ] Pass queue to ListEventsHandler

- [ ] [Update example_one to use EventQueue](https://github.com/Boomatang/crystal/issues/5)
  - [ ] Same changes as example_two

## Files Changed

| File | Change |
|------|--------|
| `internal/workflow/event_queue.go` | New file: `EventQueue` type |
| `internal/workflow/event_queue_test.go` | New file: unit tests |
| `internal/runtime/http/main.go` | Remove global slice/mutex, update handlers |
| `cmd/example_two/main.go` | Create and pass `EventQueue` instance |
| `cmd/example_one/main.go` | Same changes as example_two |

## Out of Scope

- Metrics (to be added later)
- Visualization of deduplication statistics
- DELETE event handling (currently only UPDATE events are sent)

## Change Log

| Date | Change |
|------|--------|
| 2026-03-13 | Initial design |
| 2026-03-13 | Added namespace to deduplication key, Snapshot() method, improved error handling, expanded test cases |
| 2026-03-13 | Added Implementation Tasks section with TODO checkboxes |
| 2026-03-13 | Created GitHub issues #1-#5 and linked TODOs |
