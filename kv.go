package kv

import (
	"bytes"
	"context"
	"errors"
)

var ErrNotFound = errors.New("kv: Key not found")

type Options struct {
	// MaxFileSize is the size at which the active data file is rotated.
	// 0 means "use a sensible default" (e.g. 64 MiB). Matters from stage 6.
	MaxFileSize int64
}

type DB struct {
	data map[string][]byte
}

func Open(dir string, opts Options) (*DB, error) {
	return &DB{data: make(map[string][]byte)}, nil
}

func (db *DB) Put(key, value []byte) error {
	// empty key -> error
	if len(key) == 0 {
		return errors.New("empty key")
	}

	// db.data[string(key)] = value
	// The above statement is a bit incorrect
	// Need to make sure the data is not aliased, therefore we need to copy the value
	// So when the caller mutates the value that was inserted, it doesn't affect the data in the DB
	db.data[string(key)] = bytes.Clone(value)
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	// missing -> ErrNotFound
	if val, ok := db.data[string(key)]; ok {
		// return val, nil

		// Utilize bytes.Clone to avoid aliasing
		return bytes.Clone(val), nil
	}
	return nil, ErrNotFound
}

func (db *DB) Delete(key []byte) error {
	// missing -> nil
	if _, ok := db.data[string(key)]; !ok {
		return nil
	}

	delete(db.data, string(key))
	return nil
}

func (db *DB) Close() error {
	return nil
}

// Used in stage 7
func (db *DB) Merge() error {
	return errors.New("not implemented")
}

// Used in stage 10
func (db *DB) MergeContext(ctx context.Context) error {
	return errors.New("not implemented")
}
