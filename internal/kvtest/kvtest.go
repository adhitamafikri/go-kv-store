// Package kvtest holds test helpers shared by the challenge's test files,
// in the spirit of net/http/httptest. Add it in stage 1; don't edit it.
package kvtest

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	kv "github.com/adhitamafikri/go-kv-store"
)

// MustOpen opens a DB in dir or fails the test.
func MustOpen(t testing.TB, dir string, opts kv.Options) *kv.DB {
	t.Helper()
	db, err := kv.Open(dir, opts)
	if err != nil {
		t.Fatalf("Open(%q): %v", dir, err)
	}
	return db
}

// MustPut stores k=v or fails the test.
func MustPut(t testing.TB, db *kv.DB, k, v string) {
	t.Helper()
	if err := db.Put([]byte(k), []byte(v)); err != nil {
		t.Fatalf("Put(%q): %v", k, err)
	}
}

// MustGet fails the test unless k currently maps to want.
func MustGet(t testing.TB, db *kv.DB, k, want string) {
	t.Helper()
	got, err := db.Get([]byte(k))
	if err != nil {
		t.Fatalf("Get(%q): unexpected error %v", k, err)
	}
	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("Get(%q) = %q, want %q", k, got, want)
	}
}

// MustMiss fails the test unless Get(k) returns kv.ErrNotFound.
func MustMiss(t testing.TB, db *kv.DB, k string) {
	t.Helper()
	_, err := db.Get([]byte(k))
	if !errors.Is(err, kv.ErrNotFound) {
		t.Fatalf("Get(%q): want kv.ErrNotFound, got %v", k, err)
	}
}

// MustClose closes db or fails the test.
func MustClose(t testing.TB, db *kv.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// DataFiles lists the *.data files in dir, in lexical (= write) order.
func DataFiles(t testing.TB, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*.data"))
	if err != nil {
		t.Fatal(err)
	}
	return files // filepath.Glob returns names in lexical order
}

// NewestDataFile returns the last *.data file in dir.
func NewestDataFile(t testing.TB, dir string) string {
	t.Helper()
	files := DataFiles(t, dir)
	if len(files) == 0 {
		t.Fatalf("no *.data files in %s", dir)
	}
	return files[len(files)-1]
}

// TotalDataSize sums the sizes of all *.data files in dir.
func TotalDataSize(t testing.TB, dir string) int64 {
	t.Helper()
	var total int64
	for _, f := range DataFiles(t, dir) {
		st, err := os.Stat(f)
		if err != nil {
			t.Fatal(err)
		}
		total += st.Size()
	}
	return total
}

// CountdownContext returns a context that reports "cancelled" after it has
// been checked n times (calls to Done or Err). Tests use it to cancel work
// part-way through without relying on timing.
type countdownCtx struct {
	context.Context
	mu     sync.Mutex
	left   int
	done   chan struct{}
	closed bool
}

func CountdownContext(n int) context.Context {
	return &countdownCtx{Context: context.Background(), left: n, done: make(chan struct{})}
}

func (c *countdownCtx) tick() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.left--
	if c.left <= 0 {
		c.closed = true
		close(c.done)
	}
}

func (c *countdownCtx) Done() <-chan struct{} { c.tick(); return c.done }

func (c *countdownCtx) Err() error {
	c.tick()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return context.Canceled
	}
	return nil
}
