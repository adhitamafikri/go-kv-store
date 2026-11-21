# Go KV Store

The Riak Bitcask-style KV Store built with Golang

See this artifact for details: <https://claude.ai/artifact/5ahy4QrT96jPGHeSEuqKNe>

## Project Structure

This structure will be *fully visible by stage 10*

```bash
go-kv-store/
├── go.mod                   # module kvstore  (the tests import this path)
├── kv.go                    # public API: DB, Open, Options, ErrNotFound → engine
├── db_test.go               # stage 1  ┐
├── persistence_test.go      # stage 2  │
├── index_test.go            # stage 3  │  acceptance tests for the public DB API,
├── tombstone_test.go        # stage 4  │  black-box (package kv_test), next to kv.go
├── recovery_test.go         # stage 5  │
├── segment_test.go          # stage 6  │
├── merge_test.go            # stage 7  │
├── concurrency_test.go      # stage 8  │
├── merge_cancel_test.go     # stage 10 ┘
├── server/
│   ├── server.go            # New, Serve, Shutdown: networking only
│   ├── server_test.go       # stage 9
│   └── shutdown_test.go     # stage 10
├── internal/
│   ├── kvtest/kvtest.go     # shared test helpers (given in stage 1)
│   ├── record/              # record.go   + record_test.go    ← your unit tests
│   ├── keydir/              # keydir.go   + keydir_test.go
│   ├── segment/             # segment.go  + segment_test.go
│   ├── engine/              # engine.go   + engine_test.go
│   └── protocol/            # protocol.go + protocol_test.go
└── cmd/kvserver/main.go     # runnable binary to poke at with nc
```
