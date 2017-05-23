package cc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// TODO(damnever): cooperate with flag?

// Config implements the Configer interface. The environment variable
// will shadow the value from config if they have the same name,
// Only the String/Bool/Int/Float/Duration family use environment
// variables, if any environment variable with same name is set and
// it isn't empty, the Bool/BoolOr will return true.
type Config struct {
	kv map[string]interface{}
}

// NewConfig creates a new empty Config.
func NewConfig() *Config {
	return &Config{
		kv: make(map[string]interface{}),
	}
}

// NewConfigFrom creates a new Config from map.
func NewConfigFrom(kv map[string]interface{}) *Config {
	return &Config{kv: kv}
}

// NewConfigFromJSON creates a new Config from JSON bytes.
func NewConfigFromJSON(data []byte) (*Config, error) {
	c := NewConfig()
	if err := c.MergeFromJSON(data); err != nil {
		return nil, err
	}
	return c, nil
}

// NewConfigFromYAML creates a new Config from YAML bytes.
func NewConfigFromYAML(data []byte) (*Config, error) {
	c := NewConfig()
	if err := c.MergeFromYAML(data); err != nil {
		return nil, err
	}
	return c, nil
}

// NewConfigFromFile creates a Config from a config file, extension must be
// one of ".yaml", ".yml" or ".json".
func NewConfigFromFile(fpath string) (*Config, error) {
	c := NewConfig()
	if err := c.MergeFromFile(fpath); err != nil {
		return nil, err
	}
	return c, nil
}

// MergeFromFile merges config data from file, the new config will replace
// the old. File extension must be one of ".yaml", ".yaml" or ".json".
func (c *Config) MergeFromFile(fpath string) error {
	var data map[string]interface{}
	var err error

	ext := filepath.Ext(fpath)
	switch ext {
	case ".yaml", ".yml":
		data, err = unmarshalYAMLFile(fpath)
	case ".json":
		data, err = unmarshalJSONFile(fpath)
	case "":
		err = fmt.Errorf("can not determine the config file type: %s", fpath)
	default:
		err = fmt.Errorf("unsupported config file type: %s", fpath)
	}
	if err != nil {
		return err
	}
	for k, v := range data {
		c.kv[k] = v
	}
	return nil
}

// MergeFromJSON  merges data from JSON bytes, the value from same name will be replaced.
func (c *Config) MergeFromJSON(b []byte) error {
	data, err := unmarshalJSON(b)
	if err != nil {
		return err
	}
	for k, v := range data {
		c.kv[k] = v
	}
	return nil
}

// MergeFromYAML  merges data from YAML bytes, the value from same name will be replaced.
func (c *Config) MergeFromYAML(b []byte) error {
	data, err := unmarshalYAML(b)
	if err != nil {
		return err
	}
	for k, v := range data {
		c.kv[k] = v
	}
	return nil
}

// Merge merges data from another Config, the value from same name will be replaced.
func (c *Config) Merge(config *Config) error {
	for k, v := range config.kv {
		c.kv[k] = v
	}
	return nil
}

// KV returns the Config's internal data as a string map.
func (c *Config) KV() map[string]interface{} {
	return c.kv
}

// Has returns true if the name has a value, otherwise false.
func (c *Config) Has(name string) bool {
	if env := os.Getenv(name); env != "" {
		return true
	}
	_, exists := c.kv[name]
	return exists
}

// Must creates panic if name not found.
func (c *Config) Must(name string) {
	if !c.Has(name) {
		panic(fmt.Errorf("not value found for '%s' in config", name))
	}
}

// Raw returns the raw value by name.
// No support for environment variable.
func (c *Config) Raw(name string) interface{} {
	return c.kv[name]
}

// Value returns a Valuer by name.
func (c *Config) Value(name string) Valuer {
	v, ok := c.kv[name]
	if !ok {
		return NewValue(nil)
	}
	if child, ok := v.(Configer); ok {
		return NewValue(child.KV())
	}
	return NewValue(v)
}

// Pattern returns a Patterner by name.
func (c *Config) Pattern(name string) Patterner {
	return NewPattern(c.String(name))
}

// SetDefault set the default value by name if not found.
func (c *Config) SetDefault(name string, value interface{}) {
	if !c.Has(name) {
		c.kv[name] = value
	}
}

