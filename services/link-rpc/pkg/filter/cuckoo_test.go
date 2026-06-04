package filter_test

import (
	"sync"
	"testing"

	"github.com/aliaxy/zero-link/services/link-rpc/pkg/filter"
)

func TestCodeFilter_MissBeforeInsert(t *testing.T) {
	t.Parallel()
	f := filter.NewCodeFilter(1000)
	if f.Lookup("abc123") {
		t.Fatal("Lookup must return false for unseen code")
	}
}

func TestCodeFilter_HitAfterInsert(t *testing.T) {
	t.Parallel()
	f := filter.NewCodeFilter(1000)
	f.Insert("abc123")
	if !f.Lookup("abc123") {
		t.Fatal("Lookup must return true after Insert")
	}
}

func TestCodeFilter_ConcurrentInsertLookup(t *testing.T) {
	t.Parallel()
	f := filter.NewCodeFilter(10_000)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range 1000 {
			f.Insert("code")
		}
	}()
	go func() {
		defer wg.Done()
		for range 1000 {
			f.Lookup("code")
		}
	}()
	wg.Wait()
}
