// Package filter provides a thread-safe in-process cuckoo filter for short codes.
package filter

import (
	"sync"

	cuckoo "github.com/seiflotfy/cuckoofilter"
)

// CodeFilter is a thread-safe wrapper around a cuckoo filter for short-link codes.
type CodeFilter struct {
	mu sync.RWMutex
	cf *cuckoo.Filter
}

// NewCodeFilter creates a CodeFilter with the given capacity.
func NewCodeFilter(capacity uint) *CodeFilter {
	return &CodeFilter{cf: cuckoo.NewFilter(capacity)}
}

// Insert adds code to the filter. Returns true if inserted, false if already present.
func (f *CodeFilter) Insert(code string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.cf.Insert([]byte(code))
}

// Lookup reports whether code may be in the filter.
// A false return is a definitive miss; a true return may be a false positive.
func (f *CodeFilter) Lookup(code string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.cf.Lookup([]byte(code))
}

// Count returns the number of items currently stored in the filter.
func (f *CodeFilter) Count() uint {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.cf.Count()
}
