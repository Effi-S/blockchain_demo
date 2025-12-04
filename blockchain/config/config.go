package config

// Config holds cross-cutting configuration values for the blockchain
type Config struct {
	NumWorkers int
	Difficulty int
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		NumWorkers: 12,
		Difficulty: 15,
	}
}

// Global configuration instance
var global *Config

// Init initializes the global configuration with the provided config
func Init(cfg *Config) {
	global = cfg
}

// Get returns the global configuration instance
// If not initialized, returns default configuration
func Get() *Config {
	if global == nil {
		global = Default()
	}
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
