// Package kvtest holds test helpers shared by the challenge's test files,
// in the spirit of net/http/httptest. Add it in stage 1; don't edit it.
package kvtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"

	kv "github.com/adhitamafikri/go-kv-store"
)

// TempDir creates a fresh directory for one test inside the project's tmp/
// folder (next to go.mod), so test data stays inside the repo instead of the
// system temp dir. The directory is removed when the test ends, unless the
// environment variable KVTEST_KEEP=1 is set — then it is kept so you can
// inspect the .data files, e.g. with: xxd tmp/<test>/000001.data
func TempDir(t testing.TB) string {
	t.Helper()
	root, err := projectRoot()
	if err != nil {
		t.Fatalf("kvtest.TempDir: %v", err)
	}
	base := filepath.Join(root, "tmp")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("kvtest.TempDir: %v", err)
	}
	dir, err := os.MkdirTemp(base, unsafeChars.ReplaceAllString(t.Name(), "_")+"-")
	if err != nil {
		t.Fatalf("kvtest.TempDir: %v", err)
	}
	if os.Getenv("KVTEST_KEEP") == "1" {
		t.Logf("kvtest: keeping test data in %s", dir)
	} else {
		t.Cleanup(func() { os.RemoveAll(dir) })
	}
	return dir
}

var unsafeChars = regexp.MustCompile(`[^A-Za-z0-9_.-]+`)

// projectRoot walks up from the working directory to the folder holding go.mod.
// (go test runs each package in its own directory, e.g. server/ for the server tests.)
func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found above the working directory")
		}
		dir = parent
	}
}

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
