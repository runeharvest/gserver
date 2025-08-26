// Package config provides a thread-safe, global configuration service
// that loads data from a TOML file.
package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	// mutex is a Read-Write mutex to allow multiple readers or one writer.
	// Reading is more common, so we use RWMutex for better performance.
	mutex sync.RWMutex

	// configData holds the decoded TOML data. We use a map for easy
	// manipulation of keys and values.
	configData map[string]interface{}
)

// Load decodes a TOML file into our global configData map.
// This function is write-locked, as it modifies the global state.
func Load(filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	// Read the entire file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read config file '%s': %w", filePath, err)
	}

	// Decode the TOML content into the map
	if _, err := toml.Decode(string(content), &configData); err != nil {
		return fmt.Errorf("could not decode TOML file: %w", err)
	}

	return nil
}

// getValue is an internal helper to safely access nested values.
func getValue(category, key string) (interface{}, error) {
	if configData == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	// Check if the category exists and is a map
	cat, ok := configData[category]
	if !ok {
		return nil, fmt.Errorf("category '%s' not found", category)
	}

	catMap, ok := cat.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("category '%s' is not a table", category)
	}

	// Check if the key exists within the category
	value, ok := catMap[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in category '%s'", key, category)
	}

	return value, nil
}

// --- Getters with Error Handling ---

// ValueE returns a string value or an error if not found/wrong type.
func ValueE(category, key string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := getValue(category, key)
	if err != nil {
		return "", err
	}

	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value for '%s.%s' is not a string", category, key)
	}
	return strVal, nil
}

// ValueBoolE returns a boolean value or an error.
func ValueBoolE(category, key string) (bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := getValue(category, key)
	if err != nil {
		return false, err
	}

	boolVal, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("value for '%s.%s' is not a boolean", category, key)
	}
	return boolVal, nil
}

// ValueIntE returns an integer value or an error.
func ValueIntE(category, key string) (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := getValue(category, key)
	if err != nil {
		return 0, err
	}

	// TOML library decodes all integers as int64
	intVal, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("value for '%s.%s' is not an integer", category, key)
	}
	return int(intVal), nil
}

// ValueFloatE returns a float value or an error.
func ValueFloatE(category, key string) (float64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := getValue(category, key)
	if err != nil {
		return 0.0, err
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0.0, fmt.Errorf("value for '%s.%s' is not a float", category, key)
	}
	return floatVal, nil
}

// --- Simple Getters (No Error Returned) ---
// These are convenient but will return a zero-value if the key is not found.

// Value returns a string, or an empty string if not found.
func Value(category, key string) string {
	val, _ := ValueE(category, key)
	return val
}

// ValueBool returns a boolean, or false if not found.
func ValueBool(category, key string) bool {
	val, _ := ValueBoolE(category, key)
	return val
}

// ValueInt returns an integer, or 0 if not found.
func ValueInt(category, key string) int {
	val, _ := ValueIntE(category, key)
	return val
}

// ValueFloat returns a float, or 0.0 if not found.
func ValueFloat(category, key string) float64 {
	val, _ := ValueFloatE(category, key)
	return val
}