// Set set the value by name, it will replcae the exist value.
func (c *Config) Set(name string, value interface{}) {
	c.kv[name] = value
}

// Config returns a key-value sub Configer by name, the returned Configer can consider as a reference.
func (c *Config) Config(name string) Configer {
	if v, exists := c.kv[name]; exists {
		switch x := v.(type) {
		case Configer:
			return x
		case map[string]interface{}:
			child := NewConfigFrom(x)
			c.Set(name, child)
			return child
		case map[interface{}]interface{}:
			child := NewConfig()
			child.kv = unknownMapToStringMap(x)
			c.Set(name, child)
			return child
		default:
		}
	}
	child := NewConfig()
	c.Set(name, child)
	return child
}

// String returns the string value by name, returns "" if not found.
func (c *Config) String(name string) string {
	return c.StringOr(name, "")
}

// StringOr returns the string value by name, returns the deflt if not found.
func (c *Config) StringOr(name string, deflt string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	if v, exists := c.kv[name]; exists {
		return toString(v, deflt)
	}
	return deflt
}

// StringAnd returns the (string value, true) if pattern matched,
// otherwise returns ("", false).
func (c *Config) StringAnd(name string, pattern string) (string, bool) {
	if !c.Has(name) {
		return "", false
	}
	p := NewPattern(pattern)
	if s := c.String(name); p.ValidateString(s) {
		return s, true
	}
	return "", false
}

// Bool returns the bool value by name, returns false if not found.
func (c *Config) Bool(name string) bool {
	return c.BoolOr(name, false)
}

// BoolOr returns the bool value by name, returns the deflt if not found.
func (c *Config) BoolOr(name string, deflt bool) bool {
	if env := os.Getenv(name); env != "" {
		return true
	}
	if v, exists := c.kv[name]; exists {
		return toBool(v, deflt)
	}
	return deflt
}

// Int returns the int value by name, returns 0 if not found.
func (c *Config) Int(name string) int {
	return c.IntOr(name, 0)
}

// IntOr returns the int value by name, returns the deflt if not found.
func (c *Config) IntOr(name string, deflt int) int {
	if env := os.Getenv(name); env != "" {
		if n, err := strconv.Atoi(env); err == nil {
			return n
		}
	}
	if v, exists := c.kv[name]; exists {
		return toInt(v, deflt)
	}
	return deflt
}

// IntAnd returns the (int value, true) if pattern matched,
// otherwise returns (0, false)
func (c *Config) IntAnd(name string, pattern string) (int, bool) {
	if !c.Has(name) {
		return 0, false
	}
	p := NewPattern(pattern)
	if n := c.Int(name); p.ValidateInt(n) {
		return n, true
	}
	return 0, false
}

// Float returns the float64 value by name, return 0.0 if not found.
func (c *Config) Float(name string) float64 {
	return c.FloatOr(name, 0.0)
}

// FloatOr returns the float64 value by name, return deflt if not found.
func (c *Config) FloatOr(name string, deflt float64) float64 {
	if env := os.Getenv(name); env != "" {
		if n, err := strconv.ParseFloat(env, 64); err == nil {
			return n
		}
	}
	if v, exists := c.kv[name]; exists {
		return toFloat64(v, deflt)
	}
	return deflt
}

// FloatAnd returns the (float64 value, true) if pattern matched,
// otherwise (0.0, false) returned.
func (c *Config) FloatAnd(name string, pattern string) (float64, bool) {
	if !c.Has(name) {
		return 0.0, false
	}
	p := NewPattern(pattern)
	if n := c.Float(name); p.ValidateFloat(n) {
		return n, true
	}
	return 0.0, false
}

// Duration returns the time.Duration value by name,
// return time.Duration(0) if not found.
func (c *Config) Duration(name string) time.Duration {
	return c.DurationOr(name, 0)
}

// DurationOr returns the time.Duration value by name,
// return time.Duration(deflt) if not found.
func (c *Config) DurationOr(name string, deflt int64) time.Duration {
	if env := os.Getenv(name); env != "" {
		if n, err := strconv.Atoi(env); err == nil {
			return time.Duration(n)
		}
	}
	if v, exists := c.kv[name]; exists {
		return time.Duration(toInt64(v, deflt))
	}
	return time.Duration(deflt)
}
