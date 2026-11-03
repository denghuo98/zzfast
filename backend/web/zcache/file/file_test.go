package file

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSetIfNotExist(t *testing.T) {
	ctx := context.Background()
	adapter := NewAdapterFile(t.TempDir())

	ok, err := adapter.SetIfNotExist(ctx, "k", "v1", time.Minute)
	if err != nil {
		t.Fatalf("SetIfNotExist failed: %v", err)
	}
	if !ok {
		t.Fatal("expected first SetIfNotExist to succeed")
	}

	ok, err = adapter.SetIfNotExist(ctx, "k", "v2", time.Minute)
	if err != nil {
		t.Fatalf("SetIfNotExist failed: %v", err)
	}
	if ok {
		t.Fatal("expected second SetIfNotExist to fail")
	}

	v, err := adapter.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v.String() != "v1" {
		t.Fatalf("unexpected value: %s", v.String())
	}
}

func TestSetIfNotExistConcurrent(t *testing.T) {
	ctx := context.Background()
	adapter := NewAdapterFile(t.TempDir())

	var (
		wg      sync.WaitGroup
		success int32
	)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := adapter.SetIfNotExist(ctx, "k", "v", time.Minute)
			if err != nil {
				t.Errorf("SetIfNotExist failed: %v", err)
				return
			}
			if ok {
				atomic.AddInt32(&success, 1)
			}
		}()
	}
	wg.Wait()
	if success != 1 {
		t.Fatalf("expected exactly one success, got %d", success)
	}
}
