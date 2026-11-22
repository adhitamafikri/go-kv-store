package kv

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var ErrEmptyKey = errors.New("kv: empty key")
var ErrOpenDBFile = errors.New("kv: failed opening database file")
var ErrCloseDBFile = errors.New("kv: failed closing database file")

type DB struct {
	dir  string
	file *os.File
}

func Open(dir string, opts Options) (*DB, error) {
	dbFile, err := os.OpenFile(filepath.Join(dir, "000001.data"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

	if err != nil {
		return nil, ErrOpenDBFile
	}
	return &DB{dir: dir, file: dbFile}, nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	// This method needs its own open and read file implementation, because the `db.file` is already open with append mode
	file, err := os.OpenFile(filepath.Join(db.dir, "000001.data"), os.O_RDONLY, 0o644)
	if err != nil {
		return nil, ErrOpenDBFile
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	header := make([]byte, 8)
	found := false
	var value []byte

	for {
		_, err := io.ReadFull(reader, header)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("kv: read: %w", err)
		}

		kLen := binary.BigEndian.Uint32(header[0:4]) // gets the key length
		vLen := binary.BigEndian.Uint32(header[4:8]) // gets the value length

		k := make([]byte, kLen) // prepare new key slice
		v := make([]byte, vLen) // prepare new value slice

		_, err = io.ReadFull(reader, k)
		if err != nil {
			return nil, fmt.Errorf("kv: read key: %w", err)
		}

		_, err = io.ReadFull(reader, v)
		if err != nil {
			return nil, fmt.Errorf("kv: read value: %w", err)
		}

		if bytes.Equal(k, key) {
			value = v
			found = true
		}
	}

	if found {
		return value, nil
	}

	return nil, ErrNotFound
}

func (db *DB) Put(key, value []byte) error {
	kLen := uint32(len(key))
	vLen := uint32(len(value))

	if kLen == 0 {
		return ErrEmptyKey
	}

	if db.file == nil {
		return ErrOpenDBFile
	}

	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:4], kLen) // gets the key length
	binary.BigEndian.PutUint32(header[4:8], vLen) // gets the value length

	// Construct the data row
	record := append(header, key...)
	record = append(record, value...)

	_, err := db.file.Write(record)
	if err != nil {
		return fmt.Errorf("kv: write: %w", err)
	}

	return nil
}

func (db *DB) Delete(key []byte) error {
	return nil
}

func (db *DB) Close() error {
	err := db.file.Close()
	if err != nil {
		return ErrCloseDBFile
	}
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
