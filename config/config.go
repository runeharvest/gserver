// Package config provides a thread-safe, global configuration service
// that loads data from a TOML file.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	mutex sync.RWMutex

	configData map[string]any
)

// MultiLoad loads multiple TOML files in sequence, merging their contents.
func MultiLoad(dir string, name string) error {
	path := filepath.Join(dir, "common.toml")
	err := Load(path)
	if err != nil {
		return fmt.Errorf("load common.toml: %w", err)
	}
	path = filepath.Join(dir, name+".toml")
	err = Load(path)
	if err != nil {
		return fmt.Errorf("load %s.toml: %w", name, err)
	}

	path = filepath.Join(dir, name+".override.toml")
	_, err = os.Stat(path)
	if err == nil {
		err = Load(path)
		if err != nil {
			return fmt.Errorf("load %s.override.toml: %w", name, err)
		}
	}
	return nil
}

// Load decodes a TOML file and merges it into our global configData map.
// This function is write-locked, as it modifies the global state.
func Load(filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read config file '%s': %w", filePath, err)
	}

	newData := make(map[string]any)
	if _, err := toml.Decode(string(content), &newData); err != nil {
		return fmt.Errorf("decode TOML file: %w", err)
	}

	if configData == nil {
		configData = newData
		return nil
	}
	for k, v := range newData {
		if !isSupportedType(v) {
			return fmt.Errorf("unsupported value type for key '%s': %T", k, v)
		}
	}

	return nil
}

func Close() {
	mutex.Lock()
	defer mutex.Unlock()
	configData = nil
}

// Set allows you to manually set the configuration data (for testing).
func SetConfig(data map[string]any) error {
	mutex.Lock()
	defer mutex.Unlock()
	configData = data
	err := sanitize()
	if err != nil {
		return err
	}

	return nil

}

// SetValue allows you to manually set a specific configuration value
func SetValue(category, key string, value any) error {
	if !isSupportedType(value) {
		return fmt.Errorf("unsupported value type for key '%s.%s': %T", category, key, value)
	}
	mutex.Lock()
	defer mutex.Unlock()

	if configData == nil {
		configData = make(map[string]any)
	}

	cat, ok := configData[category]
	if !ok {
		cat = make(map[string]any)
		configData[category] = cat
	}

	catMap, ok := cat.(map[string]any)
	if !ok {
		catMap = make(map[string]any)
		configData[category] = catMap
	}

	err := sanitize()
	if err != nil {
		return err
	}

	return nil
}

// value is an internal helper to safely access nested values.
func value(category, key string) (any, error) {
	if configData == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	cat, ok := configData[category]
	if !ok {
		return nil, fmt.Errorf("category '%s' not found", category)
	}

	catMap, ok := cat.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("category '%s' is not a table", category)
	}

	value, ok := catMap[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in category '%s'", key, category)
	}

	return value, nil
}

// ValueStrE returns a string value or an error if not found/wrong type.
func ValueStrE(category, key string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return "", err
	}

	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value for '%s.%s' is not a string", category, key)
	}
	return strVal, nil
}

func ValueSliceStrE(category, key string) ([]string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return nil, err
	}

	sliceVal, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("value for '%s.%s' is not a []string", category, key)
	}
	return sliceVal, nil
}

func ValueSliceBoolE(category, key string) ([]bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return nil, err
	}

	sliceVal, ok := val.([]bool)
	if !ok {
		return nil, fmt.Errorf("value for '%s.%s' is not a []bool", category, key)
	}
	return sliceVal, nil
}

func ValueSliceIntE(category, key string) ([]int64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return nil, err
	}

	switch retVal := val.(type) {
	case []int:
		intSlice := val.([]int)
		int64Slice := make([]int64, len(intSlice))
		for i, v := range intSlice {
			int64Slice[i] = int64(v)
		}
		defer SetValue(category, key, int64Slice)
		return int64Slice, nil
	case []int32:
		intSlice := val.([]int32)
		int64Slice := make([]int64, len(intSlice))
		for i, v := range intSlice {
			int64Slice[i] = int64(v)
		}
		defer SetValue(category, key, int64Slice)

		return int64Slice, nil
	case []int64:
		return retVal, nil
	default:
	}

	return nil, fmt.Errorf("value for '%s.%s' is not a []int64", category, key)
}

func ValueSliceFloatE(category, key string) ([]float64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return nil, err
	}

	sliceVal, ok := val.([]float64)
	if !ok {
		return nil, fmt.Errorf("value for '%s.%s' is not a []float64", category, key)
	}
	return sliceVal, nil
}

// ValueBoolE returns a boolean value or an error.
func ValueBoolE(category, key string) (bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
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
func ValueIntE(category, key string) (int64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return 0, err
	}

	intVal, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("value for '%s.%s' is not an integer", category, key)
	}
	return int64(intVal), nil
}

// ValueFloatE returns a float value or an error.
func ValueFloatE(category, key string) (float64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	val, err := value(category, key)
	if err != nil {
		return 0.0, err
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0.0, fmt.Errorf("value for '%s.%s' is not a float", category, key)
	}
	return floatVal, nil
}

// Value returns a string, or an empty string if not found.
func ValueStr(category, key string) string {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueStrE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

// ValueBool returns a boolean, or false if not found.
func ValueBool(category, key string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueBoolE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

// ValueInt returns an integer, or 0 if not found.
func ValueInt(category, key string) int64 {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueIntE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

// ValueFloat returns a float, or 0.0 if not found.
func ValueFloat(category, key string) float64 {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueFloatE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

func ValueSliceStr(category, key string) []string {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueSliceStrE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

func ValueSliceBool(category, key string) []bool {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueSliceBoolE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

func ValueSliceInt(category, key string) []int64 {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueSliceIntE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

func ValueSliceFloat(category, key string) []float64 {
	mutex.RLock()
	defer mutex.RUnlock()
	val, err := ValueSliceFloatE(category, key)
	if err != nil && os.Getenv("IS_PRODUCTION") == "1" {
		panic(err)
	}
	return val
}

func isSupportedType(value any) bool {
	switch value.(type) {
	case string, bool, int64, int, float64, []string, []bool, []int64, []float64:
		return true
	}
	return false
}

// sanitize reigns in weird types and validates the config data.

func sanitize() error {
	for catKey, catVal := range configData {
		catMap, ok := catVal.(map[string]any)
		if !ok {
			return fmt.Errorf("category '%s' is not a table", catKey)
		}
		for key, val := range catMap {
			switch retVal := val.(type) {
			case int:
				val = int64(retVal)
				catMap[key] = val
			case int32:
				val = int64(retVal)
				catMap[key] = val
			case float32:
				val = float64(retVal)
				catMap[key] = val
			case []int:
				intSlice := val.([]int)
				int64Slice := make([]int64, len(intSlice))
				for i, v := range intSlice {
					int64Slice[i] = int64(v)
				}
				val = int64Slice
				catMap[key] = val
			case []int32:
				intSlice := val.([]int32)
				int64Slice := make([]int64, len(intSlice))
				for i, v := range intSlice {
					int64Slice[i] = int64(v)
				}
				val = int64Slice
				catMap[key] = val
			}

			if !isSupportedType(val) {
				return fmt.Errorf("unsupported value type for key '%s.%s': %T", catKey, key, val)
			}
		}
	}
	return nil
}
