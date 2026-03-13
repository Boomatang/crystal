package workflow

import (
	"sync"
	"testing"
)

// makeTestEvent creates a test AdmissionReview with string resourceVersion
func makeTestEvent(kind, namespace, name, resourceVersion string) AdmissionReview {
	return AdmissionReview{
		Request: &AdmissionRequest{
			Kind: GroupVersionKind{
				Kind: kind,
			},
			Namespace: namespace,
			Name:      name,
			Object: map[string]any{
				"metadata": map[string]any{
					"resourceVersion": resourceVersion,
				},
			},
		},
	}
}

// makeTestEventNumericRV creates a test AdmissionReview with numeric resourceVersion (float64)
func makeTestEventNumericRV(kind, namespace, name string, resourceVersion float64) AdmissionReview {
	return AdmissionReview{
		Request: &AdmissionRequest{
			Kind: GroupVersionKind{
				Kind: kind,
			},
			Namespace: namespace,
			Name:      name,
			Object: map[string]any{
				"metadata": map[string]any{
					"resourceVersion": resourceVersion,
				},
			},
		},
	}
}

// TestAddNewEvent tests adding an event to an empty queue
func TestAddNewEvent(t *testing.T) {
	q := NewEventQueue()

	event := makeTestEvent("Deployment", "default", "my-app", "100")
	added, replaced := q.Add(event)

	if !added {
		t.Errorf("expected added=true, got added=%v", added)
	}
	if replaced {
		t.Errorf("expected replaced=false, got replaced=%v", replaced)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1, got %d", q.Len())
	}
}

// TestAddNewerEvent tests adding a newer event (higher resourceVersion) replaces the older one
func TestAddNewerEvent(t *testing.T) {
	q := NewEventQueue()

	// Add initial event
	event1 := makeTestEvent("Deployment", "default", "my-app", "100")
	q.Add(event1)

	// Add newer event
	event2 := makeTestEvent("Deployment", "default", "my-app", "200")
	added, replaced := q.Add(event2)

	if !added {
		t.Errorf("expected added=true, got added=%v", added)
	}
	if !replaced {
		t.Errorf("expected replaced=true, got replaced=%v", replaced)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1, got %d", q.Len())
	}

	// Verify the newer event is kept
	events := q.Snapshot()
	if len(events) != 1 {
		t.Fatalf("expected 1 event in snapshot, got %d", len(events))
	}
	rv, err := getResourceVersion(events[0])
	if err != nil {
		t.Fatalf("failed to get resourceVersion: %v", err)
	}
	if rv != 200 {
		t.Errorf("expected resourceVersion=200, got %d", rv)
	}
}

// TestRejectOlderEvent tests that older events are rejected
func TestRejectOlderEvent(t *testing.T) {
	q := NewEventQueue()

	// Add initial event with higher version
	event1 := makeTestEvent("Deployment", "default", "my-app", "200")
	q.Add(event1)

	// Try to add older event
	event2 := makeTestEvent("Deployment", "default", "my-app", "100")
	added, replaced := q.Add(event2)

	if added {
		t.Errorf("expected added=false, got added=%v", added)
	}
	if replaced {
		t.Errorf("expected replaced=false, got replaced=%v", replaced)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1, got %d", q.Len())
	}

	// Verify the newer event is still kept
	events := q.Snapshot()
	if len(events) != 1 {
		t.Fatalf("expected 1 event in snapshot, got %d", len(events))
	}
	rv, err := getResourceVersion(events[0])
	if err != nil {
		t.Fatalf("failed to get resourceVersion: %v", err)
	}
	if rv != 200 {
		t.Errorf("expected resourceVersion=200, got %d", rv)
	}
}

// TestMultipleResources tests that events for different resources are kept separately
func TestMultipleResources(t *testing.T) {
	q := NewEventQueue()

	// Add events for different kinds
	event1 := makeTestEvent("Deployment", "default", "app1", "100")
	event2 := makeTestEvent("Service", "default", "app1", "100")
	event3 := makeTestEvent("ConfigMap", "default", "app1", "100")

	q.Add(event1)
	q.Add(event2)
	q.Add(event3)

	if q.Len() != 3 {
		t.Errorf("expected Len()=3, got %d", q.Len())
	}
}

