package cc

import (
	"os"
	"testing"
	"time"

	"github.com/damnever/cc/assert"
)

func TestConfigNew(t *testing.T) {
	{
		c := NewConfigFrom(map[string]interface{}{"foo": "bar"})
		assert.Check(t, c.Has("foo"), true)
	}
	{
		c, err := NewConfigFromJSON([]byte(`{"foo": 123}`))
		assert.Must(t, err)
		assert.Check(t, c.Has("foo"), true)
	}
	{
		c, err := NewConfigFromYAML([]byte(`name: good`))
		assert.Must(t, err)
		assert.Check(t, c.Has("name"), true)
	}
	{
		c, err := NewConfigFromFile("./example/example.yaml")
		assert.Must(t, err)
		assert.Check(t, c.Has("name"), true)
		assert.Check(t, c.Config("map").Has("key_one"), true)
		assert.Check(t, len(c.Value("map").Map()), 3)
		assert.Check(t, len(c.Value("list").List()), 4)
	}
	{
		c, err := NewConfigFromFile("./example/example.json")
		assert.Must(t, err)
		assert.Check(t, c.Has("name"), true)
		assert.Check(t, c.Config("map").Has("key_one"), true)
		assert.Check(t, len(c.Value("map").Map()), 3)
		assert.Check(t, len(c.Value("list").List()), 4)
	}
	if _, err := NewConfigFromFile("example/main.go"); err == nil {
		t.Fatal("expected error, got nothing")
	}
}

func TestConfigBasics(t *testing.T) {
	c := NewConfig()
	assert.Check(t, c.Has("foo"), false)
	c.SetDefault("foo", "bar")
	assert.Check(t, c.Has("foo"), true)
	c.Set("foo", "baz")
	assert.Check(t, c.String("foo"), "baz")
	c.SetDefault("bar", "foo")
	assert.Check(t, c.Has("bar"), true)

	ca := NewConfig()
	assert.Check(t, ca.Has("foo"), false)
	assert.Check(t, ca.Has("baz"), false)
	ca.Set("bar", "baz")
	assert.Check(t, ca.Has("bar"), true)
	assert.Check(t, ca.String("bar"), "baz")
	ca.Merge(c)
	assert.Check(t, ca.Has("foo"), true)
	assert.Check(t, ca.String("bar"), "foo")
}

func TestConfigMust(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expect error, got nothing")
		}
	}()
	c := NewConfig()
	c.Must("must")
}

func TestConfigValue(t *testing.T) {
	c := NewConfig()
	if x, ok := c.Value("xx").(Valuer); !ok {
		t.Fatalf("expected Valuer, got %#v\n", x)
	}
}

func TestConfigPattern(t *testing.T) {
	c := NewConfig()
	if x, ok := c.Pattern("xx").(Patterner); !ok {
		t.Fatalf("expected Patterner, got %#v\n", x)
	}
}

func TestConfigGetConfig(t *testing.T) {
	c := NewConfig()
	c.Set("unknwn_map", map[interface{}]interface{}{"foo": "bar"})
	c.Set("string_map", map[string]interface{}{"foo": "bar"})
	cc := NewConfig()
	cc.Set("foo", "bar")
	c.Set("cc", cc)

	assert.Check(t, c.Has("unknwn_map"), true)
	assert.Check(t, c.Has("string_map"), true)
	assert.Check(t, c.Has("cc"), true)
	assert.Check(t, c.Config("unknwn_map").Has("foo"), true)
	assert.Check(t, c.Config("string_map").Has("foo"), true)
	assert.Check(t, c.Config("cc").Has("foo"), true)

	ccc := c.Config("non")
	ccc.Set("test", "good")
	ccc.Set("good", "bad")
	ccc = c.Config("non")
	assert.Check(t, len(c.Value("non").Map()), 2)
}

