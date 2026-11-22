# Stage 2 — Record format & how `Put` writes it

Notes on the on-disk format of the append-only log (`000001.data`).

## 1. Encode vs decode

`encoding/binary` has two directions. Don't mix them up.

```text
 Get (decode)                          Put (encode)
 bytes  ──Uint32()──▶  number          number  ──PutUint32()──▶  bytes
 [00 00 00 07]  ──▶  7                 7  ──▶  [00 00 00 07]
```

- `binary.BigEndian.PutUint32(buf, n)` writes `n` into `buf`. Use it in **`Put`**.
- `binary.BigEndian.Uint32(buf)` reads a number out of `buf`. Use it in **`Get`**.

Calling `Uint32` on a header you just made with `make` only gives `0`, because the bytes are still all zeros.

## 2. Record layout

```text
┌──────────────┬──────────────┬───────────────┬─────────────────┐
│ keyLen       │ valLen       │ key           │ value           │
│ uint32, BE   │ uint32, BE   │ keyLen bytes  │ valLen bytes    │
│ 4 bytes      │ 4 bytes      │               │                 │
└──────────────┴──────────────┴───────────────┴─────────────────┘
 ◀──── header (8 bytes) ─────▶
```

The header is always the same size, so a reader always knows how much to read first. The lengths then say how much to read next. This is what makes the format **binary-safe**: keys and values can contain `\n`, spaces or `0x00` bytes, because nothing is found by looking for a separator.

## 3. `Put` step by step

Example: `Put("guest:1", "Ayu")`

```text
 key   = "guest:1"  → len 7
 value = "Ayu"      → len 3

 1) header := make([]byte, 8)
    ┌──┬──┬──┬──┬──┬──┬──┬──┐
    │00│00│00│00│00│00│00│00│
    └──┴──┴──┴──┴──┴──┴──┴──┘
     [0:4] keyLen [4:8] valLen

 2) PutUint32(header[0:4], uint32(len(key)))
    PutUint32(header[4:8], uint32(len(value)))
    ┌──┬──┬──┬──┬──┬──┬──┬──┐
    │00│00│00│07│00│00│00│03│
    └──┴──┴──┴──┴──┴──┴──┴──┘

 3) rec := append(header, key...)
    rec  = append(rec, value...)
    ┌───────────────────────┬─────────────────────┬──────────┐
    │ 00000007 00000003     │ g u e s t : 1       │ A y u    │
    │ header (8 bytes)      │ key (7 bytes)       │ val (3)  │
    └───────────────────────┴─────────────────────┴──────────┘
                    one []byte, 18 bytes total

 4) db.file.Write(rec)     ← one write call, check the error
```

## 4. The file after two `Put`s

The file is opened with `O_APPEND`, so each record goes right after the one before it:

```text
 offset 0                    18                          37
 │◀──── record 1 (18 B) ────▶│◀──── record 2 (19 B) ────▶│
 │ 7 │ 3 │ guest:1 │ Ayu      │ 7 │ 4 │ guest:2 │ Budi     │
```

Record size = `8 + keyLen + valLen`.

## 5. How `Get` reads it back

```text
 ┌─▶ ReadFull 8 bytes ──── io.EOF? ──yes──▶ stop
 │        │ no
 │        ▼
 │   keyLen = Uint32(h[0:4]),  valLen = Uint32(h[4:8])
 │        ▼
 │   ReadFull keyLen bytes  → k
 │   ReadFull valLen bytes  → v      ← always read it, even when k doesn't match,
 │        ▼                             or the next record won't line up
 │   bytes.Equal(k, key)? → remember v (don't return yet: a later record wins)
 └────────┘
 after the loop: found ? (v, nil) : (nil, ErrNotFound)
```

- `io.EOF` on the header read means the file ended cleanly at a record boundary.
- `io.ErrUnexpectedEOF` means the file ended in the middle of a record. Stage 5 handles this.

## 6. Go notes

- **`uint32(len(x))`**: `len` returns an `int`, and `PutUint32` takes a `uint32`, so you have to convert it yourself.
- **One `Write` per record**: if the program crashes, fewer half-written records are left behind than with three separate writes.
- **Set the capacity up front**: `make([]byte, 8, 8+len(key)+len(value))` allocates the whole record once, so the two `append` calls never have to grow the slice. With a plain `make([]byte, 8)`, the capacity is 8, so the first `append` copies everything into a new, larger array.
- **No `bytes.Clone` in `Put`**: `Write` copies the bytes into the file, so later changes to the caller's slices can't affect what's stored.
- **Wrap errors**: `fmt.Errorf("kv: put: %w", err)` keeps the real cause, and you can still check it with `errors.Is`.

## 7. Check it yourself

```bash
KVTEST_KEEP=1 go test -run Stage02_SurvivesRestart -v .
xxd tmp/TestStage02_SurvivesRestart-*/000001.data
```

You should see:

```text
00000000: 0000 0007 0000 0003 6775 6573 743a 3141  ........guest:1A
00000010: 7975 0000 0007 0000 0004 6775 6573 743a  yu........guest:
00000020: 3242 7564 69                             2Budi
```

Remove the kept folders afterwards with `rm -rf tmp/TestStage02_*`.

## 8. Why the file looks like `statusbooked` (no delimiter)

In a text editor, records look stuck together: `guest:1Ayu`, `statusbooked`. **That's expected.**

### The header is the delimiter

The header in front of each record already marks the boundaries:

```text
 header says: keyLen = 6, valLen = 6
 next 12 bytes: s t a t u s b o o k e d
                └── 6 ───┘└── 6 ───┘
                  key        value
```

`Get` never searches for a separator. It reads exactly `keyLen` bytes, then exactly `valLen` bytes.

### Why you can't see the header

Header bytes such as `00 00 00 06` aren't printable characters. Editors hide them or show odd symbols, so only the key and value text stays visible. Use `xxd` to see every byte:

```text
0000 0006 0000 0006 7374 6174 7573 626f 6f6b 6564   ........statusbooked
└ keyLen=6┘└ valLen=6┘└──── status ───┘└── booked ──┘
```

### Why a delimiter would make it worse

Any separator you pick (`,`, `\n`, `\0`) could also appear **inside** a key or value. Then the reader would split the record in the wrong place. `TestStage02_BinarySafe` checks exactly this, with a key and value that contain `\n` and `0x00`. With length prefixes, the content can be any bytes at all.

### This is a common design

Length-prefixed records are used by Bitcask and by many network protocols. The tradeoff is that **people can't read the file directly**. You need a tool like `xxd` to inspect it.