// TestSameNameDifferentNamespace tests that resources with same name but different namespaces are kept separately
func TestSameNameDifferentNamespace(t *testing.T) {
	q := NewEventQueue()

	// Add events for same kind/name but different namespaces
	event1 := makeTestEvent("Deployment", "ns1", "app", "100")
	event2 := makeTestEvent("Deployment", "ns2", "app", "100")

	added1, replaced1 := q.Add(event1)
	added2, replaced2 := q.Add(event2)

	if !added1 || replaced1 {
		t.Errorf("expected first add (added=true, replaced=false), got (added=%v, replaced=%v)", added1, replaced1)
	}
	if !added2 || replaced2 {
		t.Errorf("expected second add (added=true, replaced=false), got (added=%v, replaced=%v)", added2, replaced2)
	}
	if q.Len() != 2 {
		t.Errorf("expected Len()=2, got %d", q.Len())
	}
}

// TestDrainAndCopyClears tests that DrainAndCopy returns all events and clears the queue
func TestDrainAndCopyClears(t *testing.T) {
	q := NewEventQueue()

	// Add multiple events
	event1 := makeTestEvent("Deployment", "default", "app1", "100")
	event2 := makeTestEvent("Service", "default", "svc1", "200")
	event3 := makeTestEvent("ConfigMap", "default", "cm1", "300")

	q.Add(event1)
	q.Add(event2)
	q.Add(event3)

	// Drain the queue
	events := q.DrainAndCopy()

	if len(events) != 3 {
		t.Errorf("expected 3 events in drained slice, got %d", len(events))
	}
	if q.Len() != 0 {
		t.Errorf("expected Len()=0 after drain, got %d", q.Len())
	}
}

// TestSnapshotPreserves tests that Snapshot returns events without clearing the queue
func TestSnapshotPreserves(t *testing.T) {
	q := NewEventQueue()

	// Add multiple events
	event1 := makeTestEvent("Deployment", "default", "app1", "100")
	event2 := makeTestEvent("Service", "default", "svc1", "200")
	event3 := makeTestEvent("ConfigMap", "default", "cm1", "300")

	q.Add(event1)
	q.Add(event2)
	q.Add(event3)

	// Take snapshot
	events := q.Snapshot()

	if len(events) != 3 {
		t.Errorf("expected 3 events in snapshot, got %d", len(events))
	}
	if q.Len() != 3 {
		t.Errorf("expected Len()=3 after snapshot, got %d", q.Len())
	}
}

// TestConcurrentAccess tests that the queue can handle concurrent access safely
func TestConcurrentAccess(t *testing.T) {
	q := NewEventQueue()
	var wg sync.WaitGroup

	// Number of goroutines and events per goroutine
	numGoroutines := 10
	eventsPerGoroutine := 100

	// Launch goroutines that add events concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				// Mix of same and different keys
				var event AdmissionReview
				if j%2 == 0 {
					// Same key across goroutines with increasing versions
					rv := goroutineID*eventsPerGoroutine + j
					event = makeTestEventNumericRV("Deployment", "default", "shared", float64(rv))
				} else {
					// Different keys per goroutine
					uniqueName := "app-" + string(rune('a'+goroutineID)) + "-" + string(rune('0'+j%10))
					event = makeTestEvent("Deployment", "default", uniqueName, "100")
				}
				q.Add(event)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify no panics occurred and state is consistent
	length := q.Len()
	if length < 1 {
		t.Errorf("expected at least 1 event after concurrent operations, got %d", length)
	}

	// Verify we can still operate on the queue
	snapshot := q.Snapshot()
	if len(snapshot) != length {
		t.Errorf("snapshot length %d does not match Len() %d", len(snapshot), length)
	}
}

