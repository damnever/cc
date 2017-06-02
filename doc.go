// Package cc is a very flexible configuration management library for humans,
// which is easy to use and support YAML and JSON only.
//
//
// Usage
//
//		c, _ := cc.NewConfigFromFile("./example/example.yaml")  // file must has extension
//		_ := c.MergeFromFile("./example/example.json") // do not ignore the errors
//
//		c.Must("name")  // panic if not found
//		c.String("name")
//
//		cc := c.Config("map")
//		cc.Bool("key_one")
//
//		list := c.Value("list").List()
//		list[1].Int()
//
//		// environment variables
//		os.Setenv("float_env", "11.11")
//		c.Float("float_env")
//
//      // flags
//      flag.Int("flag", 33, "usage")
//      c.Int("flag")
//
// The priorities: flags > environment variables > normal configs
//
//
// Default Configs
//
// We may write the code like this:
//
//		name := "default"
//		if c.Has("name") {
//			name = c.String("name")  // or panic
//		}
//
// Now, we can write code like this:
//
//		name := c.StringOr("name", "cc")  // or c.Must("name")
//		b := c.BoolOr("bool", true)
//		f := c.FloatOr("float", 3.14)
//		i := c.IntOr("int", 33)
//
//
// Pattern and Validation
//
// If you want to check string value whether it is matched by regexp:
//
//		s, ok := c.StringAnd("name", "^c")
//
// Or, the make the string value as a pattern:
//
//		p := c.Pattern("pattern_key_name")
//		ok := p.ValidateString("a string")
//
// For int(time.Duration) and float, cc use if-like condition to do similar work.
// Assume we have `threhold: "N>=30&&N<=80"` in config file, we can use it like this:
//
//		p := c.Pattern("threhold")
//		ok := p.ValidateInt(40)  // or ValidateFloat
//
// Or, using a pattern to validate the number:
//
//		ni, ok := c.IntAnd("int_key", "N>50")
//		nf, ok := c.FloatAnd("float_key", "N/100>=0.3")
//      // or, given a default value
//      ni = c.IntAndOr("int_key", "N>50", 51)
//      nf = c.FloatAndOr("int_key", "N/100>=0.3", 40)
//      d := c.DurationAndOr("duration", "N>20&&N<=100", 50)
//
// NOTE: bit operation is not supported.
package cc
