package kv_test

// Stage 2 — The world's simplest database: an append-only log on disk.

import (
	"testing"

	kv "github.com/adhitamafikri/go-kv-store"
	"github.com/adhitamafikri/go-kv-store/internal/kvtest"
)

func TestStage02_CreatesDataFile(t *testing.T) {
	dir := kvtest.TempDir(t)
	db := kvtest.MustOpen(t, dir, kv.Options{})
	kvtest.MustPut(t, db, "k", "v")
	kvtest.MustClose(t, db)

	if n := len(kvtest.DataFiles(t, dir)); n == 0 {
		t.Fatalf("expected at least one *.data file in %s, found none", dir)
	}
}

func TestStage02_SurvivesRestart(t *testing.T) {
	dir := kvtest.TempDir(t)

	db := kvtest.MustOpen(t, dir, kv.Options{})
	kvtest.MustPut(t, db, "guest:1", "Ayu")
	kvtest.MustPut(t, db, "guest:2", "Budi")
	kvtest.MustClose(t, db)

	db = kvtest.MustOpen(t, dir, kv.Options{})
	defer kvtest.MustClose(t, db)
	kvtest.MustGet(t, db, "guest:1", "Ayu")
	kvtest.MustGet(t, db, "guest:2", "Budi")
	kvtest.MustMiss(t, db, "guest:3")
}

func TestStage02_LatestWriteWinsAfterRestart(t *testing.T) {
	dir := kvtest.TempDir(t)

	db := kvtest.MustOpen(t, dir, kv.Options{})
	kvtest.MustPut(t, db, "status", "booked")
	kvtest.MustPut(t, db, "status", "checked-in")
	kvtest.MustPut(t, db, "status", "checked-out")
	kvtest.MustClose(t, db)

	db = kvtest.MustOpen(t, dir, kv.Options{})
	defer kvtest.MustClose(t, db)
	kvtest.MustGet(t, db, "status", "checked-out")
}

func TestStage02_AppendAcrossSessions(t *testing.T) {
	dir := kvtest.TempDir(t)

	db := kvtest.MustOpen(t, dir, kv.Options{})
	kvtest.MustPut(t, db, "a", "1")
	kvtest.MustClose(t, db)

	db = kvtest.MustOpen(t, dir, kv.Options{})
	kvtest.MustPut(t, db, "b", "2")
	kvtest.MustPut(t, db, "a", "3")
	kvtest.MustClose(t, db)

	db = kvtest.MustOpen(t, dir, kv.Options{})
	defer kvtest.MustClose(t, db)
	kvtest.MustGet(t, db, "a", "3")
	kvtest.MustGet(t, db, "b", "2")
}

func TestStage02_BinarySafe(t *testing.T) {
	dir := kvtest.TempDir(t)
	key := []byte{0x00, 0x0a, 0xff, ' ', '\n'}
	val := []byte{0x00, 0x00, '\r', '\n', 0xfe}

	db := kvtest.MustOpen(t, dir, kv.Options{})
	if err := db.Put(key, val); err != nil {
		t.Fatal(err)
	}
	kvtest.MustClose(t, db)

	db = kvtest.MustOpen(t, dir, kv.Options{})
	defer kvtest.MustClose(t, db)
	kvtest.MustGet(t, db, string(key), string(val))
}
