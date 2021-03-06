package libuconf_test

import (
	"testing"

	. "toast.cafe/x/libuconf"

	"github.com/stretchr/testify/assert"
)

func TestFindLongFlag(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		matrix = []struct {
			in  *string
			out FlagOpt
		}{
			{o.String("aflag", 0, "aval", "ahelp"), o.FindLongFlag("aflag")},
			{o.String("bflag", 0, "bval", "bhelp"), o.FindLongFlag("bflag")},
		}
	)
	for _, v := range matrix {
		vv, ok := v.out.(Getter)
		assert.True(ok)
		assert.Equal(*v.in, vv.Get().(string))
	}
}

func TestFindShortFlag(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		matrix = []struct {
			in  *string
			out FlagOpt
		}{
			{o.String("aflag", 'a', "aval", "ahelp"), o.FindShortFlag('a')},
			{o.String("bflag", 'b', "bval", "bhelp"), o.FindShortFlag('b')},
		}
	)
	for _, v := range matrix {
		vv, ok := v.out.(Getter)
		assert.True(ok)
		assert.Equal(*v.in, vv.Get().(string))
	}
}

func TestParseLongFlags(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		matrix = []struct {
			flag string
			val  string
			oval *string
		}{
			{"foo", "fooval", nil},
			{"bar", "barval", nil},
		}
	)
	for _, v := range matrix {
		v.oval = o.String(v.flag, 0, "", v.flag+"help")
		flag := "--" + v.flag

		// --a.b.c=val
		err := o.ParseFlags([]string{
			flag + "=" + v.val,
		})
		assert.Nil(err)
		assert.Equal(v.val, *v.oval)

		*v.oval = "reset"

		// --a.b.c val
		err = o.ParseFlags([]string{
			flag, v.val,
		})
		assert.Nil(err)
		assert.Equal(v.val, *v.oval)

		// --a.b.c --a.b.c: expect fail (string, we test bool in its own thing)
		err = o.ParseFlags([]string{
			flag, flag,
		})
		assert.NotNil(err)

		assert.Nil(o.Args)
	}
}

func TestParseShortFlags(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		matrix = []struct {
			flag rune
			val  string
			oval *string
		}{
			{'a', "aval", nil},
			{'b', "bval", nil},
		}
	)
	for _, v := range matrix {
		var (
			sflag = string(v.flag)
			flag  = "-" + sflag
		)
		v.oval = o.String(sflag, v.flag, "", sflag+"help")

		// -a val
		err := o.ParseFlags([]string{
			flag, v.val,
		})
		assert.Nil(err)
		assert.Equal(v.val, *v.oval)

		*v.oval = "reset"

		// -aval
		err = o.ParseFlags([]string{
			flag + v.val,
		})
		assert.Nil(err)
		assert.Equal(v.val, *v.oval)

		// -a -a
		err = o.ParseFlags([]string{
			flag, flag,
		})
		assert.NotNil(err)

		assert.Nil(o.Args)
	}
}

// ---- we already tested general parsing
// now we want to test (for bools):
// a) last flag being set to true
// b) flag followed by value being auto-true
// c) flag followed by flag being auto-true

func TestParseLongFlagsBool(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		a      = o.Bool("aa", 0, false, "aahelp")
		b      = o.Bool("bb", 0, false, "bbhelp")
	)

	err := o.ParseFlags([]string{
		"--aa",
	})
	assert.Nil(err)
	assert.Equal(true, *a)

	err = o.ParseFlags([]string{
		"--bb",
	})
	assert.Nil(err)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	err = o.ParseFlags([]string{
		"--aa", "--bb",
	})
	assert.Nil(err)
	assert.Equal(true, *a)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	err = o.ParseFlags([]string{
		"--aa", "arbitrary", "--bb", "arbitrary",
	})
	assert.Nil(err)
	assert.Equal(true, *a)
	assert.Equal(true, *b)
}

func TestParseShortFLagsBool(t *testing.T) {
	var (
		assert = assert.New(t)
		o      = &OptionSet{AppName: "test"}
		a      = o.Bool("aa", 'a', false, "aahelp")
		b      = o.Bool("bb", 'b', false, "bbhelp")
	)

	err := o.ParseFlags([]string{
		"-a",
	})
	assert.Nil(err)
	assert.Equal(true, *a)

	err = o.ParseFlags([]string{
		"-b",
	})
	assert.Nil(err)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	err = o.ParseFlags([]string{
		"-ab",
	})
	assert.Nil(err)
	assert.Equal(true, *a)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	err = o.ParseFlags([]string{
		"-a", "-b",
	})
	assert.Nil(err)
	assert.Equal(true, *a)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	err = o.ParseFlags([]string{
		"-a", "arbitrary", "-b", "arbitrary",
	})
	assert.Nil(err)
	assert.Equal(true, *a)
	assert.Equal(true, *b)

	// reset
	*a = false
	*b = false

	// *THIS* should fail
	err = o.ParseFlags([]string{
		"-aarbitrary",
	})
	assert.NotNil(err)
}
