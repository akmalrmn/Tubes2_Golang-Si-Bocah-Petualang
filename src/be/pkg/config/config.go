package config

import (
	"math"
	"runtime"
)

type Config struct {
	// Colly Boolean
	IsAsync    bool
	AllowCache bool

	// Colly Config
	MaxDepth       int
	MaxQueryThread int
	MaxQueueSize   int
	AllowedDomains []string

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
		MaxParallelism: 16,
		MaxProcessor:   50000,
		MaxThreads:     50000,
	}
}

func NewConfig(isAsync, allowCache bool, maxDepth, maxQueryThread, maxQueueSize int, allowedDomains []string, maxParallelism, maxProcessor, maxThreads int) *Config {
	return &Config{
		IsAsync:        isAsync,
		AllowCache:     allowCache,
		MaxDepth:       maxDepth,
		MaxQueryThread: maxQueryThread,
		MaxQueueSize:   maxQueueSize,
		AllowedDomains: allowedDomains,
		MaxParallelism: maxParallelism,
		MaxProcessor:   maxProcessor,
		MaxThreads:     maxThreads,
	}
}

func NewTurboConfig() *Config {
	return &Config{
		IsAsync:        true,
		AllowCache:     true,
		MaxDepth:       int(math.MaxInt32),
		MaxQueryThread: runtime.NumCPU(),
		MaxQueueSize:   10000,
		AllowedDomains: []string{"en.wikipedia.org"},
		MaxParallelism: 400,
		MaxProcessor:   runtime.NumCPU(),
		MaxThreads:     50000,
	}
}
