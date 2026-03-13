package workflow

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/boomatang/crystal/internal/logger"
)

// EventQueue manages a deduplicated queue of admission review events.
// Events are keyed by kind/namespace/name and deduplicated based on resourceVersion.
type EventQueue struct {
	mu     sync.Mutex
	events map[string]AdmissionReview // key: "kind/namespace/name"
}

// NewEventQueue creates and returns a new EventQueue.
func NewEventQueue() *EventQueue {
	return &EventQueue{
		events: make(map[string]AdmissionReview),
	}
}

// Add attempts to add an event to the queue. It deduplicates based on resourceVersion,
// keeping only the event with the higher resourceVersion for each unique resource.
//
// Returns:
//   - (true, false) - new resource added
//   - (true, true) - replaced older event with newer version
//   - (false, false) - rejected (older or invalid resourceVersion)
func (q *EventQueue) Add(event AdmissionReview) (added bool, replaced bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Generate the key for this event
	key := fmt.Sprintf("%s/%s/%s",
		event.Request.Kind.Kind,
		event.Request.Namespace,
		event.Request.Name)

	// Get resourceVersion for the incoming event
	incomingVersion, err := getResourceVersion(event)
	if err != nil {
		logger.Log.Warn("event rejected (invalid resourceVersion)",
			"kind", event.Request.Kind.Kind,
			"namespace", event.Request.Namespace,
			"name", event.Request.Name,
			"error", err)
		return false, false
	}

	// Check if we already have an event for this resource
	if existing, exists := q.events[key]; exists {
		existingVersion, err := getResourceVersion(existing)
		if err != nil {
			// Existing event is unparseable, log warning and reject incoming
			logger.Log.Warn("existing event has invalid resourceVersion, rejecting incoming event",
				"kind", event.Request.Kind.Kind,
				"namespace", event.Request.Namespace,
				"name", event.Request.Name,
				"error", err)
			return false, false
		}

		// Compare versions
		if incomingVersion > existingVersion {
			// Replace with newer version
			q.events[key] = event
			logger.Log.Info("event replaced older version",
				"kind", event.Request.Kind.Kind,
				"namespace", event.Request.Namespace,
				"name", event.Request.Name,
				"resourceVersion", incomingVersion,
				"existing_resourceVersion", existingVersion)
			return true, true
		}

		// Incoming event is older, skip it
		logger.Log.Info("event skipped (older than existing)",
			"kind", event.Request.Kind.Kind,
			"namespace", event.Request.Namespace,
			"name", event.Request.Name,
			"resourceVersion", incomingVersion,
			"existing_resourceVersion", existingVersion)
		return false, false
	}

	// New resource, add it
	q.events[key] = event
	logger.Log.Info("event added",
		"kind", event.Request.Kind.Kind,
		"namespace", event.Request.Namespace,
		"name", event.Request.Name,
		"resourceVersion", incomingVersion)
	return true, false
}

// getResourceVersion extracts and parses the resourceVersion from an AdmissionReview event.
// It handles both string and numeric representations of resourceVersion.
func getResourceVersion(event AdmissionReview) (int64, error) {
	// Access metadata using comma-ok idiom
	metadata, ok := event.Request.Object["metadata"]
	if !ok {
		return 0, fmt.Errorf("metadata field not found in event object")
	}

	// Type assert to map[string]any using comma-ok idiom
	metadataMap, ok := metadata.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("metadata is not a map[string]any")
	}

	// Extract resourceVersion field
	rvField, ok := metadataMap["resourceVersion"]
	if !ok {
		return 0, fmt.Errorf("resourceVersion field not found in metadata")
	}

	// Handle both string and numeric types
	switch rv := rvField.(type) {
	case string:
		// String type (expected): parse with strconv.ParseInt
		version, err := strconv.ParseInt(rv, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse resourceVersion string: %w", err)
		}
		return version, nil
	case float64:
		// Numeric type (from JSON): convert to int64
		return int64(rv), nil
	default:
		return 0, fmt.Errorf("resourceVersion has unexpected type: %T", rvField)
	}
}

// DrainAndCopy removes all events from the queue and returns them as a slice.
// This is a destructive operation that clears the queue.
func (q *EventQueue) DrainAndCopy() []AdmissionReview {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := len(q.events)
	logger.Log.Info("draining event queue", "count", count)

	// Copy all values to a slice
	result := make([]AdmissionReview, 0, count)
	for _, event := range q.events {
		result = append(result, event)
	}

	// Clear the map
	q.events = make(map[string]AdmissionReview)

	return result
}

// Len returns the current number of events in the queue.
// This operation is thread-safe.
func (q *EventQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.events)
}

// Snapshot returns a copy of the current events without draining the queue.
// This is a non-destructive operation that leaves the queue unchanged.
func (q *EventQueue) Snapshot() []AdmissionReview {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Copy all values to a slice
	result := make([]AdmissionReview, 0, len(q.events))
	for _, event := range q.events {
		result = append(result, event)
	}

	return result
}
