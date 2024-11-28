package attestations

import (
	"testing"
)

func TestMonitorIncrementSuccess(t *testing.T) {
	monitor := NewMonitor()
	monitor.IncrementSuccess()

	if monitor.successCount != 1 {
		t.Errorf("Expected successCount to be 1, got %d", monitor.successCount)
	}
}

func TestMonitorIncrementFailure(t *testing.T) {
	monitor := NewMonitor()
	reason := "test failure reason"
	monitor.IncrementFailure(reason)

	if monitor.failedCount != 1 {
		t.Errorf("Expected failedCount to be 1, got %d", monitor.failedCount)
	}

	if count, exists := monitor.failureReasons[reason]; !exists || count != 1 {
		t.Errorf("Expected failureReasons to contain '%s' with count 1, got %d", reason, count)
	}
}

func TestMonitorReset(t *testing.T) {
	monitor := NewMonitor()
	monitor.IncrementSuccess()
	monitor.IncrementFailure("test failure reason")
	monitor.Reset()

	if monitor.successCount != 0 {
		t.Errorf("Expected successCount to be 0 after reset, got %d", monitor.successCount)
	}

	if monitor.failedCount != 0 {
		t.Errorf("Expected failedCount to be 0 after reset, got %d", monitor.failedCount)
	}

	if len(monitor.failureReasons) != 0 {
		t.Errorf("Expected failureReasons to be empty after reset, got %d entries", len(monitor.failureReasons))
	}
}
