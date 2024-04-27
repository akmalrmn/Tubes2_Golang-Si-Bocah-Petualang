package config

import (
	"runtime"
	"time"
)

type Config struct {
	// Colly Boolean
	IsAsync    bool
	AllowCache bool

	// Colly Config
	MaxDepth       int
	MaxQueryThread int
	MaxQueueSize   int
	LimitQueue     int
	RandomDelay    time.Duration
	AllowedDomains []string
	CacheDir       string

	// Runtime + Debug
	MaxParallelism int
	MaxProcessor   int
	MaxThreads     int
}

func NewConfigDefault() *Config {
	return &Config{
		IsAsync:        false,
		AllowCache:     false,
		MaxDepth:       2,
		MaxQueryThread: 2,
		MaxQueueSize:   10000,
		AllowedDomains: []string{"en.wikipedia.org"},
		CacheDir:       "cache/",
		MaxParallelism: 2,
		MaxProcessor:   4,
		MaxThreads:     50000,
	}
}

func NewTurboConfig() *Config {
	return &Config{
		IsAsync:        false,
		AllowCache:     true,
		MaxDepth:       12,
		MaxQueryThread: runtime.NumCPU(),
		MaxQueueSize:   10000000,
		LimitQueue:     20000,
		RandomDelay:    1000 * time.Millisecond,
		AllowedDomains: []string{"en.wikipedia.org"},
		CacheDir:       "./cache",
		MaxParallelism: 400,
		MaxProcessor:   runtime.NumCPU(),
		MaxThreads:     50000,
	}
}
