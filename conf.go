// conf.go: Simple config file reader.  Config items are in the form
// foo=bar and comments can be created by placing a # at the start of
// the line.
//
// Copyright (c) 2011-2013 CloudFlare, Inc.

package conf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ConfigFileRaw is a copy of the config file

var ConfigFileRaw []byte

// The parsed contents of the config file

type parameter struct {
	s string // The parsed parameter string value
	r bool   // Set to true if this parameter had been read
	l int    // Line number at which this was found
}

// Config is a copy of the configuration file turned into a map

type Config struct {
	v map[string]parameter
}

// ReadConfigFile reads an entire config file and return a Config
// structure that can be used to extract single values.
func ReadConfigFile(file string) (c *Config, err error) {
	c = new(Config)
	c.v = make(map[string]parameter)
	l := 0

	var f *os.File
	if f, err = os.Open(file); err == nil {
		defer f.Close()
		
		scanner := bufio.NewScanner(f)
		
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			l += 1

			if len(line) > 0 && line[0] != '#' {
				parts := strings.Split(line, "=")

				if len(parts) != 2 {
					err = fmt.Errorf("Config line %d invalid: %s (missing =)", l,
						line)
					return
				}

				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])

				if k == "" {
					err = fmt.Errorf("Config line %d invalid: %s (missing key)", l,
						line)
					return
				}

				if _, found := c.v[k]; !found {
					c.v[k] = parameter{v, false, l}
				} else {
					err = fmt.Errorf("Config line %d invalid: %s (repeated parameter)", l,
						line)
					return
				}
			}
		}

		err = scanner.Err()
	}

	return
}

// GetString reads a string value from the parsed config file and
// return it (or if it is missing then return the default value passed
// in d)
func (c *Config) GetString(k string, d string) (v string) {
	s, present := c.v[k]
	if !present {
		v = d
	} else {
		v = s.s
		s.r = true
		c.v[k] = s
	}
	return
}

// GetUint reads an unsigned integer from the parsed config file and
// return it (or if it is missing then return the default value passed
// in d)
func (c *Config) GetUint(k string, d uint) (v uint) {
	if s, present := c.v[k]; !present {
		v = d
	} else {
		s.r = true
		c.v[k] = s
		if iv, err := strconv.Atoi(s.s); err != nil {
			v = d
		} else {
			v = uint(iv)
		}
	}

	return
}

// CheckUnread checks to see if any of the parameters in the file
// have not been read and returns a string containing those that have
// not been read.
func (c *Config) CheckUnread() (s string) {
	for k, v := range c.v {
		if !v.r {
			s += fmt.Sprintf("%s (%d) ", k, v.l)
		}
	}
	return
}
