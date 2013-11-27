// conf_test.go: test suite for conf.go
//
// Copyright (c) 2013 CloudFlare, Inc.

package conf

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Fail()
	}
}

// makeFile creates a file with the given contents, returns the file
// name. Returns an empty string on error.
func makeFile(c string) (fn string) {
	if f, err := ioutil.TempFile(os.TempDir(), "conf"); err == nil {
		f.WriteString(c)
		f.Close()
		return f.Name()
	} else {
		return ""
	}
}

func TestEmptyConfig(t *testing.T) {
	f := makeFile("")
	c, err := ReadConfigFile(f)
	assert(t, c != nil)
	assert(t, err == nil)
	assert(t, c.CheckUnread() == "")
}

func TestSimpleConfig(t *testing.T) {

	// Note: specific things being tested here:
	//
	// foo=1 (a normal line)
	// \n\n completely empty lines
	// \n \n lines with just whitespace
	// # comment lines that have comments in them
	//   bar =  baz  lines with embedded whitespace and not trailing \n

	f := makeFile("foo=1\n \n    \n# comment\n\n\n  bar =  baz  ")
	c, err := ReadConfigFile(f)
	assert(t, c != nil)
	assert(t, err == nil)
	assert(t, c.CheckUnread() != "")
	assert(t, c.GetString("foo", "FOO") == "1")
	assert(t, c.GetUint("foo", 2) == 1)
	assert(t, c.GetString("bar", "BAR") == "baz")
	assert(t, c.GetUint("bar", 2) == 2)
	assert(t, c.CheckUnread() == "")
	assert(t, c.GetString("baz", "BAZ") == "BAZ")
	assert(t, c.GetUint("baz", 3) == 3)
	assert(t, c.GetString("bam", "") == "")
	assert(t, c.GetUint("bam", 0) == 0)

	// Check handling of empty value

	f = makeFile("foo=1\nbar=")
	c, err = ReadConfigFile(f)
	assert(t, c != nil)
	assert(t, err == nil)
	assert(t, c.GetString("foo", "FOO") == "1")
	assert(t, c.GetUint("foo", 2) == 1)
	assert(t, c.GetString("bar", "BAR") == "")
	assert(t, c.GetUint("bar", 2) == 2)
}

func TestCheckUnread(t *testing.T) {
	f := makeFile("foo=1\nbar=baz")
	c, err := ReadConfigFile(f)
	assert(t, c != nil)
	assert(t, err == nil)
	assert(t, c.CheckUnread() != "")
	assert(t, strings.Contains(c.CheckUnread(), "foo (1)"))
	assert(t, strings.Contains(c.CheckUnread(), "bar (2)"))
	assert(t, c.GetString("bar", "BAR") == "baz")
	assert(t, strings.Contains(c.CheckUnread(), "foo (1)"))
	assert(t, !strings.Contains(c.CheckUnread(), "bar (2)"))
	assert(t, c.GetString("foo", "FOO") == "1")
	assert(t, !strings.Contains(c.CheckUnread(), "foo (1)"))
	assert(t, !strings.Contains(c.CheckUnread(), "bar (2)"))
}

func TestCheckErrors(t *testing.T) {
	f := makeFile("foo=1\nbaz")
	_, err := ReadConfigFile(f)
	assert(t, err != nil)

	f = makeFile("foo=1\n=baz")
	_, err = ReadConfigFile(f)
	assert(t, err != nil)

	f = makeFile("foo=1\nfoo=2")
	_, err = ReadConfigFile(f)
	assert(t, err != nil)
}
