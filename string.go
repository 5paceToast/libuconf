package libuconf

import "fmt"

// ensure interface compliance
var (
	_ EnvOpt  = &StringOpt{}
	_ FlagOpt = &StringOpt{}
	_ Getter  = &StringOpt{}
	_ Setter  = &StringOpt{}
	_ TomlOpt = &StringOpt{}
)

// StringOpt represents a string Option
type StringOpt struct {
	help  string
	name  string
	sname rune
	val   *string
}

// ---- integration with OptionSet

// StringVar adds a StringOpt to the OptionSet
func (o *OptionSet) StringVar(out *string, name, val, help string) {
	o.ShortStringVar(out, name, 0, help)
}

// String adds a StringOpt to the OptionSet
func (o *OptionSet) String(name, val, help string) *string {
	return o.ShortString(name, 0, val, help)
}

// ShortStringVar adds a StringOpt to the OptionSet
func (o *OptionSet) ShortStringVar(out *string, name string, sname rune, help string) {
	sopt := &StringOpt{help, name, sname, out}
	o.Var(sopt)
}

// ShortString adds a StringOpt to the Option Set
func (o *OptionSet) ShortString(name string, sname rune, val, help string) *string {
	out := &val
	o.ShortStringVar(&val, name, sname, help)
	return out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (s *StringOpt) Env() string {
	return env(s)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*StringOpt) Bool() bool {
	return false
}

// Flag returns the long-form flag for this option
func (s *StringOpt) Flag() string {
	return s.name
}

// Help returns the help string for this option
func (s *StringOpt) Help() string {
	return s.help
}

// ShortFlag returns the short-form flag for this option
func (s *StringOpt) ShortFlag() rune {
	return s.sname
}

// ---- Getter

// Get returns the internal value
func (s *StringOpt) Get() interface{} {
	return *s.val
}

// ---- Setter

// Set sets this option's value
func (s *StringOpt) Set(vv interface{}) error {
	switch v := vv.(type) {
	case string:
		*s.val = v
	default:
		*s.val = fmt.Sprint(vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (s *StringOpt) Toml() string {
	return _toml(s)
}