func TestConfigGetString(t *testing.T) {
	c := NewConfig()
	c.Set("string", "foo")
	assert.Check(t, c.Has("string"), true)
	assert.Check(t, c.String("string"), "foo")
	assert.Check(t, c.StringOr("string", "bar"), "foo")
	assert.Check(t, c.Has("foo"), false)
	assert.Check(t, c.String("foo"), "")
	assert.Check(t, c.StringOr("foo", "bar"), "bar")
	c.Set("www", "mmm")
	res, ok := c.StringAnd("www", "^m")
	assert.Check(t, ok, true)
	assert.Check(t, res, "mmm")
	res, ok = c.StringAnd("www", "^w")
	assert.Check(t, ok, false)
	assert.Check(t, res, "")
	res, ok = c.StringAnd("mmm", "^w")
	assert.Check(t, ok, false)
	assert.Check(t, res, "")

	assert.Check(t, c.Has("test_env"), false)
	os.Setenv("test_env", "string")
	defer func() { os.Unsetenv("test_env") }()
	assert.Check(t, c.Has("test_env"), true)
	assert.Check(t, c.StringOr("test_env", "XXX"), "string")
}

func TestConfigGetBool(t *testing.T) {
	c := NewConfig()
	c.Set("bool", true)
	assert.Check(t, c.Has("bool"), true)
	assert.Check(t, c.Bool("bool"), true)
	assert.Check(t, c.BoolOr("bool", false), true)
	assert.Check(t, c.Has("non"), false)
	assert.Check(t, c.Bool("non"), false)
	assert.Check(t, c.BoolOr("non", true), true)

	assert.Check(t, c.Has("test_env"), false)
	os.Setenv("test_env", "1")
	defer func() { os.Unsetenv("test_env") }()
	assert.Check(t, c.Has("test_env"), true)
	assert.Check(t, c.BoolOr("test_env", false), true)
}

func TestConfigGetInt(t *testing.T) {
	c := NewConfig()
	c.Set("int", 33)
	assert.Check(t, c.Has("int"), true)
	assert.Check(t, c.Int("int"), 33)
	assert.Check(t, c.IntOr("int", 333), 33)
	res, ok := c.IntAnd("int", "N>3")
	assert.Check(t, ok, true)
	assert.Check(t, res, 33)

	assert.Check(t, c.Has("non"), false)
	assert.Check(t, c.Int("non"), 0)
	assert.Check(t, c.IntOr("non", 333), 333)
	res, ok = c.IntAnd("non", "N>3")
	assert.Check(t, ok, false)
	assert.Check(t, res, 0)

	assert.Check(t, c.Has("test_env"), false)
	os.Setenv("test_env", "1111")
	defer func() { os.Unsetenv("test_env") }()
	assert.Check(t, c.Has("test_env"), true)
	assert.Check(t, c.IntOr("test_env", 11), 1111)
}

func TestConfigGetFloat(t *testing.T) {
	c := NewConfig()
	c.Set("float", 333.3)
	assert.Check(t, c.Has("float"), true)
	assert.Check(t, c.Float("float"), 333.3)
	assert.Check(t, c.FloatOr("float", 33.33), 333.3)
	res, ok := c.FloatAnd("float", "N*10==3333")
	assert.Check(t, ok, true)
	assert.Check(t, res, 333.3)

	assert.Check(t, c.Has("non"), false)
	assert.Check(t, c.Float("non"), 0.0)
	assert.Check(t, c.FloatOr("non", 33.33), 33.33)
	res, ok = c.FloatAnd("non", "N>0")
	assert.Check(t, ok, false)
	assert.Check(t, res, 0.0)

	assert.Check(t, c.Has("test_env"), false)
	os.Setenv("test_env", "11.11")
	defer func() { os.Unsetenv("test_env") }()
	assert.Check(t, c.Has("test_env"), true)
	assert.Check(t, c.FloatOr("test_env", 1.1), 11.11)
}

func TestConfigGetDuration(t *testing.T) {
	c := NewConfig()
	c.Set("t", 300)
	assert.Check(t, c.Has("t"), true)
	assert.Check(t, c.Duration("t"), time.Duration(300))
	assert.Check(t, c.DurationOr("t", 333), time.Duration(300))
	assert.Check(t, c.Has("tt"), false)
	assert.Check(t, c.Duration("tt"), time.Duration(0))
	assert.Check(t, c.DurationOr("tt", 333), time.Duration(333))

	assert.Check(t, c.Has("test_env"), false)
	os.Setenv("test_env", "1111")
	defer func() { os.Unsetenv("test_env") }()
	assert.Check(t, c.Has("test_env"), true)
	assert.Check(t, c.DurationOr("test_env", 11), time.Duration(1111))
}
