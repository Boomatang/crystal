# Event Deduplication Demo Design

**Date:** 2026-03-13
**Status:** Approved

## Summary

An integrated Streamlit demo page that demonstrates the EventQueue deduplication feature. Users can trigger burst events and observe how the queue deduplicates based on resourceVersion comparison.

## Goals

1. Verify deduplication works correctly (personal testing)
2. Demonstrate the feature to colleagues with a polished UI
3. Show both sequential and out-of-order event scenarios

## Design

### Page Structure

The demo is a new page in the existing Streamlit app with three areas:

**1. Controls (Sidebar)**
- Mode selector: "Sequential" / "Out-of-Order"
- Deployment name (text input, default: "demo-app")
- Namespace (text input, default: "default")
- Number of events to send (slider: 5-20, default: 10)
- Delay between events in ms (slider: 0-500ms, default: 50ms)
- "Send Burst" button

**2. Event Log (Left Column)**
- Shows each event sent with timestamp and resourceVersion
- Status indicator inferred from queue state:
  - Added (new resource)
  - Replaced (newer version)
  - Skipped (older version)
- Running count of events sent

**3. Queue State (Right Column)**
- Displays current queue contents from `/list` endpoint
- Shows resourceVersion held for each resource key
- Updates after each event is sent

**4. Summary Stats (Bottom)**
- Events sent: N
- Added: X
- Replaced: Y
- Skipped: Z
- Queue size: 1

### Demo Modes

**Mode A: Sequential**
- Sends events with incrementing resourceVersions: 100, 101, 102...
- Demonstrates deduplication via replacement
- All events except first show "Replaced"

**Mode B: Out-of-Order**
- Sends events in scrambled order: 105, 102, 108, 101, 103...
- Demonstrates rejection of older events
- Mix of "Added", "Replaced", and "Skipped" statuses

### Status Inference Logic

Since the Go server doesn't return explicit add/replace/skip status, infer it by:

1. Before sending: fetch current queue state via `/list`
2. Send the event
3. After sending: fetch queue state again
4. Compare:
   - Resource not in queue before, now present → "Added"
   - Resource in queue before with lower rv, now has new rv → "Replaced"
   - Resource in queue before with same/higher rv, unchanged → "Skipped"

## Files Changed

| File | Change |
|------|--------|
| `demo/src/crystal/__init__.py` | Add `send_event()` and `fetch_event_list()` functions |
| `demo/app.py` | Add deduplication demo page |

## New Functions

### `send_event(url, deployment_name, namespace, resource_version) -> dict`

Sends an AdmissionReview event to the server.

**Parameters:**
- `url`: Base server URL (e.g., "http://localhost:8000")
- `deployment_name`: Name of the deployment resource
- `namespace`: Kubernetes namespace
- `resource_version`: Integer resourceVersion for the event

**Returns:** Response dict from server

### `fetch_event_list(url) -> list`

Fetches current queue state from `/list` endpoint.

**Parameters:**
- `url`: Base server URL

**Returns:** List of AdmissionReview events currently in queue

## Implementation Tasks

- [x] [Add send_event() function to demo client](https://github.com/Boomatang/crystal/issues/6)
- [x] [Add fetch_event_list() function to demo client](https://github.com/Boomatang/crystal/issues/7)
- [x] [Add deduplication demo page to Streamlit app](https://github.com/Boomatang/crystal/issues/8)

## Out of Scope

- Server-side changes to return explicit status
- Log streaming from Go server
- Multiple resource types (Deployment only)
- Persistence of demo results

## Change Log

| Date | Change |
|------|--------|
| 2026-03-13 | Initial design |
| 2026-03-13 | Added Implementation Tasks section with linked GitHub issues #6-#8 |
| 2026-03-13 | Implementation complete - all tasks done |
