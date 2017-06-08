package cc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config implements the Configer interface.
// The priorities: flag > environment variables > normal configs.
// Only the String/Bool/Int/Float/Duration family use environment variables
// and flags, if any environment variable with same name is set and
// it isn't empty, the Bool/BoolOr will return true.
// NOTE: we take empty string, false boolean and zero number value as default
// value in flags, and those value has no priority.
type Config struct {
	flags map[string]interface{}
	kv    map[string]interface{}
}

func newConfig() *Config {
	return &Config{
		kv: make(map[string]interface{}),
	}
}

// NewConfig creates a new empty Config.
func NewConfig() *Config {
	c := newConfig()
	c.ParseFlags()
	return c
}

func newConfigFrom(kv map[string]interface{}) *Config {
	return &Config{kv: kv}
}

// NewConfigFrom creates a new Config from map.
func NewConfigFrom(kv map[string]interface{}) *Config {
	c := newConfigFrom(kv)
	c.ParseFlags()
	return c
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

// ParseFlags parse the flags explicitly, in genral, you don't.
func (c *Config) ParseFlags() {
	c.flags = parseFlags()
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
// Excludes the flags and environment variables.
func (c *Config) KV() map[string]interface{} {
	return c.kv
}

// Has returns true if the name has a value, otherwise false.
func (c *Config) Has(name string) bool {
	if _, in := c.flags[name]; in {
		return true
	}
	if env := os.Getenv(name); env != "" {
		return true
	}
	_, in := c.kv[name]
	return in
}

// Must creates panic if name not found.
func (c *Config) Must(name string) {
	if !c.Has(name) {
		panic(fmt.Errorf("not value found for '%s' in config", name))
	}
}

// Raw returns the raw value by name.
// Excludes the flags and environment variable.
func (c *Config) Raw(name string) interface{} {
	return c.kv[name]
}

// Value returns a Valuer by name.
// Excludes the flags and environment variable.
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
	if _, in := c.kv[name]; !in {
		c.kv[name] = value
	}
}

// Set set the value by name, it will replace the exist value.
func (c *Config) Set(name string, value interface{}) {
	c.kv[name] = value
}

// Config returns a key-value sub Configer by name, the returned Configer can consider as a reference.
// Excludes the flags and environment variable.
func (c *Config) Config(name string) Configer {
	if v, exists := c.kv[name]; exists {
		switch x := v.(type) {
		case Configer:
			return x
		case map[string]interface{}:
			child := newConfigFrom(x)
			c.Set(name, child)
			return child
		case map[interface{}]interface{}:
			child := newConfig()
			child.kv = unknownMapToStringMap(x)
			c.Set(name, child)
			return child
		default:
		}
	}
	child := newConfig()
	c.Set(name, child)
	return child
}

// String returns the string value by name, returns "" if not found.
func (c *Config) String(name string) string {
	return c.StringOr(name, "")
}

// StringOr returns the string value by name, returns the deflt if not found.
func (c *Config) StringOr(name string, deflt string) string {
	if v, ok := c.flags[name].(string); ok && v != "" {
		return v
	}
	if env := os.Getenv(name); env != "" {
		return env
	}
	if v, in := c.kv[name]; in {
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

// StringAndOr returns the string value by name if pattern matched,
// otherwise returns the deflt.
func (c *Config) StringAndOr(name string, pattern string, deflt string) string {
	if s, ok := c.StringAnd(name, pattern); ok {
		return s
	}
	return deflt
}

// Bool returns the bool value by name, returns false if not found.
func (c *Config) Bool(name string) bool {
	return c.BoolOr(name, false)
}

// BoolOr returns the bool value by name, returns the deflt if not found.
func (c *Config) BoolOr(name string, deflt bool) bool {
	if v, ok := c.flags[name].(bool); ok && v != false {
		return v
	}
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
	if v, ok := c.flags[name].(int); ok && v != 0 {
		return v
	}
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

// IntAnd returns the (int value, true) by name if pattern matched,
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

// IntAndOr returns the int value by name if pattern matched,
// otherwise returns the deflt.
func (c *Config) IntAndOr(name string, pattern string, deflt int) int {
	if n, ok := c.IntAnd(name, pattern); ok {
		return n
	}
	return deflt
}

// Int64 returns the int64 value by name, returns 0 if not found.
func (c *Config) Int64(name string) int64 {
	return c.Int64Or(name, 0)
}

// Int64Or returns the int64 value by name, returns the deflt if not found.
func (c *Config) Int64Or(name string, deflt int64) int64 {
	if v, ok := c.flags[name].(int64); ok && v != int64(0) {
		return v
	}
	if env := os.Getenv(name); env != "" {
		if n, err := strconv.ParseInt(env, 10, 64); err == nil {
			return n
		}
	}
	if v, exists := c.kv[name]; exists {
		return toInt64(v, deflt)
	}
	return deflt
}

// Int64And returns the (int64 value, true) by name if pattern matched,
// otherwise returns (0, false). NOTE: we convert all numbers into
// float64 then validate.
func (c *Config) Int64And(name string, pattern string) (int64, bool) {
	if !c.Has(name) {
		return 0, false
	}
	p := NewPattern(pattern)
	if n := c.Int64(name); p.ValidateFloat(float64(n)) {
		return n, true
	}
	return 0, false
}

// Int64AndOr returns the int64 value by name if pattern matched,
// otherwise returns the deflt. NOTE: we convert all numbers into
// float64 then validate.
func (c *Config) Int64AndOr(name string, pattern string, deflt int64) int64 {
	if n, ok := c.Int64And(name, pattern); ok {
		return n
	}
	return deflt
}

// Float returns the float64 value by name, return 0.0 if not found.
func (c *Config) Float(name string) float64 {
	return c.FloatOr(name, 0.0)
}

// FloatOr returns the float64 value by name, return deflt if not found.
func (c *Config) FloatOr(name string, deflt float64) float64 {
	if v, ok := c.flags[name].(float64); ok && v != float64(0) {
		return v
	}
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

// FloatAndOr returns the float64 value by name if pattern matched,
// otherwise returns the deflt.
func (c *Config) FloatAndOr(name string, pattern string, deflt float64) float64 {
	if n, ok := c.FloatAnd(name, pattern); ok {
		return n
	}
	return deflt
}

// Duration returns the time.Duration value by name,
// return time.Duration(0) if not found.
func (c *Config) Duration(name string) time.Duration {
	return c.DurationOr(name, 0)
}

// DurationOr returns the time.Duration value by name,
// return time.Duration(deflt) if not found.
func (c *Config) DurationOr(name string, deflt int64) time.Duration {
	return time.Duration(c.Int64Or(name, deflt))
}

// DurationAnd returns the (time.Duration(value), true) by name if pattern matched,
// otherwise (time.Duration(0), false) returned. NOTE: we convert all numbers into
// float64 then validate.
func (c *Config) DurationAnd(name string, pattern string) (time.Duration, bool) {
	n, ok := c.Int64And(name, pattern)
	return time.Duration(n), ok
}

// DurationAndOr returns the time.Duration value by name if pattern matched,
// otherwise returns the deflt. NOTE: we convert all numbers into
// float64 then validate.
func (c *Config) DurationAndOr(name string, pattern string, deflt int64) time.Duration {
	return time.Duration(c.Int64AndOr(name, pattern, deflt))
}
