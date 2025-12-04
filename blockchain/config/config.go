package config

import "sync"

var once sync.Once

// Config holds cross-cutting configuration values for the blockchain
type Config struct {
	NumWorkers int
	Difficulty int
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		NumWorkers: 20,
		Difficulty: 20,
	}
}

// Global configuration instance
var global *Config

// Get returns the global configuration instance
// If not initialized, returns default configuration
func Get() *Config {
	once.Do(func() {
		global = Default()
	})
	return global
}

// SetNumWorkers sets the number of workers globally
func SetNumWorkers(n int) {
	Get().NumWorkers = n
}

// SetDifficulty sets the difficulty globally
func SetDifficulty(d int) {
	Get().Difficulty = d
}

// NumWorkers returns the global number of workers
func NumWorkers() int {
	return Get().NumWorkers
}

// Difficulty returns the global difficulty
func Difficulty() int {
	return Get().Difficulty
}