// TestInvalidResourceVersionString tests that events with missing or unparseable string resourceVersion are rejected
func TestInvalidResourceVersionString(t *testing.T) {
	q := NewEventQueue()

	testCases := []struct {
		name          string
		event         AdmissionReview
		expectAdded   bool
		expectReplaced bool
	}{
		{
			name: "missing metadata",
			event: AdmissionReview{
				Request: &AdmissionRequest{
					Kind:      GroupVersionKind{Kind: "Deployment"},
					Namespace: "default",
					Name:      "app1",
					Object:    map[string]any{},
				},
			},
			expectAdded:   false,
			expectReplaced: false,
		},
		{
			name: "missing resourceVersion",
			event: AdmissionReview{
				Request: &AdmissionRequest{
					Kind:      GroupVersionKind{Kind: "Deployment"},
					Namespace: "default",
					Name:      "app2",
					Object: map[string]any{
						"metadata": map[string]any{},
					},
				},
			},
			expectAdded:   false,
			expectReplaced: false,
		},
		{
			name:          "unparseable resourceVersion string",
			event:         makeTestEvent("Deployment", "default", "app3", "not-a-number"),
			expectAdded:   false,
			expectReplaced: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			added, replaced := q.Add(tc.event)
			if added != tc.expectAdded {
				t.Errorf("expected added=%v, got added=%v", tc.expectAdded, added)
			}
			if replaced != tc.expectReplaced {
				t.Errorf("expected replaced=%v, got replaced=%v", tc.expectReplaced, replaced)
			}
		})
	}

	// Verify no events were added
	if q.Len() != 0 {
		t.Errorf("expected Len()=0, got %d", q.Len())
	}
}

// TestNumericResourceVersion tests that numeric (float64) resourceVersion is handled correctly
func TestNumericResourceVersion(t *testing.T) {
	q := NewEventQueue()

	// Add event with numeric resourceVersion
	event := makeTestEventNumericRV("Deployment", "default", "my-app", 100.0)
	added, replaced := q.Add(event)

	if !added {
		t.Errorf("expected added=true, got added=%v", added)
	}
	if replaced {
		t.Errorf("expected replaced=false, got replaced=%v", replaced)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1, got %d", q.Len())
	}

	// Add newer numeric resourceVersion
	event2 := makeTestEventNumericRV("Deployment", "default", "my-app", 200.0)
	added2, replaced2 := q.Add(event2)

	if !added2 {
		t.Errorf("expected added=true, got added=%v", added2)
	}
	if !replaced2 {
		t.Errorf("expected replaced=true, got replaced=%v", replaced2)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1, got %d", q.Len())
	}
}

// TestExistingUnparseable tests that when existing event has invalid resourceVersion, incoming is rejected
func TestExistingUnparseable(t *testing.T) {
	q := NewEventQueue()

	// Manually insert an event with unparseable resourceVersion
	// This simulates a corrupt state
	badEvent := AdmissionReview{
		Request: &AdmissionRequest{
			Kind: GroupVersionKind{
				Kind: "Deployment",
			},
			Namespace: "default",
			Name:      "my-app",
			Object: map[string]any{
				"metadata": map[string]any{
					"resourceVersion": map[string]string{}, // Invalid type
				},
			},
		},
	}

	// Manually add to bypass Add validation
	q.mu.Lock()
	key := "Deployment/default/my-app"
	q.events[key] = badEvent
	q.mu.Unlock()

	// Try to add a valid event
	validEvent := makeTestEvent("Deployment", "default", "my-app", "100")
	added, replaced := q.Add(validEvent)

	if added {
		t.Errorf("expected added=false when existing event is unparseable, got added=%v", added)
	}
	if replaced {
		t.Errorf("expected replaced=false when existing event is unparseable, got replaced=%v", replaced)
	}
	if q.Len() != 1 {
		t.Errorf("expected Len()=1 (bad event still there), got %d", q.Len())
	}
}
