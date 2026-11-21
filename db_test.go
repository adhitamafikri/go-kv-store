package kv_test

// Stage 1 — The API: an in-memory store.

import (
	"testing"

	kv "github.com/adhitamafikri/go-kv-store"
	"github.com/adhitamafikri/go-kv-store/internal/kvtest"
)

func TestStage01_PutGet(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	kvtest.MustPut(t, db, "guest:42", "Fikri")
	kvtest.MustPut(t, db, "stay:7", "checked-in")
	kvtest.MustGet(t, db, "guest:42", "Fikri")
	kvtest.MustGet(t, db, "stay:7", "checked-in")
}

func TestStage01_Overwrite(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	kvtest.MustPut(t, db, "k", "v1")
	kvtest.MustPut(t, db, "k", "v2")
	kvtest.MustGet(t, db, "k", "v2")
}

func TestStage01_MissingKey(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	kvtest.MustMiss(t, db, "nope")
}

func TestStage01_Delete(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	kvtest.MustPut(t, db, "k", "v")
	if err := db.Delete([]byte("k")); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	kvtest.MustMiss(t, db, "k")

	// Deleting a key that does not exist is not an error.
	if err := db.Delete([]byte("never-existed")); err != nil {
		t.Fatalf("Delete of missing key should return nil, got %v", err)
	}
}

func TestStage01_EmptyKeyRejected(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	if err := db.Put(nil, []byte("v")); err == nil {
		t.Fatal("Put with empty key should return an error")
	}
}

func TestStage01_EmptyValueAllowed(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	kvtest.MustPut(t, db, "k", "")
	kvtest.MustGet(t, db, "k", "")
}

func TestStage01_NoAliasing(t *testing.T) {
	db := kvtest.MustOpen(t, kvtest.TempDir(t), kv.Options{})
	defer kvtest.MustClose(t, db)

	key := []byte("k")
	val := []byte("original")
	if err := db.Put(key, val); err != nil {
		t.Fatal(err)
	}
	// Mutating the caller's slices after Put must not change what is stored.
	val[0] = 'X'
	key[0] = 'Z'
	kvtest.MustGet(t, db, "k", "original")

	// Mutating the slice returned by Get must not change what is stored.
	got, _ := db.Get([]byte("k"))
	if len(got) > 0 {
		got[0] = 'Y'
	}
	kvtest.MustGet(t, db, "k", "original")
}
