package kv

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("kv: Key not found")

type Options struct {
	// MaxFileSize is the size at which the active data file is rotated.
	// 0 means "use a sensible default" (e.g. 64 MiB). Matters from stage 6.
	MaxFileSize int64
}

type DB struct{}

func Open(dir string, opts Options) (*DB, error)
func (db *DB) Put(key, value []byte) error    // empty key -> error
func (db *DB) Get(key []byte) ([]byte, error) // missing -> ErrNotFound
func (db *DB) Delete(key []byte) error        // missing -> nil
func (db *DB) Close() error

// Used in stage 7
func (db *DB) Merge() error

// Used in stage 10
func (db *DB) MergeContext(ctx context.Context) error
