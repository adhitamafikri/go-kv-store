package kv

type Options struct {
	// MaxFileSize is the size at which the active data file is rotated.
	// 0 means "use a sensible default" (e.g. 64 MiB). Matters from stage 6.
	MaxFileSize int64
}
