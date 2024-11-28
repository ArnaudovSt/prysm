package attestations

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

type Monitor interface {
	IncrementSuccess()
	IncrementFailure(reason string)
	Report()
	Reset()
}

var _ Monitor = (*attMonitor)(nil)

type attMonitor struct {
	successCount   uint64
	failedCount    uint64
	failureReasons map[string]uint64
	mu             sync.RWMutex
}

func NewMonitor() *attMonitor {
	return &attMonitor{
		failureReasons: make(map[string]uint64),
	}
}

func (m *attMonitor) IncrementSuccess() {
	atomic.AddUint64(&m.successCount, 1)
}

func (m *attMonitor) IncrementFailure(reason string) {
	atomic.AddUint64(&m.failedCount, 1)

	truncatedReason := truncateString(reason, 42)

	m.mu.Lock()
	m.failureReasons[truncatedReason]++
	m.mu.Unlock()
}

func (m *attMonitor) Report() {
	var builder strings.Builder
	builder.WriteString("Epoch Summary:\n")
	builder.WriteString(fmt.Sprintf("Successful Attestations: %d\n", atomic.LoadUint64(&m.successCount)))
	builder.WriteString(fmt.Sprintf("Failed Attestations: %d\n", atomic.LoadUint64(&m.failedCount)))

	builder.WriteString("Failure Reasons:\n")
	m.mu.RLock()
	for reason, count := range m.failureReasons {
		builder.WriteString(fmt.Sprintf("%s: %d\n", reason, count))
	}
	m.mu.RUnlock()

	log.Info(builder.String())
}

func (m *attMonitor) Reset() {
	atomic.StoreUint64(&m.successCount, 0)
	atomic.StoreUint64(&m.failedCount, 0)

	m.mu.Lock()
	m.failureReasons = make(map[string]uint64)
	m.mu.Unlock()
}

func truncateString(input string, maxLen int) string {
	if len(input) > maxLen {
		return input[:maxLen] + "..."
	}
	return input
}
