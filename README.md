# Go KV Store

The Riak Bitcask-style KV Store built with Golang

See this artifact for details: <https://claude.ai/artifact/5ahy4QrT96jPGHeSEuqKNe>

## Project Structure

This structure will be *fully visible by stage 10*

```bash
kvstore/
├── go.mod
├── kv.go                 # entry point: DB, Open, Options, ErrNotFound → calls engine
├── server.go             # NewServer / Serve → calls protocol + engine
├── stage01_test.go …     # challenge tests (black-box, public API only)
├── internal/
│   ├── record/           # Record, Kind, Encode/Decode, CRC. Pure, no I/O
│   ├── keydir/           # in-memory index: key → Entry{segment, offset, size}
│   ├── segment/          # one .data file: Append, ReadAt, Scan, Truncate, naming
│   ├── engine/           # write path, recovery, rotation, merge, locking
│   └── protocol/         # parse "SET k v" → Command, format replies. Pure
└── cmd/kvserver/main.go  # runnable binary to poke at with nc
```
